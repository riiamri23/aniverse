package extractor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"aniverse/internal/crawler"
	"aniverse/pkg/common"
)

type Gogocdn struct {
	name            string
	key             []byte
	decryptionKey   []byte
	iv              []byte
	baseCrawler     *crawler.BaseCrawler
	reKeys          *regexp.Regexp
	reEncryptedData *regexp.Regexp
}

func NewGogocdn(baseCrawler *crawler.BaseCrawler) *Gogocdn {
	if baseCrawler == nil {
		baseCrawler = crawler.DefaultBaseCrawler
	}

	return &Gogocdn{
		name:          "gogocdn",
		key:           []byte("37911490979715163134003223491201"),
		decryptionKey: []byte("54674138327930866480207815084989"),
		iv:            []byte("3134003223491201"),
		baseCrawler:   baseCrawler,

		reKeys:          regexp.MustCompile(`(?:container|videocontent)-(\d+)`),
		reEncryptedData: regexp.MustCompile(`data-value="(.+?)"`),
	}
}

type gogoCdnData struct {
	Data string `json:"data"`
}

type gogoCdn struct {
	Source []struct {
		File string `json:"file"`
		Type string `json:"type"`
	} `json:"source"`
	Bkp []struct {
		File string `json:"file"`
		Type string `json:"type"`
	} `json:"source_bk"`
	Track interface{} `json:"track"`
}

func (g *Gogocdn) Extract(link string) (*common.Sources, error) {
	sources := new(common.Sources)

	parsedURL, err := url.Parse(link)
	if err != nil {
		return nil, fmt.Errorf("Gogocdn Extract: %w : url is not valid", common.ErrInvalidArgument)
	}

	contentID := parsedURL.Query().Get("id")
	if contentID == "" {
		return nil, fmt.Errorf(
			"Gogocdn Extract: %w : url does not have id query",
			common.ErrInvalidArgument,
		)
	}

	encyptedParams, err := g.parsePage(link, contentID)
	if err != nil {
		return nil, fmt.Errorf("Gogocdn Extract: %w", err)
	}

	nextHost := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	url := fmt.Sprintf("%s//encrypt-ajax.php?%s", nextHost, encyptedParams)
	headers := map[string]string{"X-Requested-With": "XMLHttpRequest"}
	response, err := g.baseCrawler.Client.Get(url, headers)
	if err != nil {
		return nil, fmt.Errorf("Gogocdn Extract: %w : %s", common.ErrRequest, err.Error())
	}

	var gogoCdnData gogoCdnData
	err = json.Unmarshal(response, &gogoCdnData)
	if err != nil {
		return nil, fmt.Errorf(
			"Gogocdn Extract: %w : gogo cdn data : %s",
			common.ErrJsonParse,
			err.Error(),
		)
	}

	decData, err := g.aesDecrypt(gogoCdnData.Data, g.decryptionKey, g.iv)
	if err != nil {
		return nil, fmt.Errorf(
			"Gogocdn Extract: %w : cannot decrypt sources : %s",
			common.ErrScraping,
			err.Error(),
		)
	}

	var dataFile gogoCdn
	err = json.Unmarshal(decData, &dataFile)
	if err != nil {
		return nil, fmt.Errorf(
			"Gogocdn Extract: %w : gogo cdn sources file : %s",
			common.ErrJsonParse,
			err.Error(),
		)
	}

	for _, s := range dataFile.Source {
		if s.File == "" {
			continue
		}
		sources.Sources = append(sources.Sources, common.Source{
			Url:    s.File,
			Type:   s.Type,
			IsM3U8: strings.Contains(s.File, ".m3u8"),
		})
	}

	switch track := dataFile.Track.(type) {
	case []interface{}:
		break
	case map[string]interface{}:
		tracks, ok := track["tracks"].([]interface{})
		if !ok {
			break
		}
		for _, ts := range tracks {
			t, ok := ts.(map[string]interface{})
			if !ok {
				break
			}
			if strings.ToLower(t["kind"].(string)) == "thumbnails" {
				sources.Thumbnail = t["file"].(string)
				sources.ThumbnailType = "Sprite"
			}
		}
	}

	for _, s := range dataFile.Bkp {
		if s.File == "" {
			continue
		}
		sources.Sources = append(sources.Sources, common.Source{
			Url:    s.File,
			Type:   s.Type,
			IsM3U8: strings.Contains(s.File, ".m3u8"),
		})
	}

	sources.Flags = []common.Flag{common.FLAG_CORS_ALLOWED}

	if len(sources.Sources) == 0 {
		return nil, fmt.Errorf("Gogocdn Extract: %w", common.ErrNoContent)
	}

	return sources, nil
}

func (g *Gogocdn) parsePage(
	link string,
	contentID string,
) (string, error) {
	var encryptionKey, iv, decryptionKey, encryptedData []byte
	response, err := g.baseCrawler.Client.Get(link, nil)
	if err != nil {
		return "", fmt.Errorf("Gogocdn parsePage: %w", common.ErrRequest)
	}

	matches := g.reKeys.FindAllSubmatch(response, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		encryptionKey = matches[0][1]
		g.key = encryptionKey
	}

	if encryptionKey == nil {
		encryptionKey = g.key
	}

	if len(matches) > 1 && len(matches[1]) > 1 {
		iv = matches[1][1]
		g.iv = iv
	}

	if iv == nil {
		iv = g.iv
	}

	if len(matches) > 2 && len(matches[2]) > 1 {
		decryptionKey = matches[2][1]
		g.decryptionKey = decryptionKey
	}

	match := g.reEncryptedData.FindSubmatch(response)
	if len(match) < 2 {
		return "", fmt.Errorf("Gogocdn parsePage: %w", common.ErrInvalidRegex)
	}
	encryptedData = match[1]

	decryptedData, err := g.aesDecrypt(string(encryptedData), encryptionKey, iv)
	if err != nil {
		return "", fmt.Errorf(
			"Gogocdn parsePage: %w : decryption error : %s",
			common.ErrScraping,
			err.Error(),
		)
	}

	encryptedContentID, err := g.aesEncrypt(
		[]byte(contentID),
		encryptionKey,
		iv,
	)
	if err != nil {
		return "", fmt.Errorf(
			"Gogocdn parsePage: %w : encryption error : %s",
			common.ErrScraping,
			err.Error(),
		)
	}

	component := fmt.Sprintf("id=%s&alias=%s&%s", encryptedContentID, contentID, decryptedData)

	return component, nil
}

func (g *Gogocdn) pad(data []byte) []byte {
	padding := 16 - (len(data) % 16)
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func (g *Gogocdn) aesEncrypt(data []byte, key []byte, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedData := g.pad(data)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (g *Gogocdn) aesDecrypt(data string, key []byte, iv []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Unpad the decrypted data
	padding := int(plaintext[len(plaintext)-1])
	return plaintext[:len(plaintext)-padding], nil
}
