package main

import (
	"bufio"
	"os"
)

func readKey() (byte, error) {
	reader := bufio.NewReader(os.Stdin)
	b, err := reader.ReadByte()

	return b, err
}
