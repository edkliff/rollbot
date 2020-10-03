package generator

import (
	"math"
	"sync"
	"time"
)

// Why not math.random? Because i can.
type Generator struct {
	m          int64 // modulo
	a          int64 // multiplier
	mut        sync.Mutex
	prevNumber int64
}

func  InitGenerator() *Generator {
	g := Generator{
		m:          1000000,
		a:          211021,
		mut:        sync.Mutex{},
		prevNumber: time.Now().Unix(),
	}
	if g.prevNumber%2 == 0 {
		g.prevNumber = g.prevNumber  + 1
	}
	return &g
}

func (gen *Generator) Next() float64 {
	gen.mut.Lock()
	defer gen.mut.Unlock()

	dec := gen.a * gen.prevNumber
	now := dec % gen.m
	gen.prevNumber = now
	num := float64(now) / float64(gen.m)
	return num
}

func (gen *Generator) Roll(count int64, dice int64) []int {
	results := make([]int, count)
	for i := 0; i < int(count); i++ {
		a := gen.Random(0, dice)
		results[i] = int(a)
	}
	return results
}

func (gen *Generator) Random(min int64, max int64) int64 {
	randomRange := max - min
	a := gen.Next() * float64(randomRange)
	result := int64(math.Ceil(a)) + min
	return result
}

func Sum(sl []int) int {
	var sum int
	for _, e := range sl {
		sum += e
	}
	return sum
}
