// $G $F.go && $L $F.$A  # don't run it - goes forever

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "flag"

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(max int, ch chan<- int) {

	// Slight optimization; after 2 we know there are no even primes so we only
	// need to consider odd values
	ch <- 2
	for i := 3; i<=max ; i += 2 {
		ch <- i // Send 'i' to channel 'ch'.
	}
	ch <- -1 // Use -1 as an indicator that we're done now
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for i := <- in; i != -1; i = <- in {

		if i % prime != 0 {
			out <- i // Send 'i' to channel 'out'.
		}
	}
	out <- -1
}

// The prime sieve: Daisy-chain Filter processes together.
func Sieve(max int) {
	ch := make(chan int) // Create a new channel.
	go Generate(max,ch)      // Start Generate() as a subprocess.
	for prime := <-ch; prime != -1; prime = <-ch {
		print(prime, "\n")
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}

func main() {

	var max *int = flag.Int("max",1,"Maximum number we wish to return")
	flag.Parse()
	Sieve(*max)
}
