package controler

import (
	"books/dao"
	"books/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

const (
	SwaggerPath = "/swagger"
	GroupPath   = "/api/v1"
	BookPath    = "/book"
	BooksPath   = "/books"

	QueryName = "name"
)

func CreateRouter(logger log15.Logger) http.Handler {
	router := gin.Default()

	//Add Swagger
	router.StaticFS(SwaggerPath, http.Dir("www/swagger-ui/dist"))

	// Add business
	apiV1 := router.Group(GroupPath)
	apiV1.GET(BooksPath, listBooks)
	apiV1.GET(BookPath, getBook)
	apiV1.POST(BookPath, createBook)
	apiV1.DELETE(BookPath, deleteBook)

	return router
}

func listBooks(c *gin.Context) {
	books, err := dao.ListBooks()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	name := c.Query(QueryName)
	if name == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	book, found, err := dao.GetBook(name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	book := dto.Book{}
	err := c.BindJSON(&book)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, found, err := dao.GetBook(book.Name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if found {
		c.Status(http.StatusConflict)
		return
	}

	err = dao.CreateBook(book)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func deleteBook(c *gin.Context) {
	name := c.Query(QueryName)
	if name == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	_, found, err := dao.GetBook(name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if !found {
		c.Status(http.StatusNotFound)
		return
	}

	err = dao.DeleteBook(name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
