package main

import (
	topics "me/gitpoc/topics"
	utils "me/gitpoc/utils"
	configuration "me/gitpoc/configuration"
)

var config = configuration.GlobalConfiguration

func main() {
	
	const componentMessage = "Main process"
	utils.PrintLogInfo("GitPoc", componentMessage, "Start listening topic " + config.Kafka.Consumertopic)
	topics.StartListening()
}