package prconf

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	//AppName is the name of the application
	AppName = "poprep"
	//AppUsage describes what you can do with it
	AppUsage = "find awesome awesomeitories"
	//AppVersion of the application
	//TODO: autoincrement
	AppVersion = "0.0.0.1"
	//AppCopyright copyright of the application
	AppCopyright = "(c) 2017 inux"
)

//AppAuthors - list of all authors of this application
var (
	//AppCompiled is the date where the binary was built
	AppCompiled = time.Now()
	AppAuthors  = []cli.Author{
		cli.Author{
			Name:  "inux",
			Email: "inux.steve@gmail.com",
		},
	}
)

//PoprepConfig is the config structure for poprep
type PoprepConfig struct {
	Github       GithubAPIconfig
	AwesomeLists map[string]AwesomeListConfig
}

//GithubAPIconfig contains the Github API configuration
type GithubAPIconfig struct {
	AccessToken string
}

//AwesomeListConfig represents a AwesomeList from github
type AwesomeListConfig struct {
	Name   string
	Author string
	URL    string
	//Format options
	StartLine                string
	NamePrefix               string
	NamePostfix              string
	URLPrefix                string
	URLPostfix               string
	CategoryIdentifier       string
	CategoryIdentifierSingle string
}

//AwesomeListData is the struct used internally to represent the awesomelist data
type AwesomeListData struct {
	Config     AwesomeListConfig
	Categories map[string]string
	RepoInfos  []*RepoInfo
}

//RepoInfo struct represents a short info about a github repo
type RepoInfo struct {
	URL      string
	User     string
	Project  string
	Category string
	Stars    int
}

//New reads PoprepConfig from toml file
func New(configpath string) (PoprepConfig, error) {
	conf := &PoprepConfig{}
	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		return *conf, errors.New("Could not find poprep config file")
	}
	if _, err = toml.Decode(string(data), &conf); err != nil {
		logrus.Info("Could not parse configuration")
		return *conf, errors.New("Could not parse poprep config file")
	}
	return *conf, nil
}
