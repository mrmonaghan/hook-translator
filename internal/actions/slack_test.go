package actions

import (
	"bytes"
	"encoding/json"
	"testing"
	"text/template"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

func TestSlackAction_GetName(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				Name: "test-action",
			},
			want: "test-action",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SlackAction{
				Name: tt.fields.Name,
			}
			if got := s.GetName(); got != tt.want {
				t.Errorf("SlackAction.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackAction_GetType(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				Name: "test-action",
			},
			want: "slack",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SlackAction{
				Name: tt.fields.Name,
			}
			if got := s.GetType(); got != tt.want {
				t.Errorf("SlackAction.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackAction_ParseConfig(t *testing.T) {
	type fields struct {
		config *viper.Viper
	}
	tests := []struct {
		name    string
		setup   func() (*viper.Viper, error)
		wantErr bool
	}{
		{
			name: "happy path",
			setup: func() (*viper.Viper, error) {
				v := viper.New()
				v.SetConfigType("yaml")
				if err := v.ReadConfig(bytes.NewBuffer(slackYAML)); err != nil {
					return v, err
				}
				return v.Sub("slack"), nil
			},
			wantErr: false,
		},
		{
			name: "config is nil",
			setup: func() (*viper.Viper, error) {
				return nil, nil
			},
			wantErr: true,
		},
		{
			name: "channels are empty",
			setup: func() (*viper.Viper, error) {
				v := viper.New()
				v.SetConfigType("yaml")
				if err := v.ReadConfig(bytes.NewBuffer(slackYAMLEmptyChannels)); err != nil {
					return v, err
				}
				return v.Sub("slack"), nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cfg, err := tt.setup()
			if err != nil {
				t.Errorf("error setting up test '%s': %s", tt.name, err.Error())
			}

			s := &SlackAction{
				config: cfg,
			}

			if err := s.ParseConfig(); (err != nil) != tt.wantErr {
				t.Errorf("SlackAction.ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlackAction_UnmarshalYAML(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid YAML",
			args: args{
				data: slackYAML,
			},
			wantErr: false,
		},
		{
			name: "invalid YAML",
			args: args{
				data: invalidYAML,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SlackAction{
				Name: tt.name,
			}
			if err := s.UnmarshalYAML(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SlackAction.UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSlackAction_Render(t *testing.T) {
	type fields struct {
		Name        string
		Templater   *template.Template
		SlackClient *slack.Client
	}
	type args struct {
		data any
	}
	tests := []struct {
		name        string
		configSetup func() (*viper.Viper, error)
		dataSetup   func() (any, error)
		fields      fields
		want        string
		wantErr     bool
	}{
		{
			name: "happy path",
			configSetup: func() (*viper.Viper, error) {
				v := viper.New()
				v.SetConfigType("yaml")
				if err := v.ReadConfig(bytes.NewBuffer(slackYAML)); err != nil {
					return nil, err
				}

				return v.Sub("slack"), nil
			},
			dataSetup: func() (any, error) {

				var tmplData = []byte(`
					{
						"testValue": "test_value",
						"key": {
							"nested_key": "nested_value"
						}
					}`)

				m := make(map[string]interface{})
				if err := json.Unmarshal(tmplData, &m); err != nil {
					return m, err
				}
				return m, nil
			},
			want:    "this is a test message test_value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cfg, err := tt.configSetup()
			if err != nil {
				t.Errorf("error setting up config for test '%s': %s", tt.name, err.Error())
			}

			s := &SlackAction{
				Name:   tt.name,
				config: cfg,
			}
			templater, err := template.New(tt.name).Parse(s.GetConfig().GetString("message"))
			if err != nil {
				t.Errorf("error setting up templater for test '%s': %s", tt.name, err.Error())
			}
			s.Templater = templater

			data, err := tt.dataSetup()
			if err != nil {
				t.Errorf("error setting up template data for test '%s': %s", tt.name, err.Error())
			}

			got, err := s.Render(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackAction.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SlackAction.Render() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlocks_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Blocks []slack.Block
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "happy path",
			fields: fields{},
			args: args{
				data: []byte(` {
					"blocks": [
					  {
						"type": "section",
						"text": {
						  "type": "mrkdwn",
						  "text": "We found *205 Hotels* in New Orleans, LA from *12/14 to 12/17* {{ .name }}"
						},
						"accessory": {
						  "type": "overflow",
						  "options": [
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option One"
							  },
							  "value": "value-0"
							},
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option Two"
							  },
							  "value": "value-1"
							},
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option Three"
							  },
							  "value": "value-2"
							},
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option Four"
							  },
							  "value": "value-3"
							}
						  ]
						}
					  }
					]
				  }`),
			},
			wantErr: false,
		},
		{
			name:   "invalid JSON",
			fields: fields{},
			args: args{
				data: []byte(` {
					blocks": [
					  {
						"something_fake": "section",
						"idk": {
						  "type": "mrkdwn",
						  "text": "We found *205 Hotels* in New Orleans, LA from *12/14 to 12/17* {{ .name }}"
						},
						"accessory": {
						  "type": "overflow",
						  "options": [
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option One"
							  },
							  "value": "value-0"
							},
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option Two"
							  },
							  "value": "value-1"
							},
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option Three"
							  },
							  "value": "value-2"
							},
							{
							  "text": {
								"type": "plain_text",
								"emoji": true,
								"text": "Option Four"
							  },
							  "value": "value-3"
							}
						  ]
						}
					  }
					]
				  }`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Blocks{
				Blocks: tt.fields.Blocks,
			}
			if err := b.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Blocks.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
