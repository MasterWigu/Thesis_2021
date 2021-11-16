package main

func main() {
	createProcessor()
	createFileHandler()
	createToolHandler()

	ginEngine := declareAPIEndpoint()

	startListening(ginEngine)

}