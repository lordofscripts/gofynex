/* *****************************************************************
 *              Copyright(C)2026 Lord of Scripts
 *                      All Rights Reserved
 * -----------------------------------------------------------------
 * Modal About Custom Dialog.
 ********************************************************************/
package fynex

import (
	"bytes"
	"errors"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/* -----------------------------------------------------------------
 *                     I N T E R F A C E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                  P U B L I C      T Y P E S
 * -----------------------------------------------------------------*/

type AboutBox struct {
	parent  fyne.Window
	logoRes fyne.Resource
	meta    fyne.AppMetadata
	text    string
}

/* -----------------------------------------------------------------
 *                  P R I V A T E    T Y P E S
 * -----------------------------------------------------------------*/

type smallTextTheme struct {
	fyne.Theme
}

/* -----------------------------------------------------------------
 *                  C O N S T R U C T O R S
 * -----------------------------------------------------------------*/
func NewAboutBox(myWin fyne.Window, logoRes fyne.Resource, meta fyne.AppMetadata, text string) *AboutBox {
	return &AboutBox{
		parent:  myWin,
		logoRes: logoRes,
		meta:    meta,
		text:    text,
	}
}

/* -----------------------------------------------------------------
 *                       M E T H O D S
 * -----------------------------------------------------------------*/

// Create an About window
func (a *AboutBox) ShowDialog() {
	if a.logoRes == nil {
		a.logoRes = DefaultBackground
	}

	//w, h, _ := GetImageDimensions(logoRes)
	logo := canvas.NewImageFromResource(a.logoRes)
	logo.FillMode = canvas.ImageFillContain
	//logo.SetMinSize(fyne.NewSize(float32(w/2), float32(h/2)))
	logo.SetMinSize(fyne.NewSize(150, 150))
	logo.Translucency = 0.2
	logo.Refresh()

	var hyperlink *widget.Hyperlink = nil
	if url, exists := a.meta.Custom["url"]; exists {
		hyperlink = widget.NewHyperlink(url, nil)
		hyperlink.SetURLFromString(url)
	}

	// Rich Text Segment with smaller Monospace font
	segment := &widget.TextSegment{
		Text:  a.text,
		Style: widget.RichTextStyleCodeBlock, // defaults to Monospace
	}
	segment.Style.Alignment = fyne.TextAlignCenter

	rt := widget.NewRichText(segment)
	smallMonospace := container.NewThemeOverride(rt, &smallTextTheme{Theme: theme.DefaultTheme()})

	container := container.NewVBox(
		container.NewCenter(logo),
		container.NewCenter(widget.NewLabel(a.meta.Name+" v"+a.meta.Version)),
	)
	if hyperlink != nil {
		container.Add(hyperlink)
	}
	container.Add(smallMonospace)
	container.Add(layout.NewSpacer())

	dialog.ShowCustom("About", "Close", container, a.parent)
}

// Override for sizes
func (m *smallTextTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 10 // force text to a smaller size
	}
	return m.Theme.Size(name) // non-override for all others
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// get the image dimensions from an (embedded) resource. NewImageFromResource()
// is lazy, it loads the resource but does not initialize the Image property
// until it is rendered for the first time.
func GetImageDimensions(res fyne.Resource) (width, height int, err error) {
	width = 0
	height = 0
	if res != nil {
		img, _, err := image.Decode(bytes.NewReader(res.Content()))
		if err == nil {
			bounds := img.Bounds()
			width = bounds.Dx()
			height = bounds.Dy()
		}
	} else {
		err = errors.New("no resource given")
	}

	return
}

/* -----------------------------------------------------------------
 *                       M A I N  |  D E M O
 * -----------------------------------------------------------------*/
