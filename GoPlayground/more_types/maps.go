package more_types

import "fmt"

type NewVertex struct {
	Lat, Long float64
}

var m map[string]NewVertex

var newMap = map[string]NewVertex{
	"Bell Labs": {
		40.68433, -74.39967,
	},
	"Google": {
		37.42202, -122.08408,
	},
}

func MapsExample() {
	m = make(map[string]NewVertex)
	m["Bell Labs"] = NewVertex{
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])
	fmt.Println(newMap)
}

func MutatingMapsExample() {
	m := make(map[string]int)

	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])

	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])

	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])

	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)
}
