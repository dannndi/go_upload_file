package file

import (
	"strconv"
	"time"

	"github.com/dannndi/go_upload_file/core/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
		log.Debug("Content Type : ", file.Header["Content-Type"]);
		milis := strconv.FormatInt(time.Now().UnixMilli(), 10)
		filePath := "./public/uploads/" + milis + "_" + file.Filename
		err = c.SaveFile(file, filePath)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.BaseResponse{
				Message: "Something wrong when saving the file",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Message: "Saved",
			Data:    filePath,
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
			err = c.SaveFile(file, filePath)
			if err != nil {
				continue
			}
			compressedFileURLs = append(compressedFileURLs, filePath)
		}

		return c.Status(fiber.StatusCreated).JSON(utils.BaseResponse{
			Message: "Saved",
			Data:    compressedFileURLs,
		})
	})
}
