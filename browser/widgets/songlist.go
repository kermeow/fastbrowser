package widgets

import (
	"fastgh3/fastbrowser/gh"
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type SongList struct {
	Theme         *material.Theme
	Charts        *[]*gh.Chart
	SelectedChart *gh.Chart

	search         *widget.Editor
	displayCharts  []*gh.Chart
	scrollPosition int
	buttons        []*SongButton
}

func NewSongList(theme *material.Theme, charts *[]*gh.Chart) *SongList {
	return &SongList{
		Theme:         theme,
		Charts:        charts,
		SelectedChart: nil,

		search:         &widget.Editor{SingleLine: true},
		displayCharts:  make([]*gh.Chart, 0),
		scrollPosition: 0,
		buttons:        make([]*SongButton, 0),
	}
}

func (sl *SongList) Layout(gtx layout.Context) layout.Dimensions {
	area := gtx.Constraints.Max
	defer clip.Rect{Max: area}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{0, 0, 0, 150})

	if len(*sl.Charts) > 0 {
		gtx.Constraints = layout.Exact(image.Pt(area.X, 32))

		searchBox := material.Editor(sl.Theme, sl.search, "Search for a song")
		searchBox.SelectionColor = color.NRGBA{21, 112, 239, 255}
		inset := layout.UniformInset(unit.Dp(4))
		inset.Layout(gtx, searchBox.Layout)

		searchTerm := strings.ToLower(strings.TrimSpace(sl.search.Text()))
		if len(searchTerm) == 0 {
			sl.displayCharts = (*sl.Charts)[:]
		} else {
			sl.displayCharts = make([]*gh.Chart, 0)
			for _, chart := range *sl.Charts {
				if !fuzzy.MatchNormalized(searchTerm, strings.ToLower(fmt.Sprintf("%s - %s", chart.Artist, chart.Name))) {
					continue
				}
				sl.displayCharts = append(sl.displayCharts, chart)
			}
		}

		defer op.Offset(image.Pt(0, 32)).Push(gtx.Ops).Pop()
		gtx.Constraints = layout.Exact(area.Sub(image.Pt(0, 32)))
	}
	sl.drawCharts(gtx)

	return layout.Dimensions{Size: area}
}

func (sl *SongList) drawCharts(gtx layout.Context) layout.Dimensions {
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
				sl.scrollPosition += x.Scroll.Round().Y
			}
		}
	}

	noCharts := len(sl.displayCharts)
	maxOnScreen := gtx.Constraints.Max.Y / 36

	if sl.scrollPosition < 0 {
		sl.scrollPosition = 0
		gtx.Execute(op.InvalidateCmd{})
	}
	if sl.scrollPosition > noCharts - maxOnScreen {
		sl.scrollPosition = noCharts - maxOnScreen
		gtx.Execute(op.InvalidateCmd{})
	}

	gtx.Constraints.Max.X -= 8
	gtx.Constraints.Max.Y -= 8
	defer op.Offset(image.Pt(4, 4)).Push(gtx.Ops).Pop()

	if noCharts > 0 {
		i := 0
		minI := int(sl.scrollPosition)
		maxI := minI + maxOnScreen + 1
		if maxI > noCharts {
			maxI = noCharts
		}
		for _, chart := range sl.displayCharts[minI:maxI] {
			if len(sl.buttons) < i+1 {
				btn := &SongButton{Chart: chart, Theme: sl.Theme}
				sl.buttons = append(sl.buttons, btn)
			}
			btn := sl.buttons[i]
			btn.Chart = chart
			btn.Layout(gtx, func(sb *SongButton) {
				sl.SelectedChart = sb.Chart
				gtx.Execute(op.InvalidateCmd{})
			})
			defer op.Offset(image.Pt(0, 36)).Push(gtx.Ops).Pop()
			i++
		}
	} else {
		layout.Center.Layout(gtx, material.Label(sl.Theme, unit.Sp(18), "No charts were found!").Layout)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
