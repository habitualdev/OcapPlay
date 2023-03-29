package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ebitenui/ebitenui"
	ebimg "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/ncruces/zenity"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"ocapPlay/SaveFile"
	"os"
	"time"
)

var ocapBytes []byte
var newData OcapData
var GoodLoad bool

func main() {
	GoodLoad = false
	jsonFile := flag.String("json", "", "Path to the json file")
	command := flag.String("command", "", "Command to run")
	flag.Parse()

	if *command == "stats" {
		GenerateStats(*jsonFile)
		os.Exit(0)
	}

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		fmt.Println(err)
	}
	ttFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		fmt.Println(err)
	}

	ttFace20, _ := loadFont(20)

	// add the button as a child of the container

	ebitenMap := NewMap(1024, 1024, newData)
	newGame := &Game{
		Map:            *ebitenMap,
		SearchText:     "",
		LoadScreen:     true,
		LastMapCommand: time.Now(),
		Controls: struct {
			FireLines  bool
			Processing bool
			PanSpeed   int
			Recording  bool
			Playing    bool
			PanRight   bool
			PanLeft    bool
			PanUp      bool
			PanDown    bool
			Reversing  bool
			FastFwd    bool
			Names      bool
			UnitTotal  bool
		}{PanSpeed: 100},
		RecordingImages: make(chan image.Image, 1000),
		Font:            ttFace,
		LargeFont:       ttFace20,
		lastFrameTick:   time.Now(),
	}

	newGame.Map.SizeX = 20480
	newGame.Map.SizeY = 20480
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(ebimg.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)
	// construct a button

	newGame.ui = &ebitenui.UI{
		Container: rootContainer,
	}

	buttonImage, _ := loadButtonImage()

	mapButton := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Select Map Directory", ttFace20, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			zipFile, _ := zenity.SelectFile(zenity.Filename("Select OCAP Play Data"), zenity.FileFilter{
				Name:     "Ocap ZIP",
				Patterns: []string{"*.ocapzip"},
				CaseFold: false,
			})
			if zipFile != "" {
				zipReader, err := SaveFile.ProcessSave(zipFile)
				if err != nil {
					fmt.Println(err)
					return
				}
				newGame.OcapSave = zipReader

				ocapFile, err := zipReader.Open("ocap.json")
				if err != nil {
					fmt.Println(err)
					return
				}
				ocapBytes, err = io.ReadAll(ocapFile)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = json.Unmarshal(ocapBytes, &newData)
				if err != nil {
					fmt.Println(err)
					return
				}
				newGame.Map.OcapMapData.Data = newData
				newGame.Map.OcapMapData.Frames = newData.EndFrame
				newGame.Map.OcapMapData.CurrentFrame = 0
				newGame.Map.PlayBackSpeed = int(newData.CaptureDelay)
				GoodLoad = true
			}
		}),
	)

	// add the button as a child of the container
	rootContainer.AddChild(mapButton)

	mapToggleButton := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Toggle to Map", ttFace20, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if GoodLoad {
				newGame.LoadScreen = false
			}
		}),
	)

	rootContainer.AddChild(mapToggleButton)

	ebiten.SetWindowSize(1024, 1024)
	ebiten.SetWindowTitle("OCAP Play")
	go newGame.Record()

	if err := ebiten.RunGame(newGame); err != nil {
		log.Fatal(err)
	}
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
