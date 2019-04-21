package main

import (
	"os"

	"github.com/tbellembois/gortrocketbot/rocket"

	// import the plugins here
	_ "github.com/tbellembois/gortrocketbot-plugins/hello"
)

func main() {

	// configure the bot
	c := rocket.Config{
		ServerHost:   os.Getenv("ROCKET_SERVERHOST"),
		ServerScheme: os.Getenv("ROCKET_SERVERSCHEME"),
		User:         os.Getenv("ROCKET_USER"),
		Email:        os.Getenv("ROCKET_EMAIL"),
		Password:     os.Getenv("ROCKET_PASSWORD"),
		Debug:        false,
	}

	// et voilà
	rocket.Run(&c)
}
