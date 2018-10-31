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

// ProgramConfig input from yaml config file
type ProgramConfig struct {
	FileFilter   string `yaml:"fileFilter"`
	SrcDirName   string `yaml:"srcDirName"`
	TgtDirName   string `yaml:"tgtDirName"`
	ConfigPrefix string `yaml:"configPrefix"`
	ConfigSuffix string `yaml:"configSuffix"`
}

// SSHEntry das
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
	//output += fmt.Sprintf("  Hostname %s\n", s.Jumpbox)
	output += fmt.Sprintf("  Port %v\n", s.Port)
	//output += fmt.Sprintf("  StrictHostKeyChecking no\n" )
	//output += `  ProxyCommand ssh -q -W %h:%p w3-boshcli`
	return output
}

var sshGlobals map[string]string
var programConfig ProgramConfig
var sshEntries []SSHEntry

func main() {
	fmt.Println("CreateSSHConfig.go")
	fmt.Println("Test Build from vscode 2")
	readConfigFile()
	sshGlobals = make(map[string]string)
	sshGlobals["proxyHostname"] = "bosh-cli-bluemix.rtp.raleigh.ibm.com"
	sshGlobals["User"] = "Stefan.Zink@de.ibm.com"
	//sshGlobals["StrictHostKeyChecking"] = "no"
	sshGlobals["ForwardAgent"] = "yes"
	for key, value := range sshGlobals {
		fmt.Println("Key:", key, "Value:", value)
	}
	getBoshcliFilesInDir(programConfig.SrcDirName)
	CreateSSHConfigFile()

}

// CreateSSHConfigFile cd
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
	//f.WriteString(generateSSHConfigHeader() )
	f.WriteString(programConfig.ConfigPrefix)
	for _, curEntry := range sshEntries {
		f.WriteString(curEntry.appendToConfig())
		f.WriteString("\n")
	}
	f.WriteString(programConfig.ConfigSuffix)
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
		sshEntryName := strings.Split(filepath.Base(file), filepath.Ext(file))
		newSSHEntryName := strings.Replace(sshEntryName[0], "boshcli_", "bm+", 1)
		log.Print("EntryName =", sshEntryName[0], " - replaced EntryName=", newSSHEntryName)
		newSSHEntry := SSHEntry{Name: newSSHEntryName}
		readBoshcliSHFile(file, &newSSHEntry)
		sshEntries = append(sshEntries, newSSHEntry)
	}
	fmt.Printf("Number of entries found %v \n", len(sshEntries))
}

func readBoshcliSHFile(shFileName string, newSSHEntry *SSHEntry) {
	f, err := os.Open(shFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Filename: %v \n Error: %v\n", shFileName, err)
	}
	input := bufio.NewScanner(f)
	for input.Scan() {
		curStr := input.Text()
		if strings.Contains(curStr, "ssh -o ServerAliveInterval=60 -o ServerAliveCountMax=15") {
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
	wrkdir, err := os.Getwd()
	fmt.Printf("workingDir: %s \n", wrkdir)
	filename, _ := filepath.Abs(path.Join(wrkdir, "CreateSSHConfig.yml"))
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
