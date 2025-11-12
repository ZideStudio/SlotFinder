package account

import (
	"app/config"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

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
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", filepath.Base(fileName))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	writer.Close()

	config := config.GetConfig()

	url := fmt.Sprintf("https://api.imgbb.com/1/upload?key=%s&name=%s", config.ImgBBApiKey, fileName)
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create HTTP request for imgbb")
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload image to imgbb")
		return "", err
	}
	defer resp.Body.Close()

	// Log the response body for debugging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return "", err
	}

	// Reset the response body for JSON decoding
	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	var imgbbResp ImgbbResponse
	if err := json.NewDecoder(resp.Body).Decode(&imgbbResp); err != nil {
		log.Error().Err(err).Msg("Failed to decode imgbb response")
		return "", err
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
	normalized := strings.TrimSpace(strings.ToLower(username))

	hash := md5.Sum([]byte(normalized))
	hashStr := hex.EncodeToString(hash[:])

	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=retro", hashStr)
}
