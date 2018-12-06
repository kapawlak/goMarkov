package main
import(
)

//For enforcing Periodic Boundary Conditions
func CyclicBoundary(x int, L int) int {
	return ((x % L) + L) % L
}
