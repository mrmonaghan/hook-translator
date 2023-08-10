package templates

import (
	"fmt"
	"os"

	"github.com/mrmonaghan/hook-translator/internal/actions"
	"github.com/mrmonaghan/hook-translator/internal/utils"
	"gopkg.in/yaml.v3"
)

type Template struct {
	Name       string                   `yaml:"name"`
	RawActions []map[string]interface{} `yaml:"actions"`
	Actions    []actions.Action         `yaml:"-"`
}

func (t *Template) UnmarshalYAML(data []byte) error {

	if err := yaml.Unmarshal(data, t); err != nil {
		return fmt.Errorf("unable to unmarshal template: %w", err)
	}

	for _, rawAction := range t.RawActions {
		actionType, err := actions.GetActionTypeFromInterface(rawAction)
		if err != nil {
			return fmt.Errorf("unable to unmarshal template '%s': %w", t.Name, err)
		}

		switch actionType {
		case actions.SLACK_TYPE:

			b, err := yaml.Marshal(rawAction)
			if err != nil {
				return fmt.Errorf("unable to unmarshal template '%s' actions: %w", t.Name, err)
			}
			a, err := actions.NewSlackAction(os.Getenv("SLACK_TOKEN"), b)
			if err != nil {
				return fmt.Errorf("unable to unmarshal template '%s' action: %w", t.Name, err)

			}

			if err := a.ParseConfig(); err != nil {
				return fmt.Errorf("unable to unmarshal template '%s': %w", t.Name, err)
			}

			t.Actions = append(t.Actions, &a)
			break
		case actions.HTTP_TYPE:
			fmt.Println("HTTP_TYPE")
			break
		default:
			fmt.Println("unknown type")
		}

	}

	return nil
}

func LoadTemplates(dir string) ([]Template, error) {
	tmplFiles, err := utils.GetYamlFileNamesFromDir(dir)
	if err != nil {
		return []Template{}, err
	}

	var templates []Template

	for _, file := range tmplFiles {
		b, err := os.ReadFile(file)
		if err != nil {
			return templates, fmt.Errorf("error reading template file '%s': %w", file, err)
		}

		var t Template

		if err := t.UnmarshalYAML(b); err != nil {
			return templates, err
		}

		templates = append(templates, t)

	}

	return templates, nil
}
