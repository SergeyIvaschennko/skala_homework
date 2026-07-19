package main

func Map(data []int, action func(int) int) []int {
	if data == nil {
		return nil
	}

	newData := make([]int, len(data))
	for i := range data {
		newData[i] = action(data[i])
	}

	return newData
}

func Filter(data []int, action func(int) bool) []int {
	if data == nil {
		return nil
	}

	newData := make([]int, 0, len(data))
	for i := range data {
		if action(data[i]) {
			newData = append(newData, data[i])
		}
	}

	return newData
}

func Reduce(data []int, initial int, action func(int, int) int) int {
	for i := range data {
		initial = action(initial, data[i])
	}

	return initial
}
