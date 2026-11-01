package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Blob struct {
	//size byte
	content string
	hash    [20]byte
}

func newBlob(data string) *Blob {
	sum := sha1.Sum([]byte(data))
	blob := &Blob{content: data, hash: sum}
	return blob
}

func find_index(input string, sub string) int {
	//fmt.Println("the string is:", input, " the sub=", sub)
	n := strings.Index(input, sub)
	if n == -1 {
		return -1
	}

	r := []rune(input[:n]) // start to end
	return len(r)
}

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	objects := make(map[string]string)
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
		fmt.Println("cmd=", cmd, "flag=", flag, "param=", param)

		if cmd == "hash-object" && flag == "-w" {

			b, err := os.ReadFile(param)
			if err != nil {
				panic(err)
			}
			content := string(b)
			bytes := len(content)

			blob := newBlob(content)

			hx := hex.EncodeToString(blob.hash[:])
			hx = strings.TrimSpace(hx)
			objects[hx] = content

			// store data to blob object
			blob_folder := ".git/objects"
			sub_folder := hx[:2]
			file_name := hx[2:]
			zlib_format := fmt.Sprintf("blob %d0x00%s", bytes, content)

			path := filepath.Join(blob_folder, sub_folder, file_name)
			dir := filepath.Dir(path)

			if err = os.MkdirAll(dir, 0755); err != nil {
				panic(err)
			}

			file, err := os.Create(path)

			if err != nil {
				panic(err)
			}
			defer file.Close()

			// write content to file
			_, err = file.Write([]byte(zlib_format))
			if err != nil {
				panic(err)
			}
			fmt.Println(hx)
		}

	case "cat-file":
		flag := os.Args[2]
		hash := strings.TrimSpace(os.Args[3])
		//fmt.Println("hash=", hash)

		if flag == "-p" {
			// fetch from folder
			blob_folder := ".git/objects"
			sub_folder := hash[:2]
			file_name := hash[2:]
			file_slice := []string{blob_folder, sub_folder, file_name}
			full_path := strings.Join(file_slice, "/")

			if _, err := os.Stat(full_path); err != nil {
				panic(err)
			} else {
				if data, err := os.ReadFile(full_path); err != nil {
					panic(err)
				} else {

					//zlib decompress
					b := bytes.NewReader(data)

					r, err := zlib.NewReader(b)
					if err != nil {
						panic(err)
					}

					buf := new(bytes.Buffer)
					buf.ReadFrom(r)
					p := buf.String()
					defer r.Close() // close zlib readerr

					//p, err := io.ReadAll(r)
					if err != nil {
						panic(err)
					} else {
						res := string(p)
						//res_digit_str := fmt.Sprintf("%d", len(res))
						null_byte_len := len("\x00")
						sub_idx := find_index(res, "\x00")
						//fmt.Println("the sub index is ", sub_idx, "len of res is: ", res_len)
						fmt.Print(res[sub_idx+null_byte_len:])
					}

				}

			}
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
