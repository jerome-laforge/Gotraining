package dto

import (
	"bytes"

	"github.com/ugorji/go/codec"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

var (
	mh codec.MsgpackHandle
)

type Book struct {
	Name   string
	Author string
	Price  uint
}

func (b Book) Marshal() ([]byte, error) {
	var data bytes.Buffer
	err := codec.NewEncoder(&data, &mh).Encode(b)
	return data.Bytes(), err
}

func (b *Book) Unmarshal(data []byte) error {
	return codec.NewDecoderBytes(data, &mh).Decode(&b)
}
