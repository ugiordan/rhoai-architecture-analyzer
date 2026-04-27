package main

// Complexity 1: no decision points
func simpleFunc() int {
	return 42
}

// Complexity 4: if + else if + for + &&
func complexFunc(x int, items []string) int {
	if x > 0 && x < 100 {
		return 1
	} else if x >= 100 {
		return 2
	}
	for _, item := range items {
		_ = item
	}
	return 0
}

// Complexity 5: switch with 4 cases
func switchFunc(op string) int {
	switch op {
	case "add":
		return 1
	case "sub":
		return 2
	case "mul":
		return 3
	default:
		return 0
	}
}

// Complexity 3: nested if + ||
func nestedFunc(a, b bool) string {
	if a {
		if b || a {
			return "both"
		}
		return "a"
	}
	return "none"
}
