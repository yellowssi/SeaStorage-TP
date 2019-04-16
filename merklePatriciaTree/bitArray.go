package merklePatriciaTree

import (
	"bytes"
	"errors"
	"log"
	"math"
)

type BitArray struct {
	Value []bool
}

func NewBitArray(length int) *BitArray {
	value := make([]bool, length)
	return &BitArray{Value: value[:]}
}

func (b *BitArray) GetLength() int {
	return len(b.Value)
}

func (b *BitArray) GetBit(index int) bool {
	if index >= len(b.Value) {
		log.Fatal("index out of range")
	}
	return b.Value[index]
}

func (b *BitArray) SetBit(index int) error {
	if index >= len(b.Value) {
		return errors.New("index out of range")
	}
	b.Value[index] = true
	return nil
}

func FromByteArray(bytes []byte) *BitArray {
	result := NewBitArray(len(bytes) * 8)
	for i := 0; i < len(bytes); i++ {
		value := bytes[i]
		for j := 0; j < 8; j++ {
			if value%2 != 0 {
				_ = result.SetBit(i*8 + (8 - j - 1))
			}
			value /= 2
		}
	}
	return result
}

func (b *BitArray) ToByteArray() []byte {
	value := b.Value
	length := len(b.Value) / 8
	diff := len(b.Value) % 8
	if diff != 0 {
		temp := NewBitArray(b.GetLength() + (8 - diff))
		for i := 0; i < len(value); i++ {
			if value[i] {
				_ = temp.SetBit(i + diff)
			}
		}
		value = temp.Value
		length++
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		v := 0
		for j := 1; j <= 8; j++ {
			if value[(i*8)+j-1] {
				v += int(math.Pow(2, float64(8-j)))
			}
		}
		result[i] = byte(v)
	}
	return result[:]
}

func CompareBitArray(a BitArray, b BitArray) int {
	if a.GetLength() != b.GetLength() {
		panic("cannot compare in different length")
	}
	return bytes.Compare(a.ToByteArray(), b.ToByteArray())
}

func (b *BitArray) SubBitArray(start int, end int) *BitArray {
	if start < 0 {
		panic("index out of range")
	}
	if end == -1 {
		return &BitArray{Value: b.Value[start:]}
	} else if start > end {
		panic("start index greater than end index")
	} else if end < -1 || end > b.GetLength() {
		panic("index out of range")
	}
	return &BitArray{Value: b.Value[start:end]}
}
