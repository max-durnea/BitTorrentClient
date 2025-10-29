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
	fmt.Printf("%v\n", decoded)
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	data_byte := decoder.InfoBytes
	hash := sha1.Sum(data_byte)
	torrent, err := CreateTorrent(decoded)
	torrent.info_hash = hash
	if err != nil {
		fmt.Println(fmt.Errorf("%v\n", err))
		return
	}
	fmt.Printf("%v\n", torrent)
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
func CreateTorrent(data bencode.BValue) (TorrentFile, error) {
	m, ok := data.(map[string]bencode.BValue)
	if !ok {
		return TorrentFile{}, fmt.Errorf("top-level is not a map")
	}

	// announce
	announceVal, ok := m["announce"]
	if !ok {
		return TorrentFile{}, fmt.Errorf("missing announce")
	}
	announce, ok := announceVal.(string)
	if !ok {
		return TorrentFile{}, fmt.Errorf("announce wrong type")
	}

	// comment (optional)
	comment := ""
	if commentVal, ok := m["comment"]; ok {
		if c, ok := commentVal.(string); ok {
			comment = c
		}
	}

	// creation date (optional)
	creationDate := 0
	if cdVal, ok := m["creation date"]; ok {
		if cd, ok := cdVal.(int); ok {
			creationDate = cd
		}
	}

	// info dict
	infoVal, ok := m["info"]
	if !ok {
		return TorrentFile{}, fmt.Errorf("missing info")
	}
	infoMap, ok := infoVal.(map[string]bencode.BValue)
	if !ok {
		return TorrentFile{}, fmt.Errorf("info is not a dict")
	}

	// info fields
	length, ok := infoMap["length"].(int)
	if !ok {
		return TorrentFile{}, fmt.Errorf("info.length missing or wrong type")
	}

	name, ok := infoMap["name"].(string)
	if !ok {
		return TorrentFile{}, fmt.Errorf("info.name missing or wrong type")
	}

	pieceLength, ok := infoMap["piece length"].(int)
	if !ok {
		return TorrentFile{}, fmt.Errorf("info.piece length missing or wrong type")
	}

	piecesRaw, ok := infoMap["pieces"].(string)
	if !ok {
		return TorrentFile{}, fmt.Errorf("info.pieces missing or wrong type")
	}
	pieces := []byte(piecesRaw)

	if len(pieces)%20 != 0 {
		return TorrentFile{}, fmt.Errorf("pieces length not multiple of 20")
	}

	var hashesList [][]byte
	for i := 0; i < len(pieces); i += 20 {
		hashesList = append(hashesList, pieces[i:i+20])
	}

	infoStruct := Info{
		length:        length,
		name:          name,
		piece_length:  uint(pieceLength),
		pieces_hashes: hashesList,
	}

	torrentFile := TorrentFile{
		announce:      announce,
		comment:       comment,
		creation_date: creationDate,
		info:          infoStruct,
	}

	return torrentFile, nil
}
