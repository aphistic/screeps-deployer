package screepsapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

type branchesResponse struct {
	OK   int               `json:"ok"`
	List []*branchesBranch `json:"list"`
}

type branchesBranch struct {
	ID          string `json:"_id"`
	Name        string `json:"branch"`
	ActiveSim   bool   `json:"activeSim"`
	ActiveWorld bool   `json:"activeWorld"`
}

type Branch struct {
	Name        string
	ActiveSim   bool
	ActiveWorld bool
}

func (c *Client) Branches() ([]*Branch, error) {
	reqURL := path.Join(c.baseURL(), "api/user/branches")

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	branchRes := &branchesResponse{}
	err = json.Unmarshal(resBody, &branchRes)
	if err != nil {
		return nil, err
	}

	if branchRes.OK != 1 {
		return nil, fmt.Errorf("could not get branches")
	}

	var branches []*Branch
	for _, branch := range branchRes.List {
		branches = append(branches, &Branch{
			Name: branch.Name,
		})
	}

	return branches, nil
}

type cloneBranchRequest struct {
	Branch  string `json:"branch"`
	NewName string `json:"newName"`
}

func (c *Client) CloneBranch(sourceBranch string, targetBranch string) error {
	reqURL := path.Join(c.baseURL(), "api/user/clone-branch")

	reqBody := &cloneBranchRequest{
		Branch:  sourceBranch,
		NewName: targetBranch,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewBuffer(reqData)
	req, err := http.NewRequest(http.MethodPost, reqURL, bodyReader)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cloneRes := &basicResponse{}
	err = json.Unmarshal(bodyData, &cloneRes)
	if err != nil {
		return err
	}

	if cloneRes.OK != 1 {
		return fmt.Errorf("error cloning branch")
	}

	return nil
}
