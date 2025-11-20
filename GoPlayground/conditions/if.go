package conditions

import (
	"fmt"
	"math"
)

func Sqrt(x float64) string {
	if x < 0 {
		return Sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func IfExample(x int) string {
	if x%2 == 0 {
		return "Even"
	} else {
		return "Odd"
	}
}

func IfWithShortStatement(x int) string {
	if y := x * 2; y > 10 {
		return "Greater than 10"
	} else {
		return "10 or less"
	}
}

func NestedIfExample(x int) string {
	if x > 0 {
		if x%2 == 0 {
			return "Positive Even"
		} else {
			return "Positive Odd"
		}
	} else if x < 0 {
		return "Negative"
	} else {
		return "Zero"
	}
}

func IfWithoutElse(x int) string {
	if x > 0 {
		return "Positive"
	}
	return "Non-positive"
}

func IfMultipleConditions(x int) string {
	if x < 0 {
		return "Negative"
	} else if x == 0 {
		return "Zero"
	} else if x > 0 && x < 10 {
		return "Positive Single Digit"
	} else {
		return "Positive Multi Digit"
	}
}

func IfWithLogicalOperators(x int) string {
	if x > 0 && x < 100 {
		return "Between 1 and 99"
	} else if x <= 0 || x >= 100 {
		return "100 or less than or equal to 0"
	}
	return "Undefined"
}

func IfWithReturnEarly(x int) string {
	if x < 0 {
		return "Negative"
	}
	if x == 0 {
		return "Zero"
	}
	return "Positive"
}

func IfWithMultipleReturns(x int) string {
	if x%2 == 0 {
		return "Even"
	}
	return "Odd"
}

func IfWithComplexCondition(x, y int) string {
	if (x > 0 && y > 0) || (x < 0 && y < 0) {
		return "Same Sign"
	}
	return "Different Signs"
}

func IfWithFunctionCall(x int) string {
	if isEven(x) {
		return "Even"
	}
	return "Odd"
}

func isEven(x int) bool {
	return x%2 == 0
}

func IfWithTypeAssertion(i interface{}) string {
	if str, ok := i.(string); ok {
		return "String: " + str
	}
	return "Not a string"
}

func IfWithNilCheck(ptr *int) string {
	if ptr == nil {
		return "Pointer is nil"
	}
	return "Pointer is not nil"
}

func IfWithSliceCheck(slice []int) string {
	if len(slice) == 0 {
		return "Slice is empty"
	}
	return "Slice has elements"
}

func IfWithMapCheck(m map[string]int, key string) string {
	if val, ok := m[key]; ok {
		return "Key found with value: " + fmt.Sprint(val)
	}
	return "Key not found"
}

func IfWithChannelCheck(ch chan int) string {
	select {
	case val := <-ch:
		return "Received value: " + fmt.Sprint(val)
	default:
		return "No value received"
	}
}

func IfWithStructFieldCheck(s struct{ Name string }) string {
	if s.Name != "" {
		return "Name is set to: " + s.Name
	}
	return "Name is not set"
}

func IfWithPointerDereference(ptr *int) string {
	if ptr != nil && *ptr > 0 {
		return "Pointer points to a positive integer"
	}
	return "Pointer is nil or points to a non-positive integer"
}

func IfWithMultipleVariables(x, y int) string {
	if a, b := x+y, x-y; a > b {
		return "Sum is greater than difference"
	}
	return "Difference is greater than or equal to sum"
}

func IfWithBooleanVariable(flag bool) string {
	if flag {
		return "Flag is true"
	}
	return "Flag is false"
}

func IfWithStringComparison(s string) string {
	if s == "hello" {
		return "Greeting detected"
	}
	return "No greeting"
}

func IfWithFloatComparison(f float64) string {
	if f > 0.0 {
		return "Positive float"
	} else if f < 0.0 {
		return "Negative float"
	}
	return "Zero"
}

func IfWithArrayCheck(arr [3]int) string {
	if arr[0] == 0 && arr[1] == 0 && arr[2] == 0 {
		return "Array is all zeros"
	}
	return "Array has non-zero elements"
}

func IfWithStructComparison(s1, s2 struct{ ID int }) string {
	if s1.ID == s2.ID {
		return "Structs have the same ID"
	}
	return "Structs have different IDs"
}

func IfWithInterfaceCheck(i interface{}) string {
	if _, ok := i.(int); ok {
		return "It's an integer"
	}
	return "It's not an integer"
}

func IfWithErrorCheck(err error) string {
	if err != nil {
		return "An error occurred: " + err.Error()
	}
	return "No error"
}

func IfWithTimeComparison(t1, t2 int64) string {
	if t1 < t2 {
		return "t1 is before t2"
	} else if t1 > t2 {
		return "t1 is after t2"
	}
	return "t1 is equal to t2"
}

func IfWithPointerComparison(p1, p2 *int) string {
	if p1 == p2 {
		return "Pointers are equal"
	}
	return "Pointers are not equal"
}

func IfWithByteComparison(b1, b2 byte) string {
	if b1 == b2 {
		return "Bytes are equal"
	}
	return "Bytes are not equal"
}

func IfWithRuneComparison(r1, r2 rune) string {
	if r1 == r2 {
		return "Runes are equal"
	}
	return "Runes are not equal"
}

func IfWithComplexNumberComparison(c1, c2 complex128) string {
	if c1 == c2 {
		return "Complex numbers are equal"
	}
	return "Complex numbers are not equal"
}

func IfWithPointerArithmetic(ptr *int, offset int) string {
	if ptr != nil && *ptr+offset > 0 {
		return "Pointer arithmetic result is positive"
	}
	return "Pointer is nil or result is non-positive"
}

func IfWithSliceLengthCheck(slice []int) string {
	if len(slice) > 5 {
		return "Slice has more than 5 elements"
	}
	return "Slice has 5 or fewer elements"
}

func IfWithMapLengthCheck(m map[string]int) string {
	if len(m) > 3 {
		return "Map has more than 3 entries"
	}
	return "Map has 3 or fewer entries"
}

func IfWithChannelLengthCheck(ch chan int) string {
	if len(ch) > 0 {
		return "Channel has buffered elements"
	}
	return "Channel is empty"
}

func IfWithStructFieldLengthCheck(s struct{ Name string }) string {
	if len(s.Name) > 5 {
		return "Name is longer than 5 characters"
	}
	return "Name is 5 or fewer characters"
}

func IfWithPointerNilCheck(ptr *int) string {
	if ptr == nil {
		return "Pointer is nil"
	}
	return "Pointer is not nil"
}

func IfWithMultipleConditionsAndReturns(x int) string {
	if x < 0 {
		return "Negative"
	} else if x == 0 {
		return "Zero"
	} else if x > 0 && x < 10 {
		return "Positive Single Digit"
	} else {
		return "Positive Multi Digit"
	}
}

func IfWithDeferExample(x int) string {
	if x < 0 {
		defer fmt.Println("Exiting negative check")
		return "Negative"
	}
	defer fmt.Println("Exiting non-negative check")
	return "Non-negative"
}

func IfWithPanicRecovery(x int) string {
	if x < 0 {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()
		panic("Negative value encountered")
	}
	return "Non-negative"
}

func IfWithLoggingExample(x int) string {
	if x < 0 {
		fmt.Println("Logging: Negative value")
		return "Negative"
	}
	fmt.Println("Logging: Non-negative value")
	return "Non-negative"
}

func IfWithMetricsExample(x int) string {
	if x < 0 {
		// Simulate metric increment
		fmt.Println("Metric: negative_value_count incremented")
		return "Negative"
	}
	// Simulate metric increment
	fmt.Println("Metric: non_negative_value_count incremented")
	return "Non-negative"
}

func IfWithTracingExample(x int) string {
	if x < 0 {
		// Simulate tracing span
		fmt.Println("Tracing: negative_value_span started")
		return "Negative"
	}
	// Simulate tracing span
	fmt.Println("Tracing: non_negative_value_span started")
	return "Non-negative"
}

func IfWithConfigurationCheck(config map[string]bool, key string) string {
	if enabled, ok := config[key]; ok && enabled {
		return "Feature is enabled"
	}
	return "Feature is disabled or not found"
}

func IfWithEnvironmentVariableCheck(env map[string]string, key string) string {
	if value, ok := env[key]; ok && value != "" {
		return "Environment variable is set to: " + value
	}
	return "Environment variable is not set"
}

func IfWithCommandLineArgumentCheck(args []string, arg string) string {
	for _, a := range args {
		if a == arg {
			return "Argument found: " + arg
		}
	}
	return "Argument not found: " + arg
}

func IfWithFileExistenceCheck(files map[string]bool, filename string) string {
	if exists, ok := files[filename]; ok && exists {
		return "File exists: " + filename
	}
	return "File does not exist: " + filename
}

func IfWithNetworkConnectivityCheck(connected bool) string {
	if connected {
		return "Network is connected"
	}
	return "Network is disconnected"
}

func IfWithDatabaseConnectionCheck(connected bool) string {
	if connected {
		return "Database is connected"
	}
	return "Database is disconnected"
}

func IfWithServiceAvailabilityCheck(available bool) string {
	if available {
		return "Service is available"
	}
	return "Service is unavailable"
}

func IfWithUserAuthenticationCheck(authenticated bool) string {
	if authenticated {
		return "User is authenticated"
	}
	return "User is not authenticated"
}

func IfWithUserAuthorizationCheck(authorized bool) string {
	if authorized {
		return "User is authorized"
	}
	return "User is not authorized"
}
