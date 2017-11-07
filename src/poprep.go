package main

import (
	"poprep/src/pr"
)

const (
	configfullname = "poprep_config.toml"
)

func main() {
	//init and run poprep
	pr.Init(configfullname)
	pr.RunCli()
}
