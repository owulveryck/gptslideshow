package ngrok

import (
	"context"
	"image"
	"image/color"
	"image/draw"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)
	return img
}

func TestNewUploader(t *testing.T) {
	uploader := NewUploader()
	assert.NotNil(t, uploader)
	assert.NotNil(t, uploader.imageStore)
}

func TestUpload(t *testing.T) {
	uploader := NewUploader()
	ctx := context.Background()

	name := "test-image"
	img := createTestImage()

	url, err := uploader.Upload(ctx, name, img)
	assert.NoError(t, err)
	assert.Contains(t, url, name)
	assert.Equal(t, img, uploader.imageStore[name])
}

func TestImageHandler(t *testing.T) {
	uploader := NewUploader()

	// Add a test image to the uploader's store
	name := "test-image"
	img := createTestImage()
	uploader.imageStore[name] = img

	req := httptest.NewRequest(http.MethodGet, "/"+name, nil)
	w := httptest.NewRecorder()

	uploader.imageHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/png", resp.Header.Get("Content-Type"))
}

func TestImageHandlerNotFound(t *testing.T) {
	uploader := NewUploader()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent-image", nil)
	w := httptest.NewRecorder()

	uploader.imageHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestImageHandlerBadRequest(t *testing.T) {
	uploader := NewUploader()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	uploader.imageHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestStartServer(t *testing.T) {
	// Note: Testing StartServer requires integration tests and might need
	// mocking for ngrok.Listen. This test is a placeholder.
	t.Skip("StartServer requires integration testing with ngrok")
}
