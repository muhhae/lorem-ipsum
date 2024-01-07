package util

import "fmt"

func Format(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	} else if n < 1000000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	} else if n < 1000000000 {
		return fmt.Sprintf("%.1fm", float64(n)/1000000)
	} else {
		return fmt.Sprintf("%.1fb", float64(n)/1000000000)
	}
}
