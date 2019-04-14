package screepsapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Module struct {
	Name   string
	Data   []byte
	Binary bool
}

type uploadRequest struct {
	Hash    int64                  `json:"_hash"`
	Branch  string                 `json:"branch"`
	Modules map[string]interface{} `json:"modules"`
}

type uploadResponse struct {
	OK    int    `json:"ok"`
	Error string `json:"error"`
}

func (c *Client) UploadBranch(branch string, modules []*Module) error {

	reqBody := &uploadRequest{
		Hash:    time.Now().UnixNano(),
		Branch:  branch,
		Modules: map[string]interface{}{},
	}

	for _, module := range modules {
		if module.Binary {
			reqBody.Modules[module.Name] = map[string]string{
				"binary": base64.StdEncoding.EncodeToString(module.Data),
			}
		} else {
			reqBody.Modules[module.Name] = string(module.Data)
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

	req.Header.Set("X-Token", c.token)
	req.Header.Set("X-Username", c.token)
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
