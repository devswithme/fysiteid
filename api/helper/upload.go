package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UploadHelper interface {
	initialize() error
	Save(c *fiber.Ctx, key string, ID string) (string, error)
	Delete(path string) error
}

type uploadHelper struct {
	parent     string
	extensions []string
	size       int64
	logger     *zap.Logger
}

func NewUploadHelper(parent string, size int64, extensions []string, logger *zap.Logger) UploadHelper {
	u := &uploadHelper{
		parent:     parent,
		extensions: extensions,
		size:       size,
		logger:     logger,
	}
	u.initialize()
	return u
}

func (u *uploadHelper) initialize() error {
	if err := os.MkdirAll(u.parent, 0755); err != nil {
		u.logger.Error("failed initialize parent directory", zap.Error(err))
		return err
	}

	return nil
}

func (u *uploadHelper) Save(c *fiber.Ctx, key string, ID string) (string, error) {
	file, err := c.FormFile(key)

	if err == nil {
		if file.Size > u.size {
			return "", fmt.Errorf("file too large")
		}

		if !slices.Contains(u.extensions, filepath.Ext(file.Filename)) {
			return "", fmt.Errorf("invalid file extension")
		}

		if err := os.MkdirAll(fmt.Sprintf("%s/%s", u.parent, key), 0755); err != nil {
			u.logger.Error(fmt.Sprintf("failed making directory %s", key), zap.Error(err))
			return "", err
		}

		path := fmt.Sprintf("%s/%s/%s%s", u.parent, key, ID, filepath.Ext(file.Filename))

		if err := c.SaveFile(file, path); err != nil {
			u.logger.Error("failed saving file", zap.Error(err))
			return "", err
		}

		return path, nil
	}

	return "", nil
}

func (u *uploadHelper) Delete(path string) error {
	if err := os.Remove(path); err != nil {
		u.logger.Error("failed deletings file", zap.Error(err))
		return err
	}

	return nil
}
