package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func redirectToFreshman(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s", r.URL.Host)
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

	err = r.ParseForm()
	if err != nil {
		fmt.Println("Error Parsing Form")
		return
	}

	dateValue := r.PostFormValue("date")

	date, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		fmt.Println(err)
		return
	}

	addDate := date.AddDate(0, 0, 13)
	nolaterDate := date.AddDate(0, 0, 16)
	prcDate := date.AddDate(0, 0, 20)

	formatDate := date.Format("02-01-2006")
	fmtAddDate := addDate.Format("02-01-2006")
	fmtnolaterDate := nolaterDate.Format("02-01-2006")
	fmtprcDate := prcDate.Format("02-01-2006")

	fmt.Printf("Paycycle Start, %s!\n", formatDate)
	fmt.Printf("Paycycle End: %s!\n", fmtAddDate)
	fmt.Printf("No Later than Wednesday: %s!\n", fmtnolaterDate)
	fmt.Printf("Process Date: %s!\n", fmtprcDate)

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

	// fmt.Printf("\n\nUpload successful\n\n")
	// fmt.Printf("\n\nDo stuff Here \n\n")

	// fmt.Printf("\n\nNow delete the same file \n\n")

	e := os.Remove(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if e != nil {
		log.Fatal(e)
	}

	// fmt.Printf("\n\nFile has been deleted \n\n")

	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "success.html")

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}
func closeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, "exit.html")
	// fmt.Println("Server Closed")
	os.Exit(1)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/exit", closeHandler)

	http.ListenAndServe(":8000", mux)

}
