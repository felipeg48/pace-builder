package initialize

import (
	"fmt"
	"os"

	"github.com/Pivotal-Field-Engineering/pace-builder/resources"
)

func InitCmd() {

	fmt.Println("Generating default pace config.json")
	if err := createDefaultConfig(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Generating default cf push manifest.yml")
	if err := createDefaultManifest(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Config and Manifest have been generated. Edit the config and manifest to your desire. Run `pace serve` to run a local version of your application!")
}

func createDefaultConfig() error {
	f, err := os.OpenFile("config.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error creating config.json")
	}
	defer f.Close()
	_, err = f.WriteString(resources.DefaultConfig)
	if err != nil {
		return fmt.Errorf("error writing default config to config.json")
	}

	return nil
}

func createDefaultManifest() error {
	f, err := os.OpenFile("manifest.yml", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error creating manifest.yml")
	}
	defer f.Close()
	_, err = f.WriteString(resources.DefaultManifest)
	if err != nil {
		return fmt.Errorf("error writing default manifest to manifest.yml")
	}

	return nil
}
