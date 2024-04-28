package main

func main() {

	server := InitWireServer()

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
