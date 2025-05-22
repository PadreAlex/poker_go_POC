package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type LookupTable struct {
	// default amount of integers in .dat file
	TABLE [32487834]int32
}

func (l *LookupTable) getTable() {
	file, err := os.Open("./HandRanks.dat")
	check(err)
	defer file.Close()

	stat, err := file.Stat()
	check(err)

	expectedSize := int64(32487834 * 4)
	if stat.Size() != expectedSize {
		panic(fmt.Sprintf("HandRanks.dat is the wrong size: expected %d bytes", expectedSize))
	}

	err = binary.Read(file, binary.LittleEndian, &l.TABLE)
	check(err)
}
