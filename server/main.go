package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:4000)/test")
	if err != nil {
		log.Fatalln(err)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/ST_Distance", func(c *gin.Context) {
		x1 := c.PostForm("x1")
		y1 := c.PostForm("y1")
		x2 := c.PostForm("x2")
		y2 := c.PostForm("y2")

		var result float64
		queryString := fmt.Sprintf("SELECT ST_Distance(POINT(%s, %s), POINT(%s, %s));", x1, y1, x2, y2)
		fmt.Println(queryString)
		row := db.QueryRow(queryString)
		err = row.Scan(&result)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(result)
		c.JSON(200, gin.H{
			"result": result,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
