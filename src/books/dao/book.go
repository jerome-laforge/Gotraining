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
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketBooksName))
		if err != nil {
			return err
		}

		book := dto.Book{
			Name:   "Le Rouge et le Noir",
			Author: "Stendhal",
			Price:  5,
		}

		buf, _ := book.MarshalBinary()
		err = bucket.Put([]byte(book.Name), buf)
		if err != nil {
			return err
		}

		book = dto.Book{
			Name:   "En attendant Godot",
			Author: "Samuel Beckett",
			Price:  0,
		}

		buf, _ = book.MarshalBinary()
		err = bucket.Put([]byte(book.Name), buf)
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
		bucket := tx.Bucket([]byte(bucketBooksName))

		c := bucket.Cursor()
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

func GetBook(name string) (book dto.Book, found bool, err error) {
	db, err := getDB()
	if err != nil {
		return book, found, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		v := bucket.Get([]byte(name))
		if v == nil {
			return err
		}

		found = true

		err = book.UnmarshalBinary(v)
		return err
	})

	return
}

func CreateBook(book dto.Book) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		buf, _ := book.MarshalBinary()
		return bucket.Put([]byte(book.Name), buf)
	})
}

func DeleteBook(name string) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		return bucket.Delete([]byte(name))
	})
}
