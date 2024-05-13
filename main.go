package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/projects", projectsHandler)
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	defer func() {
		if err := server.Shutdown(nil); err != nil {
			log.Fatalf("Shutdown error: %v", err)
		}
	}()
	select {} // Block indefinitely
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	githubUsername := "your-username"
	reposURL := fmt.Sprintf("https://api.github.com/users/%s/repos", githubUsername)
	reposResponse, err := http.Get(reposURL)
	if err != nil {
		http.Error(w, "Error fetching GitHub repositories", http.StatusInternalServerError)
		return
	}
	defer reposResponse.Body.Close()

	var repos []map[string]interface{}
	if err := json.NewDecoder(reposResponse.Body).Decode(&repos); err != nil {
		http.Error(w, "Error decoding GitHub repositories response", http.StatusInternalServerError)
		return
	}

	var projectsHTML string
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(repo map[string]interface{}) {
			defer wg.Done()

			repoName := repo["name"].(string)
			repoLanguage := repo["language"].(string)

			readmeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", githubUsername, repoName)
			readmeResponse, err := http.Get(readmeURL)
			if err != nil {
				return
			}
			defer readmeResponse.Body.Close()

			if readmeResponse.StatusCode != http.StatusOK {
				return
			}

			readmeContentBytes, err := ioutil.ReadAll(readmeResponse.Body)
			if err != nil {
				return
			}
			readmeContent, err := base64.StdEncoding.DecodeString(string(readmeContentBytes))
			if err != nil {
				return
			}

			projectsHTML += fmt.Sprintf(`
                <div class="project-box">
                    <h2>%s</h2>
                    <h3>%s</h3>
                    <p>%s</p>
                </div>
            `, repoName, repoLanguage, string(readmeContent))
		}(repo)
	}
	wg.Wait()

	fmt.Fprintf(w, projectsHTML)
}
