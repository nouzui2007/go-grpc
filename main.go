package main

import (
	"goexample/router"
)

func main() {

	// init router
	r := router.New()

	// サーバー起動
	r.Start(":80")

}
