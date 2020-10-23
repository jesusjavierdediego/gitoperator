package main

/**
The GIT Operator!
It matches to a single GIT repository specified in the config.
Only 1 running instance of the GIT Operator can run simultaneously.
It listens to a given topic that matches to the GIT repository.
*/

import (
	configuration "xqledger/gitoperator/configuration"
	topics "xqledger/gitoperator/kafka"
	utils "xqledger/gitoperator/utils"
)
const componentMessage = "GIT Operator Main process"
var config = configuration.GlobalConfiguration

func main() {
	utils.PrintLogInfo("GitOperator", componentMessage, "Start listening topic "+config.Kafka.Consumertopic)
	topics.StartListening()
}
