/* -----------------------------------------------------------------
 *              L o r d  O f   S c r i p t s (tm)
 *             Copyright (C)2026 Dídimo Grimaldo T.
 *                           APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Package/Module/Application version information
 *-----------------------------------------------------------------*/
package gofynex

import (
	"fmt"
	"runtime"
	"strings"
)

/* ----------------------------------------------------------------
 *                       G L O B A L S
 *-----------------------------------------------------------------*/

// These must be adapted for the current module
const (
	// in case vcsVersion not injected during link phase
	MANUAL_VERSION string = "1.2.0"
	// Change these values accordingly
	NAME     string    = "Go-Fynex"
	DESC     string    = "An custom widget and helper library for Fyne"
	STATUS   devStatus = statusReleased
	REVISION           = 1
)

// These are for all modules where this template is used
const (
	// Useful Unicode Characters
	CHR_COPYRIGHT       = '\u00a9'      // ©
	CHR_REGISTERED      = '\u00ae'      // ®
	CHR_GUILLEMET_L     = '\u00ab'      // «
	CHR_GUILLEMET_R     = '\u00bb'      // »
	CHR_TRADEMARK       = '\u2122'      // ™
	CHR_SAMARITAN       = '\u214f'      // ⅏
	CHR_PLACEOFINTEREST = '\u2318'      // ⌘
	CHR_HIGHVOLTAGE     = '\u26a1'      // ⚡
	CHR_TRIDENT         = rune(0x1f531) // 🔱
	CHR_SPLATTER        = rune(0x1fadf)
	CHR_WARNING         = '\u26a0' // ⚠
	CHR_EXCLAMATION     = '\u2757'
	CHR_SKULL           = '\u2620' // ☠

	CO1 = "odlamirG omidiD 5202)C("
	CO2 = "stpircS fO droL 5202)C("
	CO3 = "gnitirwnitsol"
)

var Version PackageVersion

/* ----------------------------------------------------------------
 *                       L O C A L S
 *-----------------------------------------------------------------*/

const (
	// don't change
	statusAlpha    devStatus = "Alpha"
	statusBeta     devStatus = "Beta"
	statusRC       devStatus = "RC" // Release Candidate
	statusReleased devStatus = ""
)

var (
	vcsVersion  string // automatically injected with linker
	vcsCommit   string
	vcsDate     string
	vcsBuildNum string
	//NOT USEDgo:embed version.txt
)

/* ----------------------------------------------------------------
 *                     I N T E R F A C E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                         T Y P E S
 *-----------------------------------------------------------------*/

type devStatus = string

/* ----------------------------------------------------------------
 *                   P U B L I C    T Y P E S
 *-----------------------------------------------------------------*/

// Package/Module/Application version descriptor
type PackageVersion struct {
	n  string    // name
	v  string    // version tag
	s  devStatus // status
	sv int       // Alpha/Beta/RC-### sequence
}

/* ----------------------------------------------------------------
 *                   P R I V A T E    T Y P E S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                     I N I T I A L I Z E R
 *-----------------------------------------------------------------*/
func init() {
	Version = PackageVersion{
		n:  NAME,
		v:  MANUAL_VERSION,
		s:  STATUS,
		sv: REVISION,
	}
}

/* ----------------------------------------------------------------
 *                    C O N S T R U C T O R S
 *-----------------------------------------------------------------*/

func NewPackageVersion(name, description string, verStr string, status devStatus) PackageVersion {
	return PackageVersion{
		n:  name,
		v:  verStr,
		s:  status,
		sv: 1,
	}
}

/* ----------------------------------------------------------------
 *                        M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                 P U B L I C    M E T H O D S
 *-----------------------------------------------------------------*/

func (pv PackageVersion) BuildMeta() string {
	ver := vcsVersion
	if len(vcsVersion) == 0 {
		ver = "v" + MANUAL_VERSION
	}
	return fmt.Sprintf("\t\t%s-%s %s", ver, vcsBuildNum, vcsDate)
}

func (pv PackageVersion) CommitInfo() string {
	return fmt.Sprintf("Build #%s (%s)", vcsBuildNum, vcsCommit)
}

func (pv PackageVersion) Short() string {
	var ver string

	if len(vcsVersion) != 0 {
		pv.v = vcsVersion
	}
	var buildInfo string = ""
	if vcsBuildNum != "" {
		buildInfo = fmt.Sprintf(" build #%s", vcsBuildNum)
	}

	switch pv.s {
	case statusAlpha:
		fallthrough
	case statusBeta:
		fallthrough
	case statusRC:
		ver = fmt.Sprintf("v%s-%s-%d%s", pv.v, pv.s, pv.sv, buildInfo)
	default:
		ver = fmt.Sprintf("v%s %s", pv.v, buildInfo)
	}
	return ver
}

func (pv PackageVersion) String() string {
	var ver string

	if len(vcsVersion) != 0 {
		pv.v = vcsVersion
	}
	var buildInfo string = ""
	if vcsBuildNum != "" {
		buildInfo = fmt.Sprintf(" build #%s", vcsBuildNum)
	}

	switch pv.s {
	case statusAlpha:
		fallthrough
	case statusBeta:
		fallthrough
	case statusRC:
		ver = fmt.Sprintf("%s v%s-%s-%d%s", pv.n, pv.v, pv.s, pv.sv, buildInfo)
	default:
		ver = fmt.Sprintf("%s v%s %s", pv.n, pv.v, buildInfo)
	}
	return ver
}

/* ----------------------------------------------------------------
 *                 P R I V A T E    M E T H O D S
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *                       F U N C T I O N S
 *-----------------------------------------------------------------*/

// Funny LordOfScripts logo
func Logo() string {
	const (
		whiteStar rune = '\u269d' // ⚝
		unisex    rune = '\u26a5' // ⚥
		hotSpring rune = '\u2668' // ♨
		leftConv  rune = '\u269e' // ⚞
		rightConv rune = '\u269f' // ⚟
		eye       rune = '\u25d5' // ◕
		mouth     rune = '\u035c' // ͜	‿ \u203f
		skull     rune = '\u2620' // ☠
	)
	return fmt.Sprintf("%c%c%c %c%c", leftConv, eye, mouth, eye, rightConv)
	//fmt.Sprintf("(%c%c %c)", eye, mouth, eye)
}

// Hey! My time costs money too!
func BuyMeCoffee(coffee4 ...string) {
	const (
		coffee rune = '\u2615' // ☕
	)

	var recipient string
	if len(coffee4) == 0 {
		recipient = Reverse(CO3)
	} else {
		recipient = coffee4[0]
	}

	fmt.Printf("\t%c Buy me a Coffee? https://www.buymeacoffee/%s\n", coffee, recipient)
}

func Copyright(owner string, withLogo bool) {
	fmt.Printf("\t%c %s %s %c\n", CHR_TRIDENT, Version, Reverse(owner), CHR_TRIDENT)
	fmt.Println("\t\t\t\t", Logo())
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// get the current GO language version
func GoVersion() string {
	ver := strings.Replace(runtime.Version(), "go", "", -1)
	return ver
}

// retrieve the current GO language version and compare it
// to the minimum required. It returns the current version
// and whether the condition current >= min is fulfilled or not.
func GoVersionMin(min string) (string, bool) {
	current := GoVersion()
	ok := current >= min
	return current, ok
}
