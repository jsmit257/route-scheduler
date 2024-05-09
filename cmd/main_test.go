package main

import (
	"os"
	"testing"
)

func Test_Main(t *testing.T) {
	os.Args[1] = "../data/problem1.txt"
	main()
}
