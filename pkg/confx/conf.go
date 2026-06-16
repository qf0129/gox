package confx

import (
	"os"

	"github.com/qf0129/gox/pkg/jsonx"
)

func LoadJsonConfig(jsonFile string, target any) error {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	err = jsonx.Unmarshal(data, target)
	if err != nil {
		return err
	}
	return nil
}
