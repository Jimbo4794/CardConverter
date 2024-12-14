package main

import (
	"fmt"
	"os"
	"path"

	"math/rand"
)

// run the binary with the mana cost "one, two, etc"
func main() {
	mana := os.Args[1]
	files, _ := os.ReadDir(path.Join("output", mana))
	file := files[rand.Intn(len(files))]
	fmt.Println(file.Name())
}
