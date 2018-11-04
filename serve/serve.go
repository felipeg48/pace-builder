package serve

import (
	"fmt"
	"os"
	"os/exec"
)

func ServeCmd() {

	if err := checkDependencies(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Checking workshopGen")
	if err := os.Chdir("workshopGen/"); err != nil {
		fmt.Println("Please `build` before `serve` to create the content. Error:" + err.Error())
		return
	}

	fmt.Println("Running hugo serve.  Check your content at http://localhost:1313 ...")
	if err := serveHugo(); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}
}

func checkDependencies() error {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+"hugo")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("It looks like Hugo is not installed! You need to install Hugo to run a local instance of the web site...")
	}
	return nil
}

func serveHugo() error {
	cmd := exec.Command("hugo", "serve")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error running Hugo: %s", err)
	}
	return nil
}
