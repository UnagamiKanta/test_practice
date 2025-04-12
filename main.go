package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Memo struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Context string `json:"context"`
}

func main() {
	e := echo.New()
	e.Debug = true

	e.GET("/memo", getAllMemos)
	e.POST("/memo", createMemo)
	e.PUT("/memo/:id", updateMemo)
	e.DELETE("/memo/:id", deleteMemo)

	e.Logger.Fatal(e.Start(":8080"))
}

func getAllMemos(c echo.Context) error {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/go-college")
	if err != nil {
		log.Println("Error connecting to database:", err)
		return c.String(500, "Error connecting to database")
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM memo")
	if err != nil {
		return c.String(500, "Error fetching memos")
	}
	defer rows.Close()

	var memos []Memo
	for rows.Next() {
		var memo Memo
		if err := rows.Scan(&memo.ID, &memo.Title, &memo.Context); err != nil {
			return c.String(500, "Error scanning memo")
		}
		memos = append(memos, memo)
	}

	return c.JSON(200, memos)
}

func createMemo(c echo.Context) error {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/go-college")
	if err != nil {
		return c.String(500, "Error connecting to database")
	}

	defer db.Close()

	memo := Memo{}
	if err := c.Bind(&memo); err != nil {
		return c.String(400, "Invalid input")
	}

	result, err := db.Exec("INSERT INTO memo (title, context) VALUES (?, ?)", memo.Title, memo.Context)
	if err != nil {
		return c.String(500, "Error creating memo")
	}

	id, _ := result.LastInsertId()
	memo.ID = int(id)

	return c.JSON(201, memo)
}

func updateMemo(c echo.Context) error {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/go-college")
	if err != nil {
		return c.String(500, "Error connecting to database")
	}

	defer db.Close()

	id := c.Param("id")
	memo := Memo{}
	if err := c.Bind(&memo); err != nil {
		return c.String(400, "Invalid input")
	}

	_, err = db.Exec("UPDATE memo SET title = ?, context = ? WHERE id = ?", memo.Title, memo.Context, id)
	if err != nil {
		return c.String(500, "Error updating memo")
	}

	return c.JSON(200, memo)
}

func deleteMemo(c echo.Context) error {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/go-college")
	if err != nil {
		return c.String(500, "Error connecting to database")
	}

	defer db.Close()

	id := c.Param("id")

	_, err = db.Exec("DELETE FROM memo WHERE id = ?", id)
	if err != nil {
		return c.String(500, "Error deleting memo")
	}

	return c.NoContent(204)
}
