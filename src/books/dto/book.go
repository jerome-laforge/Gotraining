package dto

import "encoding/json"

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

type Book struct {
	Name   string
	Author string
	Price  uint
}

func (b Book) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *Book) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, b)
}
