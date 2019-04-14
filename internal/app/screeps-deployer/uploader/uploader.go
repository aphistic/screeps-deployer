package uploader

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/aphistic/screeps-deployer/internal/app/screeps-deployer/consts"
	"github.com/aphistic/screeps-deployer/internal/app/screeps-deployer/deploy"
	"github.com/aphistic/screeps-deployer/internal/app/screeps-deployer/screepsapi"
)

type uploadRequest struct {
	Hash    int64                  `json:"_hash"`
	Branch  string                 `json:"branch"`
	Modules map[string]interface{} `json:"modules"`
}

type uploadResponse struct {
	OK    int    `json:"ok"`
	Error string `json:"error"`
}

type UploadOption func(*Uploader)

func WithClient(client *screepsapi.Client) UploadOption {
	return func(u *Uploader) {
		u.client = client
	}
}

type Uploader struct {
	client *screepsapi.Client
}

func NewUploader(opts ...UploadOption) *Uploader {
	u := &Uploader{
		client: screepsapi.NewClient(),
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func (u *Uploader) Upload(branch string, workspace string, dep *deploy.Deployment) error {
	lowerBranch := strings.ToLower(branch)

	// Get a list of the existing branches so we know if we need to clone a new one or not
	branches, err := u.client.Branches()
	if err != nil {
		return err
	}

	hasBranch := false
	for _, branch := range branches {
		if lowerBranch == strings.ToLower(branch.Name) {
			hasBranch = true
			break
		}
	}

	if !hasBranch {
		// We don't have the branch yet, we need to clone it
		err = u.client.CloneBranch(consts.DefaultBranch, branch)
		if err != nil {
			return fmt.Errorf("could not clone branch for upload: %s", err)
		}
	}

	var modules []*screepsapi.Module
	for _, module := range dep.Manifest.Modules {
		modulePath := path.Join(workspace, module.File)
		data, err := ioutil.ReadFile(modulePath)
		if err != nil {
			return fmt.Errorf(
				"could not read module '%s' file '%s': %s",
				module.Name, module.File, err,
			)
		}

		modules = append(modules, &screepsapi.Module{
			Name:   module.Name,
			Data:   data,
			Binary: module.Binary,
		})
	}

	err = u.client.UploadBranch(branch, modules)
	if err != nil {
		return err
	}

	return nil
}
