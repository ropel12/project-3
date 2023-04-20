package pkg

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

type StorageGCP struct {
	ClG        *storage.Client
	ProjectID  string
	BucketName string
	Path       string
}

func (s *StorageGCP) UploadFile(file multipart.File, fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	wc := s.ClG.Bucket(s.BucketName).Object(s.Path + fileName).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}
