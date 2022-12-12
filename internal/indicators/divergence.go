package indicators

type stockPoint struct {
	index int
	price float64
}

func findLocalMins(data []float64) []stockPoint {
	stockPoints := []stockPoint{}
	negativeAreas := [][]float64{}
	negativeAreasTemp := []float64{}

	for i, val := range data {
		if val < 0 {
			negativeAreasTemp = append(negativeAreasTemp, val)
		} else if x >= 0 {
			negativeAreas = append(negativeAreas, negativeAreasTemp)
			negativeAreasTemp = nil
		}
	}

	if len(negativeAreasTemp) != 0 {
		negativeAreas = append(negativeAreas, negativeAreasTemp)
	}

	for _, area := range negativeAreas {

	}

	return stockPoints
}


func Divergence(data []float64) (bool, bool) {
	return false, false
}
