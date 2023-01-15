package main

func main() {
	defer initLogger()()
	defer initConfig()()

	application, cleanup, err := NewApplication()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	application.Run()
}
