package main

import (
	"testing"
)

func TestConversion(t *testing.T) {
	{
		index := 2
		p := Convert1d2(index, 13)
		expectedPair := Pair{2, 0}
		if p != expectedPair {
			t.Fatalf("Convert1d2 with %d returns %+v, not equal to %+v", index, p, expectedPair)
		}
	}
	{
		index := 23
		p := Convert1d2(index, 13)
		expectedPair := Pair{10, 1}
		if p != expectedPair {
			t.Fatalf("Convert1d2 with %d returns %+v, not equal to %+v", index, p, expectedPair)
		}
	}
}
