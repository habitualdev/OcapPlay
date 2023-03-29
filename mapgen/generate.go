package mapgen

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
)

func ParseMap(dir string) (image.Image, error) {
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, err
	}
	dirDirs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	totalLength := len(dirDirs)
	fmt.Println(totalLength)
	tileCollection := map[int]image.RGBA{}
	for i := 0; i < totalLength; i++ {

		tileBytes := []image.Image{}

		for i2 := 0; i2 < totalLength; i2++ {
			pngBytes, err := os.ReadFile(dir + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(i2) + ".png")
			if err != nil {
				return nil, err
			}
			pngTile, _ := png.Decode(bytes.NewReader(pngBytes))
			tileBytes = append(tileBytes, pngTile)
		}
		subTileSize := tileBytes[0].Bounds().Size()
		subTileSum := image.NewRGBA(image.Rect(0, 0, subTileSize.X, subTileSize.Y*int(totalLength)))
		//Build the subtile in a single vertical column
		for I, tile := range tileBytes {
			for x := 0; x < subTileSize.X; x++ {
				for y := 0; y < subTileSize.Y; y++ {
					subTileSum.Set(x, y+I*subTileSize.Y, tile.At(x, y))
				}
			}
		}

		tileCollection[i] = *subTileSum
	}

	//Build the final map from the subtile cloumns in a single horizontal row
	testTile := tileCollection[0]
	tileSize := testTile.Bounds().Size()

	finalMap := image.NewRGBA(image.Rect(0, 0, tileSize.X*int(totalLength), tileSize.Y))
	finalCounter := 0
	for range tileCollection {
		for x := 0; x < tileSize.X; x++ {
			for y := 0; y < tileSize.Y; y++ {
				curTile := tileCollection[finalCounter]
				finalMap.Set(x+finalCounter*tileSize.X, y, curTile.At(x, y))
			}
		}
		finalCounter++
	}

	finalMapFile, err := os.Create("map.png")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = png.Encode(finalMapFile, finalMap)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return finalMap, nil
}
