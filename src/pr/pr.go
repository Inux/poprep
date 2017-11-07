package pr

import (
	"os"
	"poprep/src/prconf"
	"poprep/src/prgithub"
	"poprep/src/prlog"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var config = prconf.PoprepConfig{}
var cliapp = &cli.App{}

var logger = &logrus.Logger{}

//Init initializes poprep config full name has to be provided
func Init(awesomeListConfig, githubConfig string) {
	//initializing logrus
	prlog.New()
	logger = prlog.Get()

	logger.Info("initialize poprep - started")
	//read configuration
	logger.Info("reading config - started")
	tc, err := prconf.New(awesomeListConfig, githubConfig)
	if err != nil {
		logger.Fatal("couldn't read configuration: ", err.Error())
	}
	config = tc
	logger.Info("reading config - done")

	//init github client
	logger.Info("initialize github client - started")
	prgithub.New(config.Github.AccessToken)
	logger.Info("initialize github client - done")

	logger.Info("initialize cli application - started")
	cliapp = initCliApp()
	logger.Info("initialize cli application - done")

	logger.Info("initialize poprep - done")
}

//RunCli starts poprep as a CLi application
func RunCli() {
	logger.Info("running cli - started")
	logger.Info("os args: ", os.Args)
	cliapp.Run(os.Args)
	logger.Info("running cli - done")
}
