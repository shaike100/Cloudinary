package bl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nfnt/resize"
)

// InputImage - A struct that holds the input parameters
type InputImage struct {
	Height string
	Width  string
	URL    string
}

// ValidationError - Custom validation struct for the JSON in the http error response
type ValidationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Implements Error interface, return a JSON as error message
func (ve *ValidationError) Error() string {
	buf, _ := json.Marshal(ve)
	return string(buf)
}

// Register - Needed for registering the image formats in order to support the various image types
func Register() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
}

// ValidateInput - Validates the given query parameters
func ValidateInput(inputImage InputImage) *ValidationError {
	if len(inputImage.URL) < 1 ||
		len(inputImage.Width) < 1 ||
		len(inputImage.Height) < 1 {
		return &ValidationError{Code: http.StatusBadRequest, Message: "input parameter missing"}
	}

	wi, err := strconv.Atoi(inputImage.Width)
	if err != nil {
		return &ValidationError{Code: http.StatusBadRequest, Message: "width is not a number"}
	}

	if wi <= 0 {
		return &ValidationError{Code: http.StatusBadRequest, Message: "width must be greater than 0"}
	}

	he, err := strconv.Atoi(inputImage.Height)
	if err != nil {
		return &ValidationError{Code: http.StatusBadRequest, Message: "height is not a number"}
	}

	if he <= 0 {
		return &ValidationError{Code: http.StatusBadRequest, Message: "height must be greater than 0"}
	}

	_, err = url.ParseRequestURI(inputImage.URL)
	if err != nil {
		return &ValidationError{Code: http.StatusBadRequest, Message: fmt.Sprintf("Bad URL format: %s", inputImage.URL)}
	}

	return nil
}

// ProcessImage - Image proccessing
func ProcessImage(inputImage InputImage) (*ValidationError, []byte) {
	resp, err := http.Get(inputImage.URL)
	if err != nil {
		return &ValidationError{Code: http.StatusNotFound, Message: "Unable to get image from URL"}, nil
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	originalImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return &ValidationError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	width, _ := strconv.Atoi(inputImage.Width)
	height, _ := strconv.Atoi(inputImage.Height)

	// If the requested width and height are the same as the original, we just jpeg encode
	if width == originalImage.Bounds().Max.X && height == originalImage.Bounds().Max.Y {
		buffer := new(bytes.Buffer)
		e := jpeg.Encode(buffer, originalImage, nil)
		if e != nil {
			return &ValidationError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}

		return nil, buffer.Bytes()
	}

	var backgroundImage image.Image
	var yOffset, xOffset = 0, 0
	var finalImage image.Image = nil

	// If one of the given size parameters are smaller than the original then we need scale down and maintain the aspect ratio
	if width < originalImage.Bounds().Max.X || height < originalImage.Bounds().Max.Y {
		thumbnailImage := resize.Thumbnail(uint(width), uint(height), originalImage, resize.Lanczos3) // Scale down and maintain aspect ratio
		finalImage = thumbnailImage
	} else { // One of the given size parameters is larger than the original so no need to scale up
		finalImage = originalImage
	}

	if width > finalImage.Bounds().Max.X {
		xOffset = (width - finalImage.Bounds().Max.X) / 2
	}
	if height > finalImage.Bounds().Max.Y {
		yOffset = (height - finalImage.Bounds().Max.Y) / 2
	}

	backgroundImage = CreateBackground(width, height)

	offset := image.Pt(xOffset, yOffset)
	b := backgroundImage.Bounds()
	outputImage := image.NewRGBA(b)
	point := image.ZP

	draw.Draw(outputImage, b, backgroundImage, point, draw.Src)
	draw.Draw(outputImage, finalImage.Bounds().Add(offset), finalImage, point, draw.Over)

	buffer := new(bytes.Buffer)
	e := jpeg.Encode(buffer, outputImage, nil)
	if e != nil {
		return &ValidationError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}

	return nil, buffer.Bytes()
}

// CreateBackground - Create the black background image
func CreateBackground(width, height int) image.Image {
	backgroundImage := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			backgroundImage.Set(x, y, color.Black)
		}
	}

	return backgroundImage
}
