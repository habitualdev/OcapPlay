package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/png"
	"strings"
	"time"
)

//go:embed resources/blueCircle.png
var blueCircle []byte

//go:embed resources/redCircle.png
var redCircle []byte

//go:embed resources/whiteCircle.png
var whiteCircle []byte

//go:embed resources/blackCircle.png
var blackCircle []byte

//go:embed resources/unitBackground.png
var unitBackground []byte

var mapBaseDir string
var mapSizeGlobal int

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type Game struct {
	Map             Map
	OcapSave        *zip.ReadCloser
	LastMapCommand  time.Time
	ui              *ebitenui.UI
	MapDirInput     *widget.TextInput
	JsonFileInput   *widget.TextInput
	count           int
	LoadScreen      bool
	RecordingImages chan image.Image
	Font            font.Face
	LargeFont       font.Face
	DrawingOption   *ebiten.DrawImageOptions
	SearchText      string
	TypeMode        bool
	Controls        struct {
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
	}
	lastFrameTick time.Time
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		if g.OcapSave != nil {
			g.LoadScreen = !g.LoadScreen
		}
	}
	if g.LoadScreen {
		g.ui.Update()
		return nil
	}
	if !g.TypeMode {
		g.PanChecks()
		g.ArrowChecks()
		g.WindowOptions()

		if g.Controls.Recording {
			ebiten.SetWindowTitle("OCAP Play - Recording")
		} else {
			ebiten.SetWindowTitle("OCAP Play")
		}
	} else {
		g.TypingCheck()
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			for _, entity := range g.Map.OcapMapData.Data.Entities {
				if strings.Contains(strings.ToLower(entity.Name), strings.ToLower(g.SearchText)) {
					posData := entity.Positions[g.Map.OcapMapData.CurrentFrame-entity.StartFrameNum]
					pos := posData[0].([]interface{})
					y := pos[1].(float64)
					x := pos[0].(float64)
					fmt.Println("Panning to", x, y)
					g.Map.PanMap([]int{g.Map.OcapMapData.WindowTopLeftX + int(x) - g.Map.WindowWidth/2,
						(g.Map.SizeY - int(y)) + g.Map.OcapMapData.WindowTopLeftY - g.Map.WindowHeight/2})
				}
			}
			g.TypeMode = false
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if time.Since(g.LastMapCommand) > 200*time.Millisecond {
			g.TypeMode = !g.TypeMode
			g.SearchText = ""
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.LoadScreen {
		g.ui.Draw(screen)
	} else {
		screen.Clear()

		numTilesY := g.Map.WindowHeight / 256
		numTilesX := g.Map.WindowWidth / 256
		numTilesY++
		numTilesX++

		topLeftTile := image.Point{g.Map.OcapMapData.WindowTopLeftX / 256,
			g.Map.OcapMapData.WindowTopLeftY / 256}
		tileOffset := image.Point{g.Map.OcapMapData.WindowTopLeftX % 256,
			g.Map.OcapMapData.WindowTopLeftY % 256}
		bottomRightTile := image.Point{topLeftTile.X + numTilesX, topLeftTile.Y + numTilesY}

		for y := topLeftTile.Y; y < bottomRightTile.Y; y++ {
			for x := topLeftTile.X; x < bottomRightTile.X; x++ {
				tileFile, err := g.OcapSave.Open(fmt.Sprintf("tiles/%d/%d.png", x, y))
				if err == nil {
					tileImage, _, _ := image.Decode(tileFile)
					if tileImage.Bounds().Dx() != 256 || tileImage.Bounds().Dy() != 256 {
						tileImage = resize.Resize(256, 256, tileImage, resize.NearestNeighbor)
					}
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64((x-topLeftTile.X)*256-tileOffset.X), float64((y-topLeftTile.Y)*256-tileOffset.Y))
					screen.DrawImage(ebiten.NewImageFromImage(tileImage), op)
				}

			}
		}

		// Draw the background map, cutting off the edges that are outside the window
		op := &ebiten.DrawImageOptions{}

		onScreenUnits := []string{}
		for _, entity := range g.Map.OcapMapData.Data.Entities {
			if entity.StartFrameNum > g.Map.OcapMapData.CurrentFrame {
				continue
			}

			if len(entity.Positions) < g.Map.OcapMapData.CurrentFrame-entity.StartFrameNum+1 {
				continue
			}

			posData := entity.Positions[g.Map.OcapMapData.CurrentFrame-entity.StartFrameNum]
			if len(posData) < 3 {
				continue
			}
			pos := posData[0].([]interface{})
			y := pos[1].(float64)
			x := pos[0].(float64)
			point := image.Point{
				X: int(x),
				Y: g.Map.OcapMapData.MapSize - int(y),
			}

			if posData[2] == float64(0) {
				continue
			}

			if point.X < g.Map.OcapMapData.WindowTopLeftX || point.X > g.Map.OcapMapData.WindowTopLeftX+g.Map.WindowWidth {
				continue
			}
			if point.Y < g.Map.OcapMapData.WindowTopLeftY || point.Y > g.Map.OcapMapData.WindowTopLeftY+g.Map.WindowHeight {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(point.X-g.Map.OcapMapData.WindowTopLeftX), float64(point.Y-g.Map.OcapMapData.WindowTopLeftY))
			if entity.Type == "unit" {
				if entity.Side == "EAST" {
					screen.DrawImage(g.Map.Resources["redCircle"], op)
					for _, shot := range entity.FramesFired {
						if int(shot[0].(float64)) == g.Map.OcapMapData.CurrentFrame {
							shotPos := shot[1].([]interface{})
							shotPoint := image.Point{
								X: int(shotPos[0].(float64)),
								Y: g.Map.OcapMapData.MapSize - int(shotPos[1].(float64)),
							}
							if g.Controls.FireLines {
								vector.StrokeLine(screen,
									float32(point.X-g.Map.OcapMapData.WindowTopLeftX),
									float32(point.Y-g.Map.OcapMapData.WindowTopLeftY),
									float32(shotPoint.X-g.Map.OcapMapData.WindowTopLeftX),
									float32(shotPoint.Y-g.Map.OcapMapData.WindowTopLeftY), 1,
									color.RGBA{255, 0, 0, 255}, true)
							}
						}
					}
				} else if entity.Side == "WEST" {
					screen.DrawImage(g.Map.Resources["blueCircle"], op)
					if g.Controls.Names {
						text.Draw(screen, entity.Name, g.Font, point.X-g.Map.OcapMapData.WindowTopLeftX,
							(point.Y - g.Map.OcapMapData.WindowTopLeftY - 2), color.Black)
					}
					for _, shot := range entity.FramesFired {
						if int(shot[0].(float64)) == g.Map.OcapMapData.CurrentFrame {
							shotPos := shot[1].([]interface{})
							shotPoint := image.Point{
								X: int(shotPos[0].(float64)),
								Y: g.Map.OcapMapData.MapSize - int(shotPos[1].(float64)),
							}
							if g.Controls.FireLines {
								vector.StrokeLine(screen,
									float32(point.X-g.Map.OcapMapData.WindowTopLeftX),
									float32(point.Y-g.Map.OcapMapData.WindowTopLeftY),
									float32(shotPoint.X-g.Map.OcapMapData.WindowTopLeftX),
									float32(shotPoint.Y-g.Map.OcapMapData.WindowTopLeftY), 1,
									color.RGBA{0, 0, 255, 255}, true)
							}
						}
					}
				} else if entity.Side == "CIV" {
					drawCircle(g.Map.Resources["whiteCircle"], 3, 3, 3, color.RGBA{255, 0, 255, 255})
					screen.DrawImage(g.Map.Resources["whiteCircle"], op)
				}
				onScreenUnits = append(onScreenUnits, entity.Side+" - "+entity.Name)
			} else if entity.Type == "vehicle" {
				op.Filter = ebiten.FilterNearest
				text.Draw(screen, entity.Name, g.Font, point.X-g.Map.OcapMapData.WindowTopLeftX, point.Y-g.Map.OcapMapData.WindowTopLeftY, color.Black)
				screen.DrawImage(g.Map.Resources["blackCircle"], op)

			}
		}
		if g.Controls.UnitTotal {
			if len(onScreenUnits) > 0 {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(g.Map.WindowWidth-100), 0)
				screen.DrawImage(g.Map.Resources["unitBackground"], op)
			}
			for i, unit := range onScreenUnits {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(g.Map.WindowWidth-100), float64(10+i*10))
				screen.DrawImage(g.Map.Resources["unitBackground"], op)
				text.Draw(screen, unit, g.Font, g.Map.WindowWidth-100, 10+i*10, color.White)
			}
		}
		if g.Controls.Recording {
			recImage := image.NewRGBA(image.Rect(0, 0, g.Map.WindowWidth, g.Map.WindowHeight))
			screenBounds := screen.Bounds()
			for x := screenBounds.Min.X; x < screenBounds.Max.X; x++ {
				for y := screenBounds.Min.Y; y < screenBounds.Max.Y; y++ {
					recImage.Set(x, y, screen.At(x, y))
				}
			}
			g.RecordingImages <- recImage
		}

		if g.TypeMode {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(g.Map.WindowWidth/2), 50)
			screen.DrawImage(g.Map.Resources["unitBackground"], op)
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(g.Map.WindowWidth/2), 40)
			screen.DrawImage(g.Map.Resources["unitBackground"], op)
			text.Draw(screen, g.SearchText, g.Font, g.Map.WindowWidth/2, 50, color.White)
		}

		if !g.Controls.Recording {
			text.Draw(screen, fmt.Sprint("N : Toggle Names -- L: Toggle Fire Lines -- U: Unit's in View  -- R: Toggle Recording -- ESC: Unit Search"), g.LargeFont, 10, g.Map.WindowHeight-30, color.Black)
			text.Draw(screen, fmt.Sprint("WASD: Pan Map -- Left and Right Arrow Keys: Rewind/Fast Forward -- PgDwn/PgUp: Change Pan Speed"), g.LargeFont, 10, g.Map.WindowHeight-10, color.Black)
		}
		ebitenutil.DebugPrint(screen, fmt.Sprint(g.Map.OcapMapData.CurrentFrame, "\n", g.Map.OcapMapData.WindowTopLeftX, g.Map.OcapMapData.WindowTopLeftY))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.Map.WindowWidth = outsideWidth
	g.Map.WindowHeight = outsideHeight
	return outsideWidth, outsideHeight
}

