package mxf

import (
	"math/rand"
	"testing"
)

type TestImap struct {
	id      int
	payload uint64
}

func NewTestImap(id int) *TestImap {
	m := TestImap{
		id:      id,
		payload: rand.Uint64(),
	}
	return &m
}

func (m *TestImap) GetId() int {
	return m.id
}

func TestNewImap(t *testing.T) {
	m := NewImap(10)
	for i := 0; i < 10; i++ {
		m.Set(i, NewTestImap(i))
	}
}

func BenchmarkImap_Set(b *testing.B) {
	m := NewImap(10000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Set(i, NewTestImap(i+1))
	}
}

func BenchmarkImap_Get(b *testing.B) {
	m := NewImap(10000)
	for i := 0; i < 1e7; i++ {
		m.Set(i, NewTestImap(i+1))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Get(int(rand.Int31n(1e7)))
	}
}
