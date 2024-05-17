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
	"strings"
)

var (
	ErrReducedContentNotPresent = errors.New("reduced content not present")
)

func ReduceContentHandler(u *url.URL) error {

	force := u.Query().Has("force")

	return ReduceContent(force)
}

func ReduceContent(force bool) error {

	rca := nod.NewProgress("reducing content...")
	defer rca.End()

	scd, err := pathways.GetAbsDir(paths.SourceContent)
	if err != nil {
		return rca.EndWithError(err)
	}

	skv, err := kvas.ConnectLocal(scd, kvas.HtmlExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	rcd, err := pathways.GetAbsDir(paths.ReducedContent)
	if err != nil {
		return rca.EndWithError(err)
	}

	rkv, err := kvas.ConnectLocal(rcd, kvas.JsonExt)
	if err != nil {
		return rca.EndWithError(err)
	}

	rca.TotalInt(len(skv.Keys()))

	for _, slug := range skv.Keys() {
		if err := getSetReducedContent(slug, skv, rkv, force); err != nil {
			return rca.EndWithError(err)
		}
		rca.Increment()
	}

	rca.EndWithResult("done")

	return nil
}

func getSetReducedContent(slug string, skv, rkv kvas.KeyValues, force bool) error {

	if rkv.Has(slug) && !force {
		return nil
	}

	sc, err := skv.Get(slug)
	if err != nil {
		return err
	}
	defer sc.Close()

	body, err := html.Parse(sc)
	if err != nil {
		return err
	}

	if nextDataNode := match_node.Match(body, &nextDataMatcher{}); nextDataNode != nil && nextDataNode.FirstChild != nil {
		return rkv.Set(slug, strings.NewReader(nextDataNode.FirstChild.Data))
	}

	return ErrReducedContentNotPresent
}

type nextDataMatcher struct{}

func (ndm *nextDataMatcher) Match(node *html.Node) bool {
	if node.DataAtom != atom.Script ||
		len(node.Attr) == 0 {
		return false
	}

	return match_node.AttrVal(node, "id") == "__NEXT_DATA__"
}
