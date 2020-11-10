package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

func TestParse(t *testing.T) {
	t.Run("Header", func(t *testing.T) {
		testInclude := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Include) {
			must, doc := parseAsReader(t, s)
			must.Equal(1, len(doc.Headers))
			inc, ok := doc.Headers[0].(ast.Include)
			must.True(ok)
			return must, doc, inc
		}
		t.Run("Include", func(t *testing.T) {
			t.Run("succ", func(t *testing.T) {
				t.Run("comment in middle", func(t *testing.T) {
					must, doc, inc := testInclude(t, `include /*/*/ "abc"`)

					must.Equal(1, len(doc.Comments))
					bcmt, ok := doc.Comments[0].(ast.BlockComment)
					must.True(ok)
					checkTextSpan(must, bcmt, 0, 8, 0, 13, "/*/*/")

					checkTextSpan(must, inc.Keyword, 0, 0, 0, 7, "include")
					checkTextSpan(must, inc.Literal, 0, 14, 0, 19, `"abc"`)
					checkTextSpan(must, inc, 0, 0, 0, 19, `include "abc"`)
				})
				t.Run("consecutive", func(t *testing.T) {
					must, _, inc := testInclude(t, `include"abc"`)
					checkTextSpan(must, inc.Keyword, 0, 0, 0, 7, "include")
					checkTextSpan(must, inc.Literal, 0, 7, 0, 12, `"abc"`)
					checkTextSpan(must, inc, 0, 0, 0, 12, `include "abc"`)
				})
			})
		})
		t.Run("Cpp Include", func(t *testing.T) {
			t.Run("succ", func(t *testing.T) {
				must, doc, inc := testInclude(t, `cpp_include /*/*/ "abc"`)

				must.Equal(1, len(doc.Comments))
				bcmt, ok := doc.Comments[0].(ast.BlockComment)
				must.True(ok)
				checkTextSpan(must, bcmt, 0, 12, 0, 17, "/*/*/")

				checkTextSpan(must, inc.Keyword, 0, 0, 0, 11, "cpp_include")
				checkTextSpan(must, inc.Literal, 0, 18, 0, 23, `"abc"`)
				checkTextSpan(must, inc, 0, 0, 0, 23, `cpp_include "abc"`)
			})
			t.Run("consecutive", func(t *testing.T) {
				must, _, inc := testInclude(t, `cpp_include"abc"`)
				checkTextSpan(must, inc.Keyword, 0, 0, 0, 11, "cpp_include")
				checkTextSpan(must, inc.Literal, 0, 11, 0, 16, `"abc"`)
				checkTextSpan(must, inc, 0, 0, 0, 16, `cpp_include "abc"`)
			})
		})
		t.Run("namespace", func(t *testing.T) {
			testNamespace := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Namespace) {
				must, doc := parseAsReader(t, s)
				must.Equal(1, len(doc.Headers))
				ns, ok := doc.Headers[0].(ast.Namespace)
				must.True(ok)
				return must, doc, ns
			}
			t.Run("specified", func(t *testing.T) {
				t.Run("succ", func(t *testing.T) {
					t.Run("full", func(t *testing.T) {
						must, doc, ns := testNamespace(t, `namespace/*/*/go/*/*/go/*/*/(/*/*/a="b"/*/*/)`)
						must.Equal(5, len(doc.Comments))

						must.NotNil(ns.Annotations)
						must.Equal(1, len(ns.Annotations.Annotations))
						checkTextSpan(must, ns.Annotations.Annotations[0], 0, 34, 0, 39, `a = "b"`)

						checkTextSpan(must, ns.Language, 0, 14, 0, 16, "go")
						checkTextSpan(must, ns.Namespace, 0, 21, 0, 23, "go")
						checkTextSpan(must, ns, 0, 0, 0, 45, `namespace go go ( a = "b" )`)
					})
					t.Run("no annotations", func(t *testing.T) {
						t.Run("w parentheses", func(t *testing.T) {
							must, _, ns := testNamespace(t, `namespace go go()`)

							must.NotNil(ns.Annotations)
							must.Equal(0, len(ns.Annotations.Annotations))

							checkTextSpan(must, ns.Language, 0, 10, 0, 12, "go")
							checkTextSpan(must, ns.Namespace, 0, 13, 0, 15, "go")
							checkTextSpan(must, ns, 0, 0, 0, 17, `namespace go go`)
						})
						t.Run("wo parentheses", func(t *testing.T) {
							must, _, ns := testNamespace(t, `namespace go go`)

							must.Nil(ns.Annotations)

							checkTextSpan(must, ns.Language, 0, 10, 0, 12, "go")
							checkTextSpan(must, ns.Namespace, 0, 13, 0, 15, "go")
							checkTextSpan(must, ns, 0, 0, 0, 15, `namespace go go`)
						})
					})
				})
			})
			t.Run("glob", func(t *testing.T) {
				t.Run("succ", func(t *testing.T) {
					must, doc := parseAsReader(t, `namespace go go(a="b")
namespace * go`)
					must.Equal(2, len(doc.Headers))

					ns, ok := doc.Headers[1].(ast.Namespace)
					must.True(ok)
					must.Nil(ns.Annotations)

					checkTextSpan(must, ns.Language, 1, 10, 1, 11, "*")
					checkTextSpan(must, ns.Namespace, 1, 12, 1, 14, "go")
					checkTextSpan(must, ns, 1, 0, 1, 14, "namespace * go")
				})
				t.Run("error", func(t *testing.T) {
					must := require.New(t)
					_, err := ParseString("namespace * go()")
					must.Error(err)
				})
			})
		})
		t.Run("definition", func(t *testing.T) {
			testDefinition := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Definition) {
				must, doc := parseAsReader(t, s)
				must.Equal(1, len(doc.Definitions))
				return must, doc, doc.Definitions[0]
			}
			t.Run("const", func(t *testing.T) {
				testConst := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Const) {
					must, doc, def := testDefinition(t, s)
					cst, ok := def.(ast.Const)
					must.True(ok)
					return must, doc, cst
				}
				t.Run("consecutive", func(t *testing.T) {
					/* must,doc := */ parseAsReader(t, "const i32 a=1const i32 a=2")
				})
				t.Run("comma", func(t *testing.T) {
					/* must,doc := */ parseAsReader(t, "const i32 a=1,const i32 a=2;")
				})
				t.Run("identifier", func(t *testing.T) {
					must, _, cst := testConst(t, "const abc def = -1.2e+3")

					checkPos(must, cst.Start, 0, 0)
					typ, ok := cst.FieldType.(ast.Identifier)
					must.True(ok)
					checkTextSpan(must, typ, 0, 6, 0, 9, "abc")
					checkTextSpan(must, cst.Key, 0, 10, 0, 13, "def")
					checkPos(must, cst.Eq, 0, 14)
					num, ok := cst.Value.(ast.Number)
					must.True(ok)
					checkTextSpan(must, num, 0, 16, 0, 23, "-1.2e+3")
				})
				t.Run("base type", func(t *testing.T) {
					t.Run("wo annotation", func(t *testing.T) {
						t.Run("wo parenthese", func(t *testing.T) {
							must, _, cst := testConst(t, "const i32 a = 1")

							checkPos(must, cst.Start, 0, 0)
							typ, ok := cst.FieldType.(ast.BaseType)
							must.True(ok)
							checkTextSpan(must, typ.Identifier, 0, 6, 0, 9, "i32")
							must.Nil(typ.Annotations)
							checkTextSpan(must, cst.Key, 0, 10, 0, 11, "a")
							checkPos(must, cst.Eq, 0, 12)
							checkTextSpan(must, cst.Value, 0, 14, 0, 15, "1")
						})
						t.Run("w parenthese", func(t *testing.T) {
							must, _, cst := testConst(t, "const i32() a = 1")
							checkPos(must, cst.Start, 0, 0)
							typ, ok := cst.FieldType.(ast.BaseType)
							must.True(ok)
							checkTextSpan(must, typ.Identifier, 0, 6, 0, 9, "i32")
							must.NotNil(typ.Annotations)
							must.Equal(0, len(typ.Annotations.Annotations))
							checkSpan(must, typ.Annotations, 0, 9, 0, 11)
							checkTextSpan(must, cst.Key, 0, 12, 0, 13, "a")
							checkPos(must, cst.Eq, 0, 14)
							checkTextSpan(must, cst.Value, 0, 16, 0, 17, "1")
						})
					})
					t.Run("with annotations", func(t *testing.T) {
						must, _, cst := testConst(t, `const i32(k="v") a = 1`)
						checkPos(must, cst.Start, 0, 0)
						typ, ok := cst.FieldType.(ast.BaseType)
						must.True(ok)
						checkTextSpan(must, typ.Identifier, 0, 6, 0, 9, "i32")
						must.NotNil(typ.Annotations)
						must.Equal(1, len(typ.Annotations.Annotations))
						anno := typ.Annotations.Annotations[0]
						checkTextSpan(must, anno.Key, 0, 10, 0, 11, "k")
						checkPos(must, anno.Value.Start, 0, 11)
						checkTextSpan(must, anno.Value.Value, 0, 12, 0, 15, `"v"`)

						checkSpan(must, typ.Annotations, 0, 9, 0, 16)
						checkTextSpan(must, cst.Key, 0, 17, 0, 18, "a")
						checkPos(must, cst.Eq, 0, 19)
						checkTextSpan(must, cst.Value, 0, 21, 0, 22, "1")
					})
				})
				t.Run("container type", func(t *testing.T) {
					t.Run("map type", func(t *testing.T) {
						must, _, cst := testConst(t, `const map cpp_type"a"<i32,i32>a=1`)
						checkPos(must, cst.Start, 0, 0)
					})
					t.Run("set type", func(t *testing.T) {
						must, _, cst := testConst(t, `const set<bool>a=1`)
						checkPos(must, cst.Start, 0, 0)
					})
					t.Run("list type", func(t *testing.T) {
						must, _, cst := testConst(t, `const list<double>a=1`)
						checkPos(must, cst.Start, 0, 0)
					})
					t.Run("nested", func(t *testing.T) {
						must, _, cst := testConst(t, `const map cpp_type"a"<i32(k="v"),list <bool(k="v")>(k="v")>(k="v")a=1`)
						checkPos(must, cst.Start, 0, 0)
					})
				})
				t.Run("value type", func(t *testing.T) {
					t.Run("number", func(t *testing.T) {
						must, _, cst := testConst(t, `const i32 a=1`)
						checkPos(must, cst.Start, 0, 0)
						num, ok := cst.Value.(ast.Number)
						must.True(ok)
						checkTextSpan(must, num, 0, 12, 0, 13, "1")
					})
					t.Run("literal", func(t *testing.T) {
						must, _, cst := testConst(t, `const string a="a"`)
						checkPos(must, cst.Start, 0, 0)
						l, ok := cst.Value.(ast.Literal)
						must.True(ok)
						checkTextSpan(must, l, 0, 15, 0, 18, `"a"`)
					})
					t.Run("identifier", func(t *testing.T) {
						must, _, cst := testConst(t, `const i32 b=a`)
						checkPos(must, cst.Start, 0, 0)
						id, ok := cst.Value.(ast.Identifier)
						must.True(ok)
						checkTextSpan(must, id, 0, 12, 0, 13, "a")
					})
					t.Run("const list", func(t *testing.T) {
						must, _, cst := testConst(t, `const list<i32> a=[1 2+2,3;]`)
						checkPos(must, cst.Start, 0, 0)
						lst, ok := cst.Value.(ast.ConstList)
						must.True(ok)
						must.Equal(4, len(lst.Content))
					})
					t.Run("nested list", func(t *testing.T) {
						must, _, cst := testConst(t, `const list<list<i32>>a=[[][1];[1 2][1+3] [0;4]]`)
						checkPos(must, cst.Start, 0, 0)
						lst, ok := cst.Value.(ast.ConstList)
						must.True(ok)
						must.Equal(5, len(lst.Content))

						lst2, ok := lst.Content[0].ConstValue.(ast.ConstList)
						must.True(ok)
						must.Equal(0, len(lst2.Content))

						lst2, ok = lst.Content[3].ConstValue.(ast.ConstList)
						must.True(ok)
						must.Equal(2, len(lst2.Content))
					})
					t.Run("const map", func(t *testing.T) {
						must, _, cst := testConst(t, `const map<i32,i32> a={1:2+2:3,3:4 4:5;5:6}`)
						checkPos(must, cst.Start, 0, 0)
						lst, ok := cst.Value.(ast.ConstMap)
						must.True(ok)
						must.Equal(5, len(lst.Content))
					})
					t.Run("nested map", func(t *testing.T) {
						must, _, cst := testConst(t, `const map<i32,map<i32,i32>> a={1:{2:3 4:5}}`)
						checkPos(must, cst.Start, 0, 0)
						mp, ok := cst.Value.(ast.ConstMap)
						must.True(ok)
						must.Equal(1, len(mp.Content))

						mp2, ok := mp.Content[0].Value.(ast.ConstMap)
						must.True(ok)
						must.Equal(2, len(mp2.Content))
					})
				})
			})
			t.Run("type definition", func(t *testing.T) {
				t.Run("typedef", func(t *testing.T) {
					testTypeDef := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.TypeDef) {
						must, doc, def := testDefinition(t, s)
						td, ok := def.(ast.TypeDef)
						must.True(ok)
						return must, doc, td
					}
					t.Run("", func(t *testing.T) {
						must, _, td := testTypeDef(t, `typedef i32 abc(k="v");`)
						checkPos(must, td.Start, 0, 0)
					})
				})
				t.Run("enum", func(t *testing.T) {
					testEnum := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Enum) {
						must, doc, def := testDefinition(t, s)
						enum, ok := def.(ast.Enum)
						must.True(ok)
						return must, doc, enum
					}
					t.Run("", func(t *testing.T) {
						must, _, enum := testEnum(t, `enum a{A=1 B,C;D(k="v")E}(k="v")`)
						checkPos(must, enum.Start, 0, 0)
						must.Equal(5, len(enum.List))
					})
				})
				t.Run("senum", func(t *testing.T) {
					testSenum := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Senum) {
						must, doc, def := testDefinition(t, s)
						senum, ok := def.(ast.Senum)
						must.True(ok)
						return must, doc, senum
					}
					t.Run("", func(t *testing.T) {
						must, _, senum := testSenum(t, `senum a{"a""b" "c","d";"e"}(k="v")`)
						checkPos(must, senum.Start, 0, 0)
						must.Equal(5, len(senum.List))
					})
				})
				t.Run("struct", func(t *testing.T) {
					testStruct := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Struct) {
						must, doc, def := testDefinition(t, s)
						st, ok := def.(ast.Struct)
						must.True(ok)
						return must, doc, st
					}
					t.Run("", func(t *testing.T) {
						must, _, st := testStruct(t, `union A xsd_all{}(k="v")`)
						checkTextSpan(must, st.Head, 0, 0, 0, 5, "union")
					})
					t.Run("", func(t *testing.T) {
						must, _, st := testStruct(t, `struct A{
	i32& F1 xsd_attrs {
		i8 A = 1
	},
	1: required bool F2 = true xsd_optional
	2: optional i8 F3 xsd_nillable
}`)
						checkTextSpan(must, st.Head, 0, 0, 0, 6, "struct")
						must.Equal(3, len(st.List))
					})
				})
				t.Run("exception", func(t *testing.T) {
					testException := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Exception) {
						must, doc, def := testDefinition(t, s)
						st, ok := def.(ast.Exception)
						must.True(ok)
						return must, doc, st
					}
					t.Run("", func(t *testing.T) {
						must, _, st := testException(t, `exception A {
	i32 a
} (k="v")`)
						checkPos(must, st.Start, 0, 0)
						checkTextSpan(must, st.Identifier, 0, 10, 0, 11, "A")
					})
				})
			})
			t.Run("service", func(t *testing.T) {
				testService := func(t *testing.T, s string) (*require.Assertions, *ast.Document, ast.Service) {
					must, doc, def := testDefinition(t, s)
					svc, ok := def.(ast.Service)
					must.True(ok)
					return must, doc, svc
				}
				t.Run("", func(t *testing.T) {
					must, _, svc := testService(t, `service A{}(k="v")`)
					checkPos(must, svc.Start, 0, 0)
					checkTextSpan(must, svc.Identifier, 0, 8, 0, 9, "A")
				})
				t.Run("", func(t *testing.T) {
					must, _, svc := testService(t, `service B extends A{
	oneway void Foo (i32 A = 1; 2: bool B) (k="v"),
    i32 Bar () throws () (k="v");
}(k="v")`)
					checkPos(must, svc.Start, 0, 0)
					checkTextSpan(must, svc.Identifier, 0, 8, 0, 9, "B")

					must.Equal(2, len(svc.FunctionList))
				})
			})
		})
	})
}
