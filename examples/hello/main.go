package main

func main() {
	app, cancel, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cancel()
	app.Run()
}
