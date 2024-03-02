package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Writes the .env file. It will truncate the file and recreate it every time.
func writeDotEnvFile(envdata string) error {
	// This receives a string that's ready to be writed to the .env file
	file, err := os.OpenFile(".env", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprint(file, envdata)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	// It will display a message and the data writed to the .env file.
	fmt.Println(".env has been updated!\n\n" + envdata)
	return nil
}

// Adds the .env and json.env files to .gitignore if they are not have been added already.
func addToGitIgnore(s string) error {
	// Opens the file in append mode.
	file, err := os.OpenFile(".gitignore", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var gitignore []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		gitignore = append(gitignore, line)
	}

	// Checks if the files we want to add to .gitignore are already there and skip writing if they are.
	toadd := []string{".env", "env.json"}
	for _, r := range toadd {
		if slices.Contains(gitignore, r) {
			fmt.Printf("'%v' is already in .gitignore!\n", r)
		} else {
			fmt.Printf("'%v' has been added .gitignore!\n", r)
			_, err = fmt.Fprint(file, "\n"+r)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Deletes the .env and env.json files from the current directory.
func cleanProject(s string) error {
	cleanhouse := []string{".env", "env.json"}

	for _, file := range cleanhouse {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Error removing '%v': %v\n", file, err)
		} else {
			fmt.Printf("%v has been successfully deleted\n", file)
		}
	}
	os.Exit(0)
	return nil
}

// Displays a visual representation of the project's info.
func displayProject(s string) error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	fmt.Printf(`Current Project Information:

Vault:		%v
Item:		%v
Vars:		%v
Stage:		%v
Stage Key:	%v
tfvars:		%v
		
Current .env state:

%v
`, projectData.Vault,
		projectData.Item,
		strings.Join(projectData.Vars, ", "),
		projectData.Stage, projectData.StageKey,
		projectData.tfvars,
		string(data))

	os.Exit(0)
	return nil
}

// Prints help.
func printHelp(s string) error {
	help := `Welcome to SetEnv. 
This CLI will help you create and update .env files for variables pointing to 1Password.

Usage: 
  setenv -[flag]=<value>

Examples:
  Starting a new project:
    setenv -vault="my vault" -item=project -addvar=DB_USER,DB_PASSWORD -stagekey=item -stage=test
  Prepending TF_VAR_ to every variable:
    setenv -tfvars
  Changing the vault's name:
    setenv -vault="new project"
  Removing variables:
    setenv -rmvar=DB_USER,DB_PASSWORD
`
	fmt.Println(help)
	fmt.Println("Options:")
	flag.PrintDefaults()

	os.Exit(0)
	return nil
}

// Takes a comma separated string of vars, splits into slice and trim the stirngs.
func sanitinzeVars(vars string) ([]string, error) {
	var varslice []string
	for _, v := range strings.Split(vars, ",") {
		varslice = append(varslice, strings.TrimSpace(v))
	}
	return varslice, nil
}
