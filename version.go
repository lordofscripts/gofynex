/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024-2026 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *	PACKAGE VERSION
 *-----------------------------------------------------------------*/
package gofynex

import (
	"fmt"
	"time"

	"github.com/lordofscripts/goapp"
	"github.com/lordofscripts/goapp/app"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	// Change these values accordingly
	NAME           string = "GoFynex"
	DESC           string = "Fyne is crippled, here are the crutches."
	MANUAL_VERSION string = "1.3.2"
	//_RELEASE_CANDIDATE int    = 0
	//_BETA_VERSION int = 0
)

// NOTE: Change these values accordingly
var (
	ModuleVersion app.PackageVersion = app.NewReleaseVersion(NAME, DESC, MANUAL_VERSION)
)

/* ----------------------------------------------------------------
 *					I n i t i a l i z e r
 *-----------------------------------------------------------------*/

func init() {
	goapp.RegisterModule(ModuleVersion)
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// get a corrected copyright
func Copyright() string {
	const CHR_TRIDENT rune = rune(0x1f531) // 🔱
	creator := fmt.Sprintf("Copyright (C)2024-%d Dídimo Grimaldo", time.Now().Year())
	return creator
}

/*
  This commented part is used by the Makefile target that extracts it
  to a temporary file, compiles it and uses it to print out the full
  Module Version that includes attributes like Alpha,Beta,RC which the
  plain old "make version" didn't

//>>>BEGIN Versioner
package main

import (
    "os"
    "fmt"
    "strings"
    "github.com/lordofscripts/gofynex"
)

func main() {
    if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "short") {
        fmt.Println(gofynex.ModuleVersion.Short())
    } else {
        fmt.Println(gofynex.ModuleVersion)
    }
}
//>>>END Versioner
*/
