package dao

import (
	"books/dto"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/inconshreveable/log15"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

const (
	bucketBooksName = "books"
)

var (
	db     *bolt.DB
	err    error
	dbOnce sync.Once
)

func getDB() (*bolt.DB, error) {
	dbOnce.Do(func() {
		// Open the my.db data file in your current directory.
		// It will be created if it doesn't exist.
		db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log15.Error(err.Error())
		}
	})

	return db, err
}

func Close() {
	db, err := getDB()
	if err != nil {
		return
	}

	db.Close()
	dbOnce = sync.Once{}
}

func InitBucketBooks() error {
	db, err := getDB()
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
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

	return err
}

func ListBooks() (list []dto.Book, err error) {
	db, err := getDB()
	if err != nil {
		return list, err
	}

	err = db.View(func(tx *bolt.Tx) error {
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
