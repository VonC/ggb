package main

import "testing"

type testMain struct {
	name string
}

func TestMain(t *testing.T) {
	// No argument means help
	tests := []*testMain{
		&testMain{name: "No argument means help"},
	}
	for range tests {
		main()
	}
}
