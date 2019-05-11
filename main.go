package main

import (
	"github.com/drill-shishamo-alliance/asotech_server/interface/router"
)

func main() {
	routerHandler := router.NewRouterHandler()
	r := routerHandler.SetUpRouter()
	err := r.Run(":3001")
	if err != nil {
		panic(err)
	}
}