package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mrmonaghan/hook-translator/internal/rules"
	"go.uber.org/zap"
)

type RuleHandler struct {
	Rules  map[string]rules.Rule
	Logger *zap.SugaredLogger
}

func (h *RuleHandler) HandleWebhooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	requestedRule, err := h.getRequestedRule(r)
	if err != nil {
		h.Logger.Errorw("unable to determine rule for request", "url", r.URL.RequestURI(), "headers", r.Header, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := h.Rules[requestedRule]; !ok {
		h.Logger.Debugw("request rule not found", "rule", requestedRule)
		http.Error(w, fmt.Sprintf("rule '%s' not found", requestedRule), http.StatusNotFound)
		return
	}
	rule, _ := h.Rules[requestedRule]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorw("error reading request body", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	if err := json.Unmarshal(b, &data); err != nil {
		h.Logger.Errorw("error unmarshalling request body", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	for _, template := range rule.GetTemplates() {
		for _, action := range template.Actions {
			rendered, err := action.Render(data)
			if err != nil {
				h.Logger.Errorw("error rendering template action", "template", template.Name, "action", action.GetName(), err)
				http.Error(w, fmt.Sprintf("unable to render template '%s' action '%s'", template.Name, action.GetName()), http.StatusBadRequest)
			}
			h.Logger.Debugw("successfully rendered template action", "template", template.Name, "action", action.GetName())

			if err := action.Execute(rendered); err != nil {
				h.Logger.Errorf("error executing action", "action", action.GetName(), err)
			}
			h.Logger.Debugw("successfully executed action", "template", template.Name, "action", action.GetName())
		}

	}

	http.Error(w, "", http.StatusOK)
	return
}

func (h *RuleHandler) HandleRules(w http.ResponseWriter, r *http.Request) {
	var rules []rules.Rule

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

func (h *RuleHandler) getRequestedRule(r *http.Request) (string, error) {
	if requestedRule := r.URL.Query().Get("rule"); requestedRule != "" {
		return requestedRule, nil
	} else if requestedRule := r.Header.Get("rule"); requestedRule == "" {
		return requestedRule, nil
	} else {
		return "", errors.New("no 'stitch_rule' parameter included in URL or request headers")
	}
}
