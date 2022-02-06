package main

import (
	"testing"
	"math/rand"
	"time"
	"fmt"
)

func BenchmarkRule(b *testing.B) {
	rand.Seed(time.Now().Unix())
	reasons := []string{
		"pepinos",
	}
	rules := BuildRules()
	n := rand.Int() % len(reasons)
	fmt.Printf("%s --> %t\n",reasons[n],FindRule(reasons[n], rules))
}
