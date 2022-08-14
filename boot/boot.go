package boot

func Run() {
	initFlag()
	generateConfig()
	initHttpServer()
}
