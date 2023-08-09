package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrmonaghan/hook-translator/internal/stitch"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type RuleHandler struct {
	Rules  map[string]stitch.Rule
	Logger *zap.SugaredLogger
	Slack  *slack.Client
}

func (h *RuleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/rules" {

		var rules []stitch.Rule
		for _, rule := range h.Rules {
			rules = append(rules, rule)
		}
		b, err := json.Marshal(rules)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		http.Error(w, string(b), http.StatusOK)
		return
	}

	requestedRule := r.URL.Query().Get("rule")
	if requestedRule == "" {
		h.Logger.Debugw("no 'rule' query parameter included in URL", "url", r.URL.RequestURI())
		requestedRule = r.Header.Get("rule")
		if requestedRule == "" {
			h.Logger.Debugw("no 'rule' header included in request", "headers", r.Header)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "request must include the 'rule' query parameter or 'rule' header", http.StatusBadRequest)
		}
	}

	if _, ok := h.Rules[requestedRule]; !ok {
		h.Logger.Debugw("request rule not found", "rule", requestedRule)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf("rule '%s' not found", requestedRule), http.StatusNotFound)
	}

	rule, _ := h.Rules[requestedRule]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorw("error reading request body", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "", http.StatusInternalServerError)
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(b, &data); err != nil {
		h.Logger.Errorw("error unmarshalling request body", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "", http.StatusInternalServerError)
	}
	h.Logger.Debugw("template data", data)

	for _, template := range rule.Templates() {

		rendered, err := template.Render(data)
		if err != nil {
			fmt.Println(fmt.Errorf("error rendering template '%s': %w", template.Name, err))
		}
		h.Logger.Debugw("rendered template", "name", template.Name, "rendered_data", rendered)

		switch template.Type {
		case "slack":

			var slackOpts slack.MsgOption
			if template.Config.Get("blocks") == true {
				h.Logger.Debugw("template has blocks enabled", "template", template.Name, err)
				var blocks stitch.Blocks

				fmt.Println(rendered)
				if err := blocks.UnmarshalJSON([]byte(rendered)); err != nil {
					h.Logger.Errorw("unable to process blocks", "template", template.Name, err)
				}

				slackOpts = slack.MsgOptionBlocks(blocks.Blocks...)
			} else {
				slackOpts = slack.MsgOptionText(rendered, false)
			}

			for _, channel := range template.Config.GetStringSlice("channels") {
				messageID, _, err := h.Slack.PostMessage(channel, slackOpts)
				if err != nil {
					h.Logger.Errorw("error sending slack template", "template", template.Name, "channel", channel, err)
				}

				h.Logger.Debugw("successfully posted slack message", "template", template.Name, "channel", channel, "messageID", messageID)
			}
		default:
			h.Logger.Debugw("template type is not recognized", "template", template.Name, "type", template.Type)
		}
	}

}
