package main

import (
	//"time"
	configuration "me/gitoperator/configuration"
	topicconsumer "me/gitoperator/topicconsumer"
	//mediator "me/gitoperator/mediator"
	//utils "me/gitoperator/utils"
)
const componentMessage = "Main process"
var config = configuration.GlobalConfiguration

func main() {
	//utils.PrintLogInfo("GitPoc", componentMessage, "Start listening topic "+config.Kafka.Consumertopic)
	//topicconsumer.StartListening()
	topicconsumer.StartListeningBatches()
/* 	if err != nil {
		utils.PrintLogError(err, componentMessage, "Main", "Main function failed")
	} */
}

/* func startScheduledTasks(c configuration.Configuration){
	methodMessage := "startScheduledTasks"
	for true {
		time.Sleep(time.Duration(config.Microbatchfrequency) * time.Hour)
		utils.PrintLogInfo(componentMessage, methodMessage, "Scheduled action to run micro batches from received events: %s")
		mediator.ProcessMicroBatch()
	}
} */
