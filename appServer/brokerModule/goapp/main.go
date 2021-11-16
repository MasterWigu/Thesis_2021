package main


func main() {


	initLedger()
	initModules()
	createProcessor()
	createUsersBase()


	ginEngine := declareAPIEndpoint()

	startListening(ginEngine)

}
