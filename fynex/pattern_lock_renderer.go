/* *****************************************************************
 *					Copyright(C)2026 Lord of Scripts
 *						All Rights Reserved
 * -----------------------------------------------------------------
 * A custom renderer for the custom PatternLock widget.
 ********************************************************************/
package fynex

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

/* -----------------------------------------------------------------
 *                       G L O B A L S
 * -----------------------------------------------------------------*/

const (
	// a fixed offset to ensure both the custom widget and renderer
	// properly identify where the user touches is where the dots are.
	statusLABEL_OFFSET = 40
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ fyne.WidgetRenderer = (*patternRenderer)(nil)

/* -----------------------------------------------------------------
 *                  P R I V A T E    T Y P E S
 * -----------------------------------------------------------------*/

type patternRenderer struct {
	p           *PatternLock
	statusLabel *widget.Label
	background  *canvas.Image
	fadeOverlay *canvas.Rectangle
	objects     []fyne.CanvasObject
	lastSize    fyne.Size // the last size used by Layout()
}

/* -----------------------------------------------------------------
 *                  C O N S T R U C T O R S
 * -----------------------------------------------------------------*/

func newPatternRenderer(p *PatternLock) *patternRenderer {
	// Create a rectangle to act as the "fade" layer
	// Black with alpha 180 creates a dark fade; use White for a light fade
	fade := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 150})

	var img *canvas.Image = nil
	if p.backgroundRsrc != nil {
		//img := canvas.NewImageFromFile("/home/didi/go/patternlock/widget_background_lotr.png")
		img = canvas.NewImageFromResource(p.backgroundRsrc)
		img.FillMode = canvas.ImageFillStretch
	}

	r := &patternRenderer{
		p: p,
		statusLabel: widget.NewLabelWithStyle("Design your unlock pattern",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true}),
		background:  img,
		fadeOverlay: fade,
	}
	// ensure the first Refresh() has a valid size greater than 0,0
	r.lastSize = r.MinSize()
	return r
}

/* -----------------------------------------------------------------
 *                       M E T H O D S
 * -----------------------------------------------------------------*/

func (r *patternRenderer) Layout(size fyne.Size) {
	// status label at the top
	r.lastSize = size
	// if there is a widget background image, lay it out
	if r.background != nil {
		// if we don't do this, it stays at size 0x0 and remains invisible!
		r.background.Resize(size)
		r.background.Move(fyne.NewPos(0, 0))
	}
	// also resize the overlay
	r.fadeOverlay.Resize(size)
	r.fadeOverlay.Move(fyne.NewPos(0, 0))

	r.statusLabel.Resize(fyne.NewSize(size.Width, 30))
	r.statusLabel.Move(fyne.NewPos(0, 5))
}

func (r *patternRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 340)
}

func (r *patternRenderer) Destroy() {}

func (r *patternRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *patternRenderer) Refresh() {
	// 1. Synchronize the background canvas object with the widget's resource
	if r.p.backgroundRsrc != nil {
		if r.background == nil {
			// Create the canvas object if it didn't exist at startup
			r.background = canvas.NewImageFromResource(r.p.backgroundRsrc)
			r.background.FillMode = canvas.ImageFillStretch
			// Manually resize and move because Layout() won't be called automatically
			r.background.Resize(r.lastSize)
			r.background.Move(fyne.NewPos(0, 0))
		} else if r.background.Resource != r.p.backgroundRsrc {
			// Update the resource if it changed and trigger an image refresh
			r.background.Resource = r.p.backgroundRsrc
			r.background.Refresh()
		}
	} else {
		r.background = nil
	}

	// 2. Update the status label text
	r.statusLabel.SetText(r.p.Status)

	// 3. Rebuild the objects slice
	// We must rebuild this list so Fyne knows to include the background if it was previously nil
	r.objects = []fyne.CanvasObject{}
	if r.background != nil {
		r.objects = append(r.objects, r.background)
	}
	r.objects = append(r.objects, r.fadeOverlay)
	r.objects = append(r.objects, r.statusLabel)

	// 4. Drawing logic for lines and dots
	renderSize := r.lastSize
	gridAreaTop := float32(40)
	gridHeight := renderSize.Height - gridAreaTop
	colWidth := renderSize.Width / float32(r.p.GridSize)
	rowHeight := gridHeight / float32(r.p.GridSize)
	dotRadius := fyne.Min(colWidth, rowHeight) / 8

	lineColor := r.p.selectedColor
	if r.p.designing {
		lineColor = color.NRGBA{R: 0x9d, G: 0, B: 0xff, A: 255}
	}

	// Draw lines
	if len(r.p.Sequence) > 1 {
		for i := 0; i < len(r.p.Sequence)-1; i++ {
			r.objects = append(r.objects, r.createLine(r.p.Sequence[i], r.p.Sequence[i+1], colWidth, rowHeight, gridAreaTop, lineColor))
		}
	}

	// Draw active drag line
	if r.p.active && len(r.p.Sequence) > 0 {
		lastIdx := r.p.Sequence[len(r.p.Sequence)-1]
		lastX := colWidth*float32(lastIdx%r.p.GridSize) + colWidth/2
		lastY := rowHeight*float32(lastIdx/r.p.GridSize) + rowHeight/2 + gridAreaTop

		line := canvas.NewLine(lineColor)
		line.Position1 = fyne.NewPos(lastX, lastY)
		line.Position2 = r.p.hover
		line.StrokeWidth = 4
		r.objects = append(r.objects, line)
	}

	// Draw dots
	for i := 0; i < r.p.GridSize*r.p.GridSize; i++ {
		x := colWidth*float32(i%r.p.GridSize) + colWidth/2
		y := rowHeight*float32(i/r.p.GridSize) + rowHeight/2 + gridAreaTop

		dotColor := color.Color(color.Gray{Y: 150})
		for _, seqID := range r.p.Sequence {
			if seqID == i {
				dotColor = lineColor
				break
			}
		}

		dot := canvas.NewCircle(dotColor)
		dot.Resize(fyne.NewSize(dotRadius*2, dotRadius*2))
		dot.Move(fyne.NewPos(x-dotRadius, y-dotRadius))
		r.objects = append(r.objects, dot)
	}

	// Canvas Refresh to ensure the background is painted
	if r.background != nil {
		r.background.Refresh()
	}
}

/* -----------------------------------------------------------------
 *                  P R I V A T E    M E T H O D S
 * -----------------------------------------------------------------*/

func (r *patternRenderer) createLine(idx1, idx2 int, w, h, top float32, clr color.Color) *canvas.Line {
	l := canvas.NewLine(clr)
	l.StrokeWidth = 6
	l.Position1 = fyne.NewPos(w*float32(idx1%r.p.GridSize)+w/2, h*float32(idx1/r.p.GridSize)+h/2+top)
	l.Position2 = fyne.NewPos(w*float32(idx2%r.p.GridSize)+w/2, h*float32(idx2/r.p.GridSize)+h/2+top)
	return l
}
