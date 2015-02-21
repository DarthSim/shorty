package main

func main() {
	prepareConfig()

	initLogger()
	defer closeLogger()

	initDB(true)
	defer closeDB()

	startServer()
}
