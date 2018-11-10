package build

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var languages = [...]string{"en", "es", "fr", "pt"}

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

	fmt.Println("Cleaning up existing workshopGen...")
	if err := os.RemoveAll("workshopGen/"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := cloneBaseRepo("https://github.com/Pivotal-Field-Engineering/pace-workshop-base", "workshopGen"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := setWorkshopTitle(config); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	if err := setWorkshopContent(config); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

	fmt.Println("Cleaning up temp workshopContent...")
	if err := os.RemoveAll("workshopContent/"); err != nil {
		fmt.Println("Error " + err.Error())
		return
	}

}

func setWorkshopContent(config *WorkshopConfig) error {

	if err := cloneBaseRepo("https://github.com/Pivotal-Field-Engineering/pace-workshop-content", "workshopContent"); err != nil {
		return err
	}

	for _, module := range config.Modules {
		if (strings.Compare(module.Type, "concepts")) == 0 {
			setWorkshopConcepts(module.Content)
		} else if module.Type == "demos" {
			setWorkshopDemos(module.Content)
		} else {
			fmt.Printf("Config contains a module (%s) that is not of type demos or concepts. This is not allowed! \n", module.Type)
		}
	}
	return nil
}

func setWorkshopDemos(contents []ContentConfig) error {
	for _, content := range contents {
		err := setWorkshopExtras(content, "demos")
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, language := range languages {
			fileName := strings.Split(content.Filename, "/")
			pageFile := "workshopGen/content/demos/" + fileName[len(fileName)-1] + "." + language + ".md"
			err := createPage(pageFile, content.Name)

			if err != nil {
				return err
			}

			contentPath := "workshopContent/" + content.Filename
			err = addMarkdown(pageFile, contentPath+"."+language+".md")
			if err != nil {
				fmt.Printf("cannot add specified demo markdown to file, %s, %+v", fileName[len(fileName)-1]+"."+language+".md", err)
			}
		}
	}
	return nil
}

func setWorkshopConcepts(contents []ContentConfig) error {
	for _, content := range contents {
		err := setWorkshopExtras(content, "contents")
		if err != nil {
			return err
		}
		for _, language := range languages {
			fileName := strings.Split(content.Filename, "/")
			pageFile := "workshopGen/content/concepts/" + fileName[len(fileName)-1] + "." + language + ".md"
			err := createPage(pageFile, content.Name)

			if err != nil {
				return err
			}

			contentPath := "workshopContent/" + content.Filename
			err = addMarkdown(pageFile, contentPath+"."+language+".md")
			if err != nil {
				fmt.Printf("cannot add specified content markdown to file, %s, %+v", fileName[len(fileName)-1]+"."+language+".md", err)
			}
		}
	}
	return nil
}

func setWorkshopExtras(curContent ContentConfig, contType string) error {

	var (
		destination string
		source      string
	)

	contentPath := strings.Split(curContent.Filename, "/")
	folders := contentPath[:len(contentPath)-1]
	folderPath := strings.Join(folders, "/")

	source = "workshopContent/" + folderPath + "/"

	if contType == "demos" {
		destination = "workshopGen/content/demos/" + contentPath[len(contentPath)-1] + "/"
		_ = os.MkdirAll(destination, os.FileMode(0777))
	} else if contType == "concepts" {
		destination = "workshopGen/content/concepts/" + contentPath[len(contentPath)-1] + "/"
		_ = os.MkdirAll(destination, os.FileMode(0777))
	} else {
		return fmt.Errorf("content is not of demos or concepts types")
	}

	fds, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(source, fd.Name())
		dstfp := path.Join(destination, fd.Name())

		if !fd.IsDir() {
			if filepath.Ext(strings.TrimSpace(fd.Name())) != ".md" {

				srcfd, err := os.Open(srcfp)
				if err != nil {
					return err
				}
				defer srcfd.Close()

				dstfd, err := os.Create(dstfp)
				if err != nil {
					return err
				}
				defer dstfd.Close()

				if _, err = io.Copy(dstfd, srcfd); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func addMarkdown(existingFile string, additionalMarkDown string) error {
	additionalMarkDownWriter, err := os.Open(additionalMarkDown)
	if err != nil {
		fmt.Printf("%s not found!\n", additionalMarkDown)
		os.Remove(existingFile)
		return nil
	}
	defer additionalMarkDownWriter.Close()
	existingFileWriter, err := os.OpenFile(existingFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing %s", err)
	}
	defer existingFileWriter.Close()
	_, err = io.Copy(existingFileWriter, additionalMarkDownWriter)
	if err != nil {
		log.Fatalln("failed to append files:", err)
	}

	return nil
}

func createPage(file string, title string) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("cannot create file, %s, %+v", file, err)
	}
	header := fmt.Sprintf("+++\ntitle = \"\"\nmenuTitle = \"%s\"\nchapter = false\nweight = 3\ndescription = \"\"\ndraft = false\n+++\n", title)
	_, err = f.WriteString(header)
	if err != nil {
		return fmt.Errorf("cannot write string %s, %+v", header, err)
	}
	return nil
}

func setWorkshopTitle(config *WorkshopConfig) error {
	workshopTitle := fmt.Sprintf("%s Workshop", config.WorkshopSubject)
	workshopToml := fmt.Sprintf("+++\ntitle = \"%s\"\nchapter = true\nweight = 1\n+++\n\n", workshopTitle)
	workshopHomepageContent := workshopToml
	if config.WorkshopHomepage != "" {
		homepageContent, err := ioutil.ReadFile(config.WorkshopHomepage)
		if err != nil {
			fmt.Printf("%s not found!\n", config.WorkshopHomepage)
			return err
		}
		workshopHomepageContent = workshopHomepageContent + string(homepageContent)

	} else {
		workshopHomepageContent = workshopHomepageContent + `<div class="container" style="
    padding-top: 3rem;
    max-width: 80rem;
    background-color: #058673;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    padding-left: 5rem;
    padding-right: 5rem;
    padding-bottom: 3rem;
    width: 100%;
">
    <a href="/images/pa-logo.png" data-featherlight="image"><img src="/images/pa-logo.png" alt="Platform Architecture" style="height: 10rem;"></a>
    <p style="max-width: 40rem; padding-left: 2rem;">
        Pivotal Platform Architecture helps companies learn how to
        solve IT
        and engineering challenges. We encourage you to explore our
        workshops. Build the future with Pivotal!
    </p>
</div>
<br>
<br>

<p style="font-family: Novacento Sans Wide, Helvetica, Tahoma, Geneva, Arial, sans-serif;
    text-align: center;
    text-transform: uppercase;
    color: #222;
    font-weight: 200;
	font-size: 3rem;">` + workshopTitle + `</p>
<div style="border-top: 4px solid #F0F2F4; border-top-left-radius: 5em; border-top-right-radius: 5em; min-height: 5em;"></div>
	
{{< toc >}}`
	}

	workshop, err := os.OpenFile("workshopGen/content/workshop/_index.en.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}

	workshop, err = os.OpenFile("workshopGen/content/workshop/_index.es.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}

	workshop, err = os.OpenFile("workshopGen/content/workshop/_index.fr.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}

	workshop, err = os.OpenFile("workshopGen/content/workshop/_index.pt.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("cannot open nav workshop file")
	}

	defer workshop.Close()

	if _, err = workshop.WriteString(workshopHomepageContent); err != nil {
		return fmt.Errorf("cannot write to workshop file")
	}
	return nil
}

func cloneBaseRepo(repoPath string, destinationPath string) error {
	data, err := base64.StdEncoding.DecodeString("NTA5OWU2MjE2ZjU0NjExYjNkY2UwNTc0MTk2MzVkOTg4ZDNlMmNjZgo=")
	if err != nil {
		return err
	}

	gitToken := strings.TrimSpace(string(data))

	fmt.Println("git clone " + repoPath)

	_, err = git.PlainClone(
		destinationPath,
		false,
		&git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "doesnotmatter",
				Password: gitToken,
			},
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
	cmd := exec.Command("/bin/sh", "-c", "command -v "+"hugo")
	if err := cmd.Run(); err != nil {
		fmt.Println("It looks like Hugo is not installed! You need to install Hugo to run a local instance of the web site...")
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
	WorkshopHomepage string `json:"workshopHomepage"`
	WorkshopSubject  string `json:"workshopSubject"`
	Modules          []struct {
		Type    string          `json:"type"`
		Content []ContentConfig `json:"content"`
	} `json:"modules"`
}

type ContentConfig struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}
