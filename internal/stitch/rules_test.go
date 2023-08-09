package stitch

import (
	"reflect"
	"testing"
)

func TestRule_Templates(t *testing.T) {
	type fields struct {
		Name          string
		Enabled       bool
		TemplateNames []string
		templates     []Template
	}
	tests := []struct {
		name   string
		fields fields
		want   []Template
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Rule{
				Name:          tt.fields.Name,
				Enabled:       tt.fields.Enabled,
				TemplateNames: tt.fields.TemplateNames,
				templates:     tt.fields.templates,
			}
			if got := r.Templates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rule.Templates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRule(t *testing.T) {
	type args struct {
		path      string
		templates []Template
	}
	tests := []struct {
		name    string
		args    args
		want    Rule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRule(tt.args.path, tt.args.templates)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadRules(t *testing.T) {
	type args struct {
		dir       string
		templates []Template
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]Rule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadRules(tt.args.dir, tt.args.templates)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
