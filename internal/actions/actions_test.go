package actions

import (
	"testing"

	"gopkg.in/yaml.v3"
)

var slackYAML = []byte(`
name: test-slack-action
slack:
  channels: [test-slack-channel]
  message: |
    this is a test message {{ .testValue }}`)

var slackYAMLEmptyChannels = []byte(`
name: test-slack-action
slack:
  channels: []
`)

var httpYAML = []byte(`
name: test-http-action
http:
  method: GET
  path: https://fake-url.com
`)

var invalidActionYAML = []byte(`
name: test-invalid-action
invalid:
  bad-option: true
`)

var invalidYAML = []byte(`
'name: test-invalid-action
    invalid:
  bad-option: true
`)

func TestGetActionTypeFromInterface(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name    string
		args    func() (args, error)
		want    string
		wantErr bool
	}{
		{
			name: "happy path - slack",
			args: func() (args, error) {
				var m map[string]interface{}
				if err := yaml.Unmarshal(slackYAML, &m); err != nil {
					return args{m}, err
				}

				return args{m}, nil
			},
			want:    "slack",
			wantErr: false,
		},
		{
			name: "happy path - http",
			args: func() (args, error) {
				var m map[string]interface{}
				if err := yaml.Unmarshal(httpYAML, &m); err != nil {
					return args{m}, err
				}

				return args{m}, nil
			},
			want:    "http",
			wantErr: false,
		},
		{
			name: "invalid action type",
			args: func() (args, error) {
				var m map[string]interface{}
				if err := yaml.Unmarshal(invalidActionYAML, &m); err != nil {
					return args{m}, err
				}

				return args{m}, nil
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			args, err := tt.args()
			if err != nil {
				t.Errorf("error setting up args for test '%s': %s", tt.name, err.Error())
			}

			got, err := GetActionTypeFromInterface(args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActionTypeFromInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetActionTypeFromInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}
