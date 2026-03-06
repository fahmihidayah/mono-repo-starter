package utilities

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintBanner() {
	banner := `
   ____       ____
  / ___| ___ / ___| ___ _ __ ___   __ _
 | |  _ / _ \\___ \/ _ \ '_ ' _ \ / _' |
 | |_| | (_) |___) |  __/ | | | | | (_| |
  \____|\___/|____/ \___|_| |_| |_|\__,_|

  Go Code Generator for Domain-Driven Design
`
	color.Cyan(banner)
}

func PrintSuccess() {
	fmt.Println()
	color.Green("╔════════════════════════════════════════╗")
	color.Green("║   Code Generation Completed! 🎉       ║")
	color.Green("╔════════════════════════════════════════╗")
	fmt.Println()
}
