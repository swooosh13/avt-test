package s3

import (
	"fmt"
	"log"

	"github.com/minio/minio-go"
)

type Client struct {
	*minio.Client
	Endpoint   string
	BucketName string
}

func New(accessKey, secKey, endpoint, bucketName string) *Client {
	ssl := true

	if accessKey == "" || secKey == "" || endpoint == "" || bucketName == "" {
		log.Fatal("Must provide Cloud_KEY and Cloud_SECRET and Cloud_ENDPOINT and Cloud_BUCKET_NAME environment variables!")
	}

	// Подключиться к VK Cloud S3.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{client, endpoint, bucketName}
}

func (c *Client) UploadFile(objectName string, filePath string) error {
	contentType := "text/csv"

	_, err := c.FPutObject(c.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("upload file: %w", err)
	}

	return nil
}
