/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           Go Fynex
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package dlg

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Shows a Directory picker attached to the parent window w. If startFolder is
// empty it begins at the user's Home directory. The filter function may be
// empty if there will be no file filtering.
// Returns the directory path or empty if nothing chosen/cancelled.
func ShowDirectoryPicker(w fyne.Window, startFolder string, filterFunc storage.FileFilter) string {
	chosenDir := ""
	directoryDlg := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil || list == nil {
			return
		}
		chosenDir = list.Path()
	}, w)

	if filterFunc != nil {
		directoryDlg.SetFilter(filterFunc)
	}

	var uri fyne.URI
	if len(startFolder) == 0 {
		home, _ := os.UserHomeDir()
		uri = storage.NewFileURI(home)
	} else {
		uri = storage.NewFileURI(startFolder)
	}
	if listable, err := storage.ListerForURI(uri); err == nil {
		directoryDlg.SetLocation(listable)
	}

	directoryDlg.Show()
	return chosenDir
}

// Show a directory picker but instead of returning the chosen folder, it sets it
// into the bound string variable (it reloads.)
func ShowDirectoryPickerBind(w fyne.Window, startFolder string, filterFunc storage.FileFilter, extStr binding.ExternalString) {
	chose := ShowDirectoryPicker(w, startFolder, filterFunc)
	if len(chose) != 0 {
		extStr.Set(chose)
		extStr.Reload()
	}
}
