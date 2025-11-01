package main

import (
	"bytes"
	"crypto/rand"
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
	//fmt.Printf("%v\n", decoded)
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	data_byte := decoder.InfoBytes
	//fmt.Printf("DataByte:%v\n", data_byte)
	hash := sha1.Sum(data_byte)
	torrent, err := CreateTorrent(decoded)
	torrent.info_hash = hash
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	buf := make([]byte, 20)
	rand.Read(buf)
	port := 6882
	url, err := torrent.CreateURL(buf, port)
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return

	}
	resp, err := sendAnnounce(url)
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	fmt.Println(resp)
	//fmt.Printf("Announce: %s\n", torrent.announce)
	//fmt.Printf("Comment: %s\n", torrent.comment)
	//fmt.Printf("Creation date: %d\n", torrent.creation_date)
	//fmt.Printf("Info:\n  Name: %s\n  Length: %d\n  Piece length: %d\n",
	//	torrent.info.name, torrent.info.length, torrent.info.piece_length)

	//fmt.Println("  Pieces hashes:")
	//for i, h := range torrent.info.pieces_hashes {
	//	fmt.Printf("    piece %d: %x\n", i, h)
	//}

	//fmt.Printf("Info hash: %x\n", torrent.info_hash)
}

/*
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
*/
