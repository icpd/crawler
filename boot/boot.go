package boot

func Init() {
	initFlag()
	generateConfig()
	initHttpServer()
}
