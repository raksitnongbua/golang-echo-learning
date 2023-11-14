package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func uploadFileHandler(c echo.Context) error {
	_, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	defer file.Close()

	tempFile, err := os.CreateTemp("./uploads", "temp-*.pdf")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	//TODO: implement receive password for locked pdf to unlock pdf file.
	// Serve the temporary file.
	return c.File(tempFile.Name())
}

func main() {
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating uploads directory:", err)
		return
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/api/unlock-pdf", uploadFileHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
