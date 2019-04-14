package deploy

import (
	"io/ioutil"
)

type Deployment struct {
	Manifest *Manifest
}

func LoadDeploymentFromFile(file string) (*Deployment, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return LoadDeployment(data)
}

func LoadDeployment(manifest []byte) (*Deployment, error) {
	m, err := ParseManifest(manifest)
	if err != nil {
		return nil, err
	}

	err = m.Validate()
	if err != nil {
		return nil, err
	}

	dep := &Deployment{
		Manifest: m,
	}

	return dep, nil
}
