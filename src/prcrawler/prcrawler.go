package prcrawler

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"poprep/src/prconf"
	"strings"

	"io"

	. "github.com/ahmetalpbalkan/go-linq"

	"github.com/google/go-github/github"
)

var (
	githubIdentifiers = []string{"https://github.com/", "http://github.com/"}
	repoInfos         []prconf.RepoInfo
)

//GetRepoInfos returns repo infos if cat string is specified only this category will be updated
func GetRepoInfos(ctx context.Context, client *github.Client, al prconf.AwesomeListConfig, cat string) []*prconf.RepoInfo {
	response, err := http.Get(al.URL)
	check(err)
	defer response.Body.Close()

	repoInfos := parseFile(response.Body, al, cat)
	repoInfos = getGithubInfos(ctx, client, repoInfos)

	return repoInfos
}

func getGithubInfos(ctx context.Context, client *github.Client, repinfs []*prconf.RepoInfo) []*prconf.RepoInfo {
	repoInfos := []*prconf.RepoInfo{}
	for _, ri := range repinfs {
		repo, _, err := client.Repositories.Get(ctx, ri.User, ri.Project)
		if err == nil {
			ri.Stars = *repo.StargazersCount
		} else {
			fmt.Println("Github API ERROR: ", string(err.Error()))
		}
	}
	//Sorting slice by stars
	repoQuery := From(repinfs).OrderBy(
		func(inf interface{}) interface{} {
			return inf.(*prconf.RepoInfo).Stars > 0
		}).ThenByDescending(func(inf interface{}) interface{} {
		return inf.(*prconf.RepoInfo).Stars
	})
	repoQuery.ToSlice(&repoInfos)

	return repoInfos
}

func parseFile(f io.ReadCloser, al prconf.AwesomeListConfig, cat string) []*prconf.RepoInfo {
	var repInfs []*prconf.RepoInfo
	startParsing := false
	category := ""

	//Check if only single category (depends on cat argument)
	startline := al.StartLine
	singleCat := false
	if cat != "" {
		startline = al.CategoryIdentifier + " " + cat
		category = cat
		singleCat = true
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if startParsing {
			s := scanner.Text()
			//check if category change
			if strings.Contains(s, al.CategoryIdentifier) {
				if singleCat {
					break
				}
				category = s
				category = strings.Replace(category, al.CategoryIdentifierSingle, "", -1)
				category = strings.Replace(category, " ", "", -1)
			} else {
				repInf := parseLine(s, category, al)
				if repInf != nil {
					repInfs = append(repInfs, repInf)
				}
			}
		} else {
			if strings.Contains(scanner.Text(), startline) {
				startParsing = true
			}
		}

	}
	return repInfs
}

func parseLine(line string, category string, al prconf.AwesomeListConfig) *prconf.RepoInfo {
	name := ""
	url := ""

	namestart := strings.Index(line, al.NamePrefix)
	nameend := strings.Index(line, al.NamePostfix)
	if namestart != -1 && nameend != -1 {
		name = line[namestart+1 : nameend]
	}
	urlstart := strings.Index(line, al.URLPrefix)
	urlend := strings.Index(line, al.URLPostfix)
	if urlstart != -1 && urlend != -1 {
		url = line[urlstart+1 : urlend]
	}
	if name != "" && url != "" {
	}

	//TODO implement way to find repositories from other informations
	//currently only github links are recognized in awesome lists
	if ContainsAnyOf(url, githubIdentifiers) && name != "" {
		reinf := &prconf.RepoInfo{}
		reinf.URL = url
		reinf.User = getUserFromURL(url)
		reinf.Project = getProjectFromURL(reinf.User, url)
		reinf.Category = category
		return reinf
	}
	return nil
}

func getUserFromURL(s string) string {
	rv := ""
	for _, v := range githubIdentifiers {
		if strings.Contains(s, v) {
			rv := ""
			userstart := strings.Index(s, v)
			if userstart == -1 {
				return rv
			}
			substring := s[userstart+len(v):]
			if substring == "" {
				return rv
			}
			userend := strings.Index(substring, "/")
			if userend == -1 {
				userend = len(substring)
			}
			rv = substring[:userend]
			return rv
		}
	}
	return rv
}

func getProjectFromURL(user, url string) string {
	rv := ""
	projectstart := strings.Index(url, user+"/")
	if projectstart == -1 {
		return rv
	}
	substring := url[projectstart+len(user+"/"):]
	if substring == "" {
		return rv
	}
	rv = substring
	return rv
}

//GetCategories return all categories from a AwesomList
func GetCategories(ctx context.Context, al prconf.AwesomeListConfig) []string {
	response, err := http.Get(al.URL)
	check(err)

	categories := parseFileForCategories(response.Body, al)

	return categories
}

func parseFileForCategories(rc io.ReadCloser, al prconf.AwesomeListConfig) []string {
	categories := []string{}
	startParsing := false

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		if startParsing {
			s := scanner.Text()
			//check if category change
			if strings.Contains(s, al.CategoryIdentifier) {
				category := s
				category = strings.Replace(category, al.CategoryIdentifierSingle, "", -1)
				categories = append(categories, category)
			}
		} else {
			if strings.Contains(scanner.Text(), al.StartLine) {
				startParsing = true
			}
		}

	}

	return categories
}

//ContainsCategory returns true if a certain awesomelist contains a category, else false
func ContainsCategory(al prconf.AwesomeListConfig, cat string) bool {
	response, err := http.Get(al.URL)
	check(err)
	startParsing := false

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		if startParsing {
			s := scanner.Text()
			//check if category change
			if strings.Contains(s, al.CategoryIdentifier) {
				if strings.Contains(s, cat) {
					return true
				}
			}
		} else {
			if strings.Contains(scanner.Text(), al.StartLine) {
				startParsing = true
			}
		}
	}
	return false
}

//ContainsAnyOf returns true if a string contains any string of a array
func ContainsAnyOf(in string, toCheck []string) bool {
	for _, v := range toCheck {
		if strings.Contains(in, v) {
			return true
		}
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
