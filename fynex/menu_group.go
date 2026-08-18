/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *              github.com/lordofscripts/gofynex
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A menu group administrator accepts a top-level menu as reference
 * to enable/disable all children. With the group, the user can then
 * use Check, Enable, Disable, DisableAll, EnableAll, Clear.
 *-----------------------------------------------------------------*/
package fynex

import "fyne.io/fyne/v2"

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

var _ IMenuGroup = (*menuGroup)(nil)

// Methods to control a Menu Group helper. A menu group is a top-level
// menu and a (sub)set of enrolled MenuItems (not necessarily all children)
// that will be controlled dynamically.
type IMenuGroup interface {
	Enroll(tag MenuItemID, mitem *fyne.MenuItem) *menuGroup
	EnableAll()
	DisableAll()
	Check(tag MenuItemID)
	Enable(tag MenuItemID) *menuGroup
	Disable(tag MenuItemID) *menuGroup
	DisableItem(tag MenuItemID, disable bool) *menuGroup
	ClearAll()
}

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// Empty enumeration for MenuItem identifiers for use in MenuGroup.
// The actual enumeration iota should be done in the end-user application
// using this as type, that in fact makes it customizable to any app.
type MenuItemID byte

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

type menuGroup struct {
	Items  map[MenuItemID]*fyne.MenuItem
	parent *fyne.Menu
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

// (Ctor) a menuGroup is a group of fyne.MenuItems that belong to the
// same fyne.Menu top-level menu as a child and that is of interest to
// us for dynamic control at runtime. If you don't need to change their
// Enabled/Disabled state, then do NOT enroll the item.
// NOTE: the returned value implements the `fynex.IMenuGroup` interface.
func NewMenuGroup(parent *fyne.Menu) *menuGroup {
	return &menuGroup{
		Items:  make(map[MenuItemID]*fyne.MenuItem),
		parent: parent,
	}
}

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

// Enroll a menu item in a group.
func (mg *menuGroup) Enroll(tag MenuItemID, mitem *fyne.MenuItem) *menuGroup {
	if _, exists := mg.Items[tag]; !exists {
		mg.Items[tag] = mitem
	}
	return mg
}

// Enable ALL menu item in this group individually. Fyne v2.8 is still very
// crippled and the useful feature of Enabling/Disabling a top-level menu
// is not yet implemented, we have to do lengthy workarounds
func (mg *menuGroup) EnableAll() {
	for _, mitem := range mg.parent.Items {
		mitem.Disabled = false
	}
}

// Disable ALL menu item in this group. Same crippled Fyne stuff.
func (mg *menuGroup) DisableAll() {
	for _, mitem := range mg.parent.Items {
		mitem.Disabled = true
	}
}

// Cjeck the named menu item and uncheck all others. This is only
// meaninful if the entire Menu is comprised of mutually-exclusive
// options.
func (mg *menuGroup) Check(tag MenuItemID) {
	// find it if it is ours
	target := mg.locate(tag)
	if target != nil {
		// first disable all siblings to prevent spurious states
		for _, mitem := range mg.parent.Items {
			mitem.Checked = false
		}
		// enable the one and only
		target.Checked = true
	}
}

// Enable the menu item
func (mg *menuGroup) Enable(tag MenuItemID) *menuGroup {
	return mg.DisableItem(tag, false)
}

// Disable the menu item
func (mg *menuGroup) Disable(tag MenuItemID) *menuGroup {
	return mg.DisableItem(tag, true)
}

// Disable the specific menu item
func (mg *menuGroup) DisableItem(tag MenuItemID, disable bool) *menuGroup {
	if _, exists := mg.Items[tag]; exists {
		mg.Items[tag].Disabled = disable
	}
	return mg
}

// Removes all enrolled menu items leaving only the parent Menu. It
// only removes it from this object, not from the application lifecycle!
func (mg *menuGroup) ClearAll() {
	clear(mg.Items)
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

// given its tag/ID find the menu item instance, nil if not found.
func (mg *menuGroup) locate(tag MenuItemID) *fyne.MenuItem {
	if mitem, found := mg.Items[tag]; found {
		return mitem
	} else {
		return nil
	}
}

/* ----------------------------------------------------------------
 *                  M A I N    |    D E M O
 *-----------------------------------------------------------------*/

/*
// Demonstrates using the `fynex.MenuGroup` in an application.
func demoMenuGroup()  {
	const (
		// 1. Create your own custom Menu Item identification enumeration
		// using the type provided by the library (MenuItemID)
		MenuFileClose fynex.MenuItemID = iota
		MenuFileOpen
		MenuFileQuit
	)

	// create the menu items to be attached to a top-level menu (File)
	mnuOpen := fyne.NewMenuItem("Open", nil)
	mnuClose := fyne.NewMenuItem("Close", nil)
	mnuQuit := fyne.NewMenuItem("Quit", nil)

	// create the top-level menu to be attached to the menu bar
	topFileMenu := fyne.NewMenu("File",
		mnuOpen,
		mnuClose,
		fyne.NewMenuItemSeparator(),
		mnuQuit,
	)

	// 2. Create a MenuGroup for each top-level and
	// enroll each of the enumerated menu items you intend
	// to control
	mgrpFile := fynex.NewMenuGroup(topFileMenu)
	mgrpFile.Enroll(MenuFileOpen, mnuOpen)
	mgrpFile.Enroll(MenuFileClose, mnuClose)
	mgrpFile.Enroll(MenuFileQuit, mnuQuit)
	mgrpFile.EnableAll() // Fyne has no way to control a TL menu
	mgrpFile.DisableAll() // idem
	mgrpFile.Enable(MenuFileClose)
	mgrpFile.Disable(MenuFileOpen)
	mgrpFile.Check(MenuFileOpen)

	fyne.NewMainMenu(topFileMenu)

	mgrpFile.ClearAll()
}
*/
