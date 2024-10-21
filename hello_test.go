package main

import (
	"log"
	"os"
	"testing"
)

func TestConversion(t *testing.T) {
	{
		index := 2
		p := Convert1d2d(index, 13)
		expectedPair := Pair{2, 0}
		if p != expectedPair {
			t.Fatalf("Convert1d2 with %d returns %+v, not equal to %+v", index, p, expectedPair)
		}
	}
	{
		index := 23
		p := Convert1d2d(index, 13)
		expectedPair := Pair{10, 1}
		if p != expectedPair {
			t.Fatalf("Convert1d2 with %d returns %+v, not equal to %+v", index, p, expectedPair)
		}
	}

	{
		world := World{}

		worldSize := Pair{14, 4}
		world.init(worldSize)
		if len(world.Cells) != worldSize.X*worldSize.Y {
			t.Fatalf("WorldSize %+v is not %+v", len(world.Cells), (worldSize.X * worldSize.Y))
		}
		for i, cell := range world.Cells {
			index := world.Convert2d1d(cell.WorldPosition)
			if i != index {
				t.Fatalf("Convert2d1d with %+v returns %+v, not equal to %+v", cell, index, i)
			}
		}
	}

	{
		world := World{}
		world.init(Pair{4, 4})

		for i, cell := range world.Cells {
			index := world.Convert2d1d(cell.WorldPosition)
			if i != index {
				t.Fatalf("Convert2d1d with %+v returns %+v, not equal to %+v", cell, index, i)
			}
		}
	}

	{
		world := World{}
		world.init(Pair{12, 7})

		for i, cell := range world.Cells {
			index := Convert2d1d(cell.WorldPosition, world.MapSize.X)
			if i != index {
				t.Fatalf("Convert2d1d with %+v returns %+v, not equal to %+v", cell, index, i)
			}
		}
	}
}

func TestGetNeigbours(t *testing.T) {

	file, err := os.OpenFile("test_logs.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	world := World{}
	world.init(Pair{4, 4})
	world.Cells[0].IsAlive = true
	neighbours := world.Cells[1].getNeighbours(world)

	expectedPair := Pair{1, 4}
	if neighbours != expectedPair {
		t.Fatalf("getNeighbours with %+v returns %+v, not equal to %+v", world.Cells[0], neighbours, expectedPair)
	}

}
