package utils

import (
	"fmt"

	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	// every day every hour
	c.AddFunc("0 * * * *", runEngine)
	c.Start()

	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs \n", cronCount)
}

func runEngine() {
	fmt.Println("starting up engine")
}
