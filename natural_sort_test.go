package main

import (
	"testing"
)

func Test_isDigit(t *testing.T) {
	if isDigit('1') == true && isDigit('0') == true && isDigit('a') == false && isDigit('.') == false {
		t.Log("isdigit test passed")
	} else {
		t.Error("isdigit test failed")
	}
}

func Test_getNumberAndLength(t *testing.T) {
	checkOutput := func(p1 string, p2 int, r1 uint64, r2 int) bool {
		a, b := getNumberAndLength(p1, p2)
		return a == r1 && b == r2
	}
	if checkOutput("123", 0, 123, 3) && checkOutput("5a", 0, 5, 1) && checkOutput("abcd", 0, 0, 0) && checkOutput("abcd509", 4, 509, 3) {
		t.Log("getNumberAndLength test passed")
	} else {
		t.Error("getNumberAndLength test failed")
	}
}

func Test_compareNatural(t *testing.T) {
	if compareNatural("a50b10", "a50b2") == false && compareNatural("abc", "abc") == false && compareNatural("10xcv", "11xcv") == true {
		t.Log("compareNatural test passed")
	} else {
		t.Error("compareNatural test failed")
	}
}

func Benchmark_compareNatural(b *testing.B) {
	for i := 0; i < b.N; i++ {
		compareNatural("sdfga50b10", "sdfga50b2")
		compareNatural("abcdefghijklmnopqrtuvwxyz", "abcdefghijklmnopqrtuvwxyz")
		compareNatural("5000000000000000x10xcv", "5000000000000000x11xcv")
	}
}