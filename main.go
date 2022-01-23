package main

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
)

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

type Posts []Post

type Comments []Comment

func connect() *gorm.DB {
	dsn := "user:ucUeB2dlViI48Gfk%@tcp(127.0.0.1:3306)/publish_part2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect databases")
	}
	return db
}

func main() {
	err := connect().AutoMigrate(&Posts{}, &Comments{})
	if err != nil {
		panic("can't migrate the schema")
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/comments", getAllComments)
	e.POST("/comments", createComment)
	e.DELETE("/comments/:id", deleteComment)

	e.Logger.Fatal(e.Start(":1323"))
}

func createComment(context echo.Context) error {
	comment := &Comment{}

	if err := context.Bind(comment); err != nil {
		return err
	}
	connect().Create(&Comment{
		Id:     comment.Id,
		PostId: comment.PostId,
		Name:   comment.Name,
		Email:  comment.Email,
		Body:   comment.Body,
	})

	return context.JSON(http.StatusCreated, comment)
}

func getAllComments(context echo.Context) error {
	var c []Comment

	connect().Raw("SELECT * FROM comments").Scan(&c)

	return context.JSON(http.StatusOK, c)
}
func deleteComment(context echo.Context) error {
	id, _ := strconv.Atoi(context.Param("id"))
	connect().Delete(&Comment{}, id)
	return context.NoContent(http.StatusNoContent)
}
