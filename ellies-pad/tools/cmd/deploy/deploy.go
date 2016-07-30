package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var elliesPath = flag.String("elliesPath", "ellies-pad", "Path to the root of the ellies-pad directory")
var keyFile = flag.String("keyFile", "", "Service account key to use when deploying to App Engine")

func main() {
	flag.Parse()

	if _, err := os.Stat(*elliesPath); err != nil {
		log.Println("deploy: could not find the ellies-pad directory", err)
		os.Exit(1)
	}

	if err := buildWebApp(*elliesPath); err != nil {
		log.Println("deploy: could not build the web app", err)
		os.Exit(1)
	}

	if err := authWithAppEngine(*elliesPath, *keyFile); err != nil {
		log.Println("deploy: could not authenticate with App Engine", err)
		os.Exit(1)
	}

	if err := deployToAppEngine(*elliesPath); err != nil {
		log.Println("deploy: could not deploy to App Engine", err)
		os.Exit(1)
	}
}

func buildWebApp(elliesPath string) error {
	log.Println("deploy: building web app")
	npm := exec.Command("npm", "run", "build")
	npm.Dir = path.Join(elliesPath, "web")
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	return npm.Run()
}

func authWithAppEngine(elliesPath, keyFile string) error {
	if keyFile == "" {
		log.Println("deploy: no keyFile provided, assuming login authentication")
		return nil
	}

	log.Println("deploy: authenticating with App Engine")
	gc := exec.Command("gcloud", "auth", "activate-service-account")
	keyFile, err := filepath.Abs(keyFile)
	if err != nil {
		return err
	}
	gc.Args = append(gc.Args, "--key-file", keyFile)
	gc.Dir = path.Join(elliesPath, "api", "main")
	gc.Stdout = os.Stdout
	gc.Stderr = os.Stderr
	return gc.Run()
}

func deployToAppEngine(elliesPath string) error {
	log.Println("deploy: deploying to App Engine")
	gc := exec.Command(
		"go", "run", "../../../vendor/google.golang.org/appengine/cmd/aedeploy/aedeploy.go",
		"gcloud", "app", "deploy",
		"--project", "ellies-pad",
		"--verbosity", "info",
		"--version", "1",
		"--quiet",
	)
	gc.Dir = path.Join(elliesPath, "api", "main")
	gc.Stdout = os.Stdout
	gc.Stderr = os.Stderr
	return gc.Run()
}
