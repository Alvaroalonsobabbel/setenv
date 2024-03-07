package main

import (
	"slices"
	"strings"
	"testing"
)

func mockProject() *Project {
	return &Project{
		Vault:    "example-vault",
		Item:     "example-item",
		Vars:     []string{"EXISTING_VAR1", "EXISTING_VAR2"},
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

		if !slices.Contains(p.Vars, newVar) {
			t.Errorf("%v was not added to Vars", newVar)
		}
	})

	t.Run("addin multiple new variables", func(t *testing.T) {
		p := mockProject()
		newVars := "NEW1,NEW2"
		err := p.addVar(newVars)
		if err != nil {
			t.Errorf("addVar returned an error: %v", err)
		}

		if !slices.Contains(p.Vars, "NEW1") || !slices.Contains(p.Vars, "NEW2") {
			t.Errorf("New variables were not found in Vars")
		}
	})

	t.Run("adding an existing variable", func(t *testing.T) {
		p := mockProject()
		initialLength := len(p.Vars)
		err := p.addVar("EXISTING_VAR1")
		if err != nil {
			t.Errorf("addVar returned an error: %v", err)
		}

		if len(p.Vars) != initialLength {
			t.Errorf("Duplicated variable has been added to Vars")
		}
	})

	t.Run("adding a mix of existing and new vars", func(t *testing.T) {
		p := mockProject()
		err := p.addVar("NEW_VAR3,EXISTING_VAR1, NEW_VAR4 ")
		if err != nil {
			t.Errorf("addVar returned an error: %v", err)
		}
		if !slices.Contains(p.Vars, "NEW_VAR3") || !slices.Contains(p.Vars, "NEW_VAR4") {
			t.Errorf("New variables were not added to Vars")
		}
		if strings.Count(strings.Join(p.Vars, ","), "EXISTING_VAR1") > 1 {
			t.Errorf("Duplicate variable EXISTING_VAR1 was added to Vars")
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

		if slices.Contains(p.Vars, rmVar) {
			t.Errorf("%v was not removed from Vars", rmVar)
		}
	})

	t.Run("removing multiple variables", func(t *testing.T) {
		p := mockProject()
		rmVar := "EXISTING_VAR1,EXISTING_VAR2"
		err := p.rmVar(rmVar)
		if err != nil {
			t.Errorf("rmvar returned an error: %v", err)
		}

		if slices.Contains(p.Vars, "EXISTING_VAR1") || slices.Contains(p.Vars, "EXISTING_VAR2") {
			t.Errorf("%v was not removed from Vars", rmVar)
		}
	})

	t.Run("removing a non existing var", func(t *testing.T) {
		p := mockProject()
		rmVar := "EXISTING_VAR3"
		err := p.rmVar(rmVar)
		if err != nil {
			t.Errorf("rmvar returned an error: %v", err)
		}

		if slices.Contains(p.Vars, "EXISTING_VAR3") {
			t.Errorf("%v was not removed from Vars", rmVar)
		}
	})

	t.Run("removing a mix of new and existing vars with spaces", func(t *testing.T) {
		p := mockProject()
		rmVar := " EXISTING_VAR1, EXISTING_VAR3"
		err := p.rmVar(rmVar)
		if err != nil {
			t.Errorf("rmvar returned an error: %v", err)
		}

		if slices.Contains(p.Vars, "EXISTING_VAR3") || slices.Contains(p.Vars, "EXISTING_VAR1") {
			t.Errorf("%v was not removed from Vars", rmVar)
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
