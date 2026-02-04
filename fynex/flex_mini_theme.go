/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A mini-theme where you can inject custom colors and sizes when
 * the default theme is too restrictive.
 *-----------------------------------------------------------------*/
package fynex

import (
	"image/color"

	"fyne.io/fyne/v2"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

// A mini-theme that only overrides text names and color names
// but the overrides are specified using Include() and IncludeColor()
type FlexMiniTheme struct {
	fyne.Theme
	sizeOverrides  map[fyne.ThemeSizeName]float32
	colorOverrides map[fyne.ThemeColorName]color.Color
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewFlexMiniTheme(defaultTheme fyne.Theme) *FlexMiniTheme {
	return &FlexMiniTheme{
		Theme:          defaultTheme,
		sizeOverrides:  make(map[fyne.ThemeSizeName]float32),
		colorOverrides: make(map[fyne.ThemeColorName]color.Color),
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// fluent API method to include one or more ThemeSizeName in the
// theme's text size override.
func (m *FlexMiniTheme) Include(name fyne.ThemeSizeName, size float32) *FlexMiniTheme {
	m.sizeOverrides[name] = size
	return m
}

func (m *FlexMiniTheme) IncludeColor(name fyne.ThemeColorName, col color.Color) *FlexMiniTheme {
	m.colorOverrides[name] = col
	return m
}

// Override for sizes
func (m *FlexMiniTheme) Size(name fyne.ThemeSizeName) float32 {
	if size, exists := m.sizeOverrides[name]; exists {
		return size
	}

	return m.Theme.Size(name) // non-override for all others
}

func (m *FlexMiniTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if colour, exists := m.colorOverrides[name]; exists {
		return colour
	}

	return m.Theme.Color(name, variant) // non-override for all others
}
