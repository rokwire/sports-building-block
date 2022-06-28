package main

type application struct {
	started               bool
	simulationTimeSeconds int
	exit                  chan bool
	Games                 []game
}

func main() {
	startApp()
}

func startApp() {
	app := createApplication()
	routerHandler := newRouter(app)
	routerHandler.setRoutes()
}

func createApplication() *application {
	app := application{exit: make(chan bool)}
	return &app
}
