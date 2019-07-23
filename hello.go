package main

import (
	"context"
	"database/sql"
	"fmt"
	"hello/models"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

func GetDataTodo(ctx context.Context, c *gin.Context) {
	var b models.Todo
	if err := c.Bind(&b); err != nil {
		fmt.Errorf("%#v", err)
	}
	b.Status = 0
	err := b.InsertG(ctx, boil.Infer())

	todos, err := models.Todos().AllG(ctx)
	if err != nil {
		fmt.Errorf("Get todo error: %v", err)
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"todo": todos,
	})
}

func GetDoneTodo(ctx context.Context, c *gin.Context) {
	var b models.Todo
	if err := c.Bind(&b); err != nil {
		fmt.Errorf("%#v", err)
	}

	if b.Status == 0 {
		println("haiteru")
		b.Status = 1
	} else {
		b.Status = 0
	}

	af, err := b.UpdateG(ctx, boil.Whitelist("status", "updated_at"))
	println("Get todo error: %v", af)
	if err != nil {
		fmt.Errorf("Get todo error: %v", err)
	}

	todos, err := models.Todos().AllG(ctx)
	if err != nil {
		fmt.Errorf("Get todo error: %v", err)
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"todo": todos,
	})

}

func main() {
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/todo?parseTime=true")
	if err != nil {
		log.Fatalf("Cannot connect database: %v", err)
	}
	boil.SetDB(db)

	r := gin.Default()
	r.LoadHTMLFiles("./tmpl/index.html")
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "world",
		})
	})
	r.GET("/todo", func(c *gin.Context) {

		todos, err := models.Todos(OrderBy("update_at desc")).AllG(ctx)
		if err != nil {
			fmt.Errorf("Get todo error: %v", err)
		}

		c.HTML(http.StatusOK, "index.html", map[string]interface{}{
			"todo": todos,
		})
	})
	r.GET("/yaru", func(c *gin.Context) {
		GetDataTodo(ctx, c)
	})
	r.GET("/done", func(c *gin.Context) {
		GetDoneTodo(ctx, c)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
