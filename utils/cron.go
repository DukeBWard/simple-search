package utils

import (
	"dukebward/search/search"
	"fmt"

	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	// every day every hour
	c.AddFunc("0 * * * *", search.RunEngine)
	// quarter past every hour
	c.AddFunc("15 * * * *", search.RunIndex)
	c.Start()

	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs \n", cronCount)
}
