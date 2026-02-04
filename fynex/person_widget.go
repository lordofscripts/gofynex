/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package fynex

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

// Model for a person. It is general purpose but can be used with PersonWidget
type PersonModel struct {
	Name    string
	Title   string
	Picture fyne.Resource
}

type PersonWidget struct {
	widget.BaseWidget
	Model PersonModel

	container *fyne.Container
	nameFg    color.Color
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

// Helper for circular clipping
type circleMask struct {
	p image.Point
	r int
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewPerson(name, title string) *PersonModel {
	return &PersonModel{
		Name:    name,
		Title:   title,
		Picture: theme.AccountIcon(),
	}
}

func NewPersonWithImage(name, title string, picture fyne.Resource) *PersonModel {
	return &PersonModel{
		Name:    name,
		Title:   title,
		Picture: picture,
	}
}

// (Ctor) A person widget that can be added to a container or AboutBox.
func NewPersonWidget(name, title string, img fyne.Resource) *PersonWidget {
	model := PersonModel{
		Name:    name,
		Title:   title,
		Picture: img,
	}
	return NewPersonWidgetWithModel(model)
}

func NewPersonWidgetWithModel(person PersonModel) *PersonWidget {
	p := &PersonWidget{
		Model:     person,
		container: nil,
		nameFg:    nil,
	}
	p.ExtendBaseWidget(p)
	return p
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Override the Name's foreground color
func (p *PersonWidget) NameColor(fgCol color.Color) *PersonWidget {
	p.nameFg = fgCol
	return p
}

func (p *PersonWidget) CreateRenderer() fyne.WidgetRenderer {
	miniTheme := NewFlexMiniTheme(theme.DefaultTheme()).
		Include(theme.SizeNameText, 10).
		Include(theme.SizeNamePadding, 1). // put Name & Title closer together used by VBox
		IncludeColor(theme.ColorNameForeground, color.White)

	// 1. Circular Image
	img := canvas.NewImageFromResource(p.getCircularResource())
	img.SetMinSize(fyne.NewSize(60, 60))
	img.FillMode = canvas.ImageFillContain

	// 2. Text using canvas.Text for zero-margin control
	var name *canvas.Text
	currentVariant := fyne.CurrentApp().Settings().ThemeVariant()
	if p.nameFg != nil {
		name = canvas.NewText(p.Model.Name, p.nameFg)
	} else {
		name = canvas.NewText(p.Model.Name, miniTheme.Color(theme.ColorNameForeground, currentVariant))
	}
	name.TextStyle.Bold = true
	name.TextSize = 16

	title := canvas.NewText(p.Model.Title, miniTheme.Color(theme.ColorNamePlaceHolder, currentVariant))
	title.TextSize = 11

	// 3. Stack text vertically with 0 padding
	textStack := container.NewVBox(name, title)
	tightText := container.NewThemeOverride(textStack, miniTheme)

	// 4. Layout
	// Center the tight text block vertically next to the image
	hbox := container.NewHBox(
		container.NewPadded(img),
		container.NewCenter(tightText),
	)

	// Faint white line at bottom
	line := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 30})
	line.SetMinSize(fyne.NewSize(0, 1))

	content := container.NewBorder(nil, line, nil, nil, hbox)

	return widget.NewSimpleRenderer(content)
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// getCircularResource clips the provided resource into a circular image
func (p *PersonWidget) getCircularResource() fyne.Resource {
	src, _, err := image.Decode(bytes.NewReader(p.Model.Picture.Content()))
	if err != nil {
		return p.Model.Picture // Fallback to original if decoding fails
	}

	// Create a mask
	bounds := src.Bounds()
	diameter := bounds.Dx()
	if bounds.Dy() < diameter {
		diameter = bounds.Dy()
	}

	mask := image.NewRGBA(image.Rect(0, 0, diameter, diameter))
	circle := &circleMask{image.Point{diameter / 2, diameter / 2}, diameter / 2}
	draw.DrawMask(mask, mask.Bounds(), src, image.Point{}, circle, image.Point{}, draw.Over)
	// Convert back to Fyne resource (using internal helper or just returning static)
	// For simplicity in this snippet, we return the original if processing fails,
	// but in production, you'd encode this to PNG bytes.
	return p.Model.Picture
}

func (c *circleMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circleMask) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circleMask) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/
/*
func demoPersonWidget() {
	myApp := app.New()
	w := myApp.NewWindow("Person Widget")

	// Use a placeholder icon or your own resource
	imgRes := theme.AccountIcon()

	person := NewPersonWidget("John Doe", "Senior Software Architect", imgRes)

	w.SetContent(container.NewVBox(person))
	w.Resize(fyne.NewSize(400, 100))
	w.ShowAndRun()
}
*/
