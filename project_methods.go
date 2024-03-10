package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Reads the env.json file and return an initialized struct with it's content.
// A new env.json file will be generated later when the configuration is written.
func (p *Project) readJsonEnv() error {
	byteValue, err := os.ReadFile("env.json")
	if err != nil {
		fmt.Println("No env.json file has been found. Generating a new one.")
		return nil
	} else {
		err = json.Unmarshal(byteValue, p)
		if err != nil {
			return fmt.Errorf("error parsing JSON file: %v", err)
		}
	}
	return nil
}

// Writes the env.json file with the updated information.
func (p *Project) writeJsonEnv() error {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	err = os.WriteFile("env.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}
	return nil
}

// Takes the Project data and outputs all the variables pointing to the correct vault and item.
//
// The expected format for 1Password is:
// variable_name=op://vault_name/vault_item/variable_name
func (p *Project) setDotEnv() error {
	var envdata []string

	// Sanitizes Vault and Item values in case there was a trailing space.
	p.Vault = strings.TrimSpace(p.Vault)
	p.Item = strings.TrimSpace(p.Item)

	// Populates the stage in the selected stage key.
	if p.Stage != "" {
		p.addStage()
	}

	// Iterates over the Variables in our project and create a line for each one, populating the
	// variable name, vault and item, with the environment in the chosen key.
	for k, v := range p.Vars {
		varline := fmt.Sprintf("%v=\"op://%v/%v/%v\"", k, p.Vault, p.Item, v)
		if p.tfvars {
			varline = fmt.Sprintf("TF_VAR_" + varline)
		}
		envdata = append(envdata, varline)
	}

	// Concatenate all the lines with a new line and send it to the function to write it to a file.
	err := writeDotEnvFile(strings.Join(envdata, "\n"))
	if err != nil {
		return err
	}
	return nil
}

// Adds the stage to the stage key.
func (p *Project) addStage() {
	switch p.StageKey {
	case "vault":
		p.Vault = p.Vault + "-" + p.Stage
	case "item":
		p.Item = p.Item + "-" + p.Stage
	case "vars":
		for k, v := range p.Vars {
			p.Vars[k] = v + "-" + p.Stage
		}
	}
}

// Removes vars from the Vars struct's key
func (p *Project) rmVar(toremove string) error {
	var varslice []string
	// Split the comma separated vars from the flag into a slice.
	varslice = append(varslice, strings.Split(toremove, ",")...)
	// Trim whitepsaces and remove the key from the Vars map
	for i, v := range varslice {
		varslice[i] = strings.TrimSpace(v)
		delete(p.Vars, v)
	}
	return nil
}

// Adds vars into the Vars map.
func (p *Project) addVar(toadd string) error {
	var varslice []string
	// Split the comma separated vars from the flag into a slice.
	varslice = append(varslice, strings.Split(toadd, ",")...)
	// Iterate over the vars trimming them for whitespaces and adding them in the Vars map.
	for _, v := range varslice {
		// Split each key value pair by colon, if any.
		kvslice := strings.Split(v, ":")
		// Trim whitepsaces
		for i, v := range kvslice {
			kvslice[i] = strings.TrimSpace(v)
		}
		// Adding them to the Vars map. If no value was passed, the key will also be set as the value.
		if len(kvslice) == 1 {
			p.Vars[kvslice[0]] = kvslice[0]
		} else {
			p.Vars[kvslice[0]] = kvslice[1]
		}
	}
	return nil
}

// Validates that the stage pased with the stage flag is allowed.
func (p *Project) validateStage(stage string) error {
	stageList := []string{"test", "staging", "prod"}
	if slices.Contains(stageList, stage) {
		p.Stage = stage
		return nil
	} else if stage == "aws" {
		p.Stage = "$AWS_ENV"
		return nil
	}
	return fmt.Errorf("allowed options are: %v", stageList)
}

// Validates that the stagekey pased with the stagekey flag is allowed.
func (p *Project) validateStageKey(stagekey string) error {
	stageKeyList := []string{"vault", "item", "vars"}
	if slices.Contains(stageKeyList, stagekey) {
		p.StageKey = stagekey
		return nil
	}
	return fmt.Errorf("allowed options are: %v", stageKeyList)
}
