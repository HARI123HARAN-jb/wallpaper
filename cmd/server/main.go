package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"wallgen/internal/processor"

	"github.com/disintegration/imaging"
)

const (
	uploadDir = "uploads"
	outputDir = "web/static/generated"
)

func main() {
	// Ensure directories exist
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(outputDir, 0755)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/healthz", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server started at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save uploaded file
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// Process image
	resolutionsStr := r.FormValue("resolutions")
	if resolutionsStr == "" {
		resolutionsStr = "1366x768,1920x1080,2560x1440,3840x2160,1080x2400,1440x3200"
	}

	resolutions, err := processor.ParseResolutions(resolutionsStr)
	if err != nil {
		http.Error(w, "Invalid resolutions", http.StatusBadRequest)
		return
	}

	// Create a unique output folder for this request
	reqID := fmt.Sprintf("%d", time.Now().UnixNano())
	reqOutputDir := filepath.Join(outputDir, reqID)
	os.MkdirAll(reqOutputDir, 0755)

	srcImg, err := imaging.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open image", http.StatusInternalServerError)
		return
	}

	type Result struct {
		Resolution string
		Path       string
	}
	var results []Result

	for _, res := range resolutions {
		err := processor.ResizeImage(srcImg, res.Width, res.Height, reqOutputDir, "jpg", 95)
		if err == nil {
			results = append(results, Result{
				Resolution: fmt.Sprintf("%dx%d", res.Width, res.Height),
				Path:       fmt.Sprintf("/static/generated/%s/wallpaper_%dx%d.jpg", reqID, res.Width, res.Height),
			})
		}
	}

	// Return HTML fragment with results (HTMX style or just JSON)
	// For simplicity, we'll return a JSON response and handle it in JS
	// But to keep it simple with vanilla JS, let's return JSON.
	w.Header().Set("Content-Type", "application/json")
	// Simple JSON construction
	jsonStr := `{"success": true, "images": [`
	for i, res := range results {
		if i > 0 {
			jsonStr += ","
		}
		jsonStr += fmt.Sprintf(`{"res": "%s", "url": "%s"}`, res.Resolution, res.Path)
	}
	jsonStr += `]}`
	w.Write([]byte(jsonStr))
}
