package main

func main() {
	config, err := InitConfig()
	if err != nil {
		panic(err)
	}

	if err := StartTUI(config); err != nil {
		panic(err)
	}

	if err := StoreConfig(config); err != nil {
		panic(err)
	}
}
