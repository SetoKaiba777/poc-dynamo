package main

import (
	"fmt"
	"poc-dynamo/database"
)

func main() {
	c:=database.NewDynamoDBClient("payments","us-east-1","http://localhost:4566")
	out, err := c.GetItemByID("id","payment1")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Payment: ", out)
}