package account

import (
	"app/config"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/rs/zerolog/log"
)

type AvatarService struct{}

func NewAvatarService() *AvatarService {
	return &AvatarService{}
}

const maxSize = 32 * 1024 * 1024 // 32 MB

type ImgbbResponse struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Success bool `json:"success"`
}

func compressToJPEG(img image.Image) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func uploadToImgbb(file io.Reader, fileName string) (string, error) {
	config := config.GetConfig()

	client := resty.New()

	var imgbbResp ImgbbResponse
	resp, err := client.R().
		SetFileReader("image", filepath.Base(fileName), file).
		SetQueryParam("key", config.ImgBBApiKey).
		SetQueryParam("name", fileName).
		SetResult(&imgbbResp).
		Post("https://api.imgbb.com/1/upload")

	if err != nil {
		log.Error().Err(err).Msg("Failed to upload image to imgbb")
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error().Int("status", resp.StatusCode()).Msg("ImgBB upload failed with non-200 status")
		return "", fmt.Errorf("failed to upload image to imgbb: status %d", resp.StatusCode())
	}

	if !imgbbResp.Success {
		log.Error().Msg("ImgBB upload was not successful")
		return "", fmt.Errorf("failed to upload image to imgbb")
	}

	return imgbbResp.Data.URL, nil
}

func (*AvatarService) UploadAvatar(imgURL, fileName string) (string, error) {
	// Download the image from the URL
	resp, err := http.Get(imgURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to download image")
		return "", err
	}
	defer resp.Body.Close()

	// Read the image data
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read image data")
		return "", err
	}

	var fileReader io.Reader
	if len(imgData) > maxSize {
		img, _, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			log.Error().Err(err).Msg("Failed to decode image for compression")
			return "", err
		}

		compressed, err := compressToJPEG(img)
		if err != nil {
			log.Error().Err(err).Msg("Failed to compress image")
			return "", err
		}
		fileReader = compressed
	} else {
		fileReader = bytes.NewReader(imgData)
	}

	// Upload to imgbb
	url, err := uploadToImgbb(fileReader, fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload image to imgbb")
		return "", err
	}

	return url, nil
}

func (*AvatarService) GetGravatarURL(username string) string {
	hash := sha256.Sum256([]byte(username))
	hashStr := hex.EncodeToString(hash[:])

	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=retro", hashStr)
}
