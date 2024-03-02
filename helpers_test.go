package main

import (
	"slices"
	"testing"
)

func TestSanitizeVar(t *testing.T) {
	t.Run("single variable", func(t *testing.T) {
		stringvar := "VAR1"
		vars, err := sanitinzeVars(stringvar)
		if err != nil {
			t.Errorf("Error sanitizing Var")
		}
		if !slices.Contains(vars, "VAR1") {
			t.Errorf("Var %vhas not been sanitised", stringvar)
		}
	})

	t.Run("multiple variables", func(t *testing.T) {
		stringvar := "VAR1,VAR2"
		vars, err := sanitinzeVars(stringvar)
		if err != nil {
			t.Errorf("Error sanitizing Var")
		}
		if !slices.Contains(vars, "VAR1") || !slices.Contains(vars, "VAR2") {
			t.Errorf("Var %vhas not been sanitised", stringvar)
		}
	})

	t.Run("multiple variables with trailing spaces", func(t *testing.T) {
		stringvar := "    VAR1,VAR2    "
		vars, err := sanitinzeVars(stringvar)
		if err != nil {
			t.Errorf("Error sanitizing Var")
		}
		if !slices.Contains(vars, "VAR1") || !slices.Contains(vars, "VAR2") {
			t.Errorf("Var %vhas not been sanitised", stringvar)
		}
	})
}
