package main

func main() {

	server := InitWireServer()

	err := server.Run(":8081")
	if err != nil {
		panic(err)
	}
}
