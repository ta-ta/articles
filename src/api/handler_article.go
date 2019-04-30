package main

import (
	"errors"
	"log"
	"strings"

	"github.com/labstack/echo"
	"github.com/najeira/echo/echoutil"

	"../db"
)

func handlerArticles(c echo.Context) error {
	// 記事の並び順
	created_order := echoutil.ParamString(c, "created", "desc")
	read := echoutil.ParamInt(c, "read", 0)
	// articles取得
	var articles []*db.Article
	var err error
	if articles, err = db.DB.FetchArticles(read, created_order); err != nil {
		log.Fatalf("FetchArticles Error: %v\n", err)
	}
	for _, article := range articles {
		article.Domain = strings.Split(article.URL, "/")[2]
	}
	return render(c, "article", map[string]interface{}{"articles": articles, "created_order": created_order, "read": read})
}

func handlerUpdateRead(c echo.Context) error {
	articleID := echoutil.ParamInt(c, "articleID", 0)
	read := echoutil.ParamInt(c, "read", -1)
	if articleID <= 0 || read < 0 {
		return errors.New("Invalid articleID or read")
	}
	if err := db.DB.UpdateRead(articleID, read); err != nil {
		log.Fatalf("UpdateRead Error: %v\n", err)
		return err
	}
	return nil
}

func handlerUpdatePriproty(c echo.Context) error {
	articleID := echoutil.ParamInt(c, "articleID", 0)
	priority := echoutil.ParamInt(c, "priority", 0)
	if articleID <= 0 || priority <= 0 {
		return errors.New("Invalid articleID or priority")
	}
	if err := db.DB.UpdateProprity(articleID, priority); err != nil {
		log.Fatalf("UpdateProprity Error: %v\n", err)
		return err
	}
	return nil
}
