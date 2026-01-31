/* *****************************************************************
 *              Copyright(C)2026 Lord of Scripts
 *                      All Rights Reserved
 * -----------------------------------------------------------------
 * A Lock Pattern widget similar to those used to unblock the screens
 * of smartphones. This widget supports 3x3, 4x4 and 5x5 grids.
 ********************************************************************/
package fynex

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"reflect"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

/* -----------------------------------------------------------------
 *                       G L O B A L S
 * -----------------------------------------------------------------*/

const MSG_STATUS_UNBLOCK = "Draw the pattern to unblock"
const MSG_STATUS_DEFINE = "Draw NEW Pattern"
const MSG_STATUS_WRONG = "Wrong Pattern. Try again."
const MSG_STATUS_GRANTED = "Access Granted!"

var defaultColor = color.NRGBA{R: 0, G: 200, B: 255, A: 200} // #C8FFC8

/* -----------------------------------------------------------------
 *                     I N T E R F A C E S
 * -----------------------------------------------------------------*/

var _ fyne.Widget = (*PatternLock)(nil)
var _ fyne.Tappable = (*PatternLock)(nil)
var _ fyne.Draggable = (*PatternLock)(nil)

/* -----------------------------------------------------------------
 *                          T Y P E S
 * -----------------------------------------------------------------*/

type PatternLock struct {
	widget.BaseWidget
	GridSize    int
	Sequence    []int
	OnComplete  func([]int)
	OnValidated func(bool)
	Status      string

	descriptor     *PatternInfo
	selectedColor  color.NRGBA
	backgroundRsrc *fyne.StaticResource
	active         bool // mouse is being dragged to draw a pattern
	designing      bool // entered pattern design mode (no validation)
	hover          fyne.Position
	mux            sync.Mutex
}

/* -----------------------------------------------------------------
 *                  P U B L I C      T Y P E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                  P R I V A T E    T Y P E S
 * -----------------------------------------------------------------*/

/* -----------------------------------------------------------------
 *                  C O N S T R U C T O R S
 * -----------------------------------------------------------------*/

// (ctor) a Lock Pattern widget of the selected width in dots (square matrix)
// The onComplete callback gets the final sequence and should take care of
// validating the unlock sequence. This instance has no notion of whether
// the sequence is correct, defined or wrong, the caller handles all that.
// NOTE: This is best used for letting the user a NEW Pattern to replace an
// old one
func NewPatternLock(size int, onComplete func([]int)) *PatternLock {
	p := &PatternLock{
		GridSize:       size,
		selectedColor:  defaultColor,
		OnComplete:     onComplete,
		OnValidated:    nil,
		descriptor:     nil,
		backgroundRsrc: nil,
		active:         false,
		designing:      false,
	}
	p.ExtendBaseWidget(p)
	return p
}

// (ctor) a Lock Pattern widget that is initialized with the selected
// pattern information descriptor. The onValidated callback is called
// at the end of the drawn pattern with an indication whether it matched
// the selected pattern (from the descriptor) or not.
// NOTE: This is typically used for unblocking when the pattern is already
// defined.
func NewPatternLockWith(patternDesc *PatternInfo, onValidated func(bool)) *PatternLock {
	pl := NewPatternLock(patternDesc.Size(), nil)
	pl.descriptor = patternDesc
	pl.GridSize = patternDesc.Size()
	pl.OnValidated = onValidated
	return pl
}

/* -----------------------------------------------------------------
 *                       M E T H O D S
 * -----------------------------------------------------------------*/

// Set the Pattern Lock widget's background image from an (embedded) resource.
func (p *PatternLock) SetBackground(bg *fyne.StaticResource) *PatternLock {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.backgroundRsrc = bg
	p.Refresh()
	return p
}

func (p *PatternLock) SetStatus(msg string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.Status = msg
	p.Refresh()
}

// Set the line drawing color for selected dots in the pattern. The
// default is a light-blue shade.
func (p *PatternLock) SetSelectedColor(selColor color.NRGBA) *PatternLock {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.selectedColor = selColor
	return p
}

func (p *PatternLock) ResetColor() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.selectedColor = defaultColor
}

// Remove the current pattern descriptor to enter Pattern Definition Mode.
// The next full pattern will be stored as current.
func (p *PatternLock) EnterDesignState() *PatternLock {
	p.mux.Lock()
	defer p.mux.Unlock()

	log.Print("Entered design state")
	p.descriptor = nil
	p.active = false // we WILL begin drawing later
	p.designing = true
	p.Status = MSG_STATUS_DEFINE
	p.Refresh()
	return p
}

// Normally this shouldn't be called because PatternLock will
// automatically leave design-mode if DragEnd is received.
func (p *PatternLock) leaveDesignState() *PatternLock {
	p.active = false
	p.designing = false
	p.Status = MSG_STATUS_UNBLOCK
	log.Print("Leaving design state")
	return p
}

