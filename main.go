package main

import (
	"fmt"
	"github.com/max-durnea/BitTorrentClient/bencode"

	"bytes"
	"os"
)

func main() {
	var file_name string
	fmt.Scan(&file_name)
	fmt.Printf("%v\n", file_name)
	data, err := os.ReadFile(file_name)
	if err != nil {
		fmt.Errorf("%v\n", err)
		return
	}
	reader := bytes.NewReader(data)
	decoder := bencode.NewDecoder(reader)
	decoded, err := decoder.Decode()
	if err != nil {
		fmt.Errorf("%v\n", err)
		return
	}
	fmt.Printf("%v", decoded)
}
