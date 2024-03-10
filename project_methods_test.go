package main

import (
	"strings"
	"testing"
)

func mockProject() *Project {
	return &Project{
		Vault:    "example-vault",
		Item:     "example-item",
		Vars:     map[string]string{"EXISTING_VAR1": "BLA", "EXISTING_VAR2": "EXISTING_VAR2"},
		Stage:    "dev",
		StageKey: "item",
		tfvars:   false,
	}
}

func TestAddVar(t *testing.T) {
	t.Run("new variable", func(t *testing.T) {
		p := mockProject()
		newVar := "NEW_VAR"
		err := p.addVar(newVar)
		if err != nil {
			t.Errorf("addvar returned an error: %v", err)
		}

		if p.Vars["NEW_VAR"] != "NEW_VAR" {
			t.Errorf("%v was not added to Vars", newVar)
		}
	})

	t.Run("addin multiple new variables", func(t *testing.T) {
		p := mockProject()
		newVars := "NEW1,NEW2:test"
		err := p.addVar(newVars)
		if err != nil {
			t.Errorf("addVar returned an error: %v", err)
		}

		if p.Vars["NEW1"] != "NEW1" || p.Vars["NEW2"] != "test" {
			t.Errorf("New variables were not found in Vars")
		}
	})

	t.Run("modifying an existing variable", func(t *testing.T) {
		p := mockProject()
		err := p.addVar("EXISTING_VAR1:test")
		if err != nil {
			t.Errorf("addVar returned an error: %v", err)
		}

		if p.Vars["EXISTING_VAR1"] != "test" {
			t.Errorf("Duplicated variable has been added to Vars")
		}
	})
}

func TestRmVars(t *testing.T) {
	t.Run("removing a variable", func(t *testing.T) {
		p := mockProject()
		rmVar := "EXISTING_VAR1"
		err := p.rmVar(rmVar)
		if err != nil {
			t.Errorf("rmvar returned an error: %v", err)
		}

		if _, ok := p.Vars["EXISTING_VAR1"]; ok {
			t.Errorf("The key 'EXISTING_VAR1' is still present")
		}
	})

	t.Run("removing multiple variables", func(t *testing.T) {
		to_remove := []string{"EXISTING_VAR1", "EXISTING_VAR2"}
		p := mockProject()
		rmVar := strings.Join(to_remove, ",")
		err := p.rmVar(rmVar)
		if err != nil {
			t.Errorf("rmvar returned an error: %v", err)
		}

		for _, v := range to_remove {
			if _, ok := p.Vars[v]; ok {
				t.Errorf("The key %s is still present", v)
			}
		}
	})
}

func TestValidateStage(t *testing.T) {
	t.Run("with valid stage", func(t *testing.T) {
		stages := []string{"test", "staging", "prod"}
		for _, stage := range stages {
			p := mockProject()
			err := p.validateStage(stage)
			if err != nil {
				t.Errorf("Validate stage error for: %v", stage)
			}
			if p.Stage != stage {
				t.Errorf("stage has not been updated")
			}
		}
	})

	t.Run("with aws_env as stage", func(t *testing.T) {
		stage := "aws"
		p := mockProject()
		err := p.validateStage(stage)
		if err != nil {
			t.Errorf("Validate stage error for: %v", stage)
		}
		if p.Stage != "$AWS_ENV" {
			t.Errorf("stage has not been updated")
		}
	})

	t.Run("with an invalid stage", func(t *testing.T) {
		stage := "nope"
		p := mockProject()
		err := p.validateStage(stage)
		if err == nil {
			t.Errorf("Error for stage validation was not triggered %v", stage)
		}
		if p.Stage == "nope" {
			t.Errorf("Incorrect stage has been updated")
		}
	})
}

func TestValidateStageKey(t *testing.T) {
	t.Run("with valid stagekey", func(t *testing.T) {
		stagekeys := []string{"vault", "item", "vars"}
		for _, stagekey := range stagekeys {
			p := mockProject()
			err := p.validateStageKey(stagekey)
			if err != nil {
				t.Errorf("Validate stage error for: %v", stagekey)
			}
			if p.StageKey != stagekey {
				t.Errorf("stage has not been updated")
			}
		}
	})

	t.Run("with an invalid stage", func(t *testing.T) {
		stagekey := "nope"
		p := mockProject()
		err := p.validateStageKey(stagekey)
		if err == nil {
			t.Errorf("Error for stage validation was not triggered %v", stagekey)
		}
		if p.StageKey == "nope" {
			t.Errorf("Incorrect stage has been updated")
		}
	})
}

func TestAddStage(t *testing.T) {
	t.Run("vault as stagekey", func(t *testing.T) {
		p := mockProject()
		p.StageKey = "vault"
		p.addStage()

		if !strings.Contains(p.Vault, "dev") {
			t.Errorf("Failed to modify stage in 'vault'")
		}
	})

	t.Run("item as stagekey", func(t *testing.T) {
		p := mockProject()
		p.StageKey = "item"
		p.addStage()

		if !strings.Contains(p.Item, "dev") {
			t.Errorf("Failed to modify stage in 'item'")
		}
	})

	t.Run("vars as stagekey", func(t *testing.T) {
		p := mockProject()
		p.StageKey = "vars"
		p.addStage()

		for _, v := range p.Vars {
			if !strings.Contains(v, "dev") {
				t.Errorf("Failed to modify stage in 'vars'")
			}
		}
	})

	t.Run("with no stagekey", func(t *testing.T) {
		p := mockProject()
		p.StageKey = ""
		p.addStage()

		for _, v := range p.Vars {
			if strings.Contains(v, "-") {
				t.Errorf("Vars were modified with no stagekey present")
			}
		}
	})
}
