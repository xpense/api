package main

import (
	"expense-api/database"
	"expense-api/transaction"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	expenses := database.NewDB()

	r.GET("/admin/overview/expense", func(c *gin.Context) {
		c.JSON(http.StatusOK, expenses)
	})

	r.POST("/:user/expense", func(c *gin.Context) {
		user := c.Param("user")

		var json transaction.Transaction
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		id := expenses.TotalTransactions + 1
		fmt.Println(expenses.TotalTransactions, id)
		t, err := transaction.New(id, &json.Date, json.Amount, json.Type)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		expenses.AddTransaction(user, t)

		c.JSON(http.StatusOK, gin.H{
			"id": t.ID,
		})
	})

	r.Run()
}
