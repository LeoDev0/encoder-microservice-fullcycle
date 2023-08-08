package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"io"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video 					*domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	// // Create a custom HTTP client with Transport that skips authentication.
	// httpClient := &http.Client{
	// 	Transport: &http.Transport{
	// 			// Set the endpoint directly here.
	// 			// Note: Use the correct URL for the fake-gcs-server based on your setup.
	// 			Proxy: http.ProxyFromEnvironment,
	// 			DialContext: (&net.Dialer{
	// 					Timeout:   30 * time.Second,
	// 					KeepAlive: 30 * time.Second,
	// 					DualStack: true,
	// 			}).DialContext,
	// 			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	// 			TLSHandshakeTimeout: 10 * time.Second,
	// 	},
	// }

	// client, err := storage.NewClient(context.TODO(), option.WithHTTPClient(httpClient), option.WithEndpoint("http://localhost:4443/storage/v1/"))

	// ctx := context.TODO()
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("video %v has been stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("localStoragePath") + "/" + v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
			log.Printf("=====> Output: %s\n", string(out))
	}
}