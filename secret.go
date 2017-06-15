package zenkit

import (
	"github.com/spf13/afero"
)

func LoadSecretFromFile(filename string) ([]byte, error) {
	fs := afero.NewReadOnlyFs(afero.NewOsFs())
	return afero.ReadFile(fs, filename)
}
