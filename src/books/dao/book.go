package dao

import (
	"books/dto"

	"github.com/boltdb/bolt"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

const (
	bucketBooksName = "books"
)

func InitBucketBooks(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		list, err := tx.CreateBucketIfNotExists([]byte(bucketBooksName))
		if err != nil {
			return err
		}

		book := dto.Book{
			Name:   "Le Rouge et le Noir",
			Author: "Stendhal",
			Price:  5,
		}

		buf, _ := book.MarshalBinary()
		err = list.Put([]byte(book.Name), buf)
		if err != nil {
			return err
		}

		book = dto.Book{
			Name:   "En attendant Godot",
			Author: "Samuel Beckett",
			Price:  0,
		}

		buf, _ = book.MarshalBinary()
		err = list.Put([]byte(book.Name), buf)
		if err != nil {
			return err
		}

		return nil
	})
}

func ListBooks(db *bolt.DB) (list []dto.Book) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketBooksName))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			b := dto.Book{}
			err := b.UnmarshalBinary(v)
			if err != nil {
				continue
			}

			list = append(list, b)
		}

		return nil
	})

	return
}
