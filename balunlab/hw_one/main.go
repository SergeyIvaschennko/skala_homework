package main

import "fmt"

const (
	byteMask  uint32 = 0xFF //0xFF = 11111111
	byteSize         = 8
	byteCount        = 4
)

func ToLittleEndianCycle(number uint32) uint32 {
	var result uint32

	for i := 0; i < byteCount; i++ {
		bytePart := (number >> (byteSize * i)) & byteMask
		fmt.Printf("bytePart: 0x%08X\n", bytePart)
		result |= bytePart << (byteSize * (byteCount - 1 - i))
		fmt.Printf("result: 0x%08X\n", result)
		fmt.Printf("\n")
	}

	return result
}

func ToLittleEndian(number uint32) uint32 {
	return ((number & 0x000000FF) << 24) |
		((number & 0x0000FF00) << 8) |
		((number & 0x00FF0000) >> 8) |
		((number & 0xFF000000) >> 24)
}

func main() {
	number := uint32(0x01020304)

	fmt.Printf("Original : 0x%08X\n\n", number)
	fmt.Printf("Cycle     : 0x%08X\n", ToLittleEndianCycle(number))
	fmt.Printf("Bitwise   : 0x%08X\n", ToLittleEndian(number))
}
