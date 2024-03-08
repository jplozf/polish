package main

import (
	"fmt"
	"polish/color"
)

// *****************************************************************************
//  __________      .__  .__       .__
//  \______   \____ |  | |__| _____|  |__
//   |     ___/  _ \|  | |  |/  ___/  |  \
//   |    |  (  <_> )  |_|  |\___ \|   Y  \
//   |____|   \____/|____/__/____  >___|  /
//                               \/     \/
//
//          Polish © jpl@ozf.fr 2024
//
// *****************************************************************************

// *****************************************************************************
// TYPE & CONSTANTS
// *****************************************************************************
type Error int

const (
	OK Error = iota
	UNRECOGNIZED_COMMAND
	NOT_ENOUGH_ARGS_FLOAT
	NOT_ENOUGH_ARGS_ALPHA
	ALPHA_EXTRACT_OUT_OF_BOUNDS
)

func (e Error) String() string {
	return [...]string{
		"OK",
		"Unrecognized command",
		"Not enough arguments on float stack",
		"Not enough arguments on alpha stack",
		"Alpha extract out of bounds",
	}[e]
}

// *****************************************************************************
// raiseError()
// *****************************************************************************
func raiseError(e Error) {
	fmt.Printf("[E%04d]\t%s%s%s\n", e, color.Red, e.String(), color.Reset)
}
