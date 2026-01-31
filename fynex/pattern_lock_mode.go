/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Pattern Mode enumeration
 *-----------------------------------------------------------------*/
package fynex

import "strings"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	PatternModeNone PatternMode = iota
	PatternMode3x3
	PatternMode4x4
	PatternMode5x5
)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// The size of the square PatternLock grid
type PatternMode uint8

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer by returning "None", "3x3", "4x4" or "5x5"
func (pm PatternMode) String() string {
	var result string
	switch pm {
	case PatternMode3x3:
		result = "3x3"
	case PatternMode4x4:
		result = "4x4"
	case PatternMode5x5:
		result = "5x5"
	case PatternModeNone:
		result = "None"
	default:
		result = ""
	}
	return result
}

// for a grid size (0|3|4|5) convert to enumeration value
func (pm PatternMode) Convert(gridSize int) PatternMode {
	var result PatternMode
	switch gridSize {
	case 0:
		result = PatternModeNone
	case 3:
		result = PatternMode3x3
	case 4:
		result = PatternMode4x4
	case 5:
		result = PatternMode5x5
	default:
		result = PatternModeNone
		println("Unsupported gridSize in Convert", gridSize)
	}
	return result
}

// convert s (0/None/0x0|3x3|4x4|5x5) to enumeration value
func (pm PatternMode) Parse(s string) PatternMode {
	var result PatternMode
	switch strings.ToLower(s) {
	case "0":
		fallthrough
	case "0x0":
		fallthrough
	case "none":
		result = PatternModeNone
	case "3x3":
		result = PatternMode3x3
	case "4x4":
		result = PatternMode4x4
	case "5x5":
		result = PatternMode5x5
	default:
		result = PatternModeNone
		println("Unrecognized pattern mode in Parse", s)
	}
	return result
}

// the array width for the current pattern mode value
func (pm PatternMode) Width() int {
	var result int
	switch pm {
	case PatternMode3x3:
		result = 3
	case PatternMode4x4:
		result = 4
	case PatternMode5x5:
		result = 5
	default:
		result = 0
	}
	return result
}
