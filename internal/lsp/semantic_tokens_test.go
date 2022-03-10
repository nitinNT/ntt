package lsp_test

import (
	"fmt"
	"testing"

	"github.com/nokia/ntt/internal/lsp"
	"github.com/nokia/ntt/internal/lsp/protocol"
	"github.com/nokia/ntt/internal/ntt"
	"github.com/nokia/ntt/ttcn3"
	"github.com/stretchr/testify/assert"
)

func SetRange(startLine uint32, startCol uint32, endLine uint32, endCol uint32) *protocol.Range {
	return &protocol.Range{
		Start: protocol.Position{Line: startLine, Character: startCol},
		End:   protocol.Position{Line: endLine, Character: endCol}}
}
func DefaultRange(tree *ttcn3.Tree) *protocol.Range {
	modEnd := tree.Position(tree.Root.End())
	return &protocol.Range{Start: protocol.Position{Line: 0, Character: 0}, End: protocol.Position{Line: uint32(modEnd.Line - 1), Character: uint32(modEnd.Column - 1)}}
}
func generateTokenList(t *testing.T, suite *ntt.Suite, rng *protocol.Range) *protocol.SemanticTokens {
	lsp.EmitKeywords = true
	name := fmt.Sprintf("%s_Module_0.ttcn3", t.Name())
	tree := ttcn3.ParseFile(name)

	// index
	db := &ttcn3.DB{}
	if srcs, _ := suite.Sources(); len(srcs) > 0 {
		for _, src := range srcs {
			db.Index(src)
		}
	}

	if rng == nil {
		rng = DefaultRange(tree)
	}

	return lsp.NewSemanticTokensFromCurrentModule(tree, db, suite, name, *rng)
}

func TestComponentDef(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
		import from TestComponentDef_Module_1 all;
		type component B0 extends C0, TestComponentDef_Module_1.C1 {
			var integer i := 1;
			timer t1 := 2.0;
			port P p;
		}
		function f() runs on TestComponentDef_Module_1.C0 system B0 return integer {}
	}`, `module TestComponentDef_Module_1
	{
		type component C0 {}
		type component C1 {}
		type port P message {
            inout charstring
        }
	}`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 2, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Keyword), 0,
			0, 5, 25, uint32(lsp.Namespace), 0,
			0, 26, 3, uint32(lsp.Keyword), 0, //all
			1, 2, 4, uint32(lsp.Keyword), 0,
			0, 5, 9, uint32(lsp.Keyword), 0,
			0, 10, 2, uint32(lsp.Class), uint32(lsp.Definition),
			0, 3, 7, uint32(lsp.Keyword), 0,
			0, 8, 2, uint32(lsp.Class), 0, //C0
			0, 4, 25, uint32(lsp.Namespace), 0,
			0, 26, 2, uint32(lsp.Class), 0,
			1, 3, 3, uint32(lsp.Keyword), 0,
			0, 4, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 1, uint32(lsp.Variable), uint32(lsp.Declaration),
			1, 3, 5, uint32(lsp.Keyword), 0,
			0, 6, 2, uint32(lsp.Variable), uint32(lsp.Declaration),
			1, 3, 4, uint32(lsp.Keyword), 0,
			0, 5, 1, uint32(lsp.Interface), 0,
			0, 2, 1, uint32(lsp.Variable), uint32(lsp.Declaration),

			2, 2, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Function), uint32(lsp.Definition),
			0, 4, 4, uint32(lsp.Keyword), 0,
			0, 5, 2, uint32(lsp.Keyword), 0,
			0, 3, 25, uint32(lsp.Namespace), 0,
			0, 26, 2, uint32(lsp.Class), 0,
			0, 3, 6, uint32(lsp.Keyword), 0,
			0, 7, 2, uint32(lsp.Class), 0,
			0, 3, 6, uint32(lsp.Keyword), 0,
			0, 7, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
		}, list.Data)
}

func TestFullModuleKwTypeId(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
        type record MyRec {
			integer i
		}
        type record length(0..2) of integer RoI;
        type set MySet {
			integer i
		}
        function f() return integer {}
	  }`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 8, 4, uint32(lsp.Keyword), 0,
			0, 5, 6, uint32(lsp.Keyword), 0,
			0, 7, 5, uint32(lsp.Struct), uint32(lsp.Definition),
			1, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			2, 8, 4, uint32(lsp.Keyword), 0,
			0, 5, 6, uint32(lsp.Keyword), 0,
			0, 7, 6, uint32(lsp.Keyword), 0,
			0, 13, 2, uint32(lsp.Keyword), 0,
			0, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 3, uint32(lsp.Type), uint32(lsp.Definition),
			1, 8, 4, uint32(lsp.Keyword), 0,
			0, 5, 3, uint32(lsp.Keyword), 0,
			0, 4, 5, uint32(lsp.Struct), uint32(lsp.Definition),
			1, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			2, 8, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Function), uint32(lsp.Definition),
			0, 4, 6, uint32(lsp.Keyword), 0,
			0, 7, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary)}, list.Data)
}

