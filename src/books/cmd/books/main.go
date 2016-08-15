package main

import (
	"books/dao"
	"fmt"

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
		log15.Info(fmt.Sprintf("%#v", book))
	}
}
