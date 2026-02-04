/* *****************************************************************
 *              Copyright(C)2026 Lord of Scripts
 *                      All Rights Reserved
 * -----------------------------------------------------------------
 * Modal About Custom Dialog.
 ********************************************************************/
package dlg

import (
	"bytes"
	"errors"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/gofynex/fynex"
)

/* -----------------------------------------------------------------
 *                     I N T E R F A C E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                  P U B L I C      T Y P E S
 * -----------------------------------------------------------------*/

type AboutBox struct {
	parent     fyne.Window
	logoRes    fyne.Resource
	meta       fyne.AppMetadata
	text       string
	developer  *fynex.PersonModel
	isMarkdown bool
	centered   bool
	container  *fyne.Container
}

/* -----------------------------------------------------------------
 *                  P R I V A T E    T Y P E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                  C O N S T R U C T O R S
 * -----------------------------------------------------------------*/

// In the metadata Custom map you can set the "url" and "url.text" values.
// If text is empty the R/O text entry is not rendered. If WithPerson*()
// is called, we also render the developer information custom widget.
func NewAboutBox(myWin fyne.Window, logoRes fyne.Resource, meta fyne.AppMetadata) *AboutBox {
	return &AboutBox{
		parent:     myWin,
		logoRes:    logoRes,
		meta:       meta,
		text:       "",
		developer:  nil,
		isMarkdown: false,
		centered:   false,
		container:  nil,
	}
}

/* -----------------------------------------------------------------
 *                       M E T H O D S
 * -----------------------------------------------------------------*/

// Call this method for using a Person box instead of Entry
func (a *AboutBox) WithPerson(name, title string, pic fyne.Resource) {
	personPtr := &fynex.PersonModel{
		Name:    name,
		Title:   title,
		Picture: pic,
	}
	a.developer = personPtr
}

func (a *AboutBox) WithText(info string, isMarkdown, centered bool) *AboutBox {
	a.text = info
	a.isMarkdown = isMarkdown
	a.centered = centered
	return a
}

func (a *AboutBox) WithPersonModel(person *fynex.PersonModel) *AboutBox {
	a.developer = person
	return a
}

// Create an About window
func (a *AboutBox) ShowDialog() {
	if a.container == nil {
		a.container = a.buildUI()
	}

	dialog.ShowCustom("About", "Close", a.container, a.parent)
}

func (a *AboutBox) buildUI() *fyne.Container {
	if a.logoRes == nil {
		a.logoRes = fynex.DefaultBackground
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
		urlText := url
		if uTxt, exists := a.meta.Custom["url.text"]; exists {
			urlText = uTxt
		}
		hyperlink = widget.NewHyperlink(urlText, nil)
		hyperlink.SetURLFromString(url)
	}

	var smallMonospace *container.ThemeOverride = nil
	if len(a.text) != 0 {
		var rt *widget.RichText
		if a.isMarkdown {
			rt = widget.NewRichTextFromMarkdown(a.text)
		} else {
			// Rich Text Segment with smaller Monospace font
			segment := &widget.TextSegment{
				Text:  a.text,
				Style: widget.RichTextStyleCodeBlock, // defaults to Monospace
			}
			if a.centered {
				segment.Style.Alignment = fyne.TextAlignCenter
			}
			rt = widget.NewRichText(segment)
		}
		rt.Wrapping = fyne.TextWrapWord
		// apply theme override to the PARSED rich text
		miniTheme := fynex.NewFlexMiniTheme(theme.DefaultTheme()).Include(theme.SizeNameText, 10)
		smallMonospace = container.NewThemeOverride(rt, miniTheme)
	}

	var contenedor *fyne.Container
	if hyperlink != nil {
		contenedor = container.NewVBox(
			container.NewCenter(logo),
			container.NewCenter(widget.NewLabel(a.meta.Name+" v"+a.meta.Version)),
			container.NewCenter(hyperlink),
		)
	} else {
		contenedor = container.NewVBox(
			container.NewCenter(logo),
			container.NewCenter(widget.NewLabel(a.meta.Name+" v"+a.meta.Version)),
		)
	}

	// Add the read-only Text Entry if there was "About" text (description)
	if smallMonospace != nil {
		contenedor.Add(smallMonospace)
	}
	// Add developer information if given
	if a.developer != nil {
		devWidget := fynex.NewPersonWidgetWithModel(*a.developer)
		devWidget.NameColor(color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff})
		contenedor.Add(devWidget)
	}
	contenedor.Add(layout.NewSpacer())

	return contenedor
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
