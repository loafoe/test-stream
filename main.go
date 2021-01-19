package main

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/download/:seconds/:size", streamHandler())

	e.Start(":8080")
}

func genRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

func streamHandler() echo.HandlerFunc {

	return func(c echo.Context) error {
		randomData, err := genRandomBytes(1024 * 1024)
		if err != nil {
			return fmt.Errorf("randomData: %w", err)
		}
		// seconds to keep the download going
		seconds, err := strconv.Atoi(c.Param("seconds"))
		if err != nil {
			return fmt.Errorf("minutes: %w", err)
		}
		// file size (MB)
		size, err := strconv.Atoi(c.Param("size"))
		if err != nil {
			return fmt.Errorf("rate: %w", err)
		}
		if seconds == 0 || size == 0 {
			return fmt.Errorf("minutes or size cannot be 0")
		}

		resp := c.Response()
		resp.Header().Set("Content-Type", "application/zip")
		resp.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")

		zipWriter := zip.NewWriter(c.Response())

		speed := float64(size) / float64(seconds) // MB/sec
		sleepTime := 1000 / speed                 // in millesconds
		startTime := time.Now()
		defer func() {
			fmt.Printf("Duration: %v\n", time.Now().Sub(startTime))
		}()

		fmt.Printf("[%v]: sending at %.2f MB/sec\n", startTime, speed)

		for i := 0; i < size; i++ {
			header := &zip.FileHeader{
				Name:     fmt.Sprintf("random%d.data", i),
				Method:   zip.Store,
				Modified: time.Now(),
			}
			entryWriter, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("entryWriter: %w", err)
			}
			_, err = entryWriter.Write(randomData)
			if err != nil {
				msg := fmt.Sprintf("Error writing: %v", err)
				fmt.Printf("%s\n", msg)
				return fmt.Errorf(msg)
			}
			zipWriter.Flush()
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		}
		fmt.Printf("Done sending")
		return zipWriter.Close()
	}
}
