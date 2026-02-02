# Scrollable Slider Widget

To this date I have developed using multiple GUI libraries. That I remember
in all of them a widget Tooltip to display helpful, unobtrussive information
was easy to add to your application without having to resort to extra
development. Well, Fyne doesn't have that, one of its several shortcomings.

As a GUI developer, sometimes you need to use a `Slider` widget to let
the user select a discrete value from a range. With a tooltip the user
gets immediate feedback to learn when to stop. Fyne doesn't offer that.
So, rather than having to reimplement the same again and again I decided
to include here my custom slider.

![](./assets/sample_slider.png)

Features:

* Uses the standard Fyne `widget.Slider` but customized.
* Use the mouse scrollwheel to increase/decrease the `slider.Value`
* On the left of the slider, a small label with customizable template.
* On the right of the slider an optional small label that is normally hidden.
* Supports data binding

The left label always displays the current `slider.Value` to compensate
for the lack of a tooltip. It is displayed by default as an integer,
but should you prefer, use another format like `%.2f` by calling:

> betterSlider.SetValueTemplate(format string)

Or for the right text (hidden by default):

> betterSlider.SetRightText(convertedVal)
> betterSlider.SetRightVisible(true)

A typical use case for the right label is if your slider has a value
within a range, say ASCII value for A..Z and on the left you display
the slider value (ASCII value) and on the right label the corresponding
alphabet letter.

[ ![Buy me a coffee](./assets/buymecoffee.jpg)](https://BuyMeACoffee.com/lostinwriting)