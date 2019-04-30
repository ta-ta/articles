package main

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"../db"
	"../env"

	"net/http"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

var (
	templateExtension                     = "html"
	templateSet       *pongo2.TemplateSet = pongo2.DefaultSet
)

type AssetLoader func(string) ([]byte, error)
type AssetTemplateLoader struct {
	baseDir string
	loader  AssetLoader
}

func (l *AssetTemplateLoader) Abs(base, name string) string {
	name = normalize(name)
	if len(l.baseDir) <= 0 {
		if len(base) <= 0 {
			return name
		}
		return filepath.Join(filepath.Dir(base), name)
	} else if strings.HasPrefix(name, l.baseDir+string(filepath.Separator)) {
		return name
	}
	return filepath.Join(l.baseDir, name)
}

func (l *AssetTemplateLoader) Get(path string) (io.Reader, error) {
	body, err := l.loader(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}
func normalize(name string) string {
	if !strings.HasSuffix(name, "."+templateExtension) {
		name = name + "." + templateExtension
	}
	return name
}
func init() {
	templateSet = pongo2.NewSet("default", &AssetTemplateLoader{
		baseDir: "template",
		loader:  Asset,
	})
}

func main() {
	db.Init()

	//handlers
	app := echo.New()
	setupHandlers(app)
	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", env.Port)))
}

func render(c echo.Context, name string, ctx map[string]interface{}) error {
	filename := name
	if ctx == nil {
		ctx = make(map[string]interface{})
	}
	tpl, err := templateSet.FromCache(filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	pctx := pongo2.Context(ctx)
	pctx.Update(pongo2.Context{
		"c": c,
	})
	body, err := tpl.Execute(pctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, body)
}
