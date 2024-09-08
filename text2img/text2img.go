package text2img

import (
	"bytes"
	_ "embed"
	"image/gif"
	"image/jpeg"
	"io"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed JetBrainsMono-Regular-2.ttf
var fontBytes []byte

type ImageType string

const (
	ImageTypePNG ImageType = "png"
	ImageTypeJPG ImageType = "jpg"
	ImageTypeGIF ImageType = "gif"
)

type Option struct {
	FontSize           float64
	Margin             int
	LineSpacing        float64
	BackgroundHexColor string
	FontHexColor       string
	FontReader         io.ReadCloser
	ImageType          ImageType
}

type Options func(*Option)

func WithFontSize(size float64) Options {
	return func(o *Option) {
		o.FontSize = size
	}
}

func WithMargin(margin int) Options {
	return func(o *Option) {
		o.Margin = margin
	}
}

func WithLineSpacing(spacing float64) Options {
	return func(o *Option) {
		o.LineSpacing = spacing
	}
}

func WithBackgroundHexColor(color string) Options {
	return func(o *Option) {
		o.BackgroundHexColor = color
	}
}

func WithFontHexColor(color string) Options {
	return func(o *Option) {
		o.FontHexColor = color
	}
}

func WithFontReader(reader io.ReadCloser) Options {
	return func(o *Option) {
		o.FontReader = reader
	}
}

func WithImageType(imageType ImageType) Options {
	return func(o *Option) {
		o.ImageType = imageType
	}
}

func defaultOptions() Option {
	return Option{
		FontSize:           18,
		Margin:             20,
		LineSpacing:        1.1,
		BackgroundHexColor: "#ffffff",
		FontHexColor:       "#000000",
		FontReader:         io.NopCloser(bytes.NewReader(fontBytes)),
		ImageType:          "png",
	}
}

func ToImage(text string, writer io.Writer, options ...Options) error {
	opt := defaultOptions()
	for _, o := range options {
		o(&opt)
	}
	return toImage(text, writer, opt)
}

func toImage(text string, writer io.Writer, opt Option) error {
	face, err := loadFontFace(opt.FontSize, opt.FontReader)
	if err != nil {
		return err
	}

	w, h, err := measureText(text, face, opt.LineSpacing)
	if err != nil {
		return err
	}

	width := int(w) + opt.Margin
	height := int(h) + opt.Margin

	img := gg.NewContext(width, height)
	// setBackground
	img.SetHexColor(opt.BackgroundHexColor)
	img.Clear()

	img.SetFontFace(face)
	// set font color
	img.SetHexColor(opt.FontHexColor)

	img.DrawStringWrapped(text, float64(opt.Margin/2), float64(opt.Margin/2), 0, 0, float64(width-opt.Margin), opt.LineSpacing, gg.AlignLeft)

	switch opt.ImageType {
	case ImageTypePNG:
		return img.EncodePNG(writer)
	case ImageTypeJPG:
		return jpeg.Encode(writer, img.Image(), nil)
	case ImageTypeGIF:
		return gif.Encode(writer, img.Image(), nil)
	}

	return nil
}

func measureText(text string, face font.Face, lineSpacing float64) (float64, float64, error) {
	img := gg.NewContext(1, 1)
	img.SetFontFace(face)
	w, h := img.MeasureMultilineString(text, lineSpacing)
	return w, h, nil
}

func loadFontFace(points float64, reader io.ReadCloser) (font.Face, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	f, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
	})
	return face, nil
}
