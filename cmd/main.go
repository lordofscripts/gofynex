/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/gofynex"
	"github.com/lordofscripts/gofynex/fynex"
	"github.com/lordofscripts/gofynex/fynex/dlg"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

const (
	APP_NAME      = "Pattern Lock Demo"
	APP_ID        = "com.lordofscripts.patternlock"
	APP_DEVELOPER = "Lord of Scripts™"
	APP_TITLE     = "Software Architect/Writer"
)

var (
	/* 	   3x3
		A	B	C
	  +------------
	1 |	0	1	2
	2 |	3	4	5
	3 |	6	7	8
	Default 3x3 Pattern: A1-B1-C1-C2-C3
	*/
	INTERNAL_PATTERN_3x3 = []int{0, 1, 2, 5, 8}
	PATTERN_3x3          *fynex.PatternInfo

	/* 		 4x4
		A	B	C	D
	  +-----------------
	1 |	0	1	2	3
	2 |	4	5	6	7
	3 |	8	9	10	11
	4 |	12	13	14	15
	Default 4x4 Pattern: A1-B1-C1-D1-D2-D3-D4
	*/
	INTERNAL_PATTERN_4x4 = []int{0, 1, 2, 3, 7, 11, 15}
	PATTERN_4x4          *fynex.PatternInfo

	/* 		   5x5
		A	B	C	D	E
	  +--------------------
	1 |	0	1	2	3	4
	2 |	5	6	7	8	9
	3 | 10	11	12	13	14
	4 | 15	16	17	18	19
	5 |	20	21	22	23	24
	Default 5x5 Pattern: 0 1 2 3 4 9 14 19 24
	*/
	INTERNAL_FRIENDLY_PATTERN_5x5 = "A1-B1-C1-D1-E1-E2-E3-E4-E5"
	PATTERN_5x5                   *fynex.PatternInfo

	// CLI Flags
	flgLog          bool
	flgNoBackground bool
	flgHelp         bool
)

const (
	NoBackground BackgroundType = iota
	GradientBackground
	PictureBackground
)

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

type BackgroundType uint8

// manages the entire GUI
type patternApp struct {
	app fyne.App
	win fyne.Window

	current *activeSettings
	ui      *uiStuff
}

// currently active settings
type activeSettings struct {
	useBackground bool
	pattern       *fynex.PatternInfo
	mode          fynex.PatternMode
}

// UI widgets & containers
type uiStuff struct {
	bgType BackgroundType
	// contains all widgets
	mainContainer *fyne.Container
	// The visible instance of our custom PatternLock widget
	lockWidget *fynex.PatternLock
	// IMPORTANT: This is the container where we replace the PatternLock
	// widgets in their different modes. See replaceWidget()
	lockContainer *fyne.Container
	// The radio group allows us to select a Grid mode
	radioButtons *widget.RadioGroup
	// label where we display the current VALID pattern
	patternLabel *widget.Label
	// The Status LED
	statusLED *fynex.LedLabel
}

/* ----------------------------------------------------------------
 *                     I N I T I A L I Z E R
 *-----------------------------------------------------------------*/
