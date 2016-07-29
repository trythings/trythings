package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
)

var elliesPath = flag.String("ellies-path", "ellies-pad", "Path to the root of the ellies-pad directory")
var token = flag.String("token", "", "OAuth2 refresh token to use when deploying to App Engine")

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

	if err := deployToAppEngine(*elliesPath, *token); err != nil {
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

func deployToAppEngine(elliesPath, token string) error {
	log.Println("deploy: deploying to App Engine")
	ae := exec.Command("appcfg.py", "update", ".")

	if token != "" {
		ae.Args = append(ae.Args, "--oauth2_refresh_token", token)
	}

	ae.Dir = path.Join(elliesPath, "api")
	ae.Stdout = os.Stdout
	ae.Stderr = os.Stderr

	return ae.Run()
}
