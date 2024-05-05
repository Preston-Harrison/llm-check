package go_checks

import "go/ast"

type visitor struct {
	path      []ast.Node
	filter    func(ast.Node) bool
	processor func([]ast.Node)
}

func newVisitor(filter func(ast.Node) bool, processor func([]ast.Node)) visitor {
	return visitor{nil, filter, processor}
}

// If `node` is nil, it means the visitor is going back up the tree. Otherwise
// (since Walk is depth-first), it means the visitor is going down the tree.
func (v *visitor) applyNodeToPath(node ast.Node) {
	if node == nil && len(v.path) != 0 {
		v.path = v.path[:len(v.path)-1]
	} else {
		v.path = append(v.path, node)
	}
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	v.applyNodeToPath(node)
	if v.filter(node) {
		v.processor(v.path)
	}
	return v
}
