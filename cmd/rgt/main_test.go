package main

import (
	"testing"
	"time"
)

func TestSlow(t *testing.T) {
	// test to verify the bheaviour of the progress bar
	time.Sleep(5 * time.Second)
}

func TestFast(t *testing.T) {
	// nothing to do here it has to be quick :)
}
