package render

import (
	"github.com/zbeaver/cafe/pkg/vui"
)

type Select struct{}

func (r *Select) Style(base styling, n vui.INode) styling {
	elm, ok := n.(vui.Elementary)
	if !ok {
		return base
	}

	return TransformFrom(base)(elm.Style())
}

func (r *Select) Render(n vui.INode, s styling, child string) string {
	return s.SetString(child).Render()
}