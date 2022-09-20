package main

import (
	"home/controller"

	"github.com/jasonlvhit/gocron"
)

func main() {

	gocron.Every(10).Seconds().Do(controller.CreateNews)
	<-gocron.Start()

}
