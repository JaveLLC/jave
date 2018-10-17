package timed

import (
	"fmt"
	"time"
	//	"github.com/ark-lang/ark/src/util"
)

// magic
var indent int = 0

// Timed explode
func Timed(titleColored, titleUncolored string, fn func()) {
	//var bold string
	//if indent == 0 {
	//	bold = util.TEXT_BOLD
	//}

	if titleUncolored != "" {
		titleUncolored = " " + titleUncolored
	}

	fmt.Println("")
	fmt.Printf("Started - %s: %s\n", titleColored, titleUncolored)
	//Verbose("main", strings.Repeat(" ", indent))
	//Verboseln("main", bold+util.TEXT_GREEN+"Started "+titleColored+util.TEXT_RESET+titleUncolored)
	start := time.Now()

	indent++
	fn()
	indent--

	duration := time.Since(start)
	fmt.Printf("Ended - %s: %s (took %v)\n", titleColored, titleUncolored, duration)
	//Verbose("main", strings.Repeat(" ", indent))
	//Verboseln("main", bold+util.TEXT_GREEN+"Ended "+titleColored+util.TEXT_RESET+titleUncolored+" (%.2fms)", float32(duration)/1000000)
}
