package widgets

import (
	"fastgh3/fastbrowser/gh"
	"image"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type SongList struct {
	Theme         *material.Theme
	Charts        *[]*gh.Chart
	SelectedChart *gh.Chart

	scrollPosition float32
	buttons        []*SongButton
}

func NewSongList(theme *material.Theme, charts *[]*gh.Chart) *SongList {
	return &SongList{
		Theme:         theme,
		Charts:        charts,
		SelectedChart: nil,

		scrollPosition: 0,
		buttons:        make([]*SongButton, 0),
	}
}

func (sl *SongList) Layout(gtx layout.Context) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, sl)

	for {
		ev, ok := gtx.Event(pointer.Filter{Target: sl, Kinds: pointer.Scroll, ScrollY: pointer.ScrollRange{Min: -1, Max: 1}})
		if !ok {
			break
		}

		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			case pointer.Scroll:
				sl.scrollPosition += x.Scroll.Y
			}
		}
	}

	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 150})

	gtx.Constraints.Max.X -= 8
	gtx.Constraints.Max.Y -= 8
	defer op.Offset(image.Pt(4, 4)).Push(gtx.Ops).Pop()
	defer op.Offset(image.Pt(0, int(-sl.scrollPosition*32))).Push(gtx.Ops).Pop()

	i := 0
	for _, chart := range *sl.Charts {
		if len(sl.buttons) < i+1 {
			btn := &SongButton{Chart: chart, Theme: sl.Theme}
			sl.buttons = append(sl.buttons, btn)
		}
		btn := sl.buttons[i]
		btn.Layout(gtx, func(sb *SongButton) { sl.SelectedChart = sb.Chart })
		defer op.Offset(image.Pt(0, 36)).Push(gtx.Ops).Pop()
		i++
	}

	if sl.scrollPosition < 0 {
		sl.scrollPosition = 0
		gtx.Execute(op.InvalidateCmd{})
	}
	if sl.scrollPosition*32 > 36*float32(i-1) {
		sl.scrollPosition = 36 * float32(i-1) / 32
		gtx.Execute(op.InvalidateCmd{})
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
