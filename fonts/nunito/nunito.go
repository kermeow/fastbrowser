package nunito

import _ "embed"

//go:embed Nunito-Regular.ttf
var Regular []byte

//go:embed Nunito-Italic.ttf
var Italic []byte

//go:embed Nunito-Bold.ttf
var Bold []byte

//go:embed Nunito-BoldItalic.ttf
var BoldItalic []byte
