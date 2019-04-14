package uploader

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

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

func WithToken(token string) UploadOption {
	return func(u *Uploader) {
		u.token = token
	}
}

type Uploader struct {
	client *screepsapi.Client
	token  string
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

	reqBody := &uploadRequest{
		Hash:    time.Now().UnixNano(),
		Branch:  branch,
		Modules: map[string]interface{}{},
	}

	for _, module := range dep.Manifest.Modules {
		modulePath := path.Join(workspace, module.File)
		data, err := ioutil.ReadFile(modulePath)
		if err != nil {
			return fmt.Errorf(
				"could not read module '%s' file '%s': %s",
				module.Name, module.File, err,
			)
		}

		if module.Binary {
			reqBody.Modules[module.Name] = map[string]string{
				"binary": base64.StdEncoding.EncodeToString(data),
			}
		} else {
			reqBody.Modules[module.Name] = string(data)
		}
	}

	bodyData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(bodyData)
	req, err := http.NewRequest(http.MethodPost, "https://screeps.com/api/user/code", buf)
	if err != nil {
		return err
	}

	req.Header.Set("X-Token", u.token)
	req.Header.Set("X-Username", u.token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("couldn't read body: %s\n", err)
		return err
	}

	resData := &uploadResponse{}
	err = json.Unmarshal(resBody, &resData)
	if err != nil {
		return fmt.Errorf("response unmarshal error: %s", err)
	}

	if resData.OK != 1 {
		return fmt.Errorf("upload error: %s", resData.Error)
	}

	return nil
}
