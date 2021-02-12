package main

import (
	"kakao/mappings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	mappings.CreateURLMappings()
	// Listen and server on 0.0.0.0:8000
	mappings.Router.Run(":8000")

}
