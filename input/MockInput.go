package input

import (
	"runtime"
	"time"

	"github.com/funkygao/dbus/engine"
	"github.com/funkygao/dbus/pkg/model"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

var (
	_ engine.Input     = &MockInput{}
	_ engine.Pauser    = &MockInput{}
	_ engine.Restarter = &MockInput{}
)

type MockInput struct {
	inChan <-chan *engine.Packet

	payload engine.Payloader
	sleep   time.Duration
}

func (this *MockInput) Init(config *conf.Conf) {
	this.sleep = config.Duration("sleep", 0)
	switch config.String("payload", "Bytes") {
	case "RowsEvent":
		this.payload = &model.RowsEvent{
			Log:       "mysql-bin.0001",
			Position:  498876,
			Schema:    "mydabase",
			Table:     "user_account",
			Action:    "I",
			Timestamp: 1486554654,
			Rows:      [][]interface{}{{"user", 15, "hello world"}},
		}

	default:
		this.payload = model.Bytes(`{"log":"6633343-bin.006419","pos":795931083,"db":"owl_t_prod_cd","tbl":"owl_mi","dml":"U","ts":1488934500,"rows":[[132332,"expired_keys","过期的key的个数","mock-monitor|member-mock|10.1.1.1|10489|expired_keys","mock-monitor",244526,"10489",null,"2015-12-25 17:12:00",null,null,"28571284","Ok","TypeLong","2017-03-08 08:54:01",null,null,null,null,null,null],[132332,"expired_keys","过期的key的个数","mock-monitor|member-mock|10.1.1.6|10489|expired_keys","mock-monitor",244526,"10489",null,"2015-12-25 17:12:00",null,null,"28571320","Ok","TypeLong","2017-03-08 08:55:00",null,null,null,null,null,null]]}`)
	}
}

func (this *MockInput) OnAck(pack *engine.Packet) error {
	return nil
}

func (this *MockInput) CleanupForRestart() bool {
	return true
}

func (this *MockInput) Pause(r engine.InputRunner) error {
	log.Info("[%s] paused", r.Name())
	this.inChan = nil
	return nil
}

func (this *MockInput) Resume(r engine.InputRunner) error {
	log.Info("[%s] resumed", r.Name())
	this.inChan = r.Exchange().InChan()
	return nil
}

func (this *MockInput) Run(r engine.InputRunner, h engine.PluginHelper) error {
	this.inChan = r.Exchange().InChan()
	var n int64
	for {
		select {
		case <-r.Stopper():
			log.Info("[%s] %d packets injected", r.Name(), n)
			return nil

		case pack := <-this.inChan:
			pack.Payload = this.payload
			r.Exchange().Inject(pack)
			n++

			if this.sleep > 0 {
				time.Sleep(this.sleep)
			}

		default:
			runtime.Gosched()
		}
	}

	return nil
}

func init() {
	engine.RegisterPlugin("MockInput", func() engine.Plugin {
		return new(MockInput)
	})
}
