package main

import (
	"fmt"
	"github.com/max-durnea/BitTorrentClient/bencode"
)

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
