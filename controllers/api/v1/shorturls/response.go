package shorturls

import (
	"encoding/json"
	"net/url"
	"url-shortener/models"
)

type createShortUrlResponseHelper struct {
	Host string
	models.ShortUrl
}

func (r createShortUrlResponseHelper) MarshalJSON() ([]byte, error) {
	u := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   r.Slug,
	}

	return json.Marshal(CreateShortUrlResponse{
		ShortUrl:           u.String(),
		ShortUrlReadFields: r.ShortUrl.ShortUrlReadFields,
	})
}
