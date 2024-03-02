package main

import (
	"flag"
)

// Parses all the flags and sends them to the correct variable or method.
func setFlags() {
	flag.Usage = func() {}

	flag.BoolFunc("help", "Shows this Help :)", printHelp)
	flag.StringVar(&projectData.Vault, "vault", projectData.Vault, "Sets the Vault's name for the project.")
	flag.StringVar(&projectData.Item, "item", projectData.Item, "Sets the Item's name for the project.")
	flag.Func("stagekey", "Set the Key that will have the stage name. Allowed values: vault, item, vars", projectData.validateStageKey)
	flag.Func("stage", "Set the stage of the project. Allowed values: test, staging, prod", projectData.validateStage)
	flag.Func("addvar", "Adds one or more vars, comma separated.", projectData.addVar)
	flag.Func("rmvar", "Removes one or more vars, comma separated.", projectData.rmVar)
	flag.BoolVar(&projectData.tfvars, "tfvars", false, "Prefixes 'TF_VAR_' to every variable name")
	flag.BoolFunc("view", "Prints the current Project's status and .env", displayProject)
	flag.BoolFunc("ignore", "Adds the project's files (.env, env.json) to .gitignore", addToGitIgnore)
	flag.BoolFunc("clean", "Deletes the project's files (.env, env.json)", cleanProject)
}
