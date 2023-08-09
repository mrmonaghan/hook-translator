package stitch

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/spf13/viper"
)

var testTemplates = []Template{
	{
		Name: "test-template-1",
		Type: "slack",
		Config: 
	}
}

func TestTemplate_Render(t *testing.T) {
	type fields struct {
		Name           string
		Type           string
		Config         *viper.Viper
		TemplateString string
		Template       *template.Template
	}
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Template{
				Name:           tt.fields.Name,
				Type:           tt.fields.Type,
				Config:         tt.fields.Config,
				TemplateString: tt.fields.TemplateString,
				Template:       tt.fields.Template,
			}
			got, err := tr.Render(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Template.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Template.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTemplate(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    Template
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTemplate(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetYamlFiles(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetYamlFiles(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetYamlFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetYamlFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadTemplates(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []Template
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadTemplates(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadTemplates() = %v, want %v", got, tt.want)
			}
		})
	}
}
