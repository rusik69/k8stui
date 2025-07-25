package main

import (
	"fmt"
	"os"

	"github.com/yourusername/k8stui/internal/k8s/app"
)

func main() {
	// Create a new app instance
	appInstance := app.NewApp()

	// Run the application
	if err := appInstance.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
