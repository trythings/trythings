package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var elliesPath = flag.String("ellies-path", "ellies-pad", "Path to the root of the ellies-pad directory")
var email = flag.String("email", "", "Email for the Google account to use when deploying to App Engine")
var password = flag.String("password", "", "App-specific password to use when deploying to App Engine")

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

	if err := deployToAppEngine(*elliesPath, *email, *password); err != nil {
		log.Println("deploy: could not deploy to App Engine", err)
		os.Exit(1)
	}
}

func buildWebApp(elliesPath string) error {
	log.Println("deploy: building web app")
	npm := exec.Command("npm", "run", "production")
	npm.Dir = path.Join(elliesPath, "web")
	npm.Stdout = os.Stdout
	npm.Stderr = os.Stderr
	return npm.Run()
}

func deployToAppEngine(elliesPath, email, password string) error {
	log.Println("deploy: deploying to App Engine")
	ae := exec.Command("appcfg.py", "update", "api")

	if email != "" || password != "" {
		if email == "" || password == "" {
			log.Println("deploy: -email and -password must be both provided or both empty")
			os.Exit(1)
		}

		ae.Args = append(ae.Args, fmt.Sprintf("--email=%s", email), "--passin")
		// --passin expects the password to come from stdin.
		ae.Stdin = strings.NewReader(password)
	}

	ae.Dir = elliesPath
	ae.Stdout = os.Stdout
	ae.Stderr = os.Stderr

	return ae.Run()
}
