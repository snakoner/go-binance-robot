package testing

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func ReadCloseData(filename string) ([]float64, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return []float64{}, err
	}

	ret := strings.Split(string(b), "\n")
	retFloat64 := make([]float64, len(ret))

	for i := 0; i < len(ret); i++ {
		retFloat64[i], _ = strconv.ParseFloat(ret[i], 64)
	}

	return retFloat64, nil
}

func WriteCloseData(filename string, data []float64) {
	s := ""
	for _, val := range data {
		s += fmt.Sprintf("%f\n", val)
	}

	fmt.Println(s)
	ioutil.WriteFile(filename, []byte(s), 0664)
}
