package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

type imageInfo struct {
	path  string
	image *image.NRGBA
	err   error
}

// Image processing - sequential
// Input - directory with images.
// output - thumbnail images
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to send directory path of images")
	}
	start := time.Now()

	err := setupPipeline(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

func setupPipeline(root string) error {
	done := make(chan any)
	defer close(done)

	path_ch, error_ch := walkFiles(done, root)
	thumbnails_ch := processImage(done, path_ch)

	for r := range thumbnails_ch {
		if r.err != nil {
			return r.err
		}
		saveThumbnail(r.path, r.image)
	}

	if err := <-error_ch; err != nil {
		return err
	}

	return nil
}

// walfiles - take diretory path as input
// does the file walk
// generates thumbnail images
// saves the image to thumbnail directory.
func walkFiles(done <-chan any, root string) (<-chan string, <-chan error) {
	out_ch := make(chan string)
	err_ch := make(chan error, 1)
	go func() {
		defer close(out_ch)
		err_ch <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			// filter out error
			if err != nil {
				return err
			}

			// check if it is file
			if !info.Mode().IsRegular() {
				return nil
			}

			// check if it is image/jpeg
			contentType, _ := getFileContentType(path)
			if contentType != "image/jpeg" {
				return nil
			}

			select {
			case out_ch <- path:
			case <-done:
				return fmt.Errorf("walk was canceled")
			}

			// // process the image
			// thumbnailImage, err := processImage(path)
			// if err != nil {
			// 	return err
			// }

			// // save the thumbnail image to disk
			// err = saveThumbnail(path, thumbnailImage)
			// if err != nil {
			// 	return err
			// }
			return nil
		})
	}()
	return out_ch, err_ch
}

// processImage - takes image file as input
// return pointer to thumbnail image in memory.
func processImage(done <-chan any, path_ch <-chan string) <-chan *imageInfo {
	out_ch := make(chan *imageInfo)
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			for path := range path_ch {
				imageInfo := &imageInfo{
					path:  path,
					image: nil,
					err:   nil,
				}
				// load the image from file
				srcImage, err := imaging.Open(imageInfo.path)
				if err != nil {
					imageInfo.err = err
					select {
					case out_ch <- imageInfo:
						continue
					case <-done:
						return
					}
				}

				// scale the image to 100px * 100px
				thumbnailImage := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)

				imageInfo.image = thumbnailImage
				select {
				case out_ch <- imageInfo:
				case <-done:
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out_ch)
	}()
	return out_ch
}

// saveThumbnail - save the thumnail image to folder
func saveThumbnail(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	// save the image in the thumbnail folder.
	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

// getFileContentType - return content type and error status
func getFileContentType(file string) (string, error) {

	out, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err = out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
