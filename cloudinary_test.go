package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetThumbnail(t *testing.T) {
	res := httptest.NewRecorder()

	urlKey := "url"
	urlValue := "http://www.pethealthnetwork.com/sites/default/files/cat-seizures-and-epilepsy101.png"

	widthKey := "width"
	widthValue := "50"

	heightKey := "height"
	heightValue := "50"

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/thumbnail?%s=%s&%s=%s&%s=%s", urlKey, urlValue, widthKey, widthValue, heightKey, heightValue), nil)

	getThumbnail(res, req)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if status := res.Code; status != http.StatusOK {
		bodyString := string(bodyBytes)
		t.Errorf("handler returned bad status code: got %v, JSON: %v",
			status, bodyString)
		return
	}

	resultImage, err := jpeg.Decode(bytes.NewReader(bodyBytes))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	width, _ := strconv.Atoi(widthValue)
	height, _ := strconv.Atoi(heightValue)

	if width != resultImage.Bounds().Max.X ||
		height != resultImage.Bounds().Max.Y {
		t.Errorf("result image width or height are not as expected")
	}
}
