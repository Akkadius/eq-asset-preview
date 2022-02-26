package main

import (
	"fmt"
	"github.com/adotout/pack_2d"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// spritesheet generator

type config struct {
	directory       string
	iconNamePrefix  string
	iconNameSuffix  string
	spriteImageFile string
	spriteCssFile   string
	spriteSize      uint
}

const testHtmlTemplate = `<html>
<head>
    <link rel="stylesheet" href="{{spritesheet}}">
</head>
<body>
{{sprites}}
</body>
</html>
`

func main() {
	debug := len(os.Getenv("DEBUG")) > 0

	// config
	spriteSheets := []config{
		{
			directory:       "./assets/item_icons",
			iconNamePrefix:  "", // comes from filename already
			spriteImageFile: "./assets/sprites/item-icons.png",
			spriteCssFile:   "./assets/sprites/item-icons.css",
		},
		{
			directory:       "./assets/item_icons",
			iconNamePrefix:  "", // comes from filename already
			iconNameSuffix:  "-sm",
			spriteImageFile: "./assets/sprites/item-icons-sm.png",
			spriteCssFile:   "./assets/sprites/item-icons-sm.css",
			spriteSize:      13, // needs to relatively match font size
		},
		{
			directory:       "./assets/spell_icons",
			iconNamePrefix:  "spell-",
			iconNameSuffix:  "-20",
			spriteImageFile: "./assets/sprites/spell-icons-20.png",
			spriteCssFile:   "./assets/sprites/spell-icons-20.css",
			spriteSize:      20,
		},
		{
			directory:       "./assets/spell_icons",
			iconNamePrefix:  "spell-",
			iconNameSuffix:  "-30",
			spriteImageFile: "./assets/sprites/spell-icons-30.png",
			spriteCssFile:   "./assets/sprites/spell-icons-30.css",
			spriteSize:      30,
		},
		{
			directory:       "./assets/spell_icons",
			iconNamePrefix:  "spell-",
			iconNameSuffix:  "-40",
			spriteImageFile: "./assets/sprites/spell-icons-40.png",
			spriteCssFile:   "./assets/sprites/spell-icons-40.css",
			spriteSize:      40,
		},
		{
			directory:       "./assets/objects",
			iconNamePrefix:  "object-", // comes from filename already
			spriteImageFile: "./assets/sprites/objects.png",
			spriteCssFile:   "./assets/sprites/objects.css",
		},
		{
			directory:       "./assets/spell_icons",
			iconNamePrefix:  "spell-icon-",
			spriteImageFile: "./assets/sprites/spell-icons.png",
			spriteCssFile:   "./assets/sprites/spell-icons.css",
		},
		{
			directory:       "./assets/npc_models",
			iconNamePrefix:  "race-models-",
			spriteImageFile: "./assets/sprites/race-models.png",
			spriteCssFile:   "./assets/sprites/race-models.css",
		},
	}

	// loop through spriteSheets
	for _, c := range spriteSheets {
		dirName := c.directory

		fmt.Printf("[%v] Scanning\n", c.directory)

		dir, err := ioutil.ReadDir(dirName)
		if err != nil {
			panic(err)
		}

		packer := pack_2d.Packer2d{}
		id := 0
		images := map[int]image.Image{}
		imageNames := map[int]string{}
		for _, file := range dir {
			if file.IsDir() {
				continue
			}
			name := file.Name()
			if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".png") && !strings.HasSuffix(
				name,
				".gif",
			) {
				continue
			}
			imgReader, err := os.Open(filepath.Join(dirName, file.Name()))
			if err != nil {
				panic(err)
			}
			imgDecoded, _, err := image.Decode(imgReader)

			// if sprite size specified, we are resizing the image
			if c.spriteSize > 0 {
				imgDecoded = resize.Resize(c.spriteSize, c.spriteSize, imgDecoded, resize.Lanczos3)
			}

			packer.AddNewBlock(imgDecoded.Bounds().Max.X, imgDecoded.Bounds().Max.Y, id)
			images[id] = imgDecoded
			imageNames[id] = file.Name()
			id++

			imgReader.Close()
		}

		fmt.Printf("[%v] Found [%v] images...\n", c.directory, id)

		packedImages, maxWidth, maxHeight := packer.Pack()
		outImage := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))
		css := ""
		html := ""
		for _, img := range packedImages {

			// draw image
			currentImage := images[img.Id]
			mX := currentImage.Bounds().Max.X
			mY := currentImage.Bounds().Max.Y
			draw.Draw(outImage, image.Rect(img.X, img.Y, img.X+mX, img.Y+mY), currentImage, image.ZP, draw.Src)

			// image properties
			imageFileName := imageNames[img.Id]
			imageName := slug(strings.TrimSuffix(imageFileName, filepath.Ext(imageFileName)))
			imageHeight := (img.Y + mY) - img.Y
			imageWidth := (img.X + mX) - img.X

			// print out coordinates for debugging if need be
			if debug {
				fmt.Printf(
					"[%v] X [%v] Y [%v] mX[%v] mY[%v] width [%v] height [%v]\n",
					imageName,
					img.X,
					img.Y,
					img.X+mX,
					img.Y+mY,
					imageWidth,
					imageHeight,
				)
			}

			css += fmt.Sprintf(
				".%v { background: url('./%v') -%vpx -%vpx; height: %vpx; width: %vpx; display: inline-block; }\n",
				fmt.Sprintf("%v%v%v", c.iconNamePrefix, imageName, c.iconNameSuffix),
				filepath.Base(c.spriteImageFile),
				img.X,
				img.Y,
				imageHeight,
				imageWidth,
			)
			html += fmt.Sprintf(
				"<span class=\"%v\"></span>\n",
				fmt.Sprintf("%v%v%v", c.iconNamePrefix, imageName, c.iconNameSuffix),
			)
		}

		// sprite
		imgFile, _ := os.Create(c.spriteImageFile)
		defer imgFile.Close()
		png.Encode(imgFile, outImage)

		fmt.Printf("[%v] Wrote [%v]\n", c.directory, c.spriteImageFile)

		// css
		f, err := os.Create(c.spriteCssFile)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		_, err2 := f.WriteString(css)
		if err2 != nil {
			log.Fatal(err2)
		}
		fmt.Printf("[%v] Wrote [%v]\n", c.directory, c.spriteCssFile)

		// html
		template := testHtmlTemplate
		htmlFilePath := strings.ReplaceAll(c.spriteCssFile, "css", "html")
		sheet := fmt.Sprintf("./%v", filepath.Base(c.spriteCssFile))
		template = strings.ReplaceAll(template, "{{spritesheet}}", sheet)
		template = strings.ReplaceAll(template, "{{sprites}}", html)

		f, err = os.Create(htmlFilePath)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
		_, err2 = f.WriteString(template)
		if err2 != nil {
			log.Fatal(err2)
		}
		fmt.Printf("[%v] Wrote [%v]\n", c.directory, htmlFilePath)
	}

}

var re = regexp.MustCompile("[^a-z0-9]+")

func slug(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
