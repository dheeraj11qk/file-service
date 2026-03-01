package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"file-service/internal/db"
	"file-service/internal/models"
	"file-service/proto"

	"github.com/google/uuid"
)

type FileService struct {
	proto.UnimplementedFileServiceServer
}

// UploadFile - Client Streaming
func (s *FileService) UploadFile(stream proto.FileService_UploadFileServer) error {

	var fileName string
	var fileSize int64

	// Generate file ID
	fileID := uuid.New().String()
	filePath := filepath.Join("uploads", fileID)

	// Create uploads directory if not exists
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fileName = req.Filename

		n, err := file.Write(req.ChunkData)
		if err != nil {
			return err
		}

		fileSize += int64(n)
	}

	// Save metadata to DB
	newFile := models.File{
		FileName:   fileName,
		FileSize:   fileSize,
		FilePath:   filePath,
		UploadTime: time.Now(),
		Status:     "completed",
	}

	err = db.DB.Create(&newFile).Error
	if err != nil {
		return err
	}

	return stream.SendAndClose(&proto.UploadFileResponse{
		FileId:      newFile.ID.String(),
		FileMessage: "File uploaded successfully",
	})
}

func (s *FileService) GetFileInfo(
	ctx context.Context,
	req *proto.GetFileInfoRequest,
) (*proto.GetFileInfoResponse, error) {

	var file models.File

	// Find file by UUID
	err := db.DB.First(&file, "id = ?", req.FileId).Error
	if err != nil {
		return nil, err
	}

	return &proto.GetFileInfoResponse{
		FileId:     file.ID.String(),
		FileName:   file.FileName,
		FileSize:   file.FileSize,
		UploadTime: file.UploadTime.Format(time.RFC3339),
		Status:     file.Status,
	}, nil
}
