package main

import (
	"bytes"
	"fmt"
	ebimg "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/icza/mjpeg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"
)

func removeBrackets(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "]", ""), "[", "")

}

func stringInGroup(s string, group []string) bool {
	for _, str := range group {
		if str == s {
			return true
		}
	}
	return false
}

func (g *Game) PanChecks() {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.Controls.PanLeft = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.Controls.PanRight = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.Controls.PanUp = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.Controls.PanDown = true
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		g.Controls.PanLeft = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		g.Controls.PanRight = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		g.Controls.PanUp = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		g.Controls.PanDown = false
	}

	if g.Controls.PanLeft {
		g.Map.PanMap([]int{-g.Controls.PanSpeed, 0})
	}
	if g.Controls.PanRight {
		g.Map.PanMap([]int{g.Controls.PanSpeed, 0})
	}
	if g.Controls.PanUp {
		g.Map.PanMap([]int{0, -g.Controls.PanSpeed})
	}
	if g.Controls.PanDown {
		g.Map.PanMap([]int{0, g.Controls.PanSpeed})
	}
}

func (g *Game) ArrowChecks() {

	if g.Map.OcapMapData.CurrentFrame > g.Map.OcapMapData.Data.EndFrame {
		g.Map.OcapMapData.CurrentFrame = g.Map.OcapMapData.Data.EndFrame
		if g.Controls.Playing {
			g.Controls.Playing = false
		}
	}

	if g.Map.OcapMapData.CurrentFrame < 0 {
		g.Map.OcapMapData.CurrentFrame = 0
		if g.Controls.Reversing {
			g.Controls.Reversing = false
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Controls.Playing = !g.Controls.Playing
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.Map.OcapMapData.Data.CaptureDelay = 10
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.Map.OcapMapData.Data.CaptureDelay = 9
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.Map.OcapMapData.Data.CaptureDelay = 8
	}
	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.Map.OcapMapData.Data.CaptureDelay = 7
	}
	if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.Map.OcapMapData.Data.CaptureDelay = 6
	}
	if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.Map.OcapMapData.Data.CaptureDelay = 5
	}
	if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.Map.OcapMapData.Data.CaptureDelay = 4
	}
	if inpututil.IsKeyJustPressed(ebiten.Key8) {
		g.Map.OcapMapData.Data.CaptureDelay = 3
	}
	if inpututil.IsKeyJustPressed(ebiten.Key9) {
		g.Map.OcapMapData.Data.CaptureDelay = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		g.Map.OcapMapData.Data.CaptureDelay = 1
	}

	if g.Controls.Playing {
		if time.Since(g.Map.LastFrame) > time.Duration(100*float64(g.Map.OcapMapData.Data.CaptureDelay))*time.Millisecond {
			g.Map.OcapMapData.CurrentFrame++
			g.Map.LastFrame = time.Now()
		}
	}

	if g.Controls.FastFwd {
		g.Map.OcapMapData.CurrentFrame++
	}

	if g.Controls.Reversing {
		g.Map.OcapMapData.CurrentFrame--
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if !g.Controls.Reversing {
			g.Controls.FastFwd = true
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if !g.Controls.Playing {
			g.Controls.Reversing = true
		}
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
		g.Controls.FastFwd = false
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		g.Controls.Reversing = false
	}

}

func (g *Game) TypingCheck() {
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(g.SearchText) > 0 {
			g.SearchText = g.SearchText[:len(g.SearchText)-1]
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.SearchText += "a"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.SearchText += "b"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.SearchText += "c"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.SearchText += "d"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.SearchText += "e"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.SearchText += "f"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.SearchText += "g"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.SearchText += "h"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.SearchText += "i"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) {
		g.SearchText += "j"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.SearchText += "k"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.SearchText += "l"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.SearchText += "m"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.SearchText += "n"
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		g.SearchText += "o"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.SearchText += "p"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.SearchText += "q"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.SearchText += "r"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.SearchText += "s"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.SearchText += "t"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		g.SearchText += "u"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		g.SearchText += "v"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.SearchText += "w"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.SearchText += "x"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyY) {
		g.SearchText += "y"
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.SearchText += "z"
	}
}
func (g *Game) WindowOptions() {
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.Controls.Names = !g.Controls.Names
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		g.Controls.UnitTotal = !g.Controls.UnitTotal
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) && !g.Controls.Processing {
		g.Controls.Recording = !g.Controls.Recording
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.Controls.FireLines = !g.Controls.FireLines
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		g.Controls.PanSpeed += 10
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		g.Controls.PanSpeed -= 10
		if g.Controls.PanSpeed < 0 {
			g.Controls.PanSpeed = 0
		}
	}
}

func (g *Game) Record() {
	recording := false
	_, err := os.Stat("replays")
	if os.IsNotExist(err) {
		os.Mkdir("replays", 0755)
	}
	for {
		i := 0
		frameMap := make(map[int]string)
		for g.Controls.Recording {
			recording = true
			fileName := fmt.Sprintf("replays/%s_%d.png", g.Map.OcapMapData.Data.MissionName, i)
			frameFile, _ := os.Create(fileName)
			vidFrame := <-g.RecordingImages
			png.Encode(frameFile, vidFrame)
			frameMap[i] = fileName
			i++
		}
		if recording {
			recording = false
			g.Controls.Processing = true
			g.ProcessReplayFrames(frameMap)
			g.Controls.Processing = false
		}
	}
}

func (g *Game) ProcessReplayFrames(frameMap map[int]string) {
	aw, err := mjpeg.New(fmt.Sprintf("%s_%d.avi",
		strings.ReplaceAll(g.Map.OcapMapData.Data.MissionName, " ", "_"), len(frameMap)),
		int32(g.Map.WindowWidth), int32(g.Map.WindowHeight), int32(g.Map.OcapMapData.Data.CaptureDelay))
	if err != nil {
		panic(err)
	}
	defer aw.Close()
	for i := 0; i < len(frameMap); i++ {
		frameFile, _ := os.Open(frameMap[i])
		frame, _, _ := image.Decode(frameFile)
		bufferBytes := new(bytes.Buffer)
		o := jpeg.Options{Quality: 100}

		jpeg.Encode(bufferBytes, frame, &o)
		aw.AddFrame(bufferBytes.Bytes())
	}
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(ttfFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := ebimg.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := ebimg.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := ebimg.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
