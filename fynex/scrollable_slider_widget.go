/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Slider widget that subsanates Fynes shortcomings such as
 * not having tooltips. Almost always the GUI designer needs to display
 * the current value of the slider. This custom slider displays a small
 * label on the left with that value (template defined by designer).
 * It also has the capability of showing a normally hidden label on
 * the right side.
 *-----------------------------------------------------------------*/
package fynex

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	sliderVALUE_TEMPLATE = "%3d"
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ fyne.Scrollable = (*ScrollableSlider)(nil)

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

type ScrollableSlider struct {
	widget.BaseWidget
	slider         *widget.Slider
	leftLabel      *canvas.Text
	leftTemplate   string
	rightLabel     *canvas.Text
	container      *fyne.Container
	OnConvert      func(float64) string
	OnValueChanged func(float64)
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewScrollableSlider(min, max float64) *ScrollableSlider {
	s := &ScrollableSlider{
		leftTemplate: sliderVALUE_TEMPLATE,
		slider:       widget.NewSlider(min, max),
		leftLabel:    canvas.NewText("", color.White),
		rightLabel:   canvas.NewText("    ", color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
	}
	s.ExtendBaseWidget(s)

	s.leftLabel.Alignment = fyne.TextAlignTrailing
	s.rightLabel.Hide()

	s.slider.OnChanged = func(f float64) {
		s.updateLeftLabel()
		if s.OnConvert != nil {
			s.rightLabel.Text = s.OnConvert(f)
		}
		if s.OnValueChanged != nil {
			s.OnValueChanged(f)
		}
	}

	s.slider.SetValue(min) // ensure it works for non-zero Minimum and sync
	// Calculate fixed width for left label based on Max values
	s.recalculateSpace()
	s.updateLeftLabel()
	return s
}

func NewScrollableSliderWithData(min, max float64, data binding.Float) *ScrollableSlider {
	ss := NewScrollableSlider(min, max)
	ss.slider.Bind(data)
	return ss
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// return the actual minimum size, otherwise there might be no surface
// to receive the Scrolled event.
func (s *ScrollableSlider) MinSize() fyne.Size {
	return s.container.MinSize()
}

// Let the BaseWidget (of size 0x0) propagate the Resize call to
// the internal container. Otherwise it won't receive Scroll hits.
func (s *ScrollableSlider) Resize(size fyne.Size) {
	s.BaseWidget.Resize(size)
	s.container.Resize(size)
}

// Scrolled implements fyne.Scrollable
func (s *ScrollableSlider) Scrolled(ev *fyne.ScrollEvent) {
	// @note FWD scroll gives Dx:0 Dy:25 and BACK Dx:0 Dy:-25
	// let's assume that is for all platforms
	const DIVIDER = 25
	dy := int(float64(ev.Scrolled.DY)) / DIVIDER
	newValue := s.slider.Value + float64(dy)
	if newValue > s.slider.Max {
		newValue = s.slider.Max
	} else if newValue < s.slider.Min {
		newValue = s.slider.Min
	}

	// this will trigger s.slider.OnChange which will in turn
	// trigger (if given) ScrollableSlider.OnChanged.
	s.slider.SetValue(newValue)
}

// The the slider's text format template for displaying the slider value.
func (s *ScrollableSlider) SetValueTemplate(format string) {
	s.leftTemplate = format
	s.updateLeftLabel()
}

// make the (optional) right label visible
func (s *ScrollableSlider) SetRightVisible(visible bool) {
	if visible {
		if s.OnConvert != nil {
			s.rightLabel.Text = s.OnConvert(s.slider.Value)
		}
		s.rightLabel.Show()
	} else {
		s.rightLabel.Hide()
	}
	s.container.Refresh()
}

// set the text on the right label
func (s *ScrollableSlider) SetRightText(str string) {
	s.rightLabel.Text = str
	s.rightLabel.Refresh()
}

func (s *ScrollableSlider) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// rebuild the container based on a (new or initial) text formatting
// left label template
func (s *ScrollableSlider) recalculateSpace() {
	// Calculate fixed width for left label based on Max value
	maxText := s.leftLabelString()
	var textSize fyne.Size
	var width float32 = 40 // a default just in case...
	if fyne.CurrentApp() != nil {
		textSize, _ = fyne.CurrentApp().Driver().RenderedTextSize(maxText, s.leftLabel.TextSize, s.leftLabel.TextStyle, nil)
		width = textSize.Width + 10
	} else {
		// fallback in case CurrentApp isn't ready
		textSize = fyne.NewSize(40, 20)
	}

	// Layout: Left Label (fixed width) | Slider | Right Label
	leftBox := container.NewStack(canvas.NewRectangle(color.Transparent), s.leftLabel)
	leftBox.Resize(fyne.NewSize(width, s.leftLabel.TextSize))

	s.container = container.NewBorder(nil, nil, leftBox, s.rightLabel, s.slider)
}

// returns the formatted text for the left label based on the
// currently selected format template (float or int)
func (s *ScrollableSlider) leftLabelString() string {
	text := ""
	if strings.Contains(s.leftTemplate, "f") {
		text = fmt.Sprintf(s.leftTemplate, s.slider.Value)
	} else {
		text = fmt.Sprintf(s.leftTemplate, int(s.slider.Value))
	}
	return text
}

// updates the left label with the current slider value in
// the desired format template (float or integer)
func (s *ScrollableSlider) updateLeftLabel() {
	s.leftLabel.Text = s.leftLabelString()
	s.leftLabel.Refresh()
}

/* ----------------------------------------------------------------
 *                          T E S T S
 *-----------------------------------------------------------------*/

// Main for demonstration
/*
func main() {
	myApp := app.New()
	w := myApp.NewWindow("Custom Slider")

	slider := NewCustomSlider(0, 100)

	toggle := widget.NewCheck("Show Right Label", func(b bool) {
		slider.SetRightVisible(b)
	})

	w.SetContent(container.NewVBox(
		slider,
		toggle,
	))
	w.Resize(fyne.NewSize(400, 100))
	w.ShowAndRun()
}
*/
