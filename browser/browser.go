package browser

import (
	"fastgh3/fastbrowser/browser/widgets"
	"fastgh3/fastbrowser/config"
	"fastgh3/fastbrowser/fonts"
	"fastgh3/fastbrowser/gh"
	"image/color"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

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

	loading     bool
	loadingText string
	songList    *widgets.SongList
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

	browser := &Browser{
		Window: window,

		Config: conf,
		Theme:  &theme,
		Charts: nil,

		loading: true,
	}
	browser.songList = widgets.NewSongList(&theme, &browser.Charts)
	return browser
}

func (ui *Browser) Run() error {
	// ui.songList.Charts = &ui.Charts
	go ui.getCharts()

	var ops op.Ops
	for {
		switch e := ui.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, ui.Theme.Bg)

			if ui.Charts == nil || ui.loading {
				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(gtx,
						layout.Rigid(material.H2(ui.Theme, "Scanning charts").Layout),
						layout.Rigid(material.Label(ui.Theme, unit.Sp(18), ui.loadingText).Layout),
					)
				})
			} else {
				inset := layout.UniformInset(unit.Dp(16))
				inset.Layout(gtx, ui.draw)
			}

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
	charts := make([]*gh.Chart, 0)

	{ // Scanning charts
		startTime := time.Now()
		for _, root := range ui.Config.SearchDirs {
			root = filepath.Clean(root)
			if fi, err := os.Stat(root); os.IsNotExist(err) || !fi.Mode().IsDir() {
				log.Printf("Search dir '%s' doesn't exist", root)
				continue
			}
			filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}
				// relpath, _ := filepath.Rel(root, path)
				ui.loadingText = filepath.ToSlash(path)
				ui.Invalidate()
				name := strings.ToLower(d.Name())
				ext := filepath.Ext(name)
				nameNoExt := strings.TrimSuffix(name, ext)
				// time.Sleep(time.Millisecond * 100)
				if nameNoExt == "song" ||
					nameNoExt == "album" ||
					ext == ".chart" ||
					ext == ".mid" {
					// WE HAVE A WINNER
					chart, err := gh.ReadChart(filepath.Dir(path))
					if err != nil {
						return nil
					}
					charts = append(charts, chart)
					return filepath.SkipDir
				}
				return nil
			})
		}
		finishTime := time.Now()
		log.Printf("Read %d charts in %.3fms", len(charts), float32(finishTime.Sub(startTime).Microseconds()) / 1000)
	}

	ui.loadingText = "Sorting by title"
	ui.Invalidate()

	{ // Sorting charts
		startTime := time.Now()
		slices.SortStableFunc(charts, func(a, b *gh.Chart) int {
			return strings.Compare(a.Name, b.Name)
		})
		finishTime := time.Now()
		log.Printf("Sorted %d charts in %.3fms", len(charts), float32(finishTime.Sub(startTime).Microseconds()) / 1000)
	}

	ui.Charts = charts
	ui.loading = false
	ui.Invalidate()
}
