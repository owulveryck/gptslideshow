package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strings"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

var imageStore map[string]image.Image

func init() {
	// Initialize the image store with some sample images
	imageStore = make(map[string]image.Image)

	// Example image: a red rectangle
	redImage := image.NewRGBA(image.Rect(0, 0, 200, 100))
	draw.Draw(redImage, redImage.Bounds(), &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)
	imageStore["red.png"] = redImage

	// Example image: a green rectangle
	greenImage := image.NewRGBA(image.Rect(0, 0, 200, 100))
	draw.Draw(greenImage, greenImage.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.Point{}, draw.Src)
	imageStore["green.png"] = greenImage
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the image name from the URL path
	imageName := strings.TrimLeft(r.URL.Path, "/")
	log.Println(imageName)
	if imageName == "" {
		http.Error(w, "Image name is required", http.StatusBadRequest)
		return
	}

	log.Println(imageStore)
	// Find the image in the store
	img, exists := imageStore[imageName]
	if !exists {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Set the content type to image/png
	w.Header().Set("Content-Type", "image/png")

	// Encode the image as PNG and write to the response
	err := png.Encode(w, img)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
}

var tunURL string

func startServer() {
	fmt.Println("starting server")
	// Start Ngrok session
	ctx := context.Background()
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		log.Fatal("Error starting Ngrok tunnel:", err)
	}
	tunURL = listener.URL()

	// Print the Ngrok public URL
	fmt.Printf("Ngrok tunnel created at: %s\n", listener.URL())

	// Set up the HTTP server
	http.Serve(listener, http.HandlerFunc(imageHandler))
}
