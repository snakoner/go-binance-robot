package strategies

type StrategyElement struct {
  Func []func([]float64) (bool, bool)
  Long bool
  Short bool
}

type Strategy struct {
  Elements []StrategyElement
}

func NewStrategy() *Strategy {
  return &Strategy {}
}

func (this *Strategy) Add(f func([]float64)(bool, bool)) {
  se := &StrategyElement {
    Func : f,
  }
  this.Elements = append(this.Elements, se)
}

// Calculate if long, short OK for each strategy element
func (this *Strategy) Apply(data []float64) {
  for _, se := range this.Elements {
    se.Long, se.Short = se.Func(data)
  }
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
