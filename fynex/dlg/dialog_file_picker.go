/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           Go Fynex
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package dlg

import (
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/lordofscripts/gofynex/fynex"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Menu: File | Open
// Shows a File Open picker attached to the parent window w. If startFolder is
// empty it begins at the user's Home directory. The filter function may be
// empty if there will be no file filtering.
// Returns the directory path or empty if nothing chosen/cancelled.
func ShowFilePicker(w fyne.Window, startFolder string, filterFunc storage.FileFilter) string {
	chosenFile := ""
	fileDlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
		}
		if reader == nil {
			return // user cancelled
		}

		defer reader.Close()
		chosenFile = reader.URI().Path()
	}, w)

	if filterFunc != nil {
		fileDlg.SetFilter(filterFunc)
	}

	var uri fyne.URI
	if len(startFolder) == 0 {
		home, _ := os.UserHomeDir()
		uri = storage.NewFileURI(home)
	} else {
		uri = storage.NewFileURI(startFolder)
	}
	if listable, err := storage.ListerForURI(uri); err == nil {
		fileDlg.SetLocation(listable)
	}

	fileDlg.Show()
	return chosenFile
}

// Menu: File | Open
// Show a directory picker but instead of returning the chosen folder, it sets it
// into the bound string variable (it reloads.)
func ShowFilePickerBind(w fyne.Window, startFolder string, filterFunc storage.FileFilter, extStr binding.ExternalString) {
	chose := ShowFilePicker(w, startFolder, filterFunc)
	if len(chose) != 0 {
		extStr.Set(chose)
		extStr.Reload()
	}
}

// Menu: File | Open
// Shows a File Open picker attached to the parent window w. If startFolder is
// empty it begins at the user's Home directory. The filter function may be
// empty if there will be no file filtering.
// Returns the directory path or empty if nothing chosen/cancelled.
func ShowFileSave[T string | []byte](w fyne.Window, filename string, data T) fynex.TriState {
	success := fynex.Unset
	fileDlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			success = fynex.Unchecked
			return
		}
		if writer == nil {
			success = fynex.Unset // cancelled
			return
		}

		defer writer.Close()
		if str, ok := any(data).(string); ok {
			_, err = io.WriteString(writer, str)
		} else if dataBytes, ok := any(data).([]byte); ok {
			_, err = writer.Write(dataBytes)
			if err != nil {
				dialog.ShowError(err, w)
			}
		}

		if err != nil {
			dialog.ShowError(err, w)
		} else {
			log.Printf("Data successfully written to %s", writer.URI().Path())
		}
	}, w)

	fileDlg.SetFileName(filename)
	fileDlg.Show()
	return success
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
func demoFileOpen() {
	imageFilter := storage.NewExtensionFileFilter([]string{".jpg", ".png"})
	filename := ShowFilePicker(window, "/home/Pictures", imageFilter)
}
*/
