# For Developers

## Linux (`linux`)

On Linux there is not much extra you must do other than installing Go and Fyne.

## MacOS (`darwin`)

Never got my hands on one of these. Would have been nice but it is terrible that
its file system naming is not case sensitive like its parent Linux.

## Windows (`windows`)

Poor you, but make sure you have set your VSCode/VSCodium environment for 
`LF` as line-endings. No `CRLF` in my projects please.

It is quite cumbersome, but I started this project on a Windows machine.

* Install **MSys2** package manager
* On an MSys2 terminal type `pacman -S make` to get Make installed.
* Install the TDM-GCC-64 C/C++ compiler
* Install the `fyne` tool
* Compile `fyne package -os windows -icon Icon.png -release`
* Run `.\PatternLock.exe`