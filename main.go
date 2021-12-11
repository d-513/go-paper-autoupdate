/*
  Author: dada513
  Title: Paper AutoUpdater
  Description: Automatically updates paper minecraft server

*/

package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type paperProject struct {
	ProjectName string `json:"project_name"`
	ProjectId   string `json:"project_id"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

type paperChange struct {
	Commit  string `json:"commit"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

type paperDownload struct {
	Name   string `json:"name"`
	Sha256 string `json:"sha256"`
}

type paperDownloads struct {
	Application paperDownload `json:"application"`
}

type paperBuild struct {
	ProjectName string         `json:"project_name"`
	ProjectId   string         `json:"project_id"`
	Version     string         `json:"version"`
	Build       int            `json:"build"`
	Time        string         `json:"time"`
	Changes     []paperChange  `json:"changes"`
	Downloads   paperDownloads `json:"downloads"`
}

func main() {
	// throw an error if command line arguments arent provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a project name and version")
		os.Exit(1)
	}

	paperProjectName := os.Args[1]
	paperProjectVer := os.Args[2]

	fmt.Println("Project Name:", paperProjectName)
	fmt.Println("Project Version:", paperProjectVer)
	resp, err := http.Get("https://papermc.io/api/v2/projects/" + paperProjectName + "/versions/" + paperProjectVer)
	if err != nil {
		log.Fatal(err)
		return
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	project := paperProject{}
	jsonErr := json.Unmarshal(body, &project)
	if jsonErr != nil {
		fmt.Println(err)
		return
	}
	var largestNumber, temp int
	for _, element := range project.Builds {
		if element > temp {
			temp = element
			largestNumber = temp
		}
	}
	fmt.Println("Latest build:", largestNumber)
	fExists := exists(paperProjectName + "-" + paperProjectVer + ".jar")
	redownload := false

	if !fExists {
		redownload = true
		fmt.Println("File does not exist, downloading...")
	}
	// if !buildIsLatest {
	// 	fmt.Println("Downloading build:", largestNumber)
	// 	// dlStr := "https://papermc.io/api/v2/projects/" + paperProjectName + "/versions/" + paperProjectVer + "/builds/" + strconv.Itoa(largestNumber) + "/" + paperProjectName + "-" + paperProjectVer + "-" + strconv.Itoa(largestNumber) + ".jar"

	// 	getBuildFile(paperProjectName, paperProjectVer, largestNumber)
	// }
	buildInfo, buildErr := GetBuildInfo(paperProjectName, paperProjectVer, largestNumber)
	if buildErr {
		fmt.Println("Error getting build info")
		os.Exit(1)
	}
	if fExists {
		hash := GetSha256(paperProjectName + "-" + paperProjectVer + ".jar")

		redownload = hash != buildInfo.Downloads.Application.Sha256
	}

	if redownload {
		dlStr := "https://papermc.io/api/v2/projects/" + paperProjectName + "/versions/" + paperProjectVer + "/builds/" + strconv.Itoa(largestNumber) + "/downloads/" + buildInfo.Downloads.Application.Name
		fmt.Println("Redownloading...")
		os.Remove(paperProjectName + "-" + paperProjectVer + ".jar")
		DownloadFile(paperProjectName+"-"+paperProjectVer+".jar", dlStr)
	}
	fmt.Println("JAR is updated to latest")
}

func GetBuildInfo(name string, ver string, build int) (*paperBuild, bool) {
	resp, err := http.Get("https://papermc.io/api/v2/projects/" + name + "/versions/" + ver + "/builds/" + strconv.Itoa(build))
	if err != nil {
		log.Fatal(err)
		return nil, true
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, true
	}

	ein_paperBuild := paperBuild{}
	jsonErr := json.Unmarshal(body, &ein_paperBuild)
	if jsonErr != nil {
		return nil, true
	}
	return &ein_paperBuild, false
}

func GetSha256(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
