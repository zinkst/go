package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v3"
)

// ProgramConfig input from yaml config file
type ProgramConfig struct {
	LogLevel     string `yaml:"logLevel" validate:"required,eq=info|eq=trace|eq=debug|eq=error"`
	FileFilter   string `yaml:"fileFilter" validate:"required,eq=boshcli*.sh"`
	SrcDirName   string `yaml:"srcDirName" validate:"required"`
	TgtDirName   string `yaml:"tgtDirName" validate:"required"`
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
	output += fmt.Sprintf("  Hostname %s\n", s.Jumpbox)
	output += fmt.Sprintf("  Port %v\n", s.Port)
	output += fmt.Sprintf("  StrictHostKeyChecking no\n")
	output += `  ProxyCommand ssh -q -W %h:%p w3-boshcli`
	output += "\n"
	return output
}

var sshGlobals map[string]string
var sshEntries []SSHEntry
var programConfig ProgramConfig
var log = logrus.New()

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func initLogger() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter) //default
	// log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	// log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	// log.Out = os.Stdout
	log.Level = logrus.InfoLevel
}

func main() {

	// fmt.Println("CreateSSHConfig.go")
	initLogger()
	readConfigFile()
	logLevel, err := logrus.ParseLevel(programConfig.LogLevel)
	if err != nil {
		panic(err)
	}
	log.SetLevel(logLevel)
	sshGlobals = make(map[string]string)
	sshGlobals["proxyHostname"] = "bosh-cli-bluemix.rtp.raleigh.ibm.com"
	sshGlobals["User"] = "Stefan.Zink@de.ibm.com"
	//SSHGlobals["StrictHostKeyChecking"] = "no"
	sshGlobals["ForwardAgent"] = "yes"
	for key, value := range sshGlobals {
		log.Debugln("Key:", key, "Value:", value)
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

// CreateSSHConfigFile cd
func CreateSSHConfigFile() {
	t := time.Now()
	formattedDate := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	configFileName := path.Join(programConfig.TgtDirName, "sshConfig_"+formattedDate)
	f, err := os.Create(configFileName)
	if err != nil {
		log.Errorf(" Filename: %v \n Error: %v\n", configFileName, err)
	}
	err = os.Chmod(configFileName, 0600)
	if err != nil {
		log.Errorf(" Filename: %v \n Error changing permissions: %v\n", configFileName, err)
	}
	f.WriteString(generateSSHConfigHeader())
	for _, curEntry := range sshEntries {
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
	log.Debugln(files)
	for _, file := range files {
		newEntryName := strings.Split(filepath.Base(file), filepath.Ext(file))
		//fmt.Println("newEntryName =" , newEntryName[0] )
		newSSHEntry := SSHEntry{Name: newEntryName[0]}
		readBoshcliSHFile(file, &newSSHEntry)
		sshEntries = append(sshEntries, newSSHEntry)
	}
	log.Infof("Number of entries found %v \n", len(sshEntries))
}

func readBoshcliSHFile(shFileName string, newSSHEntry *SSHEntry) {
	f, err := os.Open(shFileName)
	if err != nil {
		log.Errorf(" Filename: %v \n Error: %v\n", shFileName, err)
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
	log.Debug(newSSHEntry.print())

}

func readConfigFile() {
	// filename, _ := filepath.Abs("/home/zinks/workspace/CreateSSHConfig/src/github.com/zinkst/go/CreateSSHConfig/CreateSSHConfig.yml")
	filename := "CreateSSHConfig.yml"
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	log.Infof("---- yamlFile: %s -------\n", filename)
	log.Debugf("Content:\n%v", string(yamlFile))

	err = yaml.Unmarshal(yamlFile, &programConfig)
	if err != nil {
		panic(err)
	}

	validate = validator.New()
	// register validation for 'User'
	// NOTE: only have to register a non-pointer type for 'User', validator
	// interanlly dereferences during it's type checks.
	// validate.RegisterStructValidation(UserStructLevelValidation, ProgramConfig{})
	valerr := validate.Struct(programConfig)
	if valerr != nil {
		if _, ok := valerr.(*validator.InvalidValidationError); ok {
			log.Debugln(valerr)
			return
		}
		for _, valerr := range valerr.(validator.ValidationErrors) {
			// fmt.Println(valerr.Namespace())
			// fmt.Println(valerr.Field())
			// fmt.Println(valerr.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			// fmt.Println(valerr.StructField())     // by passing alt name to Reporterror like below
			// fmt.Println(valerr.Tag())
			// fmt.Println(valerr.ActualTag())
			// fmt.Println(valerr.Kind())
			// fmt.Println(valerr.Type())
			// fmt.Println(valerr.Value())
			// fmt.Println(valerr.Param())
			// fmt.Println()
			// fmt.Println("Parameter " + valerr.Tag() + " has value " + valerr.Value() + " which does not match required value " + valerr.Param())
			log.Debugf("Parameter %s has value %s which does not match required value %s\n", valerr.StructNamespace(), valerr.Value(), valerr.Param())

		}
		// from here you can create your own error messages in whatever language you wish
		log.Errorf(valerr.Error())
		panic("exit on validation")
	}
	log.Debugf("FileFilter: %v\n", programConfig.FileFilter)
	// fmt.Printf("Value: %v\n", programConfig)
}
