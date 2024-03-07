# SetEnv for 1Password

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Alvaroalonsobabbel/setenv) ![Test](https://github.com/Alvaroalonsobabbel/setenv/actions/workflows/go-test.yml/badge.svg)

It creates a `.env` file with Variables pointing to 1Password credentials and a `env.json` file that stores all your current project configuation.

Every 1Password credential has a name and a value. This assumes you named the 1Password credential the same as your env variable.

ie:

```bash
API_KEY="op://Project Vault/Project Item/API_KEY"
```

1Password documentation on how to reference secrets can be found [here](https://developer.1password.com/docs/cli/secret-references).

## Install

Copy the `setenv` file to your `/usr/local/bin` et voilà.

## Usage

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

## Tips

When using values with spaces, enclose them in double quotes, ie:

```bash
setenv -vault=MyVault
setenv -vault="My Vault"

setenv -addvar=DB_USER,DB_PASSWORD
setenv -addvar"DB USER, DB PASSWORD"
```
