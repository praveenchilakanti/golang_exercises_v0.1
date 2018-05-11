package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

const HASH_BYTES = 16

func main() {
	// Check for console input
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("Computing MD5 hash for std console input...")
		md5HashOut, err := stdInDataRead()
		if err != nil {
			log.Fatal("Unable to compute MD5 hash : ", err)
		}
		fmt.Printf("MD5 hash: %x\n", md5HashOut)
		return
	}

	// check for file input
	if len(os.Args) > 1 {
		fmt.Println("Computing MD5 hash for file input: ", os.Args[1])
		md5HashOut, err := readFileContent(os.Args[1])
		if err != nil {
			log.Fatal("Unable to compute MD5 hash : ", err)
		}
		fmt.Printf("MD5 hash: %x\n", md5HashOut)
		return
	}
	fmt.Println("No input provided, computing hash for MD5 itself")
	fmt.Printf("MD5 Hash: %x\n", md5.Sum([]byte("MD5")))
}

// Read from data from std input concole... cat file.txt
func stdInDataRead() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	stdIndata := make([]byte, 1024)
	hash := md5.New()
	for {
		n, err := reader.Read(stdIndata)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		stdIndata = stdIndata[:n]
		if _, err := io.Copy(hash, bytes.NewReader(stdIndata)); err != nil {
			return nil, err
		}
	}
	return hash.Sum(nil)[:HASH_BYTES], nil
}

// Read data from the input file
func readFileContent(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//Open a new hash interface to write to
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}
	return hash.Sum(nil)[:HASH_BYTES], nil
}
