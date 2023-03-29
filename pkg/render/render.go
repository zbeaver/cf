package render

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/frankenbeanies/randhex"
	"github.com/zbeaver/cafe/pkg/vui"
)

type Registry map[uint32]Executor

type Executor interface {
	Render(vui.INode, styling, string) string
	Style(styling, vui.INode) styling
}

type exec func(styling) string

const (
	TAG_BODY = iota + 1
	TAG_BR
	TAG_DIV
	TAG_HEAD
	TAG_HTML
	TAG_IMG
	TAG_TEXT
	TAG_UNKNOWN
)

type engine struct {
	registry Registry
	ctx      context.Context
	doc      vui.Documentary
}

func NewEngine(ctx context.Context, doc vui.Documentary) *engine {
	reg := Registry{
		TAG_BODY:    (*Body)(nil),
		TAG_BR:      (*Br)(nil),
		TAG_DIV:     (*Div)(nil),
		TAG_HEAD:    (*Head)(nil),
		TAG_HTML:    (*Html)(nil),
		TAG_IMG:     (*Img)(nil),
		TAG_TEXT:    (*Text)(nil),
		TAG_UNKNOWN: (*Unknown)(nil),
	}

	return &engine{
		ctx:      ctx,
		doc:      doc,
		registry: reg,
	}
}

func (e *engine) Render() string {
	res := strings.Builder{}
	for _, c := range e.doc.ChildNodes() {
		res.WriteString(e.executor(c)(NewStyling()))
	}
	return res.String()
}

func debug(s styling, block string) string {
	w, h := lipgloss.Size(block)

	return s.Copy().Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color(randhex.New().String())).
		Render(fmt.Sprintf("[%vx%v]", w, h))
}

func (e *engine) executor(node vui.INode) exec {
	var ex Executor

	switch node.(type) {
	case *vui.HtmlElm:
		ex = e.registry[TAG_HTML]
	case *vui.HeadElm:
		ex = e.registry[TAG_HEAD]
	case *vui.BrElm:
		ex = e.registry[TAG_BR]
	case *vui.BodyElm:
		ex = e.registry[TAG_BODY]
	case *vui.DivElm:
		ex = e.registry[TAG_DIV]
	case *vui.ImgElm:
		ex = e.registry[TAG_IMG]
	case *vui.Text:
		ex = e.registry[TAG_TEXT]
	default:
		ex = e.registry[TAG_UNKNOWN]
	}

	return exec(func(base styling) string {
		cs := make(cells, 0)

		for _, c := range node.ChildNodes() {
			// rescurive call
			styling := ex.Style(base, c)
			subview := strings.TrimSpace(e.executor(c)(styling))
			if subview != "" {
				cs = append(cs, NewCell(subview, WithDisplay(styling.display)))
			}
		}

		if len(cs) > 0 {
			t := NewTissue(cs, &base)
			return strings.TrimSpace(ex.Render(node, base, t.Render()))
		} else {
			return strings.TrimSpace(ex.Render(node, base, ""))
		}
	})
}
