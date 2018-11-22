package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var PaceSpaceGUID = "08dba6e1-270b-4cdc-9869-61bc44530030"

var CfDomain = "cfapps.io"
var CfDomainGUID = "fb6bd89f-2ed9-49d4-9ad1-97951a573135"

var CfAPI = "https://api.run.pivotal.io"

var CfUser = "pivotalpaceci@gmail.com"

var DefaultConfig = `{
    "workshopSubject":"PACE",
    "workshopHomepage":"",
    "modules": [
    {
        "type": "concepts",
        "content": [
            {
            "name":"example-slide",
            "filename":"example/example-slide"
            }
        ]
    },
    {
        "type": "demos",
        "content": [
            {
            "name":"example-demo",
            "filename":"example/example-demo"
            }
        ]
    }
  ]
}`

var DefaultManifest = `---
applications:
- name: my-pace-workshop
  memory: 512M
  instances: 1
  buildpacks: 
  - staticfile_buildpack
  random-route: true
  path: public/`

type WorkshopConfig struct {
	WorkshopHomepage string `json:"workshopHomepage"`
	WorkshopSubject  string `json:"workshopSubject"`
	WorkshopHostname string `json:"workshopHostname"`
	Modules          []struct {
		Type    string          `json:"type"`
		Content []ContentConfig `json:"content"`
	} `json:"modules"`
}

type ContentConfig struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

func DetermineConfig(path string) (*WorkshopConfig, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config not found")
	}
	var config WorkshopConfig
	err = json.Unmarshal(configFile, &config)
	return &config, nil
}
