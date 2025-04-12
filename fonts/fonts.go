package fonts

import (
	"fmt"
	"nylenol/bootstrapper/fonts/nunito"
	"nylenol/bootstrapper/fonts/ubuntumono"
	"sync"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
)

var (
	once       sync.Once
	collection []text.FontFace
)

func Collection() []font.FontFace {
	once.Do(func() {
		// Nunito
		register(font.Font{Weight: font.Normal, Style: font.Regular}, "Nunito", nunito.Regular)
		// register(font.Font{Weight: font.Normal, Style: font.Italic}, "Nunito", nunito.Italic)
		// register(font.Font{Weight: font.Bold, Style: font.Regular}, "Nunito", nunito.Bold)
		// register(font.Font{Weight: font.Bold, Style: font.Italic}, "Nunito", nunito.BoldItalic)
		// Ubuntu Mono
		register(font.Font{Weight: font.Normal, Style: font.Regular}, "UbuntuMono", ubuntumono.Regular)
		// register(font.Font{Weight: font.Normal, Style: font.Italic}, "UbuntuMono", ubuntumono.Italic)
		// register(font.Font{Weight: font.Bold, Style: font.Regular}, "UbuntuMono", ubuntumono.Bold)
		// register(font.Font{Weight: font.Bold, Style: font.Italic}, "UbuntuMono", ubuntumono.BoldItalic)
	})
	return collection
}

func register(fnt font.Font, typeface string, data []byte) {
	face, err := opentype.Parse(data)
	if err != nil {
		panic(fmt.Errorf("failed to parse font: %v", err))
	}
	fnt.Typeface = font.Typeface(typeface)
	collection = append(collection, font.FontFace{Font: fnt, Face: face})
}
