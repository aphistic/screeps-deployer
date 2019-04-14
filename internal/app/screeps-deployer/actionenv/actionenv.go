package actionenv

import (
	"os"
)

type EnvReader interface {
	LookupEnv(string) (string, bool)
}

type RealEnvReader struct {}

func NewRealEnvReader() *RealEnvReader {
	return &RealEnvReader{}
}

func (rer *RealEnvReader) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
