package main

import (
	configuration "me/gitoperator/configuration"
	topics "me/gitoperator/topics"
	utils "me/gitoperator/utils"
)
const componentMessage = "Main process"
var config = configuration.GlobalConfiguration

func main() {
	utils.PrintLogInfo("GitPoc", componentMessage, "Start listening topic "+config.Kafka.Consumertopic)
	topics.StartListening()
}

