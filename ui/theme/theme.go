// Package theme implements themed widgets.
// The primary design element of this package is the "theme factory", which provides constructors
// for the various widgets.
package theme

import (
	"image/color"

	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Factory manufactures styled widgets.
type Factory struct {
	// Base is the underlying theme required to use pre-fabbed material widgets.
	// If material package changes it's API, we may not need this in the future.
	Base material.Theme
	// Palette maps semantic colors to local colors. These colors will change depending on the
	// theme, eg light vs dark themes.
	Palette
	// BreakPoints maps device types to pixel widths for reponsive layouts.
	BreakPoints
}

// Palette specifies semantic colors.
type Palette struct {
	// Primary color displayed most frequently across screens and components.
	Primary        color.NRGBA
	PrimaryVariant color.NRGBA

	// Secondary color used sparingly to accent ui elements.
	Secondary        color.NRGBA
	SecondaryVariant color.NRGBA

	// Surface affects surfaces of components such as cards, sheets and menus.
	Surface color.NRGBA

	// Background appears behind scrollable content.
	Background color.NRGBA

	// Error indicates errors in components.
	Error color.NRGBA

	// On colors appear "on top" of the base color.
	// Choose contrasting colors.
	OnPrimary    color.NRGBA
	OnSecondary  color.NRGBA
	OnBackground color.NRGBA
	OnSurface    color.NRGBA
	OnError      color.NRGBA
}

// BreakPoints specifies device viewports.
// Used to decide when to change layout for responsive user interfaces.
type BreakPoints struct {
	// (phone)
	Tiny BreakPoint
	// (tablet)
	Small BreakPoint
	// (laptop)
	Medium BreakPoint
	// (desktop)
	Large BreakPoint
}

type BreakPoint struct {
	// Width of the device screen.
	Width unit.Value
	// Zero margins means "fluid".
	Margins unit.Value
	// Zero body means "fluid".
	Body unit.Value
}

// MaterialDesign is the baseline palette for material design.
// https://material.io/design/color/the-color-system.html#color-theme-creation
var MaterialDesign Palette = Palette{
	Primary:          RGB(0x6200EE),
	PrimaryVariant:   RGB(0x3700B3),
	Secondary:        RGB(0x03DAC6),
	SecondaryVariant: RGB(0x018786),
	Background:       RGB(0xFFFFFF),
	Surface:          RGB(0xFFFFFF),
	Error:            RGB(0xB00020),
	OnPrimary:        RGB(0xFFFFFF),
	OnSecondary:      RGB(0x000000),
	OnBackground:     RGB(0x000000),
	OnSurface:        RGB(0x000000),
	OnError:          RGB(0xFFFFFF),
}

// Fluid sentinel value for boolean comparison.
// Used to check if a break point is fluid.
var Fluid = unit.Dp(0)

// DefaultBreakPoints is the screen breakpoints for material design.
// https://material.io/design/layout/understanding-layout.html#layout-anatomy
var DefaultBreakPoints BreakPoints = BreakPoints{
	Tiny: BreakPoint{
		Width:   unit.Dp(599),
		Margins: unit.Dp(16),
		Body:    Fluid,
	},
	Small: BreakPoint{
		Width:   unit.Dp(1240),
		Margins: unit.Dp(32),
		Body:    Fluid,
	},
	Medium: BreakPoint{
		Width:   unit.Dp(1439),
		Margins: unit.Dp(200),
		Body:    Fluid,
	},
	Large: BreakPoint{
		Width:   unit.Dp(1439),
		Margins: Fluid,
		Body:    unit.Dp(1040),
	},
}

// NewTheme allocates a theme factory using the given fonts and palette.
func NewFactory(fonts []text.FontFace, p Palette, b BreakPoints) *Factory {
	return &Factory{
		Base:        *material.NewTheme(fonts),
		Palette:     p,
		BreakPoints: b,
	}
}

// RGB creates a color using hexidecimal notation.
func RGB(c uint32) color.NRGBA {
	return ARGB(0xff000000 | c)
}

// ARGB creates a color using hexidecimal notation.
func ARGB(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
