package push

import (
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

	if _, err := os.Stat("./public"); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("/public directory not found! Please run `pace build` first!")
		}
	}

	appName := sillyname.GenerateStupidName()
	appName = strings.Replace(appName, " ", "-", -1)
	appName = strings.ToLower(appName)

	data, err := base64.StdEncoding.DecodeString("UGl2b3RhbDEhCg==")
	if err != nil {
		return err
	}

	fmt.Printf("Cf pushing %s app to: https://%s.%s .....", appName, appName, resources.CfDomain)

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

	err = zip.ArchiveFile("public/", "appFiles.zip", nil)
	if err != nil {
		return (err)
	}
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
		Host:       appName,
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
