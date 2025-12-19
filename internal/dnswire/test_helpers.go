package dnswire

import "fmt"

func printByteArrayAsHex(data []byte) {

	for i, b := range data {
		fmt.Printf("%02x ", b)
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

}

func printByteArrayAsASCII(data []byte) {

	for _, b := range data {
		if b >= 32 && b <= 126 {
			fmt.Printf("%c", b)
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()

}
