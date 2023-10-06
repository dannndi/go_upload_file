package file

import (
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/dannndi/go_upload_file/core/utils"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

/// api/v1/file
func Route(router fiber.Router) {

	router.Post("/", func(c *fiber.Ctx) error {
		// parse multiple uploaded
		form, err := c.MultipartForm()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "Failed to parse uploaded files",
			})
		}

		images := form.File["image"]
		if len(images) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "There's no image provided",
			})
		}

		image := images[0]
		compressedImageURL, err := saveImage(image)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "Something wrong when processing the image",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Message: "Saved",
			Data:    compressedImageURL,
		})
	})

	router.Post("/bulk", func(c *fiber.Ctx) error {
		// parse multiple uploaded
		form, err := c.MultipartForm()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "Failed to parse uploaded files",
			})
		}

		// Slice to store compressed image URLs
		var compressedImageURLs []string
		images := form.File["images[]"]
		for _, image := range images {
			compressedImageURL, err := saveImage(image)
			if err != nil {
				continue
			}
			compressedImageURLs = append(compressedImageURLs, compressedImageURL)
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Message: "Saved",
			Data:    compressedImageURLs,
		})
	})
}

func saveImage(image *multipart.FileHeader) (string, error) {

	// Open and decode the uploaded image
	src, err := image.Open()
	if err != nil {
		// Handle the error (e.g., log it) and continue to the next file
		return "", err
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		// Handle the error (e.g., log it) and continue to the next file
		return "", err
	}

	// image with 80% quality
	milis := strconv.FormatInt(time.Now().UnixMilli(), 10)
	path := "./public/uploads/" + milis + "_" + image.Filename
	compressedImage, err := os.Create(path)
	if err != nil {
		// Handle the error
		return "", err
	}
	defer compressedImage.Close()

	err = imaging.Encode(compressedImage, img, imaging.JPEG, imaging.JPEGQuality(80))
	if err != nil {
		// Handle the error
		return "", err

	}

	compressedImageURL := "uploads/" + milis + "_" + image.Filename
	return compressedImageURL, nil
}
