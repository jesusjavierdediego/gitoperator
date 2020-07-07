package main

import (
	topics "me/gitpoc/topics"
)

func main() {
	
	const componentMessage = "Main process"

	topics.StartListening()
}