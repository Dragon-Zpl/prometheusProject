package prometheus

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
)

type QueryType int

const (
	Inc QueryType = iota
	Dec
	Timing
	Set
	Add

)

var vecMap  map[string]int

func init() {
	vecMap = make(map[string]int)
}

type Prom struct {
	timer   map[string]*prometheus.HistogramVec
	counter map[string]*prometheus.CounterVec
	state   map[string]*prometheus.GaugeVec
}



func New() *Prom {
	return &Prom{}
}

func (p *Prom) WithTimer(name string, labels []string) error {
	if p.timer == nil {
		p.timer = make(map[string]*prometheus.HistogramVec)
	}
	if vecMap == nil {
		vecMap = make(map[string]int)
	}
	if _, ok := vecMap[name]; ok {
		return errors.New("该指标已存在")
	}
	p.timer[name] = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: name,
		Help: name,
	}, labels)
	PrometheusRegistry.MustRegister(p.timer[name])
	vecMap[name] = len(labels)
	return nil
}

func (p *Prom) WithCounter(name string, labels []string) error {
	if p.counter == nil {
		p.counter = make(map[string]*prometheus.CounterVec)
	}
	if vecMap == nil {
		vecMap = make(map[string]int)
	}
	if _, ok := vecMap[name]; ok {
		return errors.New("该指标已存在")
	}

	p.counter[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: name,
	}, labels)
	PrometheusRegistry.MustRegister(p.counter[name])
	vecMap[name] = len(labels)
	return nil
}

func (p *Prom) WithState(name string, labels []string) error {
	if p.state == nil {
		p.state = make(map[string]*prometheus.GaugeVec)
	}
	if vecMap == nil {
		vecMap = make(map[string]int)
	}
	if _, ok := vecMap[name]; ok {
		return errors.New("该指标已存在")
	}
	p.state[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: name,
	}, labels)
	PrometheusRegistry.MustRegister(p.state[name])
	vecMap[name] = len(labels)
	return nil
}

func (p *Prom) Timing(labels []string, time int64, name string) {

	if p.timer != nil {
		if v, ok := p.timer[name]; ok {
			v.WithLabelValues(labels...).Observe(float64(time))
		}
	}
}

func (p *Prom) Incr(name string, labels []string) {
	if p.counter != nil {
		if v, ok := p.counter[name]; ok {
			v.WithLabelValues(labels...).Inc()
		}
	}

	if p.state != nil {
		if v, ok := p.state[name]; ok {
			v.WithLabelValues(labels...).Inc()
		}
	}

}

func (p *Prom) Decr(name string, labels []string) {
	if p.state != nil {
		if v, ok := p.state[name]; ok {
			v.WithLabelValues(labels...).Dec()
		}
	}
}

func (p *Prom) State(labels []string, v int64, name string) {
	if p.state != nil {
		if vec, ok := p.state[name]; ok {
			vec.WithLabelValues(labels...).Set(float64(v))
		}
	}
}

func (p *Prom) Add(labels []string, v int64, name string) {
	if p.counter != nil {
		if vec, ok := p.counter[name]; ok{
			vec.WithLabelValues(labels...).Add(float64(v))
		}
	}

	if p.state != nil {
		if vec, ok := p.state[name]; ok {
			vec.WithLabelValues(labels...).Add(float64(v))
		}
	}

}

func (p *Prom) UnRegister(typ,name string) error {
	if vecMap == nil {
		return errors.New("不存在该类型指标")
	}
	switch typ {
	case "counter":
		if p.counter == nil {
			return errors.New("不存在该类型指标")
		} else {
			if vec, ok := p.counter[name]; ok {
				PrometheusRegistry.UnRegister(vec)
			}
		}
	case "state":
		if p.state == nil {
			return errors.New("不存在该类型指标")
		} else {
			if vec, ok := p.state[name]; ok {
				PrometheusRegistry.UnRegister(vec)
			}
		}
	case "time":
		if p.timer == nil {
			return errors.New("不存在该类型指标")
		} else {
			if vec, ok := p.timer[name]; ok {
				PrometheusRegistry.UnRegister(vec)
			}
		}
	default:
		return errors.New("不存在该类型指标")
	}
	return nil
}


func PrometheusOpeartor(jobName, name string, v int64 ,lables []string, opeator QueryType) error {
	if prom, ok := RegisterPromMap[jobName]; ok {
		if lables_len, ok := vecMap[name]; !ok {
			return errors.New("该指标不存在，请前往注册")
		} else {
			if lables_len != len(lables) {
				return errors.New("labels个数有问题")
			}
		}
		prom.Lock.Lock()
		defer prom.Lock.Unlock()
		switch opeator {
		case Inc:
			prom.Vec.Incr(name, lables)
		case Dec:
			prom.Vec.Decr(name, lables)
		case Add:
			prom.Vec.Add(lables, v, name)
		case Set:
			prom.Vec.Add(lables, v, name)
		case Timing:
			prom.Vec.Timing(lables, v, name)
		default:
			return errors.New("opeartor is not exist")
		}
		return nil
	}
	return errors.New("jobName is not exist")
}

func DeleteVec(name string)  {
	delete(vecMap, name)
}