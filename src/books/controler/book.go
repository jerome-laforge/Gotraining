package controler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

func CreateRouter(logger log15.Logger) http.Handler {
	router := gin.Default()

	//Add Swagger
	router.StaticFS("/swagger", http.Dir("www/swagger-ui/dist"))

	return router
}
