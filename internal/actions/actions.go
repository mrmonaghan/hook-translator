package actions

import (
	"fmt"

	"github.com/spf13/viper"
)

type Action interface {
	GetName() string
	GetType() string
	GetConfig() *viper.Viper
	ParseConfig() error
	UnmarshalYAML([]byte) error
	Render(any) (string, error)
	Execute(string) error
}

const (
	HTTP_TYPE = "http"
)

func GetActionTypeFromInterface(config map[string]interface{}) (string, error) {

	if _, ok := config["slack"]; ok {
		return SLACK_TYPE, nil
	}

	if _, ok := config["http"]; ok {
		return HTTP_TYPE, nil
	}

	return "", fmt.Errorf("error determining action type: must be one of 'slack,http'")
}
