package state

import (
	"math"
)

func get_solution(m [][]int, i int, j int, weight func(int) int, result *[]int) []int {
	if i == 0 {
		return *result
	}
	if m[i][j] > m[i-1][j] {
		*result = append(*result, i)
		get_solution(m, i-1, j-weight(i), weight, result)
	} else {
		get_solution(m, i-1, j, weight, result)
	}

	return *result
}

//returns the index of the items of the optimal knapsack
func Knapsack(item_count int, capacity int, value func(int) int, weight func(int) int) []int {
	var m = make([][]int, item_count+1)
	for i := 0; i < item_count; i++ {
		m[i] = make([]int, capacity+1)
	}

	for j := 0; j < capacity; j++ {
		m[0][j] = 0
	}
	for i := 0; i < item_count; i++ {
		m[i][0] = 0
	}
	for i := 1; i < item_count; i++ {
		for j := 1; j < capacity; j++ {
			if weight(i) > j {
				m[i][j] = m[i-1][j]
			} else {
				m[i][j] = int(math.Max(float64(m[i-1][j]), float64(m[i-1][j-weight(i)]+value(i))))
			}
		}
	}

	var ret []int
	get_solution(m, item_count-1, capacity-1, weight, &ret)
	return ret
}
