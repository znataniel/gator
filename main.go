package main

import "github.com/znataniel/gator/internal/config"

func printConfig(cfg config.Config) {
	println("db url:\t", cfg.DbUrl)
	println("current user:\t", cfg.CurrentUser)
}

func main() {
	println("Just started")
	println("reading file before writing username")
	println("------------------------")
	cfg, _ := config.Read()
	printConfig(cfg)
	println("------------------------")
	println()
	username := "hailie_selassie"
	println("writing username:", username)
	cfg.SetUser(username)
	cfg, _ = config.Read()
	printConfig(cfg)
	println("------------------------")
	println("end")

}
