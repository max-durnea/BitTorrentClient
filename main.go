package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/max-durnea/BitTorrentClient/bencode"
	"os"
)

type Info struct {
	length        int
	name          string
	piece_length  uint
	pieces_hashes [][]byte
}
type TorrentFile struct {
	announce      string
	comment       string
	creation_date int //unix timestamp
	info          Info
	info_hash     [20]byte
	//announce_list --> optional
}

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
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	torrent, err := CreateTorrent(decoded)
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	fmt.Printf("%v\n", torrent)
}

func CreateTorrent(data bencode.BValue) (TorrentFile, error) {
	m, ok := data.(map[string]bencode.BValue)
	if !ok {
		return TorrentFile{}, fmt.Errorf("Not a map")
	}
	announce, _ := m["announce"].(string)
	comment, _ := m["comment"].(string)
	creation_date, _ := m["creation date"].(int)
	info, _ := m["info"].(map[string]bencode.BValue)

	length, _ := info["length"].(int)
	name, _ := info["name"].(string)
	piece_length, _ := info["piece length"].(int)
	pieces := []byte(info["pieces"].(string))
	var hashesList [][]byte
	for i := 0; i < len(pieces); i += 20 {
		piece := pieces[i : i+20]
		hashesList = append(hashesList, piece)
	}
	info_struct := Info{
		length:        length,
		name:          name,
		piece_length:  uint(piece_length),
		pieces_hashes: hashesList,
	}
	torrentFile := TorrentFile{
		announce:      announce,
		comment:       comment,
		creation_date: creation_date,
		info:          info_struct,
	}
	return torrentFile, nil
}
