package widgets

import (
	"fastgh3/fastbrowser/gh"
	"fmt"
	"image"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type PressCallback func(*SongButton)

type SongButton struct {
	Theme *material.Theme
	Chart *gh.Chart

	hovering bool
	pressed  bool
	pressId  pointer.ID
}

func (sb *SongButton) Layout(gtx layout.Context, cb PressCallback) layout.Dimensions {
	gtx.Constraints = layout.Exact(image.Pt(gtx.Constraints.Max.X, 32))

	bounds := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	rrect := clip.UniformRRect(bounds, 8)
	defer rrect.Push(gtx.Ops).Pop()

	event.Op(gtx.Ops, sb)

	for {
		ev, ok := gtx.Event(pointer.Filter{Target: sb, Kinds: pointer.Enter | pointer.Leave | pointer.Press | pointer.Release})
		if !ok {
			break
		}

		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			case pointer.Enter:
				sb.hovering = true
			case pointer.Leave:
				sb.hovering = false
			case pointer.Press:
				sb.pressed = true
				sb.pressId = x.PointerID
			case pointer.Release:
				sb.pressed = false
				if x.PointerID == sb.pressId {
					cb(sb)
				}
			}
		}
	}

	paint.Fill(gtx.Ops, sb.Theme.Bg)
	if sb.hovering {
		paint.FillShape(gtx.Ops, color.NRGBA{200, 200, 200, 100},
			clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op(),
		)
	}

	length := "?:??"
	if sb.Chart.SongLength != 0 {
		totalSeconds := sb.Chart.SongLength / 1000
		seconds := totalSeconds % 60
		minutes := totalSeconds / 60
		length = fmt.Sprintf("%d:%02d", minutes, seconds)
	}

	layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, material.Label(sb.Theme, unit.Sp(16),
				fmt.Sprintf("%s - %s", sb.Chart.Artist, sb.Chart.Name)).Layout,
			),
			layout.Rigid(material.Label(sb.Theme, unit.Sp(16), length).Layout),
		)
	})

	return layout.Dimensions{Size: bounds.Size()}
}
