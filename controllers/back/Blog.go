package back

import (
	"akh/blog/helpers"
	"akh/blog/models"
	"akh/blog/requests"
	"akh/blog/responses"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BlogCreate(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// var data requests.BlogRequest

	// if err := c.ShouldBind(&data); err != nil {
	// 	c.JSON(http.StatusBadRequest, responses.BlogResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: map[string]interface{}{}})
	// 	return
	// }

	_, fileHandler, ferr := c.Request.FormFile("image")
	var saveDirName string

	if ferr == nil {
		filename, err := helpers.ImageProcessing(c, "blogs", fileHandler)
		saveDirName = filename
		if err != nil {
			fmt.Println(err)
		}
		// defer file.Close()
	}

	data := requests.BlogRequest{
		Title:       c.Request.FormValue("title"),
		Description: c.Request.FormValue("description"),
		Tags:        c.Request.FormValue("tags"),
		Category:    c.Request.FormValue("category"),
		Image:       saveDirName,
	}
	fmt.Println(saveDirName)

	if validateErr := validate.Struct(&data); validateErr != nil {
		c.JSON(http.StatusBadRequest, responses.BlogResponse{Status: http.StatusBadRequest, Message: validateErr.Error(), Data: map[string]interface{}{}})
		return
	}

	// fmt.Println(c.Request)

	result, err := models.AddBlog(data)

	if !err {
		c.JSON(http.StatusInternalServerError, responses.BlogResponse{Status: http.StatusInternalServerError, Message: "Checking database", Data: map[string]interface{}{}})
		return
	}

	c.JSON(http.StatusCreated, responses.BlogResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": "Successfully created.", "insertId": result}})
}
