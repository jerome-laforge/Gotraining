package main

import (
	"books/dao"

	"github.com/inconshreveable/log15"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

func main() {
	log15.Root().SetHandler(log15.CallerStackHandler("%+v", log15.StdoutHandler))
	log15.Info("Startup ...")

	defer dao.Close()
	err := dao.InitBucketBooks()
	if err != nil {
		panic(err)
	}

	books, err := dao.ListBooks()
	if err != nil {
		panic(err)
	}

	for _, book := range books {
		b, _ := book.MarshalBinary()
		log15.Info(string(b))
	}
}