// Sets the valid pattern without changing OnValidated
func (p *PatternLock) SetValidPattern(pinfo *PatternInfo) *PatternLock {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.descriptor = pinfo
	log.Print("SetValidPattern:", pinfo.String())
	return p
}

// Sets the valid pattern and specify a callback for validation result.
func (p *PatternLock) SetValidPatternWith(pinfo *PatternInfo, onValidated func(bool)) *PatternLock {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.descriptor = pinfo
	p.OnValidated = onValidated
	return p
}

func (p *PatternLock) CreateRenderer() fyne.WidgetRenderer {
	return newPatternRenderer(p)
}

// Implements fyne.Tappable.
// Let's catch user interaction when the mouse down hits the custom widget,
// not just when the user starts to drag the finger/stylus over the widget.
func (p *PatternLock) Tapped(e *fyne.PointEvent) {
	p.active = true
	if newDot := p.checkHit(e.Position); newDot {
		relativeX := e.Position.X
		relativeY := e.Position.Y
		log.Printf("onTapped (%f,%f)", relativeX, relativeY)
	}
	p.Refresh()
}

// Implements fyne.Draggable
func (p *PatternLock) Dragged(e *fyne.DragEvent) {
	// clear status label if it is beginning
	if !p.active {
		p.SetStatus("")
	}
	p.active = true
	p.hover = e.Position
	p.checkHit(e.Position)
	p.Refresh()
}

// Implements fyne.Draggable
func (p *PatternLock) DragEnd() {
	log.Print("DragEnd")
	if len(p.Sequence) > 0 {
		// 1st call OnComplete if defined
		if p.OnComplete != nil {
			p.OnComplete(p.Sequence)
		}
		// 2nd call OnValidated if defined after updating label.
		// But only if pattern lock descriptor is defined, else
		// there is nothing to validate, just accept as new pattern.
		if !p.designing {
			if p.descriptor != nil {
				// validate pattern against current
				p.onValidating()
			}
		} else {
			// set new pattern
			mode := PatternModeNone.Convert(p.GridSize)
			p.descriptor, _ = NewPattern(p.Sequence, mode)
		}
	}

	if p.designing {
		p.leaveDesignState()
	} else {
		p.active = false
	}
	p.Sequence = []int{}
	p.Refresh()
}

/* -----------------------------------------------------------------
 *                  P R I V A T E    M E T H O D S
 * -----------------------------------------------------------------*/

// checkHit determines if a position is inside a dot's radius
func (p *PatternLock) checkHit(pos fyne.Position) bool {
	added := false
	size := p.Size()
	// The renderer adds a label at the top, but p.Size() is the whole widget.
	// We must account for the label height (approx 30-40px) or ensure coordinates align.
	colWidth := size.Width / float32(p.GridSize)
	rowHeight := (size.Height - statusLABEL_OFFSET) / float32(p.GridSize) // 40 is approx label height + padding
	radius := float64(fyne.Min(colWidth, rowHeight) / 3)                  // Increased tolerance

	col := int(pos.X / colWidth)
	row := int((pos.Y - statusLABEL_OFFSET) / rowHeight)

	if col >= 0 && col < p.GridSize && row >= 0 && row < p.GridSize {
		id := row*p.GridSize + col
		centerX := colWidth*float32(col) + colWidth/2
		centerY := rowHeight*float32(row) + rowHeight/2 + statusLABEL_OFFSET
		dist := math.Sqrt(math.Pow(float64(pos.X-centerX), 2) + math.Pow(float64(pos.Y-centerY), 2))

		if dist < radius {
			for _, v := range p.Sequence {
				if v == id {
					return false
				}
			}
			p.Sequence = append(p.Sequence, id)
			added = true
			log.Printf("Added dot: %s", p.getPos(id))
		}
	}

	return added
}

// Validates the drawn pattern against the pattern that was set
// in the descriptor, updates the status label and if the OnValidated
// callback is set, it is called.
func (p *PatternLock) onValidating() {
	var isValid bool
	if reflect.DeepEqual(p.Sequence, p.descriptor.Pattern()) {
		isValid = true
		p.SetStatus(MSG_STATUS_GRANTED)
	} else {
		isValid = false
		p.SetStatus(MSG_STATUS_WRONG)
	}

	// do the user callback if defined
	if p.OnValidated != nil {
		p.OnValidated(isValid)
	}
}

// get the column as in A,B,C,(D)
func (p *PatternLock) getColumn(index int) rune {
	column := index % p.GridSize
	columnChar := rune(65 + column)
	return columnChar
}

// get the row as 1,2,3,(4)
func (p *PatternLock) getRow(index int) int {
	row := int(index/p.GridSize) + 1
	return row
}

// get the position like A1, C3, etc.
func (p *PatternLock) getPos(index int) string {
	return fmt.Sprintf("%c%d", p.getColumn(index), p.getRow(index))
}
