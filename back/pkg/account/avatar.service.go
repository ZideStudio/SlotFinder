package account

import (
	"app/commons/lib"
	model "app/db/models"
	"app/db/repository"
	"fmt"
	"time"

	dicebear "github.com/dicebear/dicebear-go/v10"
	"github.com/dicebear/styles/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AvatarService struct {
	accountRepository *repository.AccountRepository
}

func NewAvatarService(service *AvatarService) *AvatarService {
	if service != nil {
		return service
	}

	return &AvatarService{
		accountRepository: repository.NewAccountRepository(nil),
	}
}

const defaultAvatarSize = 256

var glyphsStyle = mustLoadGlyphsStyle()

func mustLoadGlyphsStyle() *dicebear.Style {
	style, err := dicebear.NewStyle([]byte(styles.Glyphs))
	if err != nil {
		panic(fmt.Sprintf("account: failed to load Dicebear glyphs style: %v", err))
	}
	return style
}

var glyphShapeVariants = []any{
	"variant01", "variant02", "variant03", "variant04", "variant05",
	"variant06", "variant07", "variant08", "variant09", "variant10",
	"variant11", "variant12", "variant13", "variant14", "variant15",
	"variant16", "variant17", "variant18", "variant19", "variant20",
	"variant21", "variant22", "variant23", "variant24", "variant25",
	"variant27", "variant28", "variant29", "variant30", "variant31",
	"variant32", "variant33", "variant34", "variant35",
}

var glyphColors = []any{
	"76a7ff", "525fa3", "8c8c8c", "75c675", "ffaa64",
	"ff4d6f", "7a14c8", "b697ce", "3e2253", "5abfbd",
}

// GetDicebearSVG returns the deterministic Glyphs-style Dicebear avatar, as SVG markup, for seed.
func GetDicebearSVG(seed string) (string, error) {
	avatar, err := dicebear.NewAvatar(glyphsStyle, map[string]any{
		"shapeVariant":        glyphShapeVariants,
		"glyphColor":          glyphColors,
		"glyphColorFillStops": 2,
		"glyphColorAngle":     -25,
		"seed":                seed,
	})
	if err != nil {
		return "", err
	}
	return avatar.SVG(), nil
}

// FetchAndStoreDefaultAvatar generates and returns the default avatar bytes and local URL.
func (*AvatarService) FetchAndStoreDefaultAvatar(seed string, accountId uuid.UUID) ([]byte, string) {
	avatarUrl := fmt.Sprintf("/api/v1/account/%s/avatar", accountId.String())

	svg, err := GetDicebearSVG(seed)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate Dicebear avatar")
		return nil, avatarUrl
	}

	rendered, err := lib.RenderSVGToPNG([]byte(svg), defaultAvatarSize)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to rasterize Dicebear avatar")
		return nil, avatarUrl
	}

	data, err := lib.ProcessAvatar(rendered)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to process Dicebear avatar")
		return nil, avatarUrl
	}
	return data, avatarUrl
}

func (s *AvatarService) FindAvatarById(id uuid.UUID) ([]byte, *time.Time, error) {
	return s.accountRepository.FindAvatarById(id)
}

// UploadAvatar processes an image from a URL or raw bytes and returns the result as JPEG bytes.
func (*AvatarService) UploadAvatar(imgUrl *string, imgBytes []byte) ([]byte, error) {
	if imgUrl == nil && imgBytes == nil {
		return nil, fmt.Errorf("no image provided")
	}

	if imgUrl != nil {
		data, err := lib.ProcessAvatarFromURL(*imgUrl)
		if err != nil {
			log.Error().Err(err).Msg("Failed to process avatar from URL")
			return nil, err
		}
		return data, nil
	}

	data, err := lib.ProcessAvatar(imgBytes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process avatar bytes")
		return nil, err
	}
	return data, nil
}

func (s *AvatarService) UploadUserAvatar(imgBytes []byte, userId uuid.UUID) error {
	processed, err := s.UploadAvatar(nil, imgBytes)
	if err != nil {
		return fmt.Errorf("error processing avatar: %w", err)
	}

	avatarUrl := fmt.Sprintf("/api/v1/account/%s/avatar", userId.String())

	if err := s.accountRepository.Updates(model.Account{
		Id:         userId,
		AvatarUrl:  avatarUrl,
		AvatarData: processed,
	}); err != nil {
		return fmt.Errorf("error updating avatar on account: %w", err)
	}

	return nil
}
