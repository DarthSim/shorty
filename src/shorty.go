package main

func main() {
	initDB(true)
	defer closeDB()

	startServer()
}
