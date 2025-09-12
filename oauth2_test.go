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

	// Create an new Get
	ap := &application.Application{}
	ap.Name = "TestApp"
	ap.Description = "TestDesc"
	ap, err := dg.ApplicationAPI().Create(ap)
	log.Printf("ApplicationCreate: err: %+v, app: %+v\n", err, ap)

	// Get a specific Get by it's ID
	ap, err = dg.ApplicationAPI().Get(ap.ID)
	log.Printf("Get: err: %+v, app: %+v\n", err, ap)

	// Update an existing Get with new values
	ap.Description = "Whooooa"
	ap, err = dg.ApplicationAPI().Update(ap.ID, ap)
	log.Printf("ApplicationUpdate: err: %+v, app: %+v\n", err, ap)

	// create a new bot account for this application
	bot, err := dg.ApplicationAPI().BotCreate(ap.ID)
	log.Printf("BotCreate: err: %+v, bot: %+v\n", err, bot)

	// Get a list of all applications for the authenticated user
	apps, err := dg.ApplicationAPI().GetAll()
	log.Printf("GetAll: err: %+v, apps : %+v\n", err, apps)
	for k, v := range apps {
		log.Printf("GetAll: %d : %+v\n", k, v)
	}

	// Delete the application we created.
	err = dg.ApplicationAPI().Delete(ap.ID)
	log.Printf("Delete: err: %+v\n", err)

	return
}
