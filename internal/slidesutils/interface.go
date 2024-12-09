package slidesutils

import (
	"context"

	"github.com/owulveryck/gptslideshow/internal/structure"
)

// BuilderInterface defines the methods that the Builder struct must implement.
type BuilderInterface interface {
	// CreateChapter creates a chapter with the given slide.
	CreateChapter(ctx context.Context, slide structure.Slide) error

	// CreateSlideTitleSubtitleBody creates a slide with a title, subtitle, and body.
	CreateSlideTitleSubtitleBody(ctx context.Context, slide structure.Slide) error

	// CreateCover creates a cover with the given title and subtitle.
	CreateCover(ctx context.Context, title, subtitle string) error

	// InsertImage inserts an image with the given URL, dimensions, and translation offsets.
	InsertImage(ctx context.Context, imageUrl string, width, height, translateX, translateY float64) error

	// CreateNewSlide creates a new slide with the specified layout.
	CreateNewSlide(ctx context.Context, layout string) error
}
