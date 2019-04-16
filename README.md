# Screeps Deployer

Deploy code to [Screeps](https://screeps.com) from GitHub actions.

__Note:__ Screeps deployer is currently under active development! I use it for my own deployments
but if you do use it at this point be sure to have backups of your code available just in case!

## Usage

### Upload Manifest
Screeps Deployer looks for a file named `screeps.yml` in the root of your repository to determine
how to upload your code to Screeps. Currently the only thing in the file is a list of modules as
well as the files to upload for those modules.

An example file can be seen below:

```yaml
modules:
  - name: main
    file: main.js
  - name: some.name
    file: some/name.js
```

Binary modules are also supported:

```yaml
modules:
  - name: main
    file: main.js
  - name: app
    file: app.wasm
    binary: true
```

### Action Configuration

The deployer uses a [Screeps auth token](https://docs.screeps.com/auth-tokens.html) to perform actions
against the Screeps API.  To provide this token to the action it's recommended to create a secret in
your deployment repository called `SCREEPS_TOKEN` with a full access auth token.  Next, add the action
to your workflow:

```
action "Deploy to Screeps" {
  uses = "aphistic/screeps-deployer@master"
  secrets = ["SCREEPS_TOKEN"]
}
```

When added to a "push" workflow, Screeps deployer will then upload your code to a branch with the same
name as the one pushed to the repository. For example, if you push to the branch `dev` in your GitHub
repository your code will be uploaded to the `dev` branch in Screeps. The `master` branch is currently
set to upload to the `default` branch in Screeps.