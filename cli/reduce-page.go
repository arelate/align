package cli

import (
	"errors"
	"github.com/arelate/align/paths"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/match_node"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/pathways"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

const nextDataScriptId = "__NEXT_DATA__"

var (
	ErrReducedContentNotPresent = errors.New("reduced content not present")
)

func ReducePageHandler(u *url.URL) error {
	slug := u.Query().Get("slug")
	page := u.Query().Get("page")
	force := u.Query().Has("force")
	return ReducePage(slug, page, force)
}

func ReducePage(slug, page string, force bool) error {

	if page == "" {
		page = mainPage
	}

	rca := nod.Begin("reducing page %s...", path.Join(slug, page))
	defer rca.End()

	rpd, err := pathways.GetAbsDir(paths.ReducedPages)
	if err != nil {
		return rca.EndWithError(err)
	}

	rpd = filepath.Join(rpd, slug)

	rkv, err := kvas.ConnectLocal(rpd, kvas.JsonExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	if rkv.Has(page) && !force {
		return nil
	}

	spd, err := pathways.GetAbsDir(paths.SourcePages)
	if err != nil {
		return rca.EndWithError(err)
	}

	spd = filepath.Join(spd, slug)

	skv, err := kvas.ConnectLocal(spd, kvas.HtmlExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	if err := getSetReducedContent(page, skv, rkv); err != nil {
		return rca.EndWithError(err)
	}

	rca.EndWithResult("done")

	return nil
}

func getSetReducedContent(page string, skv, rkv kvas.KeyValues) error {

	sc, err := skv.Get(page)
	if err != nil {
		return err
	}
	defer sc.Close()

	body, err := html.Parse(sc)
	if err != nil {
		return err
	}

	if nextDataNode := match_node.Match(body, &nextDataMatcher{}); nextDataNode != nil && nextDataNode.FirstChild != nil {
		return rkv.Set(page, strings.NewReader(nextDataNode.FirstChild.Data))
	}

	return ErrReducedContentNotPresent
}

type nextDataMatcher struct{}

func (ndm *nextDataMatcher) Match(node *html.Node) bool {
	if node.DataAtom != atom.Script ||
		len(node.Attr) == 0 {
		return false
	}

	return match_node.AttrVal(node, "id") == nextDataScriptId
}