func TestConstAndTemplDecl(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
		type integer X;
		type X Y;
		type record MyRec {
			integer i
		}
        const float pi := 3.14;
		template MyRec a_rec1 := *;
		template (value) MyRec a_rec2(template integer p_i) := {
			i := p_i
		}
 	  }`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 2, 4, uint32(lsp.Keyword), 0,
			0, 5, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 1, uint32(lsp.Type), uint32(lsp.Definition),
			1, 2, 4, uint32(lsp.Keyword), 0,
			0, 5, 1, uint32(lsp.Type), 0,
			0, 2, 1, uint32(lsp.Type), uint32(lsp.Definition),
			1, 2, 4, uint32(lsp.Keyword), 0,
			0, 5, 6, uint32(lsp.Keyword), 0,
			0, 7, 5, uint32(lsp.Struct), uint32(lsp.Definition),
			1, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			2, 8, 5, uint32(lsp.Keyword), 0,
			0, 6, 5, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 6, 2, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			1, 2, 8, uint32(lsp.Keyword), 0,
			0, 9, 5, uint32(lsp.Type), uint32(lsp.Undefined),
			0, 6, 6, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			1, 2, 8, uint32(lsp.Keyword), 0,
			0, 10, 5, uint32(lsp.Keyword), 0,
			0, 7, 5, uint32(lsp.Type), uint32(lsp.Undefined),
			0, 6, 6, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			0, 7, 8, uint32(lsp.Keyword), 0,
			0, 9, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 3, uint32(lsp.Parameter), uint32(lsp.Declaration),
			1, 8, 3, uint32(lsp.Parameter), 0,
		}, list.Data)
}

func TestModuleIds(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
		import from TestModuleIds_Module_1 all;
		type TestModuleIds_Module_1.Byte MyByte;
		function f(TestModuleIds_Module_1.Byte pb) {
			var TestModuleIds_Module_1.MyRec r;
			r.i := pb;
		}
	}`,
		`module TestModuleIds_Module_1 {
		type integer Byte;
		type record MyRec {integer i};
	}`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 2, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Keyword), 0,
			0, 5, 22, uint32(lsp.Namespace), 0,
			0, 23, 3, uint32(lsp.Keyword), 0, //all
			1, 2, 4, uint32(lsp.Keyword), 0, // type line 3
			0, 5, 22, uint32(lsp.Namespace), 0, //
			0, 23, 4, uint32(lsp.Type), 0,
			0, 5, 6, uint32(lsp.Type), uint32(lsp.Definition),
			1, 2, 8, uint32(lsp.Keyword), 0, // line 4
			0, 9, 1, uint32(lsp.Function), uint32(lsp.Definition),
			0, 2, 22, uint32(lsp.Namespace), 0,
			0, 23, 4, uint32(lsp.Type), 0,
			0, 5, 2, uint32(lsp.Parameter), uint32(lsp.Declaration), // pb
			1, 3, 3, uint32(lsp.Keyword), 0, //var line 5
			0, 4, 22, uint32(lsp.Namespace), 0,
			0, 23, 5, uint32(lsp.Type), uint32(lsp.Undefined),
			0, 6, 1, uint32(lsp.Variable), uint32(lsp.Declaration),
			1, 3, 1, uint32(lsp.Variable), 0,
			0, 7, 2, uint32(lsp.Parameter), 0,
		}, list.Data)
}

