package boot

func Boot() error {
	initFlag()
	generateConfig()

	if err := initHttpServer(); err != nil {
		return err
	}
	return nil
}
