package main

type CircularQueue struct {
	values      []int
	frontindex  int
	backindex   int
	currentsize int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values:      make([]int, size),
		frontindex:  0,
		backindex:   0,
		currentsize: 0,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}

	if q.currentsize == 0 {
		q.values[q.frontindex] = value
		q.currentsize++
		return true
	}

	q.currentsize++
	q.backindex++

	if q.backindex == len(q.values) {
		q.backindex = 0
	}

	q.values[q.backindex] = value

	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}

	q.currentsize--

	if q.currentsize == 0 {
		q.frontindex = 0
		q.backindex = 0
		return true
	}

	q.frontindex++

	if q.frontindex == len(q.values) {
		q.frontindex = 0
	}

	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}

	return q.values[q.frontindex]
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}

	return q.values[q.backindex]
}

func (q *CircularQueue) Empty() bool {
	return q.currentsize == 0
}

func (q *CircularQueue) Full() bool {
	return len(q.values) == q.currentsize
}
