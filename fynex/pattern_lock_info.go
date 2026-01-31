/* *****************************************************************
 *              Copyright(C)2026 Lord of Scripts
 *                      All Rights Reserved
 * -----------------------------------------------------------------
 * Holds information about the correct unblock pattern.
 ********************************************************************/
package fynex

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

/* -----------------------------------------------------------------
 *                       G L O B A L S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                     I N T E R F A C E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                  P U B L I C      T Y P E S
 * -----------------------------------------------------------------*/

// Information about the Lock pattern
type PatternInfo struct {
	mode    PatternMode
	minimum uint8
	pattern []int
}

/* -----------------------------------------------------------------
 *                  P R I V A T E    T Y P E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------",
 *                  C O N S T R U C T O R S",
 * -----------------------------------------------------------------*/

// define a new pattern using its internal format where each dot corresponds
// to a 0-based index. The Lock Pattern is flattened as a single row.
func NewPattern(pattern []int, mode PatternMode) (*PatternInfo, error) {
	if mode == PatternModeNone {
		return nil, errors.New("invalid pattern mode None")
	}
	if len(pattern) < mode.Width() {
		return nil, fmt.Errorf("pattern mode %d needs at least %d dots", mode, mode.Width())
	}
	if dups := findAdjacentDuplicates(pattern); len(dups) != 0 {
		return nil, fmt.Errorf("pattern has adjacent duplicates %v", dups)
	}
	if err := validateIndices(pattern, mode); err != nil {
		return nil, err
	}

	return &PatternInfo{
		mode:    mode,
		minimum: uint8(mode.Width()),
		pattern: pattern,
	}, nil
}

// define a new lock pattern using a human-friendly notation where each dot is
// defined by its column A..E and its row 1..5. The letter and number limits
// depend on the selected mode. The row number is 1-based.
func NewPatternFromString(pattern string, mode PatternMode) (*PatternInfo, error) {
	if internalPattern, err := ParseStringPatternFor(pattern, mode); err != nil {
		return nil, err
	} else {
		return NewPattern(internalPattern, mode)
	}
}

/* -----------------------------------------------------------------
 *                       M E T H O D S
 * -----------------------------------------------------------------*/

// the minimum number of pattern dots for the selected mode
func (pi *PatternInfo) Minimum() int {
	return int(pi.minimum)
}

// the current size of the pattern
func (pi *PatternInfo) Length() int {
	return len(pi.pattern)
}

// The width or height size of the matrix for the selected mode
func (pi *PatternInfo) Size() int {
	return pi.mode.Width()
}

// the current pattern mode. It dictates the validation requirements,
// grid notation and grid size.
func (pi *PatternInfo) Mode() PatternMode {
	return pi.mode
}

// returns the current pattern in internal format (dot indices)
func (pi *PatternInfo) Pattern() []int {
	return pi.pattern
}

// implements fmt.Stringer displaying the pattern in friendly form
func (pi *PatternInfo) String() string {
	return PatternInfoString(pi.mode, pi.pattern)
}

/* -----------------------------------------------------------------
 *                  P R I V A T E    M E T H O D S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                       F U N C T I O N S
 * -----------------------------------------------------------------*/

// given a pattern mode like 3x3 or 4x4 and a sequence of dots,
// convert the sequence from indices to friendly cartesian coordinates.
func PatternInfoString(mode PatternMode, sequence []int) string {
	if mode == PatternModeNone {
		return ""
	}
	var mod = mode.Width()

	coordSlice := make([]string, len(sequence))
	for i, point := range sequence {
		column := point % mod
		row := int(point/mod) + 1 // 1-based
		coord := fmt.Sprintf("%c%d", rune(65+column), row)
		coordSlice[i] = coord
	}

	return strings.Join(coordSlice, "-")
}

