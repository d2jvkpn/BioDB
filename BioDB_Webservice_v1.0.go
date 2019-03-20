package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/d2jvkpn/gopkgs/biodb2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var (
	DB     *sql.DB
	router *gin.Engine
	port   string
)

const (
	DBuser   = "world"
	DBpasswd = ""
	DBhost   = "tcp(localhost:3306)"
	USAGE    = `BioDB web service, usage:
  $ BioDB_Web_Service  [-p port]
`
	LISENSE = `
author: d2jvkpn
version: 1.0
release: 2019-03-18
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`
)

func main() {
	defer DB.Close()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Search.html", nil)
	})

	search := router.Group("search/")
	search.GET("/", func(c *gin.Context) { biodb.Search(c, DB) })
	search.POST("/", func(c *gin.Context) { biodb.Search(c, DB) })

	api := router.Group("api/")
	api.GET("/", func(c *gin.Context) { biodb.API(c, DB) })

	download := router.Group("download/")
	download.GET("/", func(c *gin.Context) { biodb.Download(c, DB) })

	router.Run(port)
}

func init() {
	var err error

	flag.StringVar(&port, "p", ":8090", "set port")

	flag.Usage = func() {
		fmt.Println(USAGE)
		flag.PrintDefaults()
		fmt.Println(LISENSE)
		os.Exit(2)
	}

	flag.Parse()

	err = ValidPort(&port)
	ErrExit(err)

	router = gin.Default()

	router.SetFuncMap(template.FuncMap{
		"Add": Add,
	})

	router.LoadHTMLGlob("templates/*.html")

	s := fmt.Sprintf("%s:%s@%s/BioDB", DBuser, DBpasswd, DBhost)
	DB, err = sql.Open("mysql", s)
	ErrExit(err)
}

func Add(a, b int) string {
	return strconv.Itoa(a + b)
}

func ValidPort(port *string) (err error) {
	var ok bool

	if ok, _ = regexp.MatchString("^[1-9][0-9]*$", *port); ok {
		*port = ":" + *port
		return
	}

	if ok, _ = regexp.MatchString("^:[1-9][0-9]*$", *port); !ok {
		err = fmt.Errorf("invalid port \"%s\"", *port)
	}

	return
}

func ErrExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
