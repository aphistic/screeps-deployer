package deploy

import (
	"fmt"
	"strings"

	"github.com/go-yaml/yaml"
)

type Manifest struct {
	Modules []*ManifestModule `yaml:"modules"`
}

type ManifestModule struct {
	Name   string `yaml:"name"`
	File   string `yaml:"file"`
	Binary bool   `yaml:"binary"`
}

func ParseManifest(data []byte) (*Manifest, error) {
	m := &Manifest{}

	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manifest) Validate() error {
	foundModules := map[string]bool{}

	hasMain := false
	for _, module := range m.Modules {
		lowerModule := strings.ToLower(module.Name)

		_, ok := foundModules[lowerModule]
		if ok {
			return fmt.Errorf("found duplicate module name '%s'", module.Name)
		}
		foundModules[lowerModule] = true

		if module.File == "" {
			return fmt.Errorf("no file provided for module '%s'", module.Name)
		}

		if lowerModule == "main" {
			hasMain = true
		}
	}

	if !hasMain {
		return fmt.Errorf("could not find module named 'main'")
	}

	return nil
}