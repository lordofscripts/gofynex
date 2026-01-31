/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   photoQ
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Fyne widget that displays a tri-state LED that can be
 * green (checked), red (unchecked) or orange (unset/undefined)
 *-----------------------------------------------------------------*/
package fynex

import (
	_ "embed" // required for go:embed
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"

	//"fyne.io/fyne/v2/resource"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
// NOTE: Embedded LED images in embed_custom_led.go
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/**
 * IconLabel is a custom widget that displays an icon and next to it a
 * label. It is a great and flexible replacement for cases where the
 * information is more than binary. Also a great replacement for instances
 * where you need a read-only CheckBox without the eye sore Fyne constrain
 * that would force the developer to use a hard-to-see disabled CheckBox.
 *   You can set your own icon. Or use a predefined set for TriState values
 * (unset, unchecked, checked) or just plain boolean (checked, unchecked).
 */
type LedLabel struct {
	widget.BaseWidget
	icon   *widget.Icon
	label  *widget.Label
	border *canvas.Rectangle
	state  TriState
	locker sync.Mutex
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) Creates a LED Label with default Yellow LED.
 */
func NewLedLabel(labelText string) *LedLabel {
	label := widget.NewLabel(labelText)

	il := &LedLabel{
		icon:   widget.NewIcon(iconYellowResource),
		label:  label,
		border: nil,
		state:  Unset,
		locker: sync.Mutex{},
	}
	il.ExtendBaseWidget(il)

	return il
}

func NewLedLabelWith(iconResource fyne.Resource, labelText string) *LedLabel {
	var icon *widget.Icon
	if iconResource == nil {
		icon = widget.NewIcon(iconYellowResource)
	} else {
		icon = widget.NewIcon(iconResource)
	}
	label := widget.NewLabel(labelText)

	il := &LedLabel{
		icon:   icon,
		label:  label,
		border: nil,
		state:  Unset,
		locker: sync.Mutex{},
	}
	il.ExtendBaseWidget(il)

	return il
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// UpdateIcon updates the icon of the IconLabel
func (ll *LedLabel) UpdateIcon(iconResource fyne.Resource) {
	ll.icon.SetResource(iconResource)
	ll.Refresh() // Refresh the widget to update the display
}

/**
 * Sets the internal state and updates the LED visual accordingly.
 * UnChecked=RED, Checked=GREEN & UnSet=YELLOW
 */
func (ll *LedLabel) SetState(state TriState) *LedLabel {
	ll.locker.Lock()
	defer ll.locker.Unlock()

	ll.state = state
	ll.refreshState()
	//il.Refresh()
	return ll
}

func (ll *LedLabel) State() TriState {
	return ll.state
}

/**
 * Turn on LED in RED
 */
func (ll *LedLabel) Red() *LedLabel {
	ll.locker.Lock()
	defer ll.locker.Unlock()

	ll.state = Checked
	ll.refreshState()
	return ll
}

/**
 * Turn on LED in GREEN
 */
func (ll *LedLabel) Green() *LedLabel {
	ll.locker.Lock()
	defer ll.locker.Unlock()

	ll.state = Checked
	ll.refreshState()
	return ll
}

/**
 * Turn on LED in YELLOW
 */
func (ll *LedLabel) Yellow() *LedLabel {
	ll.locker.Lock()
	defer ll.locker.Unlock()

	ll.state = Unset
	ll.refreshState()
	return ll
}

// CreateRenderer creates the renderer for the IconLabel
func (ll *LedLabel) CreateRenderer() fyne.WidgetRenderer {
	if ll.border == nil {
		ll.border = canvas.NewRectangle(theme.Color(theme.ColorNameSeparator)) // Create a rectangle for the border
	}

	ll.border.SetMinSize(fyne.NewSize(0, 0)) // Set the height of the border
	hbox := container.NewHBox(ll.icon, ll.label)
	return widget.NewSimpleRenderer(container.NewVBox(hbox, ll.border))
}

// MouseIn is called when the mouse enters the widget
func (ll *LedLabel) MouseIn(event *desktop.MouseEvent) {
	ll.Refresh()
}

// MouseOut is called when the mouse leaves the widget
func (ll *LedLabel) MouseOut() {
	ll.Refresh()
}

// Layout is called to layout the widget
func (ll *LedLabel) Layout(size fyne.Size) {
	ll.ExtendBaseWidget(ll)
	ll.CreateRenderer().Layout(size)
}

// MinSize returns the minimum size of the widget
func (ll *LedLabel) MinSize() fyne.Size {
	return ll.CreateRenderer().MinSize()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Sets Red, Green or Yellow according to the internal (tri)state status
 */
func (ll *LedLabel) refreshState() {
	switch ll.state {
	case Unset:
		ll.icon.SetResource(iconYellowResource)

	case Checked:
		ll.icon.SetResource(iconGreenResource)

	case Unchecked:
		ll.icon.SetResource(iconRedResource)
	}

	ll.icon.Refresh()
	ll.Refresh()
}

/*
func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Custom IconLabel Widget")

	// Create an instance of IconLabel
	iconLabel := NewIconLabel(resource.NewStaticResource("icon1", []byte{ ...icon data... }), "Initial Label")
	iconLabel2 := NewIconLabel(iconUnset.Resource, "Touch me")
	// Update the icon after a delay (for demonstration)
	go func() {
		// Simulate icon update after some time
		// Replace with actual icon resource
		iconLabel.UpdateIcon(resource.NewStaticResource("icon2", []byte{ ...new icon data... }))
	}()

	myWindow.SetContent(container.NewVBox(iconLabel, iconLabel2))
	myWindow.Resize(fyne.NewSize(300, 100))
	myWindow.ShowAndRun()
}
*/
