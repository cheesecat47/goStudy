package echo

import (
	"os"
	"strings"

	"github.com/cheesecat47/hello/scrapper"
	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New()

	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)

	e.Logger.Fatal(e.Start(":8000"))
}

func handleScrape(c echo.Context) error {
	fileName := "jobs.csv"
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrapper(term)
	return c.Attachment(fileName, fileName)
}

func handleHome(c echo.Context) error {
	return c.File("echo/home.html")
}
