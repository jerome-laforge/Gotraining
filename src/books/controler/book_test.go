package controler_test

import (
	"books/controler"
	"books/dao"
	"books/dto"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/stretchr/testify/assert"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

var logger = log15.New(log15.Ctx{"component": "test"})

func TestCreateBookAndCheckIt(t *testing.T) {
	router := controler.CreateRouter(logger)
	server := httptest.NewServer(router)
	defer os.Remove(dao.DBName)
	defer server.Close()
	defer dao.Close()

	book := dto.Book{
		Name:   "Le Grand Meaulnes",
		Author: " Alain-Fournier",
		Price:  10,
	}

	jsonBook, err := book.MarshalBinary()
	assert.Nil(t, err, "err should be nil")

	req, err := http.NewRequest("POST", server.URL+controler.GroupPath+controler.BookPath, bytes.NewBuffer(jsonBook))
	assert.Nil(t, err, "err should be nil")

	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	req, err = http.NewRequest("GET", server.URL+controler.GroupPath+controler.BookPath, nil)
	assert.Nil(t, err, "err should be nil")

	queryValue := url.Values{controler.QueryName: []string{book.Name}}
	req.URL.RawQuery = queryValue.Encode()

	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err, "err should be nil")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	buf, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, "err should be nil")

	assert.JSONEq(t, string(jsonBook), string(buf), "json is not equivalent")
}
