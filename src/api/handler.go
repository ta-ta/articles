package main

import (
	"bytes"
	"net/http"
	"path"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"../env"
)

func setupHandlers(app *echo.Echo) error {
	oauth := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == env.USERNAME && password == env.PASSWORD {
			return true, nil
		}
		return false, nil
	})

	app.GET("/", handlerDashboard)
	app.GET("/dashboard", handlerDashboard)
	app.GET("/articles", handlerArticles, oauth)
	app.GET("/:articleID/read", handlerUpdateRead, oauth)
	app.GET("/:articleID/priority", handlerUpdatePriproty, oauth)

	files := app.Group("")
	files.GET("/static/*", handlerStatic)

	return nil
}

func handlerStatic(c echo.Context) error {
	name := path.Join("static", c.Param("*"))

	info, err := AssetInfo(name)
	if err != nil {
		return err
	}

	body, err := Asset(name)
	if err != nil {
		return err
	}
	r := bytes.NewReader(body)

	c.Response().Header().Set("Cache-Control", "max-age=60")
	http.ServeContent(c.Response(), c.Request(), info.Name(), info.ModTime(), r)
	return nil
}
