package main

import (
	"github.com/labstack/echo"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"github.com/lmnzr/simpleshop/cmd/simpleshop/handler"
	"testing"
)

func TestHello(t *testing.T) {
	router := echo.New()
	router.GET("/", handler.Hello)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	context := router.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, handler.Hello(context)) {
		status,_ := jsonparser.GetInt(rec.Body.Bytes(),"s")
		jsondata,_ := jsonparser.GetString(rec.Body.Bytes(),"d")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t,200,int(status))
		assert.Equal(t,"Hello World !!!",jsondata)
	}
}
