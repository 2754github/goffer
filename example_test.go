package goffer_test

import (
	"fmt"

	"github.com/2754github/goffer"
)

func Example() {
	// Create an empty buffer of int with size 3.
	b := goffer.New[int](3)
	defer b.Close()

	// Subscribe beforehand.
	go (func() {
		for items := range b.Subscribe() {
			fmt.Printf("Subscribe: %d\n", items)
		}
	})()

	// Publish 10 items. (NOTE: This example ignores the Publish errors.)
	for i := range 10 {
		for b.Publish(i) != nil {
		}
	}

	fmt.Printf("Pull:      %d\n", b.Pull())

	// Output:
	// Subscribe: [0 1 2]
	// Subscribe: [3 4 5]
	// Subscribe: [6 7 8]
	// Pull:      [9]
}
