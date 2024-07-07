package cli

import (
	"errors"
	"github.com/arelate/align/paths"
	"github.com/arelate/align/render/view_models"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/url"
	"path"
	"strings"
)

const nextDataScriptId = "__NEXT_DATA__"

var (
	ErrDataNotFound = errors.New("data not found")
)

func GetDataHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	force := u.Query().Has("force")

	return GetData(slug, page, force)
}

func GetData(slug, page string, force bool) error {

	if page == "" {
		page = view_models.MainPage
	}

	rca := nod.Begin("getting data for %s...", path.Join(slug, page))
	defer rca.End()

	page, err := url.PathUnescape(page)
	if err != nil {
		return rca.EndWithError(err)
	}

	dkv, err := paths.DataKeyValues(slug)
	if err != nil {
		return rca.EndWithError(err)
	}

	has, err := dkv.Has(page)
	if err != nil {
		return rca.EndWithError(err)
	}

	if has && !force {
		rca.EndWithResult("already exist")
		return nil
	}

	pkv, err := paths.PagesKeyValues(slug)
	if err != nil {
		return rca.EndWithError(err)
	}

	has, err = pkv.Has(page)
	if err != nil {
		return rca.EndWithError(err)
	}

	if !has {
		rca.EndWithResult("page not available")
		return nil
	}

	src, err := pkv.Get(page)
	if err != nil {
		return rca.EndWithError(err)
	}

	if src == nil {
		rca.EndWithResult("not found")
	}

	if _, err := getSetReducedContent(page, src, dkv); err != nil {
		return rca.EndWithError(err)
	}

	rca.EndWithResult("done")

	return nil
}

func getSetReducedContent(page string, src io.Reader, kv kevlar.KeyValues) (string, error) {

	body, err := html.Parse(src)
	if err != nil {
		return "", err
	}

	if nextDataNode := match_node.Match(body, &nextDataMatcher{}); nextDataNode != nil && nextDataNode.FirstChild != nil {
		data := fixDataProblems(nextDataNode.FirstChild.Data)
		return data, kv.Set(page, strings.NewReader(data))
	}

	return "", ErrDataNotFound
}

type nextDataMatcher struct{}

func (ndm *nextDataMatcher) Match(node *html.Node) bool {
	if node.DataAtom != atom.Script ||
		len(node.Attr) == 0 {
		return false
	}

	return match_node.AttrVal(node, "id") == nextDataScriptId
}

func fixDataProblems(data string) string {
	// 1. htmlEntities.values sometimes is a struct of
	// github.com/arelate/southern_light/ign_integration/HTMLValue type
	// and sometimes is an array of images of
	// github.com/arelate/southern_light/ign_integration/ImageValue type
	// in order to fix that - look for "values:[" and replace with "imageValues:["
	// given that at the moment no other data would match that
	fixedData := strings.Replace(data, "\"values\":[", "\"imageValues\":[", -1)

	return fixedData
}

//func cutNextData(src io.Reader) (string, error) {
//
//	sb := new(strings.Builder)
//	if _, err := io.Copy(sb, src); err != nil {
//		return "", err
//	}
//
//	if _, nd, ok := strings.Cut(sb.String(), "<script id=\"__NEXT_DATA__\" type=\"application/json\">"); ok {
//		if nd, _, ok = strings.Cut(nd, "</script>"); ok {
//			return nd, nil
//		}
//	}
//	return "", ErrDataNotFound
//}