func init() {
	var err error

	// Define 3x3 and 4x4 patterns using internal notation which is a slice
	// of 0-based sequential indices
	PATTERN_3x3, err = fynex.NewPattern(INTERNAL_PATTERN_3x3, fynex.PatternMode3x3)
	if err != nil {
		Die(2, err.Error())
	}

	PATTERN_4x4, err = fynex.NewPattern(INTERNAL_PATTERN_4x4, fynex.PatternMode4x4)
	if err != nil {
		Die(2, err.Error())
	}

	// Now let's define a 5x5 mode but using a human-friendly pattern notation
	PATTERN_5x5, err = fynex.NewPatternFromString(INTERNAL_FRIENDLY_PATTERN_5x5, fynex.PatternMode5x5)
	if err != nil {
		Die(2, err.Error())
	}

}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (ctor) a new instance of the Fyne GUI pattern demo application
func newPatternApp(pinfo *fynex.PatternInfo) *patternApp {
	return &patternApp{
		current: &activeSettings{
			useBackground: false,
			pattern:       pinfo,
			mode:          pinfo.Mode(),
		},
		ui: &uiStuff{
			bgType:        PictureBackground,
			mainContainer: nil,
			lockWidget:    nil,
			lockContainer: nil,
			radioButtons:  nil,
			patternLabel:  widget.NewLabelWithStyle(pinfo.String(), fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
		},
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

// whether to include the default widget background.
func (a *patternApp) WithBackground(enable bool) {
	a.current.useBackground = enable
}

func (a *patternApp) Define() *patternApp {
	// (Window)
	// +--------------------------------------------------------+
	// |				Pattern Lock Demo                       |
	// +--------------------------------------------------------+
	// The standard Fyne application
	a.app = app.NewWithID(APP_ID)
	a.win = a.app.NewWindow(APP_NAME)
	a.win.SetMaster()

	// (Custom Widget)
	// + ----------------
	// | Draw pattern
	// |     O   O   O
	// |     O   O   O
	// |     O   O   O
	// + ----------------
	// And this is our initial setup with a known pattern
	//lock = NewPatternLock(currentPattern.Size(), onCompleted)
	a.ui.lockWidget = fynex.NewPatternLockWith(a.current.pattern, a.callbackOnValidated)
	if a.current.useBackground {
		a.ui.lockWidget.SetBackground(fynex.DefaultBackground)
	}
	// centered, but we need to address it for replacement later on. This
	// container gives the PatterhLock widget the necessary space required by Layout()
	a.ui.lockContainer = container.NewCenter(a.ui.lockWidget)

	// (RadioGroup widget)
	// + -----------------------------------------
	// | (*) 3x3  ( ) 4x4  ( ) 5x5   ( ) Define
	// + -----------------------------------------
	const LBL_3x3 = "3x3" // PatternMode3x3.String()
	const LBL_4x4 = "4x4" // PatternMode4x4.String()
	const LBL_5x5 = "5x5" // PatternMode5x5.String()
	const LBL_DEF = "Define"
	a.ui.radioButtons = widget.NewRadioGroup(
		[]string{LBL_3x3, LBL_4x4, LBL_5x5, LBL_DEF},
		a.callbackOnRadioChanged)
	a.ui.radioButtons.Horizontal = true

	bgRadioButtons := widget.NewRadioGroup([]string{"No Background", "Gradient", "Picture"}, func(s string) {
		switch s {
		case "No Background":
			a.ui.replaceBackground(NoBackground)
		case "Gradient":
			a.ui.replaceBackground(GradientBackground)
		case "Picture":
			a.ui.replaceBackground(PictureBackground)
		}
	})
	bgRadioButtons.Horizontal = true
	bgRadioButtons.SetSelected("Picture")

	// (Label widget)
	// + -----------------------------------------
	// | Current Pattern:
	// |  A1-B1-C1-D1-E1-E2
	// + -----------------------------------------
	hintLabel := widget.NewLabelWithStyle("Current pattern:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	slider := fynex.NewScrollableSlider(65, 90) // 'A'..'Z'
	slider.OnConvert = func(f float64) string {
		ascii := int(f)
		letter := rune(ascii)
		//fmt.Printf("Slider %d %c\n", ascii, letter)
		return string(letter)
	}
	slider.SetRightVisible(true) // after OnConvert is defined

	// A sample status LED that displays access status
	a.ui.statusLED = fynex.NewLedLabel("Status")

	// define the window
	a.ui.mainContainer = container.NewVBox(
		a.ui.lockContainer,
		container.NewCenter(a.ui.radioButtons),
		container.NewCenter(bgRadioButtons),
		hintLabel,
		a.ui.patternLabel,
		slider,
		a.ui.statusLED,
	)
	a.win.SetContent(a.ui.mainContainer)

	a.win.Resize(fyne.NewSize(400, 500))

	return a
}

// Set other default values for widgets and bind data (if any)
func (a *patternApp) Setup() *patternApp {
	a.ui.radioButtons.SetSelected(a.current.mode.String())
	a.ui.statusLED.SetState(fynex.Unset)
	return a
}

func (a *patternApp) GetApp() fyne.App {
	return a.app
}

func (a *patternApp) GetWindow() fyne.Window {
	return a.win
}

func (a *patternApp) Run() {
	a.win.ShowAndRun()
}

func (a *patternApp) Quit() {
	a.app.Quit()
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// refresh the main container of the main window
func (gui *uiStuff) Refresh() {
	gui.mainContainer.Refresh()
}

func (gui *uiStuff) replaceBackground(bgtype BackgroundType) {
	switch bgtype {
	case NoBackground:
		gui.lockWidget.SetBackground(nil)
	case GradientBackground:
		gui.lockWidget.SetBackground(fynex.GradientBackground)
	case PictureBackground:
		gui.lockWidget.SetBackground(fynex.DefaultBackground)
	}
	gui.bgType = bgtype
}

func (gui *uiStuff) replaceBackgroundOn(bgtype BackgroundType, lockW *fynex.PatternLock) {
	switch bgtype {
	case NoBackground:
		lockW.SetBackground(nil)
	case GradientBackground:
		lockW.SetBackground(fynex.GradientBackground)
	case PictureBackground:
		lockW.SetBackground(fynex.DefaultBackground)
	}
	gui.bgType = bgtype
}

// Swaps the current PatternLock widget with a new instance. This
// usually happens when we switch (grid) modes or switch between
// pattern design and pattern recognition.
func (gui *uiStuff) replaceWidget(with *fynex.PatternLock, hasBackground bool) {
	log.Print("replaceWidget")
	// the PatternLock container only contains that widget
	//parent.Objects[0] = with
	gui.lockContainer.Objects = []fyne.CanvasObject{with}
	// Or
	//parent.Remove(lock)
	//parent.Add(with)
	gui.lockWidget = with
	if hasBackground {
		//gui.lockWidget.SetBackground(fynex.DefaultBackground)
		gui.replaceBackground(gui.bgType)
	}
	// refresh after the swap
	gui.lockContainer.Refresh()
}

func (a *patternApp) setupWidget(pinfo *fynex.PatternInfo, doSwap, altColor bool) {
	GREEN := color.NRGBA{R: 0, G: 0x9d, B: 0, A: 255} // #009D00
	a.current.pattern = pinfo
	a.current.mode = pinfo.Mode()
	a.ui.patternLabel.SetText(pinfo.String())
	newLock := fynex.NewPatternLockWith(pinfo, a.callbackOnValidated)
	if altColor {
		newLock.SetSelectedColor(GREEN)
	}
	if doSwap {
		a.ui.replaceWidget(newLock, a.current.useBackground)
	} else {
		a.ui.replaceBackgroundOn(a.ui.bgType, newLock)
	}
}

// [RadioGroup] (Callback:onChanged)
// The selected radio button has changed.
func (a *patternApp) callbackOnRadioChanged(s string) {
	log.Printf("Selected %s", s)
	switch s {
	case fynex.PatternMode3x3.String():
		a.setupWidget(PATTERN_3x3, true, true)

	case fynex.PatternMode4x4.String():
		a.setupWidget(PATTERN_4x4, true, true)

	case fynex.PatternMode5x5.String():
		a.setupWidget(PATTERN_5x5, true, true)

	case "Define":
		a.ui.patternLabel.SetText("None")
		// configure widget for defining a new pattern. Internally it will use a
		// different color for the drawn lines.
		newLock := fynex.NewPatternLock(a.current.mode.Width(), a.callbackOnDefined)
		newLock.EnterDesignState() // reconfigure for design-mode before we swap
		a.ui.replaceWidget(newLock, a.current.useBackground)
	}
	a.ui.Refresh()
}

// [PatternLock] (Callback:OnComplete)
// Callback when Pattern needs to be validated. This is only used if you
// want to do the validation yourself and set status on the PatternLock widget,
// or if there are other actions you wish to do.
func (a *patternApp) callbackOnCompleted(sequence []int) {
	fmt.Printf("User drew: %v\n", sequence)
	log.Printf("onCompleted DRAW %s", fynex.PatternInfoString(a.current.mode, sequence))
	log.Printf("onCompleted REQD %s", a.current.pattern.String())
	if reflect.DeepEqual(sequence, a.current.pattern.Pattern()) {
		fmt.Println("Access granted")
	} else {
		fmt.Println("Access denied")
	}
}

// [PatternLock] (Callback:OnValidated)
// For normal applications if you need to do something else when
// the pattern is recognized as valid (granted) or invalid (denied)
func (a *patternApp) callbackOnValidated(isValid bool) {
	log.Printf("OnValidated (user) is-valid: %t (%s)", isValid, a.ui.statusLED.State())
	if isValid {
		a.ui.statusLED.SetState(fynex.Checked) // or a.ui.statusLED.Green()
	} else {
		a.ui.statusLED.SetState(fynex.Unchecked) // or a.ui.statusLED.Red()
	}
}

// [PatternLock] (Callback:onDefined)
// Callback when Pattern has been defined
func (a *patternApp) callbackOnDefined(sequence []int) {
	log.Print("onDefined (user) entered")
	fmt.Printf("User drew: %v\n", sequence)

	var newPattern *fynex.PatternInfo = nil
	var err error

	// This validates both the mode as well as the minimum pattern length
	if newPattern, err = fynex.NewPattern(sequence, a.current.mode); err != nil {
		fmt.Println("ERROR", err)
		log.Print("ERROR", err)
	} else {
		/*
			// assign newly defined pattern
			a.current.pattern = newPattern                         //@todo WHY IS IT NOT GETTING ASSIGNED?
			a.ui.radioButtons.SetSelected(a.current.mode.String()) // reconfigure & swap
			a.ui.lockWidget.SetValidPattern(newPattern)
			// optional, just for demonstration
			a.ui.lockWidget.OnComplete = a.callbackOnCompleted
			// display human-readable pattern
			a.ui.patternLabel.SetText(a.current.pattern.String())
		*/
		a.ui.radioButtons.OnChanged = nil // prevent next tone from triggering
		a.setupWidget(newPattern, true, true)
		a.ui.radioButtons.SetSelected(a.current.mode.String())
		a.ui.radioButtons.OnChanged = a.callbackOnRadioChanged
	}
	log.Print("onDefined (user) left")
}

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

func Help() {
	flag.PrintDefaults()
	gofynex.Copyright(gofynex.CO1, true)
	os.Exit(0)
}

// parse CLI flags
func ParseFlags() {
	flag.BoolVar(&flgLog, "log", false, "View log output")
	flag.BoolVar(&flgNoBackground, "nobg", false, "Don't use widget background")
	flag.BoolVar(&flgHelp, "help", false, "Help!")
	flag.Parse()

	if flgHelp {
		Help()
	}

	if !flgLog {
		log.SetOutput(io.Discard)
	}

}

func FixMetadata(meta fyne.AppMetadata) fyne.AppMetadata {
	if len(meta.Name) == 0 {
		meta.Name = APP_NAME
	}
	if len(meta.ID) == 0 {
		meta.ID = APP_ID
	}
	meta.Version = gofynex.Version.Short()
	meta.Custom["url"] = "https://github.com/lordofscripts"
	meta.Custom["url.text"] = "GitHub"
	return meta
}

func AboutDemo(win fyne.Window, meta fyne.AppMetadata) {
	const About = `
Copyright (c)2026 Lord of Scripts
This is a custom for Pattern Lock.
And there are more useful widgets for
you to enjoy in the limited Fyne world.
	`

	about := dlg.NewAboutBox(win, nil, meta).
		WithPersonModel(fynex.NewPersonWithImage(APP_DEVELOPER, APP_TITLE, fynex.DeveloperIcon)).
		WithText(About, false, true)
	about.ShowDialog()
}

// die with style
func Die(exitCode int, message string) {
	fmt.Println(message)
	os.Exit(exitCode)
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

// Demonstration of custom PatternLock widget
func main() {
	gofynex.Copyright(gofynex.CO1, true)

	ParseFlags()

	app := newPatternApp(PATTERN_3x3)
	app.WithBackground(!flgNoBackground)
	app.Define().Setup()

	// All log output will be redirected to this non-modal window.
	// When run on Windows you won't see the console output anywhere.
	// because of the "-H windowsgui" linker flags.
	var logWindow fyne.Window = nil
	if flgLog {
		logWindow = fynex.NewLogWindow(app.GetApp(), 500, 600)
		win := app.GetWindow()
		win.SetCloseIntercept(func() {
			logWindow.Close()
			win.Close()
		})
		logWindow.Show()
	}

	AboutDemo(app.GetWindow(), FixMetadata(app.GetApp().Metadata()))
	app.Run()

	gofynex.BuyMeCoffee()
}
