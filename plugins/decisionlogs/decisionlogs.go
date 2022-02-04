package decisionlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/logs"

	"github.com/open-policy-agent/opa/util"
	"github.com/tidwall/gjson"
)

const PluginName = "print_decision_logs_on_failure"

func New(m *plugins.Manager, config interface{}) plugins.Plugin {

	m.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})

	return &PrintlnLogger{
		manager: m,
		config:  config.(Config),
	}
}

func Validate(_ *plugins.Manager, config []byte) (interface{}, error) {
	parsedConfig := Config{}
	return parsedConfig, util.Unmarshal(config, &parsedConfig)
}

type Config struct {
	Stdout bool `json:"stdout"` // true => stdout, false => stderr
}

type PrintlnLogger struct {
	manager *plugins.Manager
	mtx     sync.Mutex
	config  Config
}

func (p *PrintlnLogger) Start(ctx context.Context) error {
	p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateOK})
	return nil
}

func (p *PrintlnLogger) Stop(ctx context.Context) {
	p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})
}

func (p *PrintlnLogger) Reconfigure(ctx context.Context, config interface{}) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.config = config.(Config)
}

// Log is called by the decision logger when a record (event) should be emitted. The logs.EventV1 fields
// map 1:1 to those described in https://www.openpolicyagent.org/docs/latest/management-decision-logs
func (p *PrintlnLogger) Log(ctx context.Context, event logs.EventV1) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	w := os.Stdout
	if !p.config.Stdout {
		w = os.Stderr
	}
	bs, err := json.Marshal(event)
	if err != nil {
		p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
		return nil
	}
	result := gjson.Get(string(bs), "result.allowed")

	// Print the decision logs only when result.allowed is false
	if !result.Bool() {
		_, err = fmt.Fprintln(w, string(bs))
	}

	if err != nil {
		p.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
	}
	return nil
}
