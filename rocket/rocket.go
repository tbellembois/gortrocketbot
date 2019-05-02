package rocket

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
)

var (
	// Rocket.Chat.Go.SDK client and user
	rtc  *realtime.Client
	user *models.User

	// bot configuration
	config *Config
	// channels the bot has subscribed in
	channels []models.Channel
	// channels already subscribed
	channelsIds map[string]string
	// registered plugins
	plugins map[string]Plugin
	// plugins result message
	cmdResult sync.Pool

	debug bool
	e     error
)

func init() {
	plugins = make(map[string]Plugin)
	channelsIds = make(map[string]string)

	// initializing the pool
	newCmdResult := func() interface{} {
		return new(string)
	}
	cmdResult = sync.Pool{New: newCmdResult}
}

// RegisterPlugin adds the plugin p
// to the list of registered plugins
func RegisterPlugin(p Plugin) {
	fmt.Println(fmt.Printf("registered plugin: %s", p.Name))
	plugins[p.Name] = p
}

// Run connects to the Rocket server and start listening
// for requests
func Run(c *Config) {

	config = c

	// connecting to the server
	serverURL := &url.URL{Scheme: c.ServerScheme, Host: c.ServerHost}
	if rtc, e = realtime.NewClient(serverURL, debug); e != nil {
		log.Panic("can not connect to Rocket " + e.Error())
	}

	// login attempt
	if user, e = rtc.Login(&models.UserCredentials{
		Email:    c.Email,
		Name:     c.User,
		Password: c.Password,
	}); e != nil {
		log.Panic("can not login user " + e.Error())
	}

	// getting the channels the bot is member of
	if channels, e = rtc.GetChannelsIn(); e != nil {
		log.Panic("can not get channels in " + e.Error())
	}

	// subscribing to the channels
	msgChannel := make(chan models.Message, 10)

	// updating the subscribed channels regularly
	go func() {
		for {
			// getting the channels the bot is member of
			if channels, e = rtc.GetChannelsIn(); e != nil {
				log.Panic("can not get channels in " + e.Error())
			}

			for _, ch := range channels {
				if _, ok := channelsIds[ch.ID]; !ok {
					if config.Debug {
						fmt.Println(fmt.Printf("subscribed to: %v", ch))
					}
					if e = rtc.SubscribeToMessageStream(&models.Channel{ID: ch.ID}, msgChannel); e != nil {
						log.Println("can not subscribe to message stream for channel " + ch.ID + " " + e.Error())
					}
					// mark channel as subscribed
					channelsIds[ch.ID] = ch.ID
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// listening for incoming messages forever
	for {
		select {
		case m := <-msgChannel:

			if c.Debug {
				fmt.Println(fmt.Sprintf("ID: %s RoomID: %s Msg: %s User: %v", m.ID, m.RoomID, m.Msg, m.User))
			}

			go func(m models.Message) {

				// getting string from the pool
				cr := cmdResult.Get().(*string)

				// getting the command name from message
				splCmd := strings.Split(m.Msg, " ")
				if c.Debug {
					fmt.Println(fmt.Sprintf("cmd: %s cmdArgs: %v", splCmd[0], splCmd[1:]))
				}

				// should never happen
				if len(splCmd) == 0 {
					return
				}
				cmdName := splCmd[0]
				// does command has arguments?
				cmdHasArgs := false
				if len(splCmd) > 1 {
					cmdHasArgs = true
				}

				// help?
				if cmdName == "help" || cmdName == "?" || cmdName == "aide" {
					// building the help message
					helpMsg := ""
					for k, v := range plugins {
						if v.IsAllowed(*m.User) {
							if config.Debug {
								fmt.Println(fmt.Sprintf("plugin: %s help: %s", k, v.Help))
							}
							helpMsg += fmt.Sprintf("`%s` %s\n", k, v.Help)
						}
					}

					*cr = helpMsg
				} else if cmd, ok := plugins[cmdName]; ok {
					if cmd.IsAllowed(*m.User) {
						// executing the command
						if cmdHasArgs {
							*cr = cmd.CommandFunc(splCmd[1:]...)
						} else {
							*cr = cmd.CommandFunc()
						}
						if c.Debug {
							fmt.Println(fmt.Sprintf("executing command %s", cmd.Name))
							fmt.Println(fmt.Sprintf("command result: %s", *cr))
						}
					}
				}

				// sending the response
				if *cr != "" {
					rtc.SendMessage(&models.Message{
						RoomID: m.RoomID,
						Msg:    *cr,
						User:   m.User,
					})
				}

				// returning to the pool
				*cr = ""
				cmdResult.Put(cr)

			}(m)
		}
	}

}
