package push

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/Pallinder/sillyname-go"
	"github.com/Pivotal-Field-Engineering/pace-builder/resources"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/pierrre/archivefile/zip"
)

func PushCmd() error {

	if _, err := os.Stat("./publicGen"); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("./publicGen directory not found! Please run `pace build` first!")
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Pivotal Username - This is used to enable re-pushing your workshop: ")
	username, _ := reader.ReadString('\n')
	username = strings.ToLower(username)
	username = strings.Replace(username, " ", "-", -1)
	username = strings.Replace(username, "\n", "", -1)
	if len(username) < 3 {
		return fmt.Errorf("That's not your username! Try again....")
	}

	config, err := resources.DetermineConfig("config.json")
	if err != nil {
		return err
	}

	appName := config.WorkshopHostname
	if appName == "" {
		appName = sillyname.GenerateStupidName()
	}
	appName = strings.Replace(appName, " ", "-", -1)
	appName = strings.ToLower(appName)
	hostname := appName

	appName = fmt.Sprintf("%s-%s", username, appName)

	data, err := base64.StdEncoding.DecodeString("UGl2b3RhbDEhCg==")
	if err != nil {
		return err
	}

	fmt.Printf("Cf pushing %s app to: https://%s.%s .....", appName, hostname, resources.CfDomain)

	cfPass := strings.TrimSpace(string(data))

	c := &cfclient.Config{
		ApiAddress: resources.CfAPI,
		Username:   resources.CfUser,
		Password:   cfPass,
	}
	client, err := cfclient.NewClient(c)
	if err != nil {
		return err
	}

	appRequest := cfclient.AppCreateRequest{
		Name:      appName,
		SpaceGuid: resources.PaceSpaceGUID,
		Instances: 1,
		Buildpack: "staticfile_buildpack",
	}
	app, err := client.CreateApp(appRequest)
	if err != nil {
		return err
	}

	err = zip.ArchiveFile("publicGen/", "appFiles.zip", nil)
	if err != nil {
		return (err)
	}
	defer os.Remove("appFiles.zip")
	appBits, err := os.Open("appFiles.zip")
	if err != nil {
		return err
	}
	defer appBits.Close()

	client.UploadAppBits(appBits, app.Guid)

	envVars := map[string]interface{}{
		"FORCE_HTTPS": "true",
	}

	aur := cfclient.AppUpdateResource{
		Name:        app.Name,
		SpaceGuid:   app.SpaceGuid,
		Instances:   1,
		State:       "STARTED",
		Environment: envVars,
	}
	_, err = client.UpdateApp(app.Guid, aur)
	if err != nil {
		return err
	}

	routeRequest := cfclient.RouteRequest{

		DomainGuid: resources.CfDomainGUID,
		Host:       hostname,
		SpaceGuid:  resources.PaceSpaceGUID,
	}
	route, err := client.CreateRoute(routeRequest)
	if err != nil {
		return err
	}

	routeMappingRequest := cfclient.RouteMappingRequest{
		AppGUID:   app.Guid,
		RouteGUID: route.Guid,
		AppPort:   8080,
	}
	_, err = client.MappingAppAndRoute(routeMappingRequest)
	if err != nil {
		return err
	}

	fmt.Println("SUCCEEDED!")

	return nil
}
