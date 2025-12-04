package main

import (
	"fmt"
	"strings"
)

func main() {
	qname := "www.google.com"

	labels := strings.Split(qname, ".")

	byte_size := 0

	for index, value := range labels {

		byte_size = byte_size + len(labels[index])
		fmt.Printf("Index: %d, Value: %s, Length: %d\n", index, value, len(labels[index]))
	}

	fmt.Printf("QNAME Length: %d", byte_size)

}
