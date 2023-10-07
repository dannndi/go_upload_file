package file

import (
	"strconv"
	"time"

	"github.com/dannndi/go_upload_file/core/utils"
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

		files := form.File["file"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "There's no file provided",
			})
		}

		file := files[0]
		milis := strconv.FormatInt(time.Now().UnixMilli(), 10)
		filePath := "./public/uploads/" + milis + "_" + file.Filename
		fileUrl := "uploads/" + milis + "_" + file.Filename
		err = c.SaveFile(file, filePath)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "Something wrong when saving the file",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Code:    fiber.StatusCreated,
			Message: "Saved",
			Data:    fileUrl,
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
		var compressedFileURLs []string
		files := form.File["files[]"]
		for _, file := range files {
			milis := strconv.FormatInt(time.Now().UnixMilli(), 10)
			filePath := "./public/uploads/" + milis + "_" + file.Filename
			fileUrl := "uploads/" + milis + "_" + file.Filename
			err = c.SaveFile(file, filePath)
			if err != nil {
				continue
			}
			compressedFileURLs = append(compressedFileURLs, fileUrl)
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Code:    fiber.StatusCreated,
			Message: "Saved",
			Data:    compressedFileURLs,
		})
	})
}
