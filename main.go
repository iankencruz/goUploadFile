package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// truncated for brevity

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")
	fmt.Printf("\n\nDo stuff Here \n\n")

	fmt.Printf("\n\nNow delete the same file \n\n")

	e := os.Remove(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf("\n\nFile has been deleted \n\n")

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)

	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
