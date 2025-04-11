package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
)

func IsConvertible(contentType string) bool {
	supportedFormats := map[string]bool{
		"application/step":         true,
		"model/step":               true,
		"application/iges":         true,
		"application/octet-stream": true,
		// Добавьте другие поддерживаемые форматы
	}
	return supportedFormats[contentType]
}

func ConvertToGLTF(fileObject *minio.Object) ([]byte, error) {
	// Читаем файл из MinIO
	fileBytes, err := io.ReadAll(fileObject)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Создаем multipart/form-data запрос
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "model.step")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("failed to write file to form: %w", err)
	}
	writer.WriteField("output_format", "gltf")
	writer.Close()

	// Отправляем запрос во внешний сервис (например, 3DConvert)
	req, err := http.NewRequest("POST", "https://api.cloudconvert.com/v2/convert", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer YOUR_API_KEY") // Замените на ваш API-ключ

	// Выполняем запрос с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("conversion failed (%d): %s", resp.StatusCode, string(responseBody))
	}

	return io.ReadAll(resp.Body)
}
