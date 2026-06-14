package main

import (
	"fmt"
	"unsafe"
)

type COWBuffer struct {
	data        []byte
	countOfRefs *int
}

func NewCOWBuffer(data []byte) COWBuffer {
	count := 1

	return COWBuffer{
		data:        data,
		countOfRefs: &count,
	}
}

func (b *COWBuffer) Clone() COWBuffer {
	*b.countOfRefs++

	return COWBuffer{
		data:        b.data,
		countOfRefs: b.countOfRefs,
	}
}

func (b *COWBuffer) Close() {
	if b.countOfRefs == nil {
		return
	}

	*b.countOfRefs--
	b.countOfRefs = nil
	b.data = nil
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if index < 0 || index >= len(b.data) {
		return false
	}

	if *b.countOfRefs > 1 {
		*b.countOfRefs--

		newCount := 1

		newData := make([]byte, len(b.data))
		copy(newData, b.data)

		b.data = newData
		b.countOfRefs = &newCount
	}

	b.data[index] = value
	return true
}

func (b *COWBuffer) String() string {
	if len(b.data) == 0 {
		return ""
	}

	return unsafe.String(&b.data[0], len(b.data))
}

func main() {
	buffer := NewCOWBuffer([]byte("rockstar"))

	copyBuffer := buffer.Clone()

	buffer.Update(0, 'g')

	fmt.Println(buffer.String())
	fmt.Println(copyBuffer.String())

	buffer.Close()
	copyBuffer.Close()
}
