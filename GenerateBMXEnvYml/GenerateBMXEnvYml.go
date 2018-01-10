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
    "strconv"
    "time"
    "github.com/smallfish/simpleyaml"
//	  "github.com/smallfish/simpleyaml/helper/util"
)    

// programConfig input from yaml config file
type ProgramConfig struct {
    FileFilter string `yaml:"fileFilter"`
    BoshcliCmdsSrcDirName string `yaml:"boshcliCmdsSrcDirName"`
    TgtDirName string `yaml:"tgtDirName"`
    ConfigPrefix string `yaml:"configPrefix"`    
    ConfigSuffix string `yaml:"configSuffix"`    
    JmlSrcDirName string `yaml:"jmlSrcDirName"`
}


type SSHEntry struct {
	Name string
	Port int
	Jumpbox string
}

type JMLEntry struct {
	Name string
	Port int
}


func (s SSHEntry) print() string {
	output := fmt.Sprintf("Name: %s, Port: %v, Jumpbox: %s \n" , s.Name, s.Port, s.Jumpbox) 
	return output
}

func (s JMLEntry) print() string {
	output := fmt.Sprintf("Name: %s, Port: %v \n" , s.Name, s.Port) 
	return output
}


func (s SSHEntry) appendToConfig() string {
	output := fmt.Sprintf("Host %s\n" , s.Name )
	//output += fmt.Sprintf("  Hostname %s\n", s.Jumpbox)
	output += fmt.Sprintf("  Port %v\n", s.Port)
	//output += fmt.Sprintf("  StrictHostKeyChecking no\n" )
	//output += `  ProxyCommand ssh -q -W %h:%p w3-boshcli`
	return output
}



var SSHGlobals map[string]string
var programConfig  ProgramConfig
var SSHEntries [] SSHEntry
var JMLEntries [] JMLEntry

func main() {
	fmt.Println("CreateBMXEnvYml.go")
	readConfigFile()
	SSHGlobals=make(map[string]string)
	SSHGlobals["proxyHostname"] = "bosh-cli-bluemix.rtp.raleigh.ibm.com"
	SSHGlobals["User"] = "Stefan.Zink@de.ibm.com"
	//SSHGlobals["StrictHostKeyChecking"] = "no"
	SSHGlobals["ForwardAgent"] = "yes"
	for key, value := range SSHGlobals {
	    fmt.Println("Key:", key, "Value:", value)
	}
	//GetBoshcliFilesInDir( programConfig.BoshcliCmdsSrcDirName)
	getJMLEntriesInDir()
	//CreateBMXEnvYmlFile()
	
}

func CreateBMXEnvYmlFile() {
	t := time.Now()
	formattedDate := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	configFileName := path.Join(programConfig.TgtDirName, "sshConfig_" + formattedDate )
	f, err := os.Create(configFileName)
	if err != nil {
        fmt.Fprintf(os.Stderr, " Filename: %v \n Error: %v\n",configFileName, err)
    }
    err = os.Chmod(configFileName, 0600)
    if err != nil {
        fmt.Fprintf(os.Stderr, " Filename: %v \n Error changing permissions: %v\n",configFileName, err)
    }
    //f.WriteString(generateSSHConfigHeader() )
    f.WriteString(programConfig.ConfigPrefix)
    for _, curEntry := range SSHEntries {
    	f.WriteString(curEntry.appendToConfig())
    	f.WriteString("\n")
    }    
    f.WriteString(programConfig.ConfigSuffix)
    f.Close()
}

func GetBoshcliFilesInDir (boshClisSrcDir string) {
	log.Print("\nfilepattern: " + path.Join(boshClisSrcDir, programConfig.FileFilter ))
	files, err := filepath.Glob(path.Join(boshClisSrcDir, programConfig.FileFilter ))	
//	files, err := ioutil.ReadDir(boshClisSrcDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files)		
	for _, file := range files {
		sshEntryName := strings.Split(filepath.Base(file),filepath.Ext(file))
		newSshEntryName := strings.Replace(sshEntryName[0],"boshcli_","bm+",1)
		log.Print("EntryName =" , sshEntryName[0]," - replaced EntryName=" ,  newSshEntryName )
		newSSHEntry := SSHEntry{Name: newSshEntryName}
		SSHEntry.ReadBoshcliSHFile(newSSHEntry,file)
		SSHEntries = append(SSHEntries,newSSHEntry)
	}
	fmt.Printf("Number of entries found %v \n", len(SSHEntries))
}

func (s SSHEntry) ReadBoshcliSHFile (shFileName string) {
	f, err := os.Open(shFileName)
    if err != nil {
        fmt.Fprintf(os.Stderr, " Filename: %v \n Error: %v\n",shFileName, err)
    }
    input := bufio.NewScanner(f) 
    for input.Scan() {
    	curStr := input.Text()    	
        if (strings.Contains(curStr, "ssh -o ServerAliveInterval=60 -o ServerAliveCountMax=15")) {
	        //fmt.Println(curStr)
	        idx := strings.Index(curStr, "-p")
	        port := curStr[idx+3:idx+8]
	        //fmt.Printf("Port: %s", port)
	        s.Port, err = strconv.Atoi(port)
        }
        if (strings.Contains(curStr, "JUMPBOX=")) {
	        //fmt.Println(curStr)
	        idx := strings.Index(curStr, "OX=")
	        s.Jumpbox = curStr[idx+3:]
        }
        
    }       
    f.Close()
    fmt.Print(s.print()) 
       
}


func readConfigFile() {
	wrkdir, err := os.Getwd()
	fmt.Printf("workingDir: %s \n", wrkdir)
	filename, _ := filepath.Abs(path.Join(wrkdir,"GenerateBMXEnvYml.yml"))
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

func getJMLEntriesInDir () {
	log.Print("\nfilepattern: " + path.Join(programConfig.JmlSrcDirName, "shared" ))
	files, err := filepath.Glob(path.Join(programConfig.JmlSrcDirName, "shared","*" )) 	
//	files, err := ioutil.ReadDir(boshClisSrcDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(files)		
	for _, file := range files {
		//doctor_ssh_tunnel_port: 60100
		jmlFileNames, err := filepath.Glob(path.Join(file,"jml_config*"))
		if err != nil {
			log.Print(err)
		} else {	
			if len(jmlFileNames) == 0 {
				log.Print("no File with pattern jml_config* found in path " + file)
			} else {
				log.Print("jmlFileName =" , jmlFileNames[0] )
				jmlFile,err := ioutil.ReadFile(jmlFileNames[0])
				jmlFileYAML,err  := simpleyaml.NewYaml(jmlFile)
				if err != nil {
					log.Print("init yaml failed")
				}
				doctorTunnelPort, err := jmlFileYAML.Get("doctor_ssh_tunnel_port").Int()
				if err != nil {
					log.Print("get yaml key doctor_ssh_tunnel_port failed ")
				}
			    newJMLEntry := JMLEntry{Name: path.Base(file)}
			    newJMLEntry.Port = doctorTunnelPort
			    log.Print(newJMLEntry.print())
				JMLEntries = append(JMLEntries,newJMLEntry)
			}
		}	
	}
	fmt.Printf("Number of entries found %v \n", len(JMLEntries))
}
