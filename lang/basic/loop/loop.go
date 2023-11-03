package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func convertToBin(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	printFileContents(file)
}

func printFileContents(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func forever() {
	for {
		fmt.Println("abc")
	}
}

func main() {
	fmt.Println("convertToBin results:")
	fmt.Println(
		convertToBin(5),  // 101
		convertToBin(13), // 1101
		convertToBin(72387885),
		convertToBin(0),
	)

	fmt.Println("abc.txt contents:")
	printFile("lang/basic/branch/abc.txt")

	fmt.Println("printing a string:")
	s := `abc"d"
	kkkk
	123

	p`
	printFileContents(strings.NewReader(s))

	// Uncomment to see it runs forever
	// forever()
}
