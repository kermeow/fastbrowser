package browser

import (
	"fastgh3/fastbrowser/browser/widgets"
	"fastgh3/fastbrowser/config"
	"fastgh3/fastbrowser/fonts"
	"fastgh3/fastbrowser/gh"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Browser struct {
	*app.Window

	Config *config.Config
	Theme  *material.Theme
	Charts []*gh.Chart

	songList *widgets.SongList
}

func New(conf *config.Config) *Browser {
	window := new(app.Window)
	window.Option(
		app.Size(unit.Dp(800), unit.Dp(600)),
		app.Title("FastBrowser"),
	)
	window.Perform(system.ActionCenter)

	theme := material.NewTheme().WithPalette(
		material.Palette{
			Bg: color.NRGBA{18, 18, 18, 255},
			Fg: color.NRGBA{240, 240, 240, 255},
		},
	)
	theme.Shaper = text.NewShaper(text.WithCollection(fonts.Collection()))
	theme.Face = "Nunito"

	charts := make([]*gh.Chart, 0)

	browser := &Browser{
		Window: window,

		Config: conf,
		Theme:  &theme,
		Charts: charts,
	}
	browser.songList = widgets.NewSongList(&theme, &browser.Charts)
	return browser
}

func (ui *Browser) Run() error {
	// ui.songList.Charts = &ui.Charts
	ui.getCharts()

	var ops op.Ops
	for {
		switch e := ui.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, ui.Theme.Bg)

			inset := layout.UniformInset(unit.Dp(16))
			inset.Layout(gtx, ui.draw)

			e.Frame(gtx.Ops)
		}
	}
}

func (ui *Browser) draw(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, ui.songList.Layout),
		layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
	)
}

func (ui *Browser) getCharts() {
	
}
