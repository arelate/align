package cli

import (
	"errors"
	"github.com/arelate/align/paths"
	"github.com/boggydigital/kvas"
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
	ErrReducedContentNotPresent = errors.New("reduced content not present")
)

func GetDataHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	force := u.Query().Has("force")

	return GetData(slug, page, force)
}

func GetData(slug, page string, force bool) error {

	if page == "" {
		page = mainPage
	}

	rca := nod.Begin("reducing page %s...", path.Join(slug, page))
	defer rca.End()

	page, err := url.PathUnescape(page)
	if err != nil {
		return rca.EndWithError(err)
	}

	rpd, err := paths.AbsDataSlugDir(slug)
	if err != nil {
		return rca.EndWithError(err)
	}

	rkv, err := kvas.ConnectLocal(rpd, kvas.JsonExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	if rkv.Has(page) && !force {
		rca.EndWithResult("already exist")
		return nil
	}

	spd, err := paths.AbsPagesSlugDir(slug)
	if err != nil {
		return rca.EndWithError(err)
	}

	skv, err := kvas.ConnectLocal(spd, kvas.HtmlExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	src, err := skv.Get(page)
	if err != nil {
		return rca.EndWithError(err)
	}

	if _, err := getSetReducedContent(page, src, rkv); err != nil {
		return rca.EndWithError(err)
	}

	rca.EndWithResult("done")

	return nil
}

func getSetReducedContent(page string, src io.Reader, kv kvas.KeyValues) (string, error) {

	body, err := html.Parse(src)
	if err != nil {
		return "", err
	}

	if nextDataNode := match_node.Match(body, &nextDataMatcher{}); nextDataNode != nil && nextDataNode.FirstChild != nil {
		data := fixDataProblems(nextDataNode.FirstChild.Data)
		return data, kv.Set(page, strings.NewReader(data))
	}

	return "", ErrReducedContentNotPresent
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
