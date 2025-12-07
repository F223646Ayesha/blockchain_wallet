package services

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// GLOBAL scheduler instance
var Scheduler *cron.Cron

func InitScheduler() {
	Scheduler = cron.New()

	// ----------------------------------------
	// 🔥 RUN ZAKAT ON 1st OF EVERY MONTH AT 00:00
	// Cron expression: MIN HOUR DOM MONTH DOW
	// ----------------------------------------
	_, err := Scheduler.AddFunc("0 0 1 * *", func() {
		fmt.Println("⏳ Monthly Zakat Cron Triggered")
		err := RunZakatService()
		if err != nil {
			fmt.Println("❌ Zakat scheduler failed:", err)
		} else {
			fmt.Println("✅ Monthly Zakat processed successfully")
		}
	})
	// ----------------------------------------
	// 🔥 TEST MODE — Run Zakat EVERY 1 MINUTE
	// ----------------------------------------
	/*_, err = Scheduler.AddFunc("@every 1m", func() {
		fmt.Println("🧪 TEST MODE: 1-minute zakat triggered")
		if err := RunZakatService(); err != nil {
			fmt.Println("❌ Test zakat failed:", err)
		} else {
			fmt.Println("✅ Test zakat processed successfully")
		}
	})*/

	if err != nil {
		fmt.Println("❌ Failed to schedule zakat:", err)
	}

	// START THE SCHEDULER
	Scheduler.Start()

	fmt.Println("⏰ Monthly Zakat Scheduler Started.")
}