func TestDeltaModuleKwTypeId(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
        type record MyRec {
			integer i
		}
        type record length(0..2) of integer RoI;
        type set MySet {
			integer i
		}
        function f() return integer {}
	  }`)

	list := generateTokenList(t, suite, SetRange(5, 8, 8, 3))

	assert.Equal(t,
		[]uint32{
			0, 0, 4, uint32(lsp.Keyword), 0,
			0, 5, 6, uint32(lsp.Keyword), 0,
			0, 7, 6, uint32(lsp.Keyword), 0,
			0, 13, 2, uint32(lsp.Keyword), 0,
			0, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 3, uint32(lsp.Type), uint32(lsp.Definition),
			1, 8, 4, uint32(lsp.Keyword), 0,
			0, 5, 3, uint32(lsp.Keyword), 0,
			0, 4, 5, uint32(lsp.Struct), uint32(lsp.Definition),
			1, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary)}, list.Data)
}

func TestUniversalCharstring(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
        type record MyRec {
			universal charstring  uch,
			charstring ch
		}
		type universal charstring UcharT;
		function f() return universal charstring {}
	  }`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 8, 4, uint32(lsp.Keyword), 0,
			0, 5, 6, uint32(lsp.Keyword), 0,
			0, 7, 5, uint32(lsp.Struct), uint32(lsp.Definition),
			1, 3, 20, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			1, 3, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			2, 2, 4, uint32(lsp.Keyword), 0,
			0, 5, 20, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 21, 6, uint32(lsp.Type), uint32(lsp.Definition),
			1, 2, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Function), uint32(lsp.Definition),
			0, 4, 6, uint32(lsp.Keyword), 0,
			0, 7, 20, uint32(lsp.Type), uint32(lsp.DefaultLibrary)}, list.Data)
}

func TestPortTypeDeclSemTok(t *testing.T) {
	suite := buildSuite(t, `module Test
{
	import from TestPortTypeDeclSemTok_Module_1 all;
	type charstring AddrType;
	type boolean Type1;
	type port MyPort message {
		in integer, universal charstring
		out Type1, TestPortTypeDeclSemTok_Module_1.Type2
		map param (in integer pi, boolean pb);
		address AddrType
	}
}`, `module TestPortTypeDeclSemTok_Module_1
{
	type integer Type2;
}`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 1, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Keyword), 0,
			0, 5, 31, uint32(lsp.Namespace), 0,
			0, 32, 3, uint32(lsp.Keyword), 0,
			1, 1, 4, uint32(lsp.Keyword), 0,
			0, 5, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 11, 8, uint32(lsp.Type), uint32(lsp.Definition),
			1, 1, 4, uint32(lsp.Keyword), 0,
			0, 5, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 5, uint32(lsp.Type), uint32(lsp.Definition),
			1, 1, 4, uint32(lsp.Keyword), 0,
			0, 5, 4, uint32(lsp.Keyword), 0,
			0, 5, 6, uint32(lsp.Interface), uint32(lsp.Definition),
			0, 7, 7, uint32(lsp.Keyword), 0,
			1, 2, 2, uint32(lsp.Keyword), 0,
			0, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 9, 20, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			1, 2, 3, uint32(lsp.Keyword), 0,
			0, 4, 5, uint32(lsp.Type), 0,
			0, 7, 31, uint32(lsp.Namespace), 0,
			0, 32, 5, uint32(lsp.Type), 0,
			1, 2, 3, uint32(lsp.Keyword), 0,
			0, 4, 5, uint32(lsp.Keyword), 0,
			0, 7, 2, uint32(lsp.Keyword), 0,
			0, 3, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 2, uint32(lsp.Parameter), uint32(lsp.Declaration),
			0, 4, 7, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 8, 2, uint32(lsp.Parameter), uint32(lsp.Declaration),
			1, 2, 7, uint32(lsp.Keyword), 0,
			0, 8, 8, uint32(lsp.Type), 0}, list.Data)
}

