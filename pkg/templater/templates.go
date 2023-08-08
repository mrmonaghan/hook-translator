package templater

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	SLACK_TYPE string = "slack"
	HTTP_TYPE  string = "http"
)

type Executor interface {
	Execute(string) error
}

type Template struct {
	Name           string
	Type           string
	Config         *viper.Viper
	TemplateString string `yaml:"data"`
	Template       *template.Template
}

func (t *Template) Render(data any) (string, error) {

	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func NewTemplate(path string) (Template, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return Template{}, fmt.Errorf("unable to read file: %w", err)
	}

	var t Template
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return t, fmt.Errorf("unable to unmarshal template: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return t, fmt.Errorf("unable to ready config: %w", err)
	}

	config := v.Sub("config")
	if config == nil {
		return t, fmt.Errorf("error parsing template '%s': 'config' field is not present", path)
	}

	t.Config = config

	tmpl, err := template.New(t.Name).Parse(t.TemplateString)
	if err != nil {
		return t, fmt.Errorf("unable to parse template data for template '%s': %w", t.Name, err)
	}

	t.Template = tmpl

	return t, nil
}

func GetYamlFiles(dir string) ([]string, error) {

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return []string{}, fmt.Errorf("unable to get absolute filepath for dir '%s': %w", dir, err)
	}

	files := os.DirFS(absDir)

	var s []string
	err = fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			ext := filepath.Ext(d.Name())
			if ext == ".yaml" || ext == ".yml" {
				s = append(s, fmt.Sprintf("%s/%s", absDir, d.Name()))
			}
		}
		return nil
	})

	if err != nil {
		return s, err
	}

	return s, nil
}

func LoadTemplates(dir string) ([]Template, error) {
	tmplFiles, err := GetYamlFiles(dir)
	if err != nil {
		return []Template{}, err
	}

	var templates []Template

	for _, file := range tmplFiles {
		t, err := NewTemplate(file)
		if err != nil {
			return templates, err
		}
		templates = append(templates, t)
	}

	return templates, nil
}
