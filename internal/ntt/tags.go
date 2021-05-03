package ntt

import (
	"github.com/nokia/ntt/internal/ttcn3/ast"
)

// Tags will return all named nodes (defintions, enums, ...) from a file.
func (suite *Suite) Tags(file string) (*ParseInfo, []ast.Node) {
	mod := suite.Parse(file)
	if mod == nil || mod.Module == nil {
		return mod, nil
	}

	t := make([]ast.Node, 0, len(mod.Module.Defs)*2)
	ast.Inspect(mod.Module, func(n ast.Node) bool {
		if n == nil {
			return false
		}

		switch n := n.(type) {
		case *ast.Module:
			t = append(t, n)
			return true

		case *ast.ImportDecl:
			return false

		case *ast.FriendDecl:
			return false

		case *ast.Field:
			t = append(t, n)
			return true

		case *ast.PortTypeDecl:
			t = append(t, n)
			return false

		case *ast.ComponentTypeDecl:
			t = append(t, n)
			return true

		case *ast.StructTypeDecl:
			t = append(t, n)
			return true

		case *ast.EnumTypeDecl:
			t = append(t, n)
			for _, e := range n.Enums {
				t = append(t, e)
			}
			return false

		case *ast.EnumSpec:
			for _, e := range n.Enums {
				t = append(t, e)
			}
			return false

		case *ast.BehaviourTypeDecl:
			t = append(t, n)
			return false

		case *ast.Declarator:
			t = append(t, n)
			return false

		case *ast.FormalPar:
			t = append(t, n)
			return false

		case *ast.TemplateDecl:
			t = append(t, n)
			return true

		case *ast.FuncDecl:
			t = append(t, n)
			return true

		case *ast.SignatureDecl:
			t = append(t, n)
			return false

		default:
			return true
		}
	})
	return mod, t
}
