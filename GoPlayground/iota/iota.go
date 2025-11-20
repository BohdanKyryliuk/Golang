package iota

import "fmt"

type ByteSize float64

const (
	_           = iota             // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota) // 1 shifted left by 10 bits
	MB                             // 1 shifted left by 20 bits
	GB                             // 1 shifted left by 30 bits
	TB                             // 1 shifted left by 40 bits
	PB                             // 1 shifted left by 50 bits
	EB                             // 1 shifted left by 60 bits
)

type Weekday int

const (
	Sunday Weekday = iota + 1
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func PrintIotaBytes() {
	fmt.Printf("KB: %v bytes\n", KB)
	fmt.Printf("MB: %v bytes\n", MB)
	fmt.Printf("GB: %v bytes\n", GB)
	fmt.Printf("TB: %v bytes\n", TB)
	fmt.Printf("PB: %v bytes\n", PB)
	fmt.Printf("EB: %v bytes\n", EB)
}

func PrintWeekdays() {
	fmt.Printf("Sunday: %d\n", Sunday)
	fmt.Printf("Monday: %d\n", Monday)
	fmt.Printf("Tuesday: %d\n", Tuesday)
	fmt.Printf("Wednesday: %d\n", Wednesday)
	fmt.Printf("Thursday: %d\n", Thursday)
	fmt.Printf("Friday: %d\n", Friday)
	fmt.Printf("Saturday: %d\n", Saturday)
}
