package prcache

import "poprep/src/prconf"

var conf = prconf.PoprepConfig{}

//New initializes a new prcache
func New(config prconf.PoprepConfig) {
	conf = config
}

//Check returns a list of awesomelists where the cache is outdated
func Check() map[string]prconf.AwesomeListConfig {
	return conf.AwesomeLists
}

//CreateForAwesomelist creates a cache for a given Awesomelist Key
func CreateForAwesomelist(alkey string) {

}
