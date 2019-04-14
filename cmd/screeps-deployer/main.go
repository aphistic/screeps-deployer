package main

import (
	"fmt"
	"os"
	"path"

	"github.com/aphistic/screeps-deployer/internal/app/screeps-deployer/actionenv"
	"github.com/aphistic/screeps-deployer/internal/app/screeps-deployer/deploy"
	"github.com/aphistic/screeps-deployer/internal/app/screeps-deployer/uploader"
)

func main() {
	env := actionenv.NewRealEnvReader()

	workspace, ok := env.LookupEnv("GITHUB_WORKSPACE")
	if !ok {
		fmt.Fprintf(os.Stderr, "Could not find workspace path\n")
		os.Exit(1)
	}

	manifestPath := path.Join(workspace, "screeps.yml")
	dep, err := deploy.LoadDeploymentFromFile(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load deployment manifest: %s\n", err)
		os.Exit(1)
	}

	token, ok := env.LookupEnv("SCREEPS_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "Could not find screeps API token 'SCREEPS_TOKEN'\n")
		os.Exit(1)
	}

	ref, ok := env.LookupEnv("GITHUB_REF")
	if !ok {
		fmt.Fprintf(os.Stderr, "Could not find github ref\n")
		os.Exit(1)
	}

	branch, err := actionenv.ParseBranch(ref)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse branch name: %s\n", err)
		os.Exit(1)
	}

	u := uploader.NewUploader(uploader.WithToken(token))
	err = u.Upload(branch, workspace, dep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not upload branch '%s': %s\n", branch, err)
		os.Exit(1)
	}

	fmt.Printf("Uploaded to branch '%s'\n", branch)
}
