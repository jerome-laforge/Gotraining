package controler

import (
	"books/dao"
	"books/dto"
	"books/log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func CreateRouter() http.Handler {
	router := gin.Default()

	//Add Swagger
	router.StaticFS(SwaggerPath, http.Dir("www/swagger-ui/dist"))

	// Add business
	apiV1 := router.Group(GroupPath, log.SetLogger)
	apiV1.GET(BooksPath, listBooks)
	apiV1.GET(BookPath, getBook)
	apiV1.POST(BookPath, createBook)
	apiV1.DELETE(BookPath, deleteBook)

	return router
}

func listBooks(ctx *gin.Context) {
	books, err := dao.ListBooks(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, books)
}

func getBook(ctx *gin.Context) {
	name := ctx.Query(QueryName)
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	book, found, err := dao.GetBook(ctx, name)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if !found {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func createBook(ctx *gin.Context) {
	book := dto.Book{}
	err := ctx.BindJSON(&book)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	_, found, err := dao.GetBook(ctx, book.Name)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if found {
		ctx.Status(http.StatusConflict)
		return
	}

	err = dao.CreateBook(ctx, book)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusCreated)
}

func deleteBook(ctx *gin.Context) {
	name := ctx.Query(QueryName)
	if name == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	_, found, err := dao.GetBook(ctx, name)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if !found {
		ctx.Status(http.StatusNotFound)
		return
	}

	err = dao.DeleteBook(ctx, name)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
