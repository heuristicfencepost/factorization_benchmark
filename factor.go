// $G $F.go && $L $F.$A  # don't run it - goes forever

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"
import "flag"
import vector "container/vector"

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(max int, ch chan<- int) {
	//fmt.Printf("Generating primes less than or equal to %d \n",max)

	// Slight optimization; after 2 we know there are no even primes so we only
	// need to consider odd values
	ch <- 2
	for i := 3; i<=max ; i += 2 {
		ch <- i // Send 'i' to channel 'ch'.
	}
	//fmt.Printf("Sending -1");
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
	//fmt.Printf("Terminating prime %d \n",prime)
	out <- -1
}

func main() {

	var target *int = flag.Int("target",1,"Number we wish to factor")
	flag.Parse()

	t := *target
	fmt.Printf("Target: %d\n",t)

	var rv vector.IntVector

	// Retrieve a prime value and see if we can divide the target evenly by
	// that prime.  If so perform the multiplication and update the current
	// value.
	ch := make(chan int) // Create a new channel.
	go Generate(t,ch)      // Start Generate() as a subprocess.
	for prime := <-ch; prime != -1; prime = <-ch {

		for ;t % prime == 0; {
			t = t / prime
			rv.Push(prime)
		}

		// Create a goroutine for each prime number whether we use it or
		// not.  This performs the daisy chaining setup that was being
		// done by the Sieve() function in sieve.go.
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}

	fmt.Printf("Results: %s\n",fmt.Sprint(rv))
}
