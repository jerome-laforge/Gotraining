package dao

import (
	"books/dto"
	"books/log"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"golang.org/x/net/context"
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

const DBName = "my.db"

func getDB(ctx context.Context) (*bolt.DB, error) {
	dbOnce.Do(func() {
		// Open the my.db data file in your current directory.
		// It will be created if it doesn't exist.
		db, err = bolt.Open(DBName, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.GetLoggerFromContext(ctx).Error(err.Error())
			return
		}

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketBooksName))
			if err != nil {
				log.GetLoggerFromContext(ctx).Error("impossible to create bucket", "name", bucketBooksName, "err", err)
				return err
			}
			return nil
		})
	})

	return db, err
}

func Close(ctx context.Context) {
	db, err := getDB(ctx)
	if err != nil {
		return
	}

	db.Close()
	dbOnce = sync.Once{}
}

func InitBucketBooks(ctx context.Context) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		book := dto.Book{
			Name:   "Le Rouge et le Noir",
			Author: "Stendhal",
			Price:  5,
		}

		buf, _ := book.MarshalBinary()
		err = bucket.Put([]byte(book.Name), buf)
		if err != nil {
			log.GetLoggerFromContext(ctx).Error("Impossible to store the data", "err", err)
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
			log.GetLoggerFromContext(ctx).Error("Impossible to store the data", "err", err)
			return err
		}

		return nil
	})

	return err
}

func ListBooks(ctx context.Context) (list []dto.Book, err error) {
	db, err := getDB(ctx)
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
				log.GetLoggerFromContext(ctx).Error("Impossible to UnmarshalBinary", "err", err)
				continue
			}

			list = append(list, b)
		}

		return nil
	})

	return
}

func GetBook(ctx context.Context, name string) (book dto.Book, found bool, err error) {
	db, err := getDB(ctx)
	if err != nil {
		return book, found, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		v := bucket.Get([]byte(name))
		if v == nil {
			log.GetLoggerFromContext(ctx).Info("Book not found", "bookName", name)
			return err
		}

		found = true

		err = book.UnmarshalBinary(v)
		if err != nil {
			log.GetLoggerFromContext(ctx).Error("Impossible to UnmarshalBinary", "err", err, "value", string(v))
			return err
		}
		return nil
	})

	return
}

func CreateBook(ctx context.Context, book dto.Book) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		buf, _ := book.MarshalBinary()
		err = bucket.Put([]byte(book.Name), buf)
		if err != nil {
			log.GetLoggerFromContext(ctx).Error("Impossible to store the data", "err", err)
			return err
		}
		return nil
	})
}

func DeleteBook(ctx context.Context, name string) error {
	db, err := getDB(ctx)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBooksName))

		err := bucket.Delete([]byte(name))
		if err != nil {
			log.GetLoggerFromContext(ctx).Error("Impossible to delete the data", "err", err)
			return err
		}
		return nil
	})
}
