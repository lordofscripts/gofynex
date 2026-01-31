/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * The other side of Boolean because sometimes things aren't as clear
 * as Yes/No, True/False, there is Checked/Unchecked/Unset
 *-----------------------------------------------------------------*/
package fynex

import "strings"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	Unset TriState = iota
	Unchecked
	Checked
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// For those who truly know digital systems, there is High/Low/Floating
// or True/False/Unset you get the point, that 3rd value
type TriState uint8

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// implements fmt.Stringer displaying the True/False/Unset value
func (t TriState) String() string {
	var val string
	switch t {
	case Checked:
		val = "True"
	case Unchecked:
		val = "False"
	default:
		val = "Unset"
	}
	return val
}

// Parses the value of type string, rune, bool or *bool
// into a TriState value and returns it. It does NOT modify
// the current variable.
// Example: Unset.Parse("yes") returns Checked.
func (t TriState) Parse(value any) TriState {
	var result TriState = Unset
	if val, isBool := value.(bool); isBool {
		if val {
			result = Checked
		} else {
			result = Unchecked
		}
	} else if val, isBoolPtr := value.(*bool); isBoolPtr {
		if val != nil {
			if *val {
				result = Checked
			} else {
				result = Unchecked
			}
		}
	} else if val, isString := value.(string); isString {
		val = strings.ToLower(strings.TrimSpace(val))
		switch val {
		case "yes":
			fallthrough
		case "true":
			result = Checked
		case "no":
			fallthrough
		case "false":
			result = Unchecked
		case "maybe":
			result = Unset
		}
	} else if val, isRune := value.(rune); isRune {
		switch val {
		case 'H': // High
			fallthrough
		case 'Y': // Yes
			result = Checked
		case 'L': // Low
			fallthrough
		case 'N': // No
			result = Unchecked
		case '-': // Floating
			fallthrough
		case ' ':
			fallthrough
		case 'X':
			fallthrough
		case '?':
			result = Unset
		}
	}

	return result
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/
