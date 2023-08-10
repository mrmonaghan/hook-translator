package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

const SLACK_TYPE string = "slack"

type SlackAction struct {
	Name        string       `yaml:"name"`
	config      *viper.Viper `yaml:"-"`
	Templater   *template.Template
	SlackClient *slack.Client
}

func NewSlackAction(token string, data []byte) (SlackAction, error) {

	var action SlackAction

	if token == "" {
		return action, fmt.Errorf("SLACK_TOKEN environment variable is not set")
	}

	action.SlackClient = slack.New(token)

	if err := action.UnmarshalYAML(data); err != nil {
		return action, err
	}

	if err := action.initTemplater(); err != nil {
		return action, err
	}

	return action, nil
}

func (s SlackAction) GetName() string {
	return s.Name
}

func (s SlackAction) GetType() string {
	return SLACK_TYPE
}

func (s SlackAction) GetConfig() *viper.Viper {
	return s.config
}

func (s *SlackAction) ParseConfig() error {
	if s.config == nil {
		return fmt.Errorf("config for action '%s' is invalid: cannot be nil", s.Name)
	}

	if s.config.GetString("slack") != "" {
		s.config = s.config.Sub("slack")
	}

	if len(s.config.GetStringSlice("channels")) < 1 {
		return fmt.Errorf("config for action '%s' is invalid: 'slack.channels' must contain at least 1 entry", s.Name)
	}

	if err := s.initTemplater(); err != nil {
		return fmt.Errorf("config for action '%s' is invalid: %w", s.Name, err)
	}

	return nil
}

func (s *SlackAction) UnmarshalYAML(data []byte) error {
	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	s.Name = v.GetString("name")
	s.config = v.Sub("slack")

	return nil
}

func (s *SlackAction) Render(data any) (string, error) {
	var buf bytes.Buffer
	if err := s.Templater.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error templating values for action '%s': %w", s.Name, err)
	}

	return buf.String(), nil
}

func (s *SlackAction) Execute(rendered string) error {
	var slackOpts slack.MsgOption
	if s.GetConfig().Get("blocks") == true {
		var blocks Blocks
		if err := blocks.UnmarshalJSON([]byte(rendered)); err != nil {
			return fmt.Errorf("unable to process blocks for action '%s': %w", s.Name, err)
		}
		slackOpts = slack.MsgOptionBlocks(blocks.Blocks...)
	} else {
		slackOpts = slack.MsgOptionText(rendered, false)
	}

	for _, channel := range s.GetConfig().GetStringSlice("channels") {
		_, _, err := s.SlackClient.PostMessage(channel, slackOpts)
		if err != nil {
			return fmt.Errorf("unable to send slack message for action '%s': %w", s.Name, err)
		}
	}
	return nil
}

func (s *SlackAction) initTemplater() error {
	if s.Templater == nil {
		templater, err := template.New(s.Name).Parse(s.config.GetString("message"))
		if err != nil {
			return fmt.Errorf("unable to initialize templater: %w", err)
		}
		s.Templater = templater
	}
	return nil
}

// Blocks & associated methods allow easier serializing of `blocks` JSON objects returned by the Slack API or obtained from the Block Kit Builder
type Blocks struct {
	Blocks []slack.Block `json:"blocks"`
}
type blockhint struct {
	Typ string `json:"type"`
}

func (b *Blocks) UnmarshalJSON(data []byte) error {
	var proxy struct {
		Blocks []json.RawMessage `json:"blocks"`
	}
	if err := json.Unmarshal(data, &proxy); err != nil {
		return fmt.Errorf(`failed to unmarshal blocks array: %w`, err)
	}
	for _, rawBlock := range proxy.Blocks {
		var hint blockhint
		if err := json.Unmarshal(rawBlock, &hint); err != nil {
			return fmt.Errorf(`failed to unmarshal next block for hint: %w`, err)
		}
		var block slack.Block
		switch hint.Typ {
		case "actions":
			block = &slack.ActionBlock{}
		case "context":
			block = &slack.ContextBlock{}
		case "divider":
			block = &slack.DividerBlock{}
		case "file":
			block = &slack.FileBlock{}
		case "header":
			block = &slack.HeaderBlock{}
		case "image":
			block = &slack.ImageBlock{}
		case "input":
			block = &slack.InputBlock{}
		case "section":
			block = &slack.SectionBlock{}
		default:
			block = &slack.UnknownBlock{}
		}
		if err := json.Unmarshal(rawBlock, block); err != nil {
			return fmt.Errorf(`failed to unmarshal next block: %w`, err)
		}
		b.Blocks = append(b.Blocks, block)
	}
	return nil
}
