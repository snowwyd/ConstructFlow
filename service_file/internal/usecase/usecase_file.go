package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	"service-file/pkg/logger/slogger"
	"time"

	"github.com/minio/minio-go/v7"
)

type FileUsecase struct {
	fileMetadataRepo interfaces.FileMetadataRepository
	directoryRepo    interfaces.DirectoryRepository
	fileStorage      interfaces.FileStorage
	log              *slog.Logger
}

func NewFileUsecase(directoryRepo interfaces.DirectoryRepository, fileMetadataRepo interfaces.FileMetadataRepository, fileStorage interfaces.FileStorage, log *slog.Logger) *FileUsecase {
	return &FileUsecase{
		directoryRepo:    directoryRepo,
		fileMetadataRepo: fileMetadataRepo,
		fileStorage:      fileStorage,
		log:              log,
	}
}

func (u *FileUsecase) GetFileInfo(ctx context.Context, fileID uint, userID uint) (*domain.FileResponse, error) {
	const op = "usecase.file.GetFileInfo"

	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))
	log.Info("getting file info")

	file, err := u.fileMetadataRepo.GetFileInfo(ctx, fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		case errors.Is(err, domain.ErrAccessDenied):
			log.Warn("access denied", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		default:
			log.Error("failed to get file info", slogger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	response := &domain.FileResponse{
		ID:          file.ID,
		NameFile:    file.Name,
		Status:      file.Status,
		DirectoryID: file.DirectoryID,
	}

	log.Info("file info response prepared successfully")
	return response, nil
}

func (u *FileUsecase) CreateFile(ctx context.Context, directoryID uint, name string, fileData []byte, contentType string, userID uint) (err error) {
	const op = "usecase.file.CreateFile"

	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID))
	log.Info("uploading file")

	log.Debug("checking acces to directory", slog.Any("directory_id", directoryID))
	hasAccess, err := u.directoryRepo.CheckUserDirectoryAccess(ctx, userID, directoryID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDirectoryNotFound):
			log.Error("directory not found", slogger.Err(domain.ErrDirectoryNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrDirectoryNotFound)
		default:
			log.Error("failed to check directory access", slogger.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	if !hasAccess {
		log.Warn("user has no access to directory", slog.Any("directory_id", directoryID))
		return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	log.Debug("uploading file into minio")

	size := int64(len(fileData))

	objectKey := fmt.Sprintf("files/%d/%s", directoryID, name)
	err = u.fileStorage.UploadFile(ctx, "files", objectKey, fileData, contentType)
	if err != nil {
		return fmt.Errorf("minio upload failed: %w", err)
	}

	file := &domain.File{
		DirectoryID:    directoryID,
		Name:           name,
		Status:         "draft",
		MinioObjectKey: objectKey,
		Size:           size,
		ContentType:    contentType,
	}

	err = u.fileMetadataRepo.CreateFile(ctx, file.DirectoryID, file.Name, file.Status, userID, file.MinioObjectKey, file.Size, file.ContentType)
	if err != nil {
		return err
	}

	return nil
}

func (u *FileUsecase) DownloadFileDirect(ctx context.Context, fileID uint, userID uint) (*domain.File, *minio.Object, error) {
	const op = "usecase.file.DownloadFileDirect"
	log := u.log.With(slog.String("op", op), slog.Any("file_id", fileID))

	// 1. Получаем метаданные файла
	file, err := u.fileMetadataRepo.GetFileByID(ctx, fileID)
	if err != nil {
		log.Error("failed to get file metadata", slogger.Err(err))
		return nil, nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
	}

	// 2. Проверяем доступ пользователя
	hasAccess, err := u.directoryRepo.CheckUserDirectoryAccess(ctx, userID, file.DirectoryID)
	if err != nil || !hasAccess {
		log.Warn("access denied to file", slogger.Err(domain.ErrAccessDenied))
		return nil, nil, fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	// 3. Получаем объект из MinIO
	object, err := u.fileStorage.GetFile(ctx, "files", file.MinioObjectKey)
	if err != nil {
		log.Error("failed to retrieve file from MinIO", slogger.Err(err))
		return nil, nil, fmt.Errorf("%s: %w", op, domain.ErrInternal)
	}

	return file, object, nil
}

func (u *FileUsecase) UpdateFile(ctx context.Context, fileID uint, newData []byte, userID uint) error {
	const op = "usecase.file.UpdateFile"
	log := u.log.With(slog.String("op", op), slog.Any("file_id", fileID))

	// 1. Получаем текущий файл
	file, err := u.fileMetadataRepo.GetFileByID(ctx, fileID)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		}
		log.Error("failed to get file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// 2. Проверяем доступ пользователя к директории
	hasAccess, err := u.directoryRepo.CheckUserDirectoryAccess(ctx, userID, file.DirectoryID)
	if err != nil {
		log.Error("failed to check directory access", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasAccess {
		log.Warn("user has no access to directory", slogger.Err(domain.ErrAccessDenied))
		return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
	}

	newVersion := file.Version + 1

	// 4. Загружаем новую версию в MinIO
	newKey, err := u.fileStorage.UploadNewVersion(ctx, "files", file.MinioObjectKey, newData, file.ContentType, newVersion)
	if err != nil {
		log.Error("failed to upload new version", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	// 5. Обновляем метаданные файла
	file.MinioObjectKey = newKey
	file.Version = newVersion
	file.Size = int64(len(newData))
	file.UpdatedAt = time.Now()

	// 6. Сохраняем изменения в БД
	if err := u.fileMetadataRepo.UpdateFile(ctx, file); err != nil {
		log.Error("failed to update file metadata", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file updated successfully", slog.Int("new_version", file.Version))
	return nil
}

func (u *FileUsecase) DeleteFile(ctx context.Context, fileID uint, userID uint) (err error) {
	const op = "usecase.file.DeleteFile"
	log := u.log.With(slog.String("op", op), slog.Any("user_id", userID), slog.Any("file_id", fileID))

	log.Info("deleting file")
	err = u.fileMetadataRepo.DeleteFile(ctx, fileID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			log.Error("file not found", slogger.Err(domain.ErrFileNotFound))
			return fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
		case errors.Is(err, domain.ErrAccessDenied):
			log.Error("user has no access to this directory", slogger.Err(domain.ErrAccessDenied))
			return fmt.Errorf("%s: %w", op, domain.ErrAccessDenied)
		case errors.Is(err, domain.ErrCannotDeleteNonDraftFile):
			log.Error("cannot delete non draft file", slogger.Err(domain.ErrCannotDeleteNonDraftFile))
			return fmt.Errorf("%s: %w", op, domain.ErrCannotDeleteNonDraftFile)
		}
		log.Error("failed to delete file", slogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("file deleted successfully")
	return nil
}

func (u *FileUsecase) ConvertSTPToGLTF(ctx context.Context, fileID uint, userID uint) (string, error) {
	const op = "usecase.file.ConvertSTPToGLTF"

	// Получаем метаданные и файл из MinIO
	fileMeta, fileObj, err := u.DownloadFileDirect(ctx, fileID, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer fileObj.Close()

	// Сохраняем STP временно на диск
	// Изменим tmpDir, используя MinioObjectKey без расширения
	tmpDir := filepath.Join(os.TempDir(), fileMeta.MinioObjectKey)
	tmpInput := tmpDir
	tmpOutput := tmpInput + ".gltf"

	// Создаем директорию, если она не существует
	err = os.MkdirAll(filepath.Dir(tmpInput), os.ModePerm) // Создаем только родительскую директорию
	if err != nil {
		return "", fmt.Errorf("%s: failed to create temp directory: %w", op, err)
	}

	outFile, err := os.Create(tmpInput) // Создаем файл по пути tmpInput
	if err != nil {
		return "", fmt.Errorf("%s: failed to create temp file: %w", op, err)
	}
	defer os.Remove(tmpInput)
	defer outFile.Close()

	_, err = io.Copy(outFile, fileObj)
	if err != nil {
		return "", fmt.Errorf("%s: failed to copy file: %w", op, err)
	}

	// Вызываем python-скрипт
	cmd := exec.Command("python3", "/app/scripts/convert_stp_to_glb.py", tmpInput, tmpOutput)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: script error: %s\n%s", op, err, string(output))
	}

	// Можно загрузить tmpOutput обратно в MinIO и отдать ссылку
	return tmpOutput, nil
}