func TestFunctionParams(t *testing.T) {
	suite := buildSuite(t, `module Test
{
	import from TestFunctionParams_Module_1 all;

	function f(charstring p_ch) return template charstring {
		var R v_r1 := {p_ch};
		var template R t_r2 := {ch1 := p_ch};
		var charstring ch := p_ch;
		ch := p_ch;
		v_r1.ch1 := p_ch;
		t_r2 := {p_ch};
		t_r2 := {ch1 := p_ch};
		f(p_ch := p_ch);
		return R:{p_ch}
	}
}`, `module TestFunctionParams_Module_1
{
	type record R {
		charstring ch1
	};
}`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 1, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Keyword), 0,
			0, 5, 27, uint32(lsp.Namespace), 0,
			0, 28, 3, uint32(lsp.Keyword), 0, //all
			2, 1, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Function), uint32(lsp.Definition),
			0, 2, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 11, 4, uint32(lsp.Parameter), uint32(lsp.Declaration),
			0, 6, 6, uint32(lsp.Keyword), 0, // return
			0, 7, 8, uint32(lsp.Keyword), 0, //  template
			0, 9, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary), // charstring
			1, 2, 3, uint32(lsp.Keyword), 0,
			0, 4, 1, uint32(lsp.Type), 0,
			0, 2, 4, uint32(lsp.Variable), uint32(lsp.Declaration),
			0, 9, 4, uint32(lsp.Parameter), 0,
			1, 2, 3, uint32(lsp.Keyword), 0,
			0, 4, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Type), 0,
			0, 2, 4, uint32(lsp.Variable), uint32(lsp.Declaration),
			0, 16, 4, uint32(lsp.Parameter), 0, //var template R tr2
			1, 2, 3, uint32(lsp.Keyword), 0,
			0, 4, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 11, 2, uint32(lsp.Variable), uint32(lsp.Declaration),
			0, 6, 4, uint32(lsp.Parameter), 0,
			1, 2, 2, uint32(lsp.Variable), 0,
			0, 6, 4, uint32(lsp.Parameter), 0,
			1, 2, 4, uint32(lsp.Variable), 0,
			0, 12, 4, uint32(lsp.Parameter), 0, // v_r1.ch1 := p_ch;
			1, 2, 4, uint32(lsp.Variable), 0,
			0, 9, 4, uint32(lsp.Parameter), 0,
			1, 2, 4, uint32(lsp.Variable), 0,
			0, 16, 4, uint32(lsp.Parameter), 0, // t_r2 := {ch1 := p_ch};
			1, 2, 1, uint32(lsp.Function), 0,
			0, 2, 4, uint32(lsp.Parameter), 0,
			0, 8, 4, uint32(lsp.Parameter), 0, // f(p_ch := p_ch);
			1, 2, 6, uint32(lsp.Keyword), 0,
			0, 7, 1, uint32(lsp.Type), 0,
			0, 3, 4, uint32(lsp.Parameter), 0}, list.Data)
}

