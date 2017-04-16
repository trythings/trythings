package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	try := cli.NewApp()
	try.Name = "try"
	try.Usage = "try things"

	hash, err := gitHash()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	try.Version = hash

	try.Run(os.Args)
}

func gitHash() (string, error) {
	output, err := exec.Command("git", "stash", "create").Output()
	if err != nil {
		return "", err
	}

	var rev string
	if len(output) == 0 {
		rev = "HEAD"
	} else {
		rev = strings.TrimSuffix(string(output), "\n")
	}

	hash, err := exec.Command("git", "rev-parse", rev+"^{tree}").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(hash), "\n"), err
}
