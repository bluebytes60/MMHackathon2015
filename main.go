package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

var result Result

func main() {
	t := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     "9d7ad4e06ed11c70c81b",
		ClientSecret: "8d91201c4a5dbd7ab45ba41e98c91da7d1f10111",
	}
	client := github.NewClient(t.Client())

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repos, _, err := client.Repositories.List("mediamath", opt)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	languages := make(map[string]Language)

	for _, repo := range repos {
		fmt.Println(*repo.Name)
		langs, _, err := client.Repositories.ListLanguages("Mediamath", *repo.Name)

		if err != nil {
			fmt.Errorf(err.Error())
		}

		for langStr := range langs {
			if lang, ok := languages[langStr]; ok {
				lang.Repos = append(lang.Repos, Repo{*repo.Name})
				languages[langStr] = lang
			} else {
				languages[langStr] = Language{Name: langStr, Repos: make([]Repo, 1, 100)}
				languages[langStr].Repos[0] = Repo{*repo.Name}
			}
		}
	}

	langArray := make([]Language, len(languages))
	i := 0
	for _, lang := range languages {
		langArray[i] = lang
		i++
	}

	result = Result{langArray}

	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":7777", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(result)
	if err == nil {
		w.Write(data)
	} else {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
	}
}

type Language struct {
	Name  string `json:"language"`
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Name string `json:"repository"`
}

type Result struct {
	Languages []Language `json:"data"`
}
