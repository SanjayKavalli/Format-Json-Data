package main

import (
	"formatjsondata/controller"
	"formatjsondata/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//Intialze services
	ConstantService := services.ConstantserviceCtor()
	Jsonformatservice := services.JsonFormatserviceCtor(ConstantService)

	//controller
	JsonformatController := controller.JsonControllerCtor(Jsonformatservice, ConstantService)

	//routes
	r.POST("/formatdata", JsonformatController.JsonFormat)
	r.Run()
}
