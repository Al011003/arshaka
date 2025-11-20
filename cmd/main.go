package main

import (
	"backend/config"
	"fmt"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		fmt.Println("gagal coy:", err)
		return
	}

	fmt.Println("berhasilll", db)
}