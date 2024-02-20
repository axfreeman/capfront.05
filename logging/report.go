//logging.report.go

package logging

import (
	"capfront/models"
	"fmt"
)

// creates a trace that logs events on the standard output
// and also creates a user-readable record called 'Trace'
// that provides information about the progress of the simulation

//TODO distinguish between what is logged to output and what is shown to the user

func Report(level int, simulation models.Simulation, message string) {
	fmt.Println(message) //placeholder
}
