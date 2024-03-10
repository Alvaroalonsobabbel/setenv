# SetEnv for 1Password

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Alvaroalonsobabbel/setenv) ![Test](https://github.com/Alvaroalonsobabbel/setenv/actions/workflows/go-test.yml/badge.svg)

It creates a `.env` file with Variables pointing to 1Password credentials and a `env.json` file that stores all your current project configuation.

1Password documentation on how to reference secrets can be found [here](https://developer.1password.com/docs/cli/secret-references).

## Install

Copy the latest release of the `setenv` file to your `/usr/local/bin` et voilà.

## Usage Examples

See the help:

```bash
setenv -help
```

Set up a new project:

```bash
setenv -vault="My Vault" -item=Item -addvars=DB_USER,DB_PASSWORD -stagekey=item -stage=test
```

View the current project's status:

```bash
setenv -view
```

## Adding and removing variables

You can add variables using the `-addvar=` flag.
Every 1Password credential has a name and a value. When setting a new variable you'll poit the variable name to the 1Password credential name and that will render the value.
The  `-addvar=` flag takes one or multiple comma separated vars and the 1Password credential name can be optionally passed separating the variable and the name with a colon. ie:

```bash
setenv -addvar=DB_USER:my_db_user
```

This will set a .env like this:

```bash
DB_USER="op://Project Vault/Project Item/my_db_user"
```

If no value is passed, the variable name will be used as the 1Pwd credential name, ie:

```bash
setenv -addvar=DB_USER
DB_USER="op://Project Vault/Project Item/DB_USER"
```

Removing vars can be done with the `-rmvar` flag.
Simply pass the var of list of comma separated vars you want to delete. No need to pass the 1pwd credential name.

```bash
setenv -rmvar=DB_USER,DB_PASSWORD
```

## Understanding the Stage and StageKey values

Usually you'll have a set of credentials for every stage of the project: **test**, **staging**, **prod**. In order to differentiate them when working on your project, you can choose in which part of the 1Password address structure you want the stage to be added: **vault**, **item**, **vars**.

If you have a vault called *My Project* with three items inside named *Project Vars-test*, *Project Vars-staging* and *Project Vars-prod*, you can set the stagekey to `item` and then you can switch between stages using the `-stage` flag. You'll end up with a `.env` file like this:

```bash
DB_PASSWORD="op://My Project/Project Vars-test/DB_PASSWORD"
DB_USER="op://My Project/Project Vars-test/DB_USER"
```

Switching stages:

```bash
$ setenv -stage=staging
.env has been updated!

DB_PASSWORD="op://My Project/Project Vars-staging/DB_PASSWORD"
DB_USER="op://My Project/Project Vars-staging/DB_USER"
```

AWS Stage:

By using the `-stage=aws` command the stage will be set to `$AWS_ENV`.

```bash
$ setenv -stage=aws
.env has been updated!

DB_PASSWORD="op://My Project/Project Vars-$AWS_ENV/DB_PASSWORD"
DB_USER="op://My Project/Project Vars-$AWS_ENV/DB_USER"
```

## Handling values with spaces

When using values with spaces, enclose them in double quotes, ie:

```bash
setenv -vault=MyVault
setenv -vault="My Vault"

setenv -addvar=DB_USER,DB_PASSWORD
setenv -addvar"DB USER, DB PASSWORD"
```
