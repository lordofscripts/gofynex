# Flex Mini Theme

Yes, contrary to what Fyne creators think, there are times when the current
theme presents the information in disproportionate sizes that look aweful.

This module includes a `FlexMiniTheme` that you can create dynamically and
inject custom colors and custom sizes for the `fyne.ThemeSizeNa,e` and
`fyne.ThemeColorName` theme properties. 

The `FlexMiniTheme` can be used to dynamically create a theme for your
application, or just a mini-theme for use in a specific widget where you 
need fine-grained control beyond the restrictions of a theme.

First, create the mini theme instance:

> miniTheme ::= fynex.NewFlexMiniTheme(defaultTheme)

Notice that you provide the current or default theme. Any requested
color or size that has not been injected into the mini-theme will 
fallback to the provided default theme.

You can override (as many) theme sizes, for example to use smaller text:

> miniTheme.Include(theme.SizeNameText, 10)

You can also override (as many) theme colors:

> miniTheme.IncludeColor(theme.ColorNameForeground, Orange)


[ ![Buy me a coffee](./assets/buymecoffee.jpg)](https://BuyMeACoffee.com/lostinwriting)