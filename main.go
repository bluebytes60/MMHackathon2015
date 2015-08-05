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
		ClientID:     "91ebc275df096fcdc2f1",
		ClientSecret: "44df9cd0058914eafb2a5065bc8886931a085e96",
	}
	client := github.NewClient(t.Client())

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repositories, _, err := client.Repositories.List("mediamath", opt)

	if err != nil {
		fmt.Errorf(err.Error())
	}

	languages := make(map[string]Language)
	repositoriesData := make([]Repository, 0)

	for _, repo := range repositories {
		fmt.Println(*repo.Name)
		//language
		langs, _, err := client.Repositories.ListLanguages("Mediamath", *repo.Name)

		if err != nil {
			fmt.Errorf(err.Error())
		}
		for langStr := range langs {
			if lang, ok := languages[langStr]; ok {
				lang.Repos = append(lang.Repos, *repo.Name)
				languages[langStr] = lang
			} else {
				languages[langStr] = Language{Name: langStr, Repos: make([]string, 1, 100)}
				languages[langStr].Repos[0] = *repo.Name
			}
		}
		//collaborators
		users, _, err := client.Repositories.ListCollaborators("Mediamath", *repo.Name, &github.ListOptions{0, 0})
		fmt.Println(users)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		userData := make([]string, 0)
		for _, user := range users {
			userData = append(userData, *user.Login)
		}
		repositoriesData = append(repositoriesData, Repository{*repo.Name, userData})

	}

	langArray := make([]Language, len(languages))
	i := 0
	for _, lang := range languages {
		langArray[i] = lang
		i++
	}

	result = Result{langArray}

	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":5555", nil))
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

type Result struct {
	Languages []Language `json:"data"`
	//Repositories []Repository `json:"repos"`
}

type Language struct {
	Name  string   `json:"language"`
	Repos []string `json:"repos"`
}

type Repository struct {
	Name  string
	Users []string
}
