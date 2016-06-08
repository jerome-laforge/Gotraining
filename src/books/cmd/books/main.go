package main

import "github.com/inconshreveable/log15"

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

func main() {
	log15.Root().SetHandler(log15.CallerStackHandler("%+v", log15.StdoutHandler))
	log15.Info("HelloWorld")
}
