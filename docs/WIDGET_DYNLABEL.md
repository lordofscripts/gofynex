# Dynamic Label Custom Widget

Once again, oftentimes you find yourself on foot while trying to develop
a Fyne application. Themes are nice, but sometimes you need to deviate
and Fyne doesn't allow you unless you go through workarounds that take
your time. Mostly, you use custom widgets.

### Use Case 1

You need to display text for which a typical `Label` won't do. And
due to Fyne's quirky nature, a standard `Entry` widget in disabled
state is simply unreadable to the human eye. I guess the designers
just had simple use cases in mind. But life is about building 
complex things.

The `DynamicLabel` fulfills that and more in an aesthetic manner.

### Use Case 2

I needed a set of labels in my application, but not just a crippled
label that works on theme. I needed:

* ability to perform an action when I clicked on a label.
* I didn't want it to look like a button
* I needed to be able to change its color, etc.

Because basically, sometimes you want to display short text that
would fit in an `Entry` widget, but you need it to be **read-only**.
That is a perfectly valid use case and that functionality is present
in just about every GUI library I got my hands on over decades 
developing software. If you try to do that in Fyne, the disabled
state of `Entry` renders the text unreadable for most humans.

[ ![Buy me a coffee](./assets/buymecoffee.jpg)](https://BuyMeACoffee.com/lostinwriting)