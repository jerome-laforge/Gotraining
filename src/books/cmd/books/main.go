package main

import (
	"books/controler"
	"books/dao"
	"books/log"
	"net/http"

	"golang.org/x/net/context"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

func main() {
	defer dao.Close(context.Background())
	logger := log.GetLogger()
	logger.Info("Startup ...")

	err := dao.InitBucketBooks(context.Background())
	if err != nil {
		panic(err)
	}

	router := controler.CreateRouter()
	err = http.ListenAndServe(":8080", router)
	logger.Crit("An error has occurred during the startup of server", "err", err)
}