type Map struct {
	Resources     map[string]*ebiten.Image
	LastFrame     time.Time
	WindowWidth   int
	WindowHeight  int
	SizeX         int
	SizeY         int
	PlayBackSpeed int
	OcapMapData   OcapMapData
}

type OcapMapData struct {
	Frames         int
	Data           OcapData
	CurrentFrame   int
	WindowTopLeftX int
	WindowTopLeftY int
	MapSize        int
	SizeX          int
	SizeY          int
}

func NewMap(windowWidth, windowHeight int, data OcapData) *Map {
	redCircleImg, err := png.Decode(bytes.NewReader(redCircle))
	if err != nil {
		return nil
	}
	blueCircleImg, err := png.Decode(bytes.NewReader(blueCircle))
	if err != nil {
		return nil
	}

	whiteCircleImg, err := png.Decode(bytes.NewReader(whiteCircle))
	if err != nil {
		return nil
	}

	blackCircleImg, err := png.Decode(bytes.NewReader(blackCircle))
	if err != nil {
		return nil
	}

	unitBackgroundImg, err := png.Decode(bytes.NewReader(unitBackground))
	if err != nil {
		return nil
	}

	ebitenRedCircle := ebiten.NewImageFromImage(redCircleImg)
	ebitenBlueCircle := ebiten.NewImageFromImage(blueCircleImg)
	ebitenWhiteCircle := ebiten.NewImageFromImage(whiteCircleImg)
	ebitenBlackCircle := ebiten.NewImageFromImage(blackCircleImg)
	ebitenUnitBackground := ebiten.NewImageFromImage(unitBackgroundImg)

	return &Map{
		Resources: map[string]*ebiten.Image{
			"redCircle":      ebitenRedCircle,
			"blueCircle":     ebitenBlueCircle,
			"whiteCircle":    ebitenWhiteCircle,
			"blackCircle":    ebitenBlackCircle,
			"unitBackground": ebitenUnitBackground,
		},
		LastFrame:     time.Now(),
		WindowWidth:   windowWidth,
		WindowHeight:  windowHeight,
		SizeX:         mapSizeGlobal,
		SizeY:         mapSizeGlobal,
		PlayBackSpeed: int(data.CaptureDelay),
		OcapMapData: OcapMapData{
			Frames:         data.EndFrame,
			Data:           data,
			CurrentFrame:   1,
			MapSize:        20480,
			SizeX:          windowWidth,
			SizeY:          windowHeight,
			WindowTopLeftX: 0,
			WindowTopLeftY: 0,
		},
	}
}

func (m *Map) PanMap(panCoords []int) {
	m.OcapMapData.WindowTopLeftX += panCoords[0]
	m.OcapMapData.WindowTopLeftY += panCoords[1]
	if m.OcapMapData.WindowTopLeftX < 0 {
		m.OcapMapData.WindowTopLeftX = 0
	}
	if m.OcapMapData.WindowTopLeftY < 0 {
		m.OcapMapData.WindowTopLeftY = 0
	}

	if m.OcapMapData.WindowTopLeftX > m.SizeX-m.WindowWidth {
		m.OcapMapData.WindowTopLeftX = m.SizeX - m.WindowWidth
	}
	if m.OcapMapData.WindowTopLeftY > m.SizeY-m.WindowHeight {
		m.OcapMapData.WindowTopLeftY = m.SizeY - m.WindowHeight
	}
}
