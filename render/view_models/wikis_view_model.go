package view_models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/arelate/southern_light/ign_integration"
	"github.com/boggydigital/kvas"
	"html/template"
	"net/url"
	"path"
)

var (
	slugTitles        = make(map[string]string)
	slugPrimaryImages = make(map[string]string)
)

func NewWikisViewModel(wikis []string, keyValues map[string]kvas.KeyValues) (*WikisSlugViewModel, error) {
	wsvm := &WikisSlugViewModel{
		GuideTitle: "All Guides",
	}

	for _, w := range wikis {

		if _, ok := slugTitles[w]; !ok {
			title, primaryImage, err := getTitlePrimaryImage(keyValues[w])
			if err != nil {
				return nil, err
			}
			slugTitles[w] = title
			slugPrimaryImages[w] = primaryImage
		}

		piu, err := url.Parse(slugPrimaryImages[w])
		if err != nil {
			return nil, err
		}

		u := fmt.Sprintf("<a href='/wikis/%s'><img src='/%s'/><span>%s</span></a>",
			w,
			path.Join("primary_image", piu.Path),
			slugTitles[w])

		wsvm.Items = append(wsvm.Items,
			template.HTML(u))
	}

	return wsvm, nil
}

func getTitlePrimaryImage(kv kvas.KeyValues) (string, string, error) {

	if err := kv.IndexRefresh(); err != nil {
		return "", "", err
	}

	wp, err := kv.Get(MainPage)
	if err != nil {
		return "", "", err
	}
	if wp == nil {
		return "", "", errors.New("page not found: " + MainPage)
	}
	defer wp.Close()

	var wikiProps ign_integration.WikiProps
	if err := json.NewDecoder(wp).Decode(&wikiProps); err != nil {
		return "", "", err
	}

	return wikiProps.Props.PageProps.Page.Name,
		wikiProps.PrimaryImageUrl(),
		nil

}
