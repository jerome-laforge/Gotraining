package main

import (
	"books/controler"
	"books/dao"
	"net/http"

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

	router := controler.CreateRouter(log15.New("component", "router"))
	err = http.ListenAndServe(":8080", router)
	log15.Crit("An error has occurred during the startup of server", "err", err)
}
