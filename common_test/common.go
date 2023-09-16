package common_test

import (
	"math/rand"
	"time"
)

var Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomChar() byte {
	return randomCharPool[Rnd.Intn(len(randomCharPool))]
}

func RandomString() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = RandomChar()
	}
	return string(b)
}

var randomCharPool = func() (a []byte) {
	for c := 'a'; c <= 'z'; c++ {
		a = append(a, byte(c))
	}
	for c := 'A'; c <= 'Z'; c++ {
		a = append(a, byte(c))
	}
	return
}()
