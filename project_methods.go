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

	// We populate the stage in the selected stage key.
	if p.Stage != "" {
		p.addStage()
	}

	// Here we iterate over the Variables in our project and create a line for each one, populating the
	// variable name, vault and item, with the environment in the chosen key.
	for _, variable := range p.Vars {
		varline := fmt.Sprintf("%v=op://%v/%v/%v", variable, p.Vault, p.Item, variable)
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
		for i, variable := range p.Vars {
			p.Vars[i] = variable + "-" + p.Stage
		}
	}
}

// It removes vars from the Vars struct's key
func (p *Project) rmVar(toremove string) error {
	varstoremove, err := sanitinzeVars(toremove)
	if err != nil {
		return err
	}
	var updatedvars []string
	// Iterates over the current variables, checkin if it's in the list of variables to remove
	for _, currentvar := range p.Vars {
		// Adds the var to a new slice if it's not to remove.
		if !slices.Contains(varstoremove, currentvar) {
			updatedvars = append(updatedvars, currentvar)
		}
	}
	p.Vars = updatedvars
	return nil
}

// It adds vars from the Vars struct's key
func (p *Project) addVar(toadd string) error {
	varstoadd, err := sanitinzeVars(toadd)
	if err != nil {
		return err
	}
	// This function receives the vars to add in a comma separated string.
	// We split it into a slice and append it to the project's Vars.
	p.Vars = append(p.Vars, varstoadd...)

	// We sort and remove the duplicated vars after.
	slices.Sort(p.Vars)
	p.Vars = slices.Compact(p.Vars)
	return nil
}

// It validates that the stage pased with the stage flag is allowed.
func (p *Project) validateStage(stage string) error {
	stageList := []string{"test", "staging", "prod"}
	if slices.Contains(stageList, stage) {
		p.Stage = stage
		return nil
	}

	return fmt.Errorf("allowed options are: %v", stageList)
}

// It validates that the stagekey pased with the stagekey flag is allowed.
func (p *Project) validateStageKey(stagekey string) error {
	stageKeyList := []string{"vault", "item", "vars"}
	if slices.Contains(stageKeyList, stagekey) {
		p.StageKey = stagekey
		return nil
	}

	return fmt.Errorf("allowed options are: %v", stageKeyList)
}
