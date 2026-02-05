# Go Fyne Custom Widgets

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/gofynex)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/gofynex?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/gofynex)
![Build](https://github.com/lordofscripts/gofynex/actions/workflows/go-fyne.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/gofynex.svg)](https://pkg.go.dev/github.com/lordofscripts/gofynex)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/gofynex)](https://github.com/lordofscripts/gofynex/releases/latest)[![License: MIT](https://img.shields.io/badge/License-MIT-lightgrey.svg)](https://choosealicense.com/licenses/mit/)

Anybody coming from the world of Microsoft .NET and its versatile `WinForms`
GUI library will soon find out that the somewhat portability of `Fyne` 
comes at a price. Many Fyne widgets miss some basic features with valid
and common use cases. Despite that, that functionality is left out of 
many Fyne widgets. The result is that the developer usually ends up 
spending extra time adding missing functionality to existent Fyne widgets.

This (growing) `gofynex` library was born out of that frustration. And I
decided to create this library so that I, as well as others, can benefit
from it by reusing it. No need to reinvent the wheel.

## Custom Dialogs

* [About dialog](./DIALOG_ABOUT.md)

## Custom Windows

* [Log window](./WINDOW_LOG.md) for outputing `log` data.

## Custom Widgets

* [Dynamic Label](./WIDGET_DYNLABEL.md) widget
* [Pattern Lock](./WIDGET_PATTERN_LOCK.md) widget
* [Scrollable slider](./WIDGET_SLIDER.md) widget. An enhanced `widget.Slider`
* A custom [Tri-color LED](./WIDGET_LED.md) to display tri-state values
* A flexible [Flexible Mini Theme](./MINI_THEME.md)
* A [Person widget](./WIDGET_PERSON.md)

### Sponsor Me

If you like my work -which takes useful free time that you don't have to spend- please
consider [Sponsoring ️❤️ me](https://github.com/sponsors/lordofscripts). Or...

[ ![Buy me a coffee](./assets/buymecoffee.jpg)](https://BuyMeACoffee.com/lostinwriting)