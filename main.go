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
	m, ok := decoded.(map[string]bencode.BValue)
	if !ok {
		fmt.Printf("Not a map")
		return
	}
	info := m["info"].(map[string]bencode.BValue)
	pieces := []byte(info["pieces"].(string))
	var hashesList [][]byte
	for i := 0; i < len(pieces); i += 20 {
		piece := pieces[i : i+20]
		hashesList = append(hashesList, piece)
	}
	for _, hash := range hashesList {
		fmt.Printf("%v\n", hash)
	}
}
