package shorturls

import (
	"encoding/json"
	"net/url"
	"url-shortener/models"
)

type shortUrlResponseHelper struct {
	Host string
	models.ShortUrl
}
type ShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
	models.ShortUrlReadFields
}

func (r shortUrlResponseHelper) MarshalJSON() ([]byte, error) {
	u := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   r.Slug,
	}

	return json.Marshal(ShortUrlResponse{
		ShortUrl:           u.String(),
		ShortUrlReadFields: r.ShortUrl.ShortUrlReadFields,
	})
}
