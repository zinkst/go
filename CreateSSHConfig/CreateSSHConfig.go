package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// programConfig input from yaml config file
type ProgramConfig struct {
	FileFilter string `yaml:"fileFilter"`
	SrcDirName string `yaml:"srcDirName"`
	TgtDirName string `yaml:"tgtDirName"`
}

type SSHEntry struct {
	Name    string
	Port    int
	Jumpbox string
}

func (s SSHEntry) print() string {
	output := fmt.Sprintf("Name: %s, Port: %v, Jumpbox: %s \n", s.Name, s.Port, s.Jumpbox)
	return output
}

func (s SSHEntry) appendToConfig() string {
	output := fmt.Sprintf("Host %s\n", s.Name)
	output += fmt.Sprintf("  Hostname %s\n", s.Jumpbox)
	output += fmt.Sprintf("  Port %v\n", s.Port)
	output += fmt.Sprintf("  StrictHostKeyChecking no\n")
	output += `  ProxyCommand ssh -q -W %h:%p w3-boshcli`
	output += "\n"
	return output
}

var SSHGlobals map[string]string
var programConfig ProgramConfig
var SSHEntries []SSHEntry

func main() {
	fmt.Println("CreateSSHConfig.go")
	readConfigFile()
	SSHGlobals = make(map[string]string)
	SSHGlobals["proxyHostname"] = "bosh-cli-bluemix.rtp.raleigh.ibm.com"
	SSHGlobals["User"] = "Stefan.Zink@de.ibm.com"
	//SSHGlobals["StrictHostKeyChecking"] = "no"
	SSHGlobals["ForwardAgent"] = "yes"
	for key, value := range SSHGlobals {
		fmt.Println("Key:", key, "Value:", value)
	}
	getBoshcliFilesInDir(programConfig.SrcDirName)
	CreateSSHConfigFile()

}

func generateSSHConfigHeader() string {
	output := `
Host *
  StrictHostKeyChecking no
  ForwardX11    yes

Host w3-boshcli
  Hostname bosh-cli-bluemix.rtp.raleigh.ibm.com
  User Stefan.Zink@de.ibm.com
  StrictHostKeyChecking no
  ForwardAgent yes
  
`
	return output
}

func CreateSSHConfigFile() {
	t := time.Now()
	formattedDate := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	configFileName := path.Join(programConfig.TgtDirName, "sshConfig_"+formattedDate)
	f, err := os.Create(configFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Filename: %v \n Error: %v\n", configFileName, err)
	}
	err = os.Chmod(configFileName, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Filename: %v \n Error changing permissions: %v\n", configFileName, err)
	}
	f.WriteString(generateSSHConfigHeader())
	for _, curEntry := range SSHEntries {
		f.WriteString(curEntry.appendToConfig())
		f.WriteString("\n")
	}
	f.Close()
}

func getBoshcliFilesInDir(boshClisSrcDir string) {
	log.Print("\nfilepattern: " + path.Join(boshClisSrcDir, programConfig.FileFilter))
	files, err := filepath.Glob(path.Join(boshClisSrcDir, programConfig.FileFilter))
	//	files, err := ioutil.ReadDir(boshClisSrcDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files)
	for _, file := range files {
		newEntryName := strings.Split(filepath.Base(file), filepath.Ext(file))
		//fmt.Println("newEntryName =" , newEntryName[0] )
		newSSHEntry := SSHEntry{Name: newEntryName[0]}
		readBoshcliSHFile(file, &newSSHEntry)
		SSHEntries = append(SSHEntries, newSSHEntry)
	}
	fmt.Printf("Number of entries found %v \n", len(SSHEntries))
}

func readBoshcliSHFile(shFileName string, newSSHEntry *SSHEntry) {
	f, err := os.Open(shFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Filename: %v \n Error: %v\n", shFileName, err)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		curStr := input.Text()
		if strings.Contains(curStr, "ssh -p") {
			//fmt.Println(curStr)
			idx := strings.Index(curStr, "-p")
			port := curStr[idx+3 : idx+8]
			//fmt.Printf("Port: %s", port)
			newSSHEntry.Port, err = strconv.Atoi(port)
		}
		if strings.Contains(curStr, "JUMPBOX=") {
			//fmt.Println(curStr)
			idx := strings.Index(curStr, "OX=")
			newSSHEntry.Jumpbox = curStr[idx+3:]
		}

	}
	f.Close()
	fmt.Print(newSSHEntry.print())

}

func readConfigFile() {
	// filename, _ := filepath.Abs("/home/zinks/workspace/CreateSSHConfig/src/github.com/zinkst/go/CreateSSHConfig/CreateSSHConfig.yml")
	filename := "CreateSSHConfig.yml"
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	fmt.Printf("---- yamlFile: %s -------\n%v", filename, string(yamlFile))

	err = yaml.Unmarshal(yamlFile, &programConfig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("FileFilter: %v\n", programConfig.FileFilter)
	fmt.Printf("Value: %v\n", programConfig)
}
