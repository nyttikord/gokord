package gokord_test

import (
	"log"
	"os"

	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/application"
)

func ExampleApplication() {
	// Authentication Token pulled from environment variable DGU_TOKEN
	Token := os.Getenv("DGU_TOKEN")
	if Token == "" {
		return
	}

	// Create a new Discordgo session
	dg := gokord.New(Token)

	// Create an new Application
	ap := &application.Application{}
	ap.Name = "TestApp"
	ap.Description = "TestDesc"
	ap, err := dg.ApplicationAPI().ApplicationCreate(ap)
	log.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, ap)

	// Application a specific Application by it's ID
	ap, err = dg.ApplicationAPI().Application(ap.ID)
	log.Printf("Application: err: %+v, app: %+v\n", err, ap)

	// Update an existing Application with new values
	ap.Description = "Whooooa"
	ap, err = dg.ApplicationAPI().ApplicationUpdate(ap.ID, ap)
	log.Printf("ApplicationUpdate: err: %+v, app: %+v\n", err, ap)

	// create a new bot account for this application
	bot, err := dg.ApplicationAPI().BotCreate(ap.ID)
	log.Printf("BotCreate: err: %+v, bot: %+v\n", err, bot)

	// Application a list of all applications for the authenticated user
	apps, err := dg.ApplicationAPI().Applications()
	log.Printf("Applications: err: %+v, apps : %+v\n", err, apps)
	for k, v := range apps {
		log.Printf("Applications: %d : %+v\n", k, v)
	}

	// Delete the application we created.
	err = dg.ApplicationAPI().ApplicationDelete(ap.ID)
	log.Printf("Delete: err: %+v\n", err)

	return
}
