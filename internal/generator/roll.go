package generator

import (
	"errors"
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

func InitGenerator() *Generator {
	g := Generator{
		m:          1000000,
		a:          211021,
		mut:        sync.Mutex{},
		prevNumber: time.Now().Unix(),
	}
	if g.prevNumber%2 == 0 {
		g.prevNumber = g.prevNumber + 1
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

func (gen *Generator) Roll(count int64, dice int64) ([]int64, error) {
	if count > 1000 || dice > 1000 || count < 1 || dice < 1 {
		return nil, errors.New("Это мы считать не будем потому, что Малой - мудак.")
	}
	results := make([]int64, count)
	for i := 0; i < int(count); i++ {
		a := gen.Random(0, dice)
		results[i] = a
	}
	return results, nil
}

func (gen *Generator) Random(min int64, max int64) int64 {
	randomRange := max - min
	a := gen.Next() * float64(randomRange)
	result := int64(math.Ceil(a)) + min
	return result
}

func Sum(sl []int64) int64 {
	var sum int64
	for _, e := range sl {
		sum += e
	}
	return sum
}
