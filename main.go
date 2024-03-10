package main

import (
	"flag"
	"log"
	"os"
)

type Project struct {
	Vault    string            `json:"vault"`
	Item     string            `json:"item"`
	StageKey string            `json:"stage_key"`
	Vars     map[string]string `json:"vars"`
	Stage    string            `json:"stage"`
	tfvars   bool              `json:"-"`
}

var projectData = Project{
	Stage:  "",
	tfvars: false,
	Vars:   make(map[string]string),
}

func main() {
	// Load project data from env.json
	err := projectData.readJsonEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Set all the flags and methods to them
	setFlags()

	// Parse the flags
	err = flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Display help if no flags were provided
	if flag.NFlag() == 0 {
		printHelp("")
		return
	}

	// Save the env.json and write the .env file
	err = projectData.writeJsonEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Concatenates the Project's data into a OP variable and creates the .env file
	err = projectData.setDotEnv()
	if err != nil {
		log.Fatal(err)
	}
}
