package provider

import (
	"AynaLivePlayer/player"
	"fmt"
	"github.com/tidwall/gjson"
	"net/url"
	"regexp"
)

type Bilibili struct {
	InfoApi   string
	FileApi   string
	SearchApi string
	IdRegex0  *regexp.Regexp
	IdRegex1  *regexp.Regexp
}

func _newBilibili() *Bilibili {
	return &Bilibili{
		InfoApi:   "https://www.bilibili.com/audio/music-service-c/web/song/info?sid=%s",
		FileApi:   "https://api.bilibili.com/audio/music-service-c/url?device=phone&mid=8047632&mobi_app=iphone&platform=ios&privilege=2&songid=%s&quality=2",
		SearchApi: "https://api.bilibili.com/audio/music-service-c/s?search_type=music&keyword=%s&page=1&pagesize=100",
		IdRegex0:  regexp.MustCompile("^[0-9]+"),
		IdRegex1:  regexp.MustCompile("^au[0-9]+"),
	}
}

var BilibiliAPI *Bilibili

func init() {
	BilibiliAPI = _newBilibili()
	Providers[BilibiliAPI.GetName()] = BilibiliAPI
}

func (b *Bilibili) GetName() string {
	return "bilibili"
}

func (b *Bilibili) MatchMedia(keyword string) *player.Media {
	if id := b.IdRegex0.FindString(keyword); id != "" {
		return &player.Media{
			Meta: Meta{
				Name: b.GetName(),
				Id:   id,
			},
		}
	}
	if id := b.IdRegex1.FindString(keyword); id != "" {
		return &player.Media{
			Meta: Meta{
				Name: b.GetName(),
				Id:   id[2:],
			},
		}
	}
	return nil
}

func (b *Bilibili) FormatPlaylistUrl(uri string) string {
	return ""
}

func (b *Bilibili) GetPlaylist(meta Meta) ([]*player.Media, error) {
	return nil, ErrorExternalApi
}

func (b *Bilibili) Search(keyword string) ([]*player.Media, error) {
	resp := httpGetString(fmt.Sprintf(b.SearchApi, url.QueryEscape(keyword)), map[string]string{
		"user-agent": "BiliMusic/2.233.3",
	})
	if resp == "" {
		return nil, ErrorExternalApi
	}
	result := make([]*player.Media, 0)
	gjson.Get(resp, "data.result").ForEach(func(key, value gjson.Result) bool {
		result = append(result, &player.Media{
			Title:  value.Get("title").String(),
			Cover:  player.Picture{Url: value.Get("cover").String()},
			Artist: value.Get("author").String(),
			Meta: Meta{
				Name: b.GetName(),
				Id:   value.Get("id").String(),
			},
		})
		return true
	})
	return result, nil
}

func (b *Bilibili) UpdateMedia(media *player.Media) error {
	resp := httpGetString(fmt.Sprintf(b.InfoApi, media.Meta.(Meta).Id), map[string]string{
		"user-agent": "BiliMusic/2.233.3",
	})
	if resp == "" {
		return ErrorExternalApi
	}
	if gjson.Get(resp, "data.title").String() == "" {
		return ErrorExternalApi
	}
	media.Title = gjson.Get(resp, "data.title").String()
	media.Cover.Url = gjson.Get(resp, "data.cover").String()
	media.Artist = gjson.Get(resp, "data.author").String()
	media.Album = media.Title
	return nil
}

func (b *Bilibili) UpdateMediaUrl(media *player.Media) error {
	resp := httpGetString(fmt.Sprintf(b.FileApi, media.Meta.(Meta).Id), map[string]string{
		"user-agent": "BiliMusic/2.233.3",
	})

	if resp == "" {
		return ErrorExternalApi
	}
	media.Header = map[string]string{
		"user-agent": "BiliMusic/2.233.3",
	}
	uri := gjson.Get(resp, "data.cdns.0").String()
	if uri == "" {
		return ErrorExternalApi
	}
	media.Url = uri
	return nil
}
func (k *Bilibili) UpdateMediaLyric(media *player.Media) error {

	return nil
}
