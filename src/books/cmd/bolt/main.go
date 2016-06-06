package main

import (
	"books/dao"
	"time"

	"github.com/boltdb/bolt"
	"github.com/inconshreveable/log15"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

func main() {
	log15.Root().SetHandler(log15.CallerStackHandler("%+v", log15.StdoutHandler))
	log15.Info("Startup ...")

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log15.Error(err.Error())
	}
	defer db.Close()

	dao.InitBucketBooks(db)

	for _, book := range dao.ListBooks(db) {
		b, _ := book.MarshalBinary()
		log15.Info(string(b))
	}
}