func TestTemplateAndConst(t *testing.T) {
	suite := buildSuite(t, `module Test
{
	import from TestTemplateAndConst_Module_1 all;

	const charstring c_ch2 := c_ch;
	template charstring a_ch2 := a_ch;
	template R a_r2(template charstring p_ch := a_ch) modifies a_r := {
		ch1 := c_ch,
		ch2 := p_ch
	}
	template R a_r3 := a_r(p_ch := a_ch2);

}`, `module TestTemplateAndConst_Module_1
{
	type record R {
		charstring ch1,
		charstring ch2 optional
	};
	const charstring c_ch := "xx";
	template charstring a_ch := ?;
	template (omit) R a_r(template charstring p_ch) := {
		ch1 := p_ch,
		ch2 := omit
	}
}`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 1, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Keyword), 0,
			0, 5, 29, uint32(lsp.Namespace), 0,
			0, 30, 3, uint32(lsp.Keyword), 0, //all
			2, 1, 5, uint32(lsp.Keyword), 0,
			0, 6, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 11, 5, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			0, 9, 4, uint32(lsp.Variable), uint32(lsp.Readonly),
			1, 1, 8, uint32(lsp.Keyword), 0,
			0, 9, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 11, 5, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			0, 9, 4, uint32(lsp.Variable), uint32(lsp.Readonly),

			1, 1, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Type), 0,
			0, 2, 4, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			0, 5, 8, uint32(lsp.Keyword), 0,
			0, 9, 10, uint32(lsp.Type), uint32(lsp.DefaultLibrary),
			0, 11, 4, uint32(lsp.Parameter), uint32(lsp.Declaration),
			0, 8, 4, uint32(lsp.Variable), uint32(lsp.Readonly),
			0, 6, 8, uint32(lsp.Keyword), 0,
			0, 9, 3, uint32(lsp.Variable), uint32(lsp.Readonly),
			1, 9, 4, uint32(lsp.Variable), uint32(lsp.Readonly),
			1, 9, 4, uint32(lsp.Parameter), 0,

			2, 1, 8, uint32(lsp.Keyword), 0,
			0, 9, 1, uint32(lsp.Type), 0,
			0, 2, 4, uint32(lsp.Variable), uint32(lsp.Declaration | lsp.Readonly),
			0, 8, 3, uint32(lsp.Variable), uint32(lsp.Readonly),
			0, 4, 4, uint32(lsp.Parameter), 0,
			0, 8, 5, uint32(lsp.Variable), uint32(lsp.Readonly),
		}, list.Data)
}

func TestComponentVars(t *testing.T) {
	suite := buildSuite(t, `module Test
    {
		import from TestComponentVars_Module_1 all;
		testcase tc() runs on C3 system C0 {
			vC0_C1 := C0.create();
			vi_C2 := 1;
			log(vi_C2);
			p0_C0.send("hello");
			vC0_C1.start(f(p_i := vi_C2));
		}
	}`, `module TestComponentVars_Module_1
	{
		type component C0 {
			port P p0_C0;
		}
		type component C1 {
			var C0 vC0_C1;
		}
		type component C2 extends C1 {
			var integer vi_C2;
		}
		type component C3 extends C2, C0 {}
		type port P message {
            inout charstring
        }
		function f(integer p_i) runs on C0 {}
	}`)

	list := generateTokenList(t, suite, nil)

	assert.Equal(t,
		[]uint32{
			0, 0, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Namespace), uint32(lsp.Definition),
			2, 2, 6, uint32(lsp.Keyword), 0,
			0, 7, 4, uint32(lsp.Keyword), 0,
			0, 5, 26, uint32(lsp.Namespace), 0,
			0, 27, 3, uint32(lsp.Keyword), 0, //all
			1, 2, 8, uint32(lsp.Keyword), 0,
			0, 9, 2, uint32(lsp.Function), uint32(lsp.Definition),
			0, 5, 4, uint32(lsp.Keyword), 0,
			0, 5, 2, uint32(lsp.Keyword), 0,
			0, 3, 2, uint32(lsp.Class), 0, //C3
			0, 3, 6, uint32(lsp.Keyword), 0,
			0, 7, 2, uint32(lsp.Class), 0, //C0
			1, 3, 6, uint32(lsp.Property), 0, //vC0_C1 :=
			0, 10, 2, uint32(lsp.Class), 0,
			/*0, 3, 6, uint32(lsp.Keyword), 0,*/ //create not a kw. TODO: check
			1, 3, 5, uint32(lsp.Property), 0,    //vi_C2
			1, 3, 3, uint32(lsp.Function), uint32(lsp.DefaultLibrary), //log
			0, 4, 5, uint32(lsp.Property), 0, //vi_C2
			1, 3, 5, uint32(lsp.Property), 0, // p0_C0.
			/*0, 6, 4, uint32(lsp.Keyword), 0, */ // send not a kw. TODO: check
			1, 3, 6, uint32(lsp.Property), 0,
			/*0, 7, 5, uint32(lsp.Keyword), 0,*/  // start( not a kw. TODO: check
			0, 6 + 7, 1, uint32(lsp.Function), 0, // f
			0, 2, 3, uint32(lsp.Parameter), 0, // p_i :=
			0, 7, 5, uint32(lsp.Property), 0,
		}, list.Data)
}
