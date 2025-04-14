package widgets

import (
	"fastgh3/fastbrowser/gh"
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type SongDetails struct {
	Theme *material.Theme
	Chart *gh.Chart

	playButton *widget.Clickable
}

func NewSongDetails(theme *material.Theme) *SongDetails {
	return &SongDetails{
		Theme: theme,
		Chart: nil,

		playButton: &widget.Clickable{},
	}
}

func (sd *SongDetails) Layout(gtx layout.Context) layout.Dimensions {
	constraints := gtx.Constraints.Constrain(image.Pt(360, gtx.Constraints.Max.Y))

	defer clip.Rect{Max: constraints}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 150})

	if sd.Chart == nil {
		gtx := gtx
		gtx.Constraints = layout.Exact(constraints)

		layout.Center.Layout(gtx, material.Label(sd.Theme, unit.Sp(18), "No song selected").Layout)
	} else {
		gtx := gtx
		gtx.Constraints = layout.Exact(constraints.Sub(image.Pt(8, 8)))

		playBtn := material.Button(sd.Theme, sd.playButton, "Start FastGH3")
		playBtn.TextSize = unit.Sp(18)
		playBtn.Background = color.NRGBA{39, 98, 33, 255}

		defer op.Offset(image.Pt(4, 4)).Push(gtx.Ops).Pop()
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(material.Label(sd.Theme, unit.Sp(22), fmt.Sprintf("%s - %s", sd.Chart.Artist, sd.Chart.Name)).Layout),
			layout.Rigid(material.Label(sd.Theme, unit.Sp(18), fmt.Sprintf("Album: %s", sd.Chart.Album)).Layout),
			layout.Rigid(material.Label(sd.Theme, unit.Sp(18), fmt.Sprintf("Year: %s", sd.Chart.Year)).Layout),
			layout.Rigid(material.Label(sd.Theme, unit.Sp(18), fmt.Sprintf("Charter: %s", sd.Chart.Charter)).Layout),
			layout.Rigid(playBtn.Layout),
		)

	}

	return layout.Dimensions{Size: constraints}
}
