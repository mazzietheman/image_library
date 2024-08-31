package main

import (
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

// input struct for resize
// used for binding form structure from frontend
type ResizeInput struct {
	Width int `form:"width"`
}

// input struct for crop
type CropInput struct {
	Width  int `form:"width"`
	Height int `form:"height"`
}

type AdjustContrast struct {
	Percentage float64 `form:"percentage"`
}

func main() {
	r := gin.New()
	//handle request from http://localhost:8080/resize_image
	r.POST("/resize_image", func(c *gin.Context) {
		//binding post data from frontend to input variable
		var input ResizeInput
		c.Bind(&input)
		width := input.Width
		//get uploaded file from frontend
		file, _ := c.FormFile("file")
		//get image type (JPEG or PNG)
		contentType := file.Header.Get("Content-Type")
		//only JPEG and PNG image type is allowed
		if contentType == "image/png" || contentType == "image/jpeg" {
			//"C:\image_folder" is my preferred folder location
			//you can choose the folder location according to your needs.
			//strconv.Itoa is command to convert int to string
			filePath := `C:\image_folder\resize-` + strconv.Itoa(width) + file.Filename
			//put uploaded image to selected folder
			c.SaveUploadedFile(file, filePath)
			//open image with imaging library
			src, _ := imaging.Open(filePath, imaging.AutoOrientation(true))
			//resize image
			resize := imaging.Resize(src, width, 0, imaging.Lanczos)
			//replace uploaded image with new resized image
			if contentType == "image/jpeg" {
				//if image format is JPEG
				imaging.Save(resize, filePath, imaging.JPEGQuality(90))
			} else {
				//if image format is PNG
				imaging.Save(resize, filePath, imaging.PNGCompressionLevel(0))
			}
			//return JSON format data to show in frontend
			c.JSON(200, gin.H{
				"filePath":    filePath,
				"contentType": contentType,
			})
		} else {
			//uploaded file not supported
			c.JSON(415, gin.H{
				"message": "Unsupported image type",
			})
			return
		}
	})
	//handle request from http://localhost:8080/crop_image
	r.POST("/crop_image", func(c *gin.Context) {
		//binding post data from frontend to input variable
		var input CropInput
		c.Bind(&input)
		width := input.Width
		height := input.Height
		//get uploaded file from frontend
		file, _ := c.FormFile("file")
		//get image type (JPEG or PNG)
		contentType := file.Header.Get("Content-Type")
		//only JPEG and PNG image type is allowed
		if contentType == "image/png" || contentType == "image/jpeg" {
			//"C:\image_folder" is my preferred folder location
			//you can choose the folder location according to your needs.
			//strconv.Itoa is command to convert int to string
			filePath := `C:\image_folder\crop-` + strconv.Itoa(width) + file.Filename
			//put uploaded image to selected folder
			c.SaveUploadedFile(file, filePath)
			//open image with imaging library
			src, _ := imaging.Open(filePath, imaging.AutoOrientation(true))
			//crop image
			crop := imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)
			//replace uploaded image with new cropped image
			if contentType == "image/jpeg" {
				//if image format is JPEG
				imaging.Save(crop, filePath, imaging.JPEGQuality(90))
			} else {
				//if image format is PNG
				imaging.Save(crop, filePath, imaging.PNGCompressionLevel(0))
			}
			//return file path
			c.JSON(200, gin.H{
				"filePath":    filePath,
				"contentType": contentType,
			})
		} else {
			//uploaded file not supported
			c.JSON(415, gin.H{
				"message": "Unsupported image type",
			})
			return
		}
	})
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"OPTIONS", "GET", "POST"},
	})
	handler := c.Handler(r)
	http.ListenAndServe("localhost:8080", handler)

	//handle request from http://localhost:8080/resize_image
	r.POST("/adjust_contrast", func(c *gin.Context) {

		var input AdjustContrast
		c.Bind(&input)

		percentage := input.Percentage
		//get uploaded file from frontend
		file, _ := c.FormFile("file")
		//get image type (JPEG or PNG)
		contentType := file.Header.Get("Content-Type")
		//only JPEG and PNG image type is allowed
		if contentType == "image/png" || contentType == "image/jpeg" {
			//"C:\image_folder" is my preferred folder location
			//you can choose the folder location according to your needs.
			//strconv.Itoa is command to convert int to string
			filePath := `C:\image_folder\contrast-` + strconv.FormatFloat(percentage, 'f', 0, 64) + file.Filename
			//put uploaded image to selected folder
			c.SaveUploadedFile(file, filePath)
			//open image with imaging library
			src, _ := imaging.Open(filePath, imaging.AutoOrientation(true))
			//resize image
			contrast := imaging.AdjustContrast(src, -20)
			//replace uploaded image with new resized image
			if contentType == "image/jpeg" {
				//if image format is JPEG
				imaging.Save(contrast, filePath, imaging.JPEGQuality(90))
			} else {
				//if image format is PNG
				imaging.Save(contrast, filePath, imaging.PNGCompressionLevel(0))
			}
			//return JSON format data to show in frontend
			c.JSON(200, gin.H{
				"filePath":    filePath,
				"contentType": contentType,
			})
		} else {
			//uploaded file not supported
			c.JSON(415, gin.H{
				"message": "Unsupported image type",
			})
			return
		}

	})
}
