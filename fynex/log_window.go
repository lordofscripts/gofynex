/* *****************************************************************
 *              Copyright(C)2026 Lord of Scripts
 *                      All Rights Reserved
 * -----------------------------------------------------------------
 * Non-modal window where standard Log messages are redirected.
 ********************************************************************/
package fynex

import (
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/* -----------------------------------------------------------------
 *                     I N T E R F A C E S
 * -----------------------------------------------------------------*/

var _ io.Writer = (*logWriter)(nil)

/* -----------------------------------------------------------------
 *                  P R I V A T E    T Y P E S
 * -----------------------------------------------------------------*/

// logWriter implements io.Writer to redirect logs to a Fyne Entry
type logWriter struct {
	output *widget.Entry
}

/* -----------------------------------------------------------------
 *                  C O N S T R U C T O R S
 * -----------------------------------------------------------------*/

// Create a non-modal window with a multi-line entry widget
// to which the log output will be redirected. Call Show()
// method before the main windows's ShowAndRun()
func NewLogWindow(myApp fyne.App, width, height float32) fyne.Window {
	// 1. Create the Log Window (Non-modal)
	logWindow := myApp.NewWindow("System Logs")
	logWindow.Resize(fyne.NewSize(width, height))

	// 2. Create Multi-line Entry for logs
	logEntry := widget.NewMultiLineEntry()
	logEntry.Wrapping = fyne.TextTruncate
	logWindow.SetContent(container.NewStack(logEntry)) // 3. Redirect standard log output

	writer := &logWriter{output: logEntry}
	log.SetOutput(writer)

	return logWindow
}

/* -----------------------------------------------------------------
 *                       M E T H O D S
 * -----------------------------------------------------------------*/

// implements io.Writer so that Log messages can be redirected to Entry widget
func (w *logWriter) Write(p []byte) (n int, err error) {
	// Append text to the entry
	w.output.SetText(w.output.Text + string(p))
	// Scroll to the bottom
	w.output.CursorColumn = 0
	w.output.CursorRow = len(w.output.Text)
	return len(p), nil
}

/* -----------------------------------------------------------------
 *                       M A I N  |  D E M O
 * -----------------------------------------------------------------*/

/*
func demo() {
	myApp := app.New()
	mainWindow := myApp.NewWindow("Main Console")

	// 1. Create the Log Window (Non-modal)
	logWindow := NewLogWindow(500, 300)

	// UI to trigger logs
	btn := widget.NewButton("Generate Log Entry", func() {
		log.Printf("Button clicked at: %s", time.Now().Format("15:04:05"))
	})

	mainWindow.SetContent(container.NewVBox(		widget.NewLabel("Click the button to see logs in the other window"),
		btn,
	))

	// Show both windows
	logWindow.Show()
	mainWindow.ShowAndRun()
}
*/
