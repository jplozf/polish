package main

import "runtime"

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
// CONSTS
// *****************************************************************************
const (
	APP_NAME    = "Polish"
	APP_VERSION = "0.1.0"
	APP_STRING  = "Polish © jpl@ozf.fr 2024"
	APP_URL     = "https://github.com/jplozf/polish"
	APP_FOLDER  = ".polish"
	FSTACK_FILE = "fstack.gob"
	SSTACK_FILE = "sstack.gob"
)

// *****************************************************************************
// VARS
// *****************************************************************************
var (
	ICON_DISK  = "⛁"
	ICON_ARROW = "⯈"
	BUILD_TIME string
	GIT_COMMIT string
)

// *****************************************************************************
// init()
// *****************************************************************************
func init() {
	if runtime.GOOS == "windows" || runtime.GOOS == "android" {
		// These UTF8 icons are not correctly rendered in Windows or Android's Termux,
		// so we convert them to plain vanilla ASCII characters
		ICON_DISK = "#"
		ICON_ARROW = ">"
	}
}
