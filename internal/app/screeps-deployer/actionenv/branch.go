package actionenv

import (
	"fmt"
	"strings"
)

func ParseBranch(refName string) (string, error) {
	prefixes := []string{
		"refs/heads/",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(refName, prefix) {
			return refName[len(prefix):], nil
		}
	}

	return "", fmt.Errorf("could not parse branch name '%s'", refName)
}
