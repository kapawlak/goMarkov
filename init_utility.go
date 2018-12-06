package main

import (
	"fmt"
	"math/rand"
)

//Grid Manipulation
func AddGrids(g ...[]int) []int {
	gnew := make([]int, len(g[0]))
	for _, gg := range g {
		if len(gg) != len(gnew) {
			fmt.Println("Warning: Grids of inevquivalent sizes")
			break
		}
	}

	for site := range gnew {
		for _, gg := range g {
			gnew[site] += gg[site]
		}
	}
	return gnew
}

func RandomPop(minpopulation, maxpopulation, length, width int) []int {
	grid := make([]int, length*width)
	for i := range grid {
		grid[i] = rand.Intn(maxpopulation-minpopulation+1) + minpopulation
	}
	return grid
}

func Circular(xpos, ypos, radius float64, pop, length, width int) []int {
	rsq := int(radius * float64(length) * radius * float64(length))
	px := int(float64(length) * xpos)
	py := int(float64(width) * ypos)
	grid := make([]int, length*width)
	for i := 0; i < length; i++ {
		for j := 0; j < width; j++ {
			if (i-px)*(i-px)+(j-py)*(j-py) < rsq {
				grid[i*width+j] = pop
			}
		}
	}
	return grid
}

func GeneratePopulation(f interface{}, length, width int) []int {
	grid := make([]int, length*width)
	switch ff := f.(type) {
	case func(xsite int, ysite int) int:
		for i := 0; i < length; i++ {
			for j := 0; j < width; j++ {
				grid[i*width+j] = ff(i, j)
			}
		}
	case func(xsite int, ysite int) float64:
		for i := 0; i < length; i++ {
			for j := 0; j < width; j++ {
				grid[i*width+j] = int(ff(i, j))
			}
		}
	case func(xsite, ysite float64) float64:
		for i := 0; i < length; i++ {
			for j := 0; j < width; j++ {
				grid[i*width+j] = int(ff(float64(i), float64(j)))
			}
		}

	case int:
		for i := 0; i < length; i++ {
			for j := 0; j < width; j++ {
				grid[i*width+j] = ff
			}
		}

	}
	return grid
}



////////////////////////////////////// Init System ////////////////////////////////
func MakeSystem(Length, Width int, Resolution float64) *Sys {
	sys := new(Sys)
	//if res=0, auto generate
	sys.glen = Length
	sys.gwid = Width
	sys.res = Resolution

	return sys
}