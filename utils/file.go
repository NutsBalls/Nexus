// utils/file.go
package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFile(file *multipart.FileHeader, path string) error {
	// Создаем директорию, если она не существует
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	// Открываем исходный файл
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Создаем целевой файл
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Копируем содержимое файла
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// Дополнительная функция для удаления файла
func DeleteFile(path string) error {
	return os.Remove(path)
}

// Функция для проверки существования файла
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
