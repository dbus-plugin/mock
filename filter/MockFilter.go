package filter

import (
	"github.com/funkygao/dbus/engine"
	conf "github.com/funkygao/jsconf"
)

var (
	_ engine.Filter = &MockFilter{}
)

type MockFilter struct {
}

func (this *MockFilter) Init(config *conf.Conf) {
}

func (this *MockFilter) Run(r engine.FilterRunner, h engine.PluginHelper) error {
	for pack := range r.Exchange().InChan() {
		pack.Recycle()
	}

	return nil
}

func init() {
	engine.RegisterPlugin("MockFilter", func() engine.Plugin {
		return new(MockFilter)
	})
}
