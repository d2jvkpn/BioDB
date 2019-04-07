package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	DB     *sql.DB
	router *gin.Engine
	fh     *os.File
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
version: 1.1
release: 2019-04-07
project: https://github.com/d2jvkpn/BioDB
lisense: GPLv3 (https://www.gnu.org/licenses/gpl-3.0.en.html)
`
)

func main() {
	defer DB.Close()
	defer fh.Close()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Search.html", nil)
	})

	search := router.Group("search/")
	search.GET("/", func(c *gin.Context) { Search(c, DB) })

	api := router.Group("api/")
	api.GET("/", func(c *gin.Context) { API(c, DB) })

	download := router.Group("download/")
	download.GET("/", func(c *gin.Context) { Download(c, DB) })

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

	s := fmt.Sprintf("%s:%s@%s/BioDB", DBuser, DBpasswd, DBhost)
	DB, err = sql.Open("mysql", s)
	ErrExit(err)

	router = gin.New() // router = gin.Default()

	// time.Now().Format("2006-01-02"), time.Now().UnixNano()
	fn := fmt.Sprintf("log_%X.txt", time.Now().Unix())
	fh, err = os.Create(fn)
	ErrExit(err)

	log.Println("Log will be written in", fn)

	gin.DefaultWriter = io.MultiWriter(fh)

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s  %s  %d  %s %s %s %s  \"%s\"  \"%s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.StatusCode,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery())

	router.SetFuncMap(template.FuncMap{
		"Add": Add,
	})

	router.LoadHTMLGlob("templates/*.html")

	fmt.Fprintf(fh, "BioDB_Webservice start at %s\n\n", time.Now().Format(time.RFC3339))
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