// takes a pattern in human-readable form and validates it for the selected
// pattern mode. It ensures the dots are correct for the matrix size. It
// also validates for adjacent duplicates.
func ParseStringPatternFor(pattern string, mode PatternMode) ([]int, error) {
	if mode == PatternModeNone {
		return []int{}, errors.New("cannot parse for pattern mode None")
	}
	// remove whitespace
	pattern = strings.Trim(pattern, " \t")
	// normalize
	pattern = strings.ToUpper(pattern)
	// split in dots
	const SEP = "-"
	dots := strings.Split(pattern, SEP)
	result := make([]int, 0)
	if len(dots) < mode.Width() {
		return []int{}, errors.New("not enough dots in the pattern")
	}
	// validate
	for i, dot := range dots {
		// The letters and numbers in the Pattern must match the ranges
		// allowed for the maximum defined PatternMode: 3x3 ABC123,
		// 4x4 ABCD1234 and 5x5 ABCDE12345
		const PATTERN_3 = `^[A-C][1-3]$` // int: 0..9 hum: L[1..3]
		const PATTERN_4 = `^[A-D][1-4]$` // int: 0..16 hum: L[1..4]
		const PATTERN_5 = `^[A-E][1-5]$` // int: 0..24 hum: L[1..5]
		var chosenPattern string
		switch mode {
		case PatternMode3x3:
			chosenPattern = PATTERN_3
		case PatternMode4x4:
			chosenPattern = PATTERN_4
		case PatternMode5x5:
			chosenPattern = PATTERN_5
		default:
			return []int{}, errors.New("invalid pattern mode to ParseStringPattern")
		}
		re := regexp.MustCompile(chosenPattern)
		if matched := re.MatchString(dot); matched {
			// Column names are in ASCII so only 1-byte per letter
			column := int(dot[0] - 'A') // 'A' = 65
			row := int(dot[1] - '1')    // '1' = 49
			// the human-friendly string rows are 1-based
			index := row*int(mode.Width()) + column
			// append the pattern dot in internal format (a 0-based slice index)
			result = append(result, index)
		} else {
			return []int{}, fmt.Errorf("pattern dot error at %d='%s' invalid notation", i, dot)
		}
	}

	// so far, so good, ensure the same dot is not adjacent to itself
	if len(findAdjacentDuplicates(result)) != 0 {
		return []int{}, errors.New("adjancent duplicates found")
	}

	return result, nil
}

// Parses and identifies a pattern. It ensures the dots are in the correct
// human format, that the same dot does not repeat after itself and that it
// has the minimum dots allowed for a mode. It begins with 3x3, then 4x4 and
// if that also fails 5x5. Note that 3x3 patterns are perfectly valid for
// 4x4, and 4x4 valid for 5x5; therefore the PatternMode may not be properly
// identified in those circumstances. Otherwise parse specifically using the
// ParseStringPatternFor() function.
func ParseStringPattern(pattern string) ([]int, PatternMode, error) {
	var err error = nil
	var mode PatternMode = PatternMode3x3
	var internalPattern []int
	internalPattern, err = ParseStringPatternFor(pattern, mode)
	if err != nil {
		// Now let's try 4x4
		mode = PatternMode4x4
		internalPattern, err = ParseStringPatternFor(pattern, mode)
		if err != nil {
			// Now let's try 5x5
			mode = PatternMode5x5
			internalPattern, err = ParseStringPatternFor(pattern, mode)
		}
	}

	return internalPattern, mode, err
}

// finds adjacent duplicates. We can have dots reconnecting at
// some point, but never the same dot in sequence.
func findAdjacentDuplicates[T int | string](slice []T) []T {
	duplicates := make([]T, 0)

	if len(slice) < 2 {
		return duplicates
	}

	for i := 0; i < len(slice)-1; i++ {
		if slice[i] == slice[i+1] {
			duplicates = append(duplicates, slice[i])
		}
	}

	return duplicates
}

// validate the index value on the context of the selected PatternMode
func validateIndices(indices []int, mode PatternMode) error {
	var max int
	switch mode {
	case PatternMode3x3:
		max = 8
	case PatternMode4x4:
		max = 15
	case PatternMode5x5:
		max = 24
	default:
		return errors.New("cannot validate indices for None")
	}

	for i, index := range indices {
		// The internal index can never be negative
		if index < 0 {
			return fmt.Errorf("negative dot index %d at [%d]", index, i)
		}
		if index > max {
			return fmt.Errorf("dot index %d out of range (0..%d) at [%d]", index, max, i)
		}
	}

	return nil
}
