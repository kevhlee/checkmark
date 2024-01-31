package main

func main() {
	config, err := InitTaskConfig()
	if err != nil {
		panic(err)
	}

	if err := StartTUI(config); err != nil {
		panic(err)
	}

	if err := StoreTaskConfig(config); err != nil {
		panic(err)
	}
}
