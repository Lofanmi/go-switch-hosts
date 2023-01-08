package main

func main() {
	initLogger()
	initConfig()

	application, cleanup, err := NewApplication()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	application.Run()
}
