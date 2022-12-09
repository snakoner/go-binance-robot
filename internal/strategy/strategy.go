package strategy

import (
	"sync"
)

type StrategyElement struct {
	Name  string
	Func  func([]float64) (bool, bool)
	Long  bool
	Short bool
}

type Strategy struct {
	Elements []*StrategyElement
}

func New() *Strategy {
	return &Strategy{}
}

func (this *Strategy) Add(f func([]float64) (bool, bool), name string) {
	se := &StrategyElement{
		Func: f,
		Name: name,
	}
	this.Elements = append(this.Elements, se)
}

// Calculate if long, short OK for each strategy element
func (this *Strategy) Apply(data []float64) {
	wg := new(sync.WaitGroup)
	for _, se := range this.Elements {
		wg.Add(1)
		go func(se *StrategyElement) {
			se.Long, se.Short = se.Func(data)
			wg.Done()
		}(se)
	}
	wg.Wait()
}

func (this *Strategy) GetName() string {
	result := ""
	first := true
	for _, se := range this.Elements {
		name := se.Name
		if first {
			first = false
			result += name
			continue
		}
		result += ", " + name
	}

	return result
}

func (this *Strategy) IsLong() bool {
	for _, se := range this.Elements {
		if !se.Long {
			return false
		}
	}
	return true
}

func (this *Strategy) IsShort() bool {
	for _, se := range this.Elements {
		if !se.Short {
			return false
		}
	}
	return true
}
