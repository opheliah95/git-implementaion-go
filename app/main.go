package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type Blob struct {
	//size byte
	content string
	hash    [20]byte
}

func newBlob(data string) *Blob {
	sum := sha1.Sum([]byte(data))
	blob := Blob{content: data, hash: sum}
	return &blob
}

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "mygit":
		cmd := os.Args[2]
		flag := os.Args[3]
		param := os.Args[4]
		//fmt.Println("cmd=", cmd, "flag=", flag, "param=", param)

		if cmd == "hash-object" && flag == "-w" {

			b, err := os.ReadFile(param)
			if err != nil {
				panic(err)
			}
			content := string(b)

			blob := newBlob(content)

			hx := hex.EncodeToString(blob.hash[:])
			fmt.Println(hx)
		}

	case "init":
		//TODO: Uncomment the code below to pass the first stage!

		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
