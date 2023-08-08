package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mrmonaghan/hook-translator/pkg/templater"
	"go.uber.org/zap"
)

type RuleHandler struct {
	Rules  map[string]templater.Rule
	Logger *zap.SugaredLogger
}

func (h *RuleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	h.Logger.Debug("template data", data)

	for _, template := range rule.Templates {

		rendered, err := template.Render(data)
		if err != nil {
			fmt.Println(fmt.Errorf("error rendering template '%s': %w", template.Name, err))
		}
		h.Logger.Debugw("rendered template", "name", template.Name, "rendered_data", rendered)
	}

}
