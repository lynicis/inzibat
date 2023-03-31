package main

func main() {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	server := NewServer(config)
	app := server.GetFiberInstance()

	client := NewClient()
	router := NewRouter(config, app, client)
	router.CreateRoutes()

	err = server.Start()
	if err != nil {
		panic(err)
	}
}
