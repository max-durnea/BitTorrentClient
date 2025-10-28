package bencode

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type BValue interface{}

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) Decode() (BValue, error) {
	b, err := d.r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {
	case b == 'i':
		return d.decodeInt()
	case b == 'l':
		return d.decodeList()
	case b >= '0' && b <= '9':
		return d.decodeString(b)
	case b == 'd':
		return d.decodeDict()
	default:
		return nil, fmt.Errorf("Invalid bencode prefix %v\n", b)

	}
}
func (d *Decoder) decodeInt() (int, error) {
	numStr, err := d.readUntil('e')
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (d *Decoder) decodeString(b byte) (string, error) {
	lenStr, err := d.readUntil(':')
	if err != nil {
		return "", err
	}
	lenStr = string(b) + lenStr
	l, err := strconv.Atoi(lenStr)
	if err != nil {
		return "", err
	}
	buf := make([]byte, l)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
func (d *Decoder) decodeDict() (map[string]BValue, error) {
	m := make(map[string]BValue)
	for {
		b, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 'e' {
			break
		}
		key, err := d.decodeString(b)
		if err != nil {
			return nil, err
		}
		value, err := d.Decode()
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}

func (d *Decoder) decodeList() ([]BValue, error) {
	var list []BValue
	for {
		b, err := d.r.Peek(1)
		if err != nil {
			return nil, err
		}
		if b[0] == 'e' {
			d.r.ReadByte()
			break
		}
		v, err := d.Decode()
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}
func (d *Decoder) readUntil(delim byte) (string, error) {
	var str string
	for {
		b, err := d.r.ReadByte()
		if err != nil {
			return "", err
		}
		if b == delim {
			return str, nil
		}
		str += string(b)
	}
}
