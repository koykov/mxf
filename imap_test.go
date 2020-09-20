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

func (m *TestImap) GetIid() int {
	return m.id
}

func (m *TestImap) randomize() *TestImap {
	m.id = rand.Intn(100)
	m.payload = rand.Uint64()
	return m
}

func TestNewImap(t *testing.T) {
	m := NewImap(10)
	for i := 0; i < 10; i++ {
		m.Set(NewTestImap(i))
	}
}

func BenchmarkImap_Set(b *testing.B) {
	m := NewImap(10000)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.Set(NewTestImap(i + 1))
	}
}

func BenchmarkImap_Get(b *testing.B) {
	m := NewImap(10000)
	for i := 0; i < 1e7; i++ {
		m.Set(NewTestImap(i + 1))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Get(int(rand.Int31n(1e7)))
	}
}

func TestImap_Transaction(t *testing.T) {
	m := NewImap(10)

	for i := 0; i < 10; i++ {
		m.Set(NewTestImap(i))
	}

	m.Begin()
	for i := 5; i < 15; i++ {
		m.Set(NewTestImap(i))
	}
	m.Commit()
}

func BenchmarkImap_Transaction(b *testing.B) {
	m := NewImap(1000)

	// Fill up only even positions.
	for i := 0; i < 1000; i = i + 10 {
		m.Set(NewTestImap(i)).Set(NewTestImap(i + 2)).Set(NewTestImap(i + 4)).Set(NewTestImap(i + 6)).Set(NewTestImap(i + 8))
	}

	// Prepare chunk of test structs.
	c := make([]*TestImap, 0, 100)
	for i := 0; i < 100; i++ {
		c = append(c, NewTestImap(i))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m.Begin()

		for j := 0; j < 100; j++ {
			m.Set(c[j].randomize())
		}

		m.Commit()
	}

}
