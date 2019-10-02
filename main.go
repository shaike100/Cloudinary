package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"shai.com/cloudinary/bl"
)

func getThumbnail(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	inputImage := bl.InputImage{Height: keys.Get("height"), Width: keys.Get("width"), URL: keys.Get("url")}

	err := bl.ValidateInput(inputImage)
	if err != nil {
		http.Error(w, err.Error(), err.Code)
		return
	}

	err, b := bl.ProcessImage(inputImage)
	if err != nil {
		http.Error(w, err.Error(), err.Code)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=result.jpeg")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))

	if _, err := w.Write(b); err != nil {
		log.Println("unable to write image.")
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/thumbnail", getThumbnail)
	log.Fatal(http.ListenAndServe(":8080", router))
}
