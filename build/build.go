package build

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	git "gopkg.in/src-d/go-git.v4"
)

func BuildCmd() {

	if err := checkDependencies(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	config, err := determineConfig("config.json")
	if err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := cloneBaseRepo("https://github.com/Pivotal-Field-Engineering/pace-workshop-base"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := setWorkshopTitle(config); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

}

func setWorkshopTitle(config *WorkshopConfig) error {
	workshopNav := fmt.Sprintf("\n[1]\n    Name = \"%s Workshop\"\n    URL = \"/\"", config.WorkshopSubject)

	f, err := os.OpenFile("workshopGen/data/Menu.toml", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav menu file")
	}

	defer f.Close()

	if _, err = f.WriteString(workshopNav); err != nil {
		return fmt.Errorf("cannot write to nav menu file")
	}

	return nil
}

func cloneBaseRepo(repoPath string) error {

	fmt.Println("git clone " + repoPath)

	_, err := git.PlainClone(
		"workshopGen",
		false,
		&git.CloneOptions{
			URL:      repoPath,
			Progress: os.Stdout,
		},
	)

	if err != nil {
		return fmt.Errorf("cannot clone base git repo + %+v", err)
	}

	return nil
}

func checkDependencies() error {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+"git")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dependency git not found")
	}
	return nil
}

func determineConfig(path string) (*WorkshopConfig, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config not found")
	}
	var config WorkshopConfig
	err = json.Unmarshal(configFile, &config)
	return &config, nil
}

type WorkshopConfig struct {
	WorkshopSubject string `json:"workshopSubject"`
	Modules         []struct {
		Name    string `json:"name"`
		Content []struct {
			Name string `json:"name"`
			Page string `json:"page"`
		} `json:"content"`
	} `json:"modules"`
}
