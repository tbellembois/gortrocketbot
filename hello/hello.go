package hello

import (
	"github.com/tbellembois/gortrocketbot/rocket"
)

func hello(...string) string {
	return "hey!"
}

func init() {
	rocket.RegisterPlugin(rocket.Plugin{
		Name:        "hello",            // the rocket command to trigger the action
		CommandFunc: hello,              // command bind function
		Args:        []string{"a", "b"}, // command arguments (not used here, just for example)
		Help:        "Say hello to you", // command help
	})
}
