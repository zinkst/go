package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
    "path"
    "path/filepath"
    "bufio"
    "os"
    "strings"
    "log"    
)    

// programConfig input from yaml config file
type ProgramConfig struct {
    FileFilter string `yaml:"fileFilter"`
    SrcDirName string `yaml:"srcDirName"`
    TgtDirName string `yaml:"tgtDirName"`
}

var SSHGlobals map[string]string

type SSHEntry struct {
	Name string
	Port int
}

var programConfig  ProgramConfig


func main() {
	fmt.Println("CreateSSHConfig.go")
	readConfigFile()
	SSHGlobals=make(map[string]string)
	SSHGlobals["proxyHostname"] = "bosh-cli-bluemix.rtp.raleigh.ibm.com"
	SSHGlobals["User"] = "Stefan.Zink@de.ibm.com"
	//SSHGlobals["StrictHostKeyChecking"] = "no"
	SSHGlobals["ForwardAgent"] = "yes"
	for key, value := range SSHGlobals {
	    fmt.Println("Key:", key, "Value:", value)
	}
	getBoshcliFilesInDir(programConfig.SrcDirName)
	readBoshcliSHFile( path.Join( programConfig.SrcDirName , "boshcli_ukpostoffice_1.sh"))	
}

func getBoshcliFilesInDir (boshClisSrcDir string) {
	log.Print("\nfilepattern: " + path.Join(boshClisSrcDir, programConfig.FileFilter ))
	files, err := filepath.Glob(path.Join(boshClisSrcDir, programConfig.FileFilter ))	
//	files, err := ioutil.ReadDir(boshClisSrcDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files)		
	for _, file := range files {
		fmt.Println(strings.SplitAfter(filepath.Base(file),filepath.Ext(file)))
	}
}

func readBoshcliSHFile (shFileName string) {
	f, err := os.Open(shFileName)
    if err != nil {
        fmt.Fprintf(os.Stderr, " Filename: %v \n Error: %v\n",shFileName, err)
    }
    input := bufio.NewScanner(f) 
    for input.Scan() {
    	curStr := input.Text()    	
        if (strings.Contains(curStr, "ssh -p")) {
	        fmt.Println(curStr)
	        idx := strings.Index(curStr, "-p")
	        port := curStr[idx+3:idx+8]
	        fmt.Printf("Port: %s", port)
        }
    }       
    f.Close()    
}


func readConfigFile() {
	filename, _ := filepath.Abs("/home/zinks/workspace/CreateSSHConfig/src/github.com/zinkst/go/CreateSSHConfig/CreateSSHConfig.yml")
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
