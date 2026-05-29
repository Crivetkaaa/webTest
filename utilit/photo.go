package utilit

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func ProcessAndSaveAsWebP(fileHeader *multipart.FileHeader, dstDir string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}

	origName := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
	fileName := fmt.Sprintf("%d_%s.webp", time.Now().UnixNano(), origName)
	dstPath := filepath.Join(dstDir, fileName)

	outFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	options := &webp.Options{Lossless: false, Quality: 80}
	if err := webp.Encode(outFile, img, options); err != nil {
		return "", err
	}

	return dstPath, nil
}
