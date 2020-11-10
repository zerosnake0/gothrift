%{

package parser

import (
	//"log"

	"github.com/zerosnake0/gothrift/pkg/ast"
)

%}


%union {
	document        ast.Document

	headers         []ast.Header
	header          ast.Header

	definitions	[]ast.Definition
	definition	ast.Definition

	fieldType	ast.FieldType

	simpleContainerType ast.SimpleContainerType
	cppType		*ast.CppType

	constValue	ast.ConstValue

	annotations	*ast.Annotations
	annotationList  []ast.Annotation
	annotation	ast.Annotation
	annotationValue ast.AnnotationValue

	constListContent []ast.ConstListItem
	constMapContent  []ast.ConstMapItem

	enumDefList	[]ast.EnumDef
	enumDef		ast.EnumDef
	enumValue	ast.EnumValue

	senumDefList	[]ast.SenumDef
        senumDef	ast.SenumDef

        fieldList	[]ast.Field
        field		ast.Field
        fieldIdentifier *ast.FieldIdentifier
        fieldValue      *ast.FieldValue

        xsdAttrs	*ast.XsdAttrs

	extends		*ast.Extends

	functionList	[]ast.Function
	function	ast.Function
	functionType	ast.FunctionType

	throws		*ast.Throws

	literal		ast.Literal
	identifier	*ast.Identifier
	number 		ast.Number
}

%type <document> Document

%type <headers> Headers
%type <header>	Header
		Include
		Namespace

%type <definitions> Definitions
%type <definition>	Definition
			Const
			TypeDefinition
			Typedef
			Enum
			Senum
			Struct
			Exception
			Service

%type <fieldType>	FieldType
			BaseType
			ContainerType

%type <simpleContainerType>	SimpleContainerType
				MapType
				SetType
				ListType

%type <cppType> CppType

%type <constValue>	ConstValue
			ConstList
			ConstMap

%type <constListContent> ConstListContent
%type <constMapContent> ConstMapContent

%type <enumDefList>	EnumDefList
%type <enumDef>		EnumDef
%type <enumValue>	EnumValue

%type <senumDefList>	SenumDefList
%type <senumDef>	SenumDef

%type <fieldList>	FieldList
%type <field>		Field
%type <fieldIdentifier>	FieldIdentifier
%type <fieldValue>	FieldValue

%type <annotations>	TypeAnnotations
%type <annotationList>	TypeAnnotationList
%type <annotation>	TypeAnnotation
%type <annotationValue>	TypeAnnotationValue

%type <xsdAttrs>	XsdAttributes

%type <extends>		Extends

%type <functionList>	FunctionList
%type <function>	Function
%type <functionType>	FunctionType

%type <throws>		Throws

%type <identifier>	CommaOrSemicolonOptional
			FieldRequiredness
			FieldReference
			XsdAll
			XsdOptional
			XsdNillable
			Oneway

%token <identifier>	'(' ')'
			'='
			',' ';'
			'*'
			'<' '>'
			'[' ']'
			'{' '}' ':'
			'&'
			INCLUDE
			NAMESPACE
			CONST
			BASETYPE
			MAP SET LIST
			CPPTYPE
			TYPEDEF ENUM SENUM
			STRUCTHEAD
			REQUIRED OPTIONAL
			XSDALL XSDOPTIONAL XSDNILLABLE XSDNAMESPACE XSDATTRS
			EXCEPTION
			SERVICE EXTENDS
			ONEWAY VOID THROWS
			Identifier
%token <literal>	Literal
%token <number>		Number

%start Root

%%

Root:
	Document
	{
		setExprValDoc(exprlex, $1)
	}

Document:
	Headers Definitions
	{
		$$.Headers = $1
		$$.Definitions = $2
	}

Headers:
	Headers Header
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

Header:
	Include // include & cpp_include
	{
		$$ = $1
	}
|	Namespace
	{
		$$ = $1
	}

Include:
	INCLUDE Literal
	{
		$$ = ast.Include{
			Keyword: *$1,
			Literal: $2,
		}
	}

Namespace:
	NAMESPACE Identifier Identifier TypeAnnotations
	{
		$$ = ast.Namespace{
			Start:       $1.Start,
			Language:    *$2,
			Namespace:   *$3,
			Annotations: $4,
		}
	}
|	NAMESPACE '*' Identifier
	{
		$$ = ast.Namespace{
			Start:       $1.Start,
			Language:    *$2,
			Namespace:   *$3,
		}
	}

Definitions:
	Definitions Definition
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

Definition:
	Const
	{
		$$ = $1
	}
|	TypeDefinition
	{
		$$ = $1
	}
|	Service
	{
		$$ = $1
	}

Const:
	CONST FieldType Identifier '=' ConstValue CommaOrSemicolonOptional
	{
		$$ = ast.Const{
			Start:     $1.Start,
			FieldType: $2,
			Key:       *$3,
			Eq:        $4.Start,
			Value:     $5,
			End:       $6,
		}
	}

FieldType:
	Identifier
	{
		$$ = *$1
	}
|	BaseType
	{
		$$ = $1
	}
|	ContainerType
	{
		$$ = $1
	}

BaseType:
	BASETYPE TypeAnnotations
	{
		$$ = ast.BaseType{
			Identifier:  *$1,
			Annotations: $2,
		}
	}


ContainerType:
	SimpleContainerType TypeAnnotations
	{
		$$ = ast.ContainerType{
			SimpleContainerType: $1,
			Annotations:         $2,
		}
	}

SimpleContainerType:
	MapType
	{
		$$ = $1
	}
|	SetType
	{
		$$ = $1
	}
|	ListType
	{
		$$ = $1
	}

MapType:
	MAP CppType '<' FieldType ',' FieldType '>'
	{
		$$ = ast.MapType{
			Start:    $1.Start,
			CppType:  $2,
			LChevron: $3.Start,
			Key:      $4,
			Comma:    $5.Start,
                	Value:    $6,
                       	RChevron: $7.Start,
		}
	}

SetType:
	SET CppType '<' FieldType '>'
	{
		$$ = ast.SetType{
			Start:    $1.Start,
                        CppType:  $2,
                        LChevron: $3.Start,
                        Elem:     $4,
                        RChevron: $5.Start,
		}
	}

ListType:
	LIST '<' FieldType '>' CppType
	{
		$$ = ast.ListType{
			Start:    $1.Start,
			LChevron: $2.Start,
			Elem:     $3,
			RChevron: $4.Start,
			CppType:  $5,
		}
	}

CppType:
	CPPTYPE Literal
	{
		$$ = &ast.CppType{
			Start:   $1.Start,
			Literal: $2,
		}
	}
|	// empty
	{
		$$ = nil
	}

ConstValue:
	Number
	{
		$$ = $1
	}
|	Literal
	{
		$$ = $1
	}
|	Identifier
	{
		$$ = *$1
	}
|	ConstList
	{
		$$ = $1
	}
|	ConstMap
	{
		$$ = $1
	}

ConstList:
	'[' ConstListContent ']'
	{
		$$ = ast.ConstList{
			Node: ast.Node{
				Start: $1.Start,
				End:   $3.End,
			},
			Content: $2,
		}
	}

ConstListContent:
	ConstListContent ConstValue CommaOrSemicolonOptional
	{
		$$ = append($1, ast.ConstListItem{
			ConstValue: $2,
			End:        $3,
		})
	}
|	// empty
	{
		$$ = nil
	}

ConstMap:
	'{' ConstMapContent '}'
	{
		$$ = ast.ConstMap{
			Node: ast.Node{
				Start: $1.Start,
				End:   $3.End,
			},
			Content: $2,
		}
	}

ConstMapContent:
	ConstMapContent ConstValue ':' ConstValue CommaOrSemicolonOptional
	{
		$$ = append($1, ast.ConstMapItem{
			Key:   $2,
			Colon: $3.Start,
			Value: $4,
			End:   $5,
		})
	}
|	// empty
	{
		$$ = nil
	}

TypeAnnotations:
	'(' TypeAnnotationList ')'
	{
		$$ = &ast.Annotations{}
		$$.Start = $1.Start
		$$.Annotations = $2
		$$.End = $3.End
	}
|	// empty
	{
		$$ = nil
	}

TypeAnnotationList:
	TypeAnnotationList TypeAnnotation
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

TypeAnnotation:
	Identifier TypeAnnotationValue CommaOrSemicolonOptional
	{
		$$.Key = *$1
		$$.Value = $2
		$$.End = $3
	}

TypeAnnotationValue:
	'=' Literal
	{
		$$.Start = $1.Start
		$$.Value = $2
	}

CommaOrSemicolonOptional:
	','
	{
		$$ = $1
	}
|	';'
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

TypeDefinition:
	Typedef
	{
		$$ = $1
	}
|	Enum
	{
		$$ = $1
	}
|	Senum
	{
		$$ = $1
	}
|	Struct
	{
		$$ = $1
	}
|	Exception
	{
		$$ = $1
	}

Typedef:
	TYPEDEF FieldType Identifier TypeAnnotations CommaOrSemicolonOptional
	{
		$$ = ast.TypeDef{
			Start:       $1.Start,
			FieldType:   $2,
			Identifier:  *$3,
			Annotations: $4,
			End:         $5,
		}
	}

Enum:
	ENUM Identifier '{' EnumDefList '}' TypeAnnotations
	{
		$$ = ast.Enum{
			Start:       $1.Start,
			Identifier:  *$2,
			LBrace:      $3.Start,
			List:        $4,
			RBrace:      $5.Start,
			Annotations: $6,
		}
	}

EnumDefList:
	EnumDefList EnumDef
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

EnumDef:
	EnumValue TypeAnnotations CommaOrSemicolonOptional
	{
		$$ = ast.EnumDef{
			Value:       $1,
			Annotations: $2,
			End:         $3,
		}
	}

EnumValue:
	Identifier '=' Number
	{
		$$ = ast.EnumValueWithNumber{
			Identifier: *$1,
			Eq:         $2.Start,
			Number:     $3,
		}
	}
|	Identifier
	{
		$$ = *$1
	}

Senum:
	SENUM Identifier '{' SenumDefList '}' TypeAnnotations
	{
		$$ = ast.Senum{
			Start:       $1.Start,
			Identifier:  *$2,
			LBrace:      $3.Start,
			List:        $4,
			RBrace:      $5.Start,
			Annotations: $6,
		}
	}

SenumDefList:
	SenumDefList SenumDef
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

SenumDef:
	Literal CommaOrSemicolonOptional
	{
		$$ = ast.SenumDef{
			Literal: $1,
			End:     $2,
		}
	}

Struct:
	STRUCTHEAD Identifier XsdAll '{' FieldList '}' TypeAnnotations
	{
		$$ = ast.Struct{
			Head:        *$1,
			Identifier:  *$2,
			XsdAll:      $3,
			LBrace:      $4.Start,
			List:        $5,
			RBrace:      $6.Start,
			Annotations: $7,
		}
	}

FieldList:
	FieldList Field
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

Field:
	FieldIdentifier FieldRequiredness FieldType FieldReference Identifier FieldValue XsdOptional XsdNillable XsdAttributes TypeAnnotations CommaOrSemicolonOptional
	{
		$$ = ast.Field{
			FieldIdentifier: $1,
			Requiredness:    $2,
			FieldType:       $3,
			Reference:       $4,
			Identifier:      *$5,
			FieldValue:      $6,
			XsdOptional:     $7,
			XsdNillable:     $8,
			XsdAttrs:        $9,
			Annotations:     $10,
			End:             $11,
		}
	}

FieldIdentifier:
	Number ':'
	{
		$$ = &ast.FieldIdentifier{
			Number: $1,
			Colon:  $2.Start,
		}
	}
|	// empty
	{
		$$ = nil
	}

FieldRequiredness:
	REQUIRED
	{
		$$ = $1
	}
|	OPTIONAL
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

FieldReference:
	'&'
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

FieldValue:
	'=' ConstValue
	{
		$$ = &ast.FieldValue{
			Start: $1.Start,
			Value: $2,
		}
	}
|	// empty
	{
		$$ = nil
	}

XsdAll:
	XSDALL
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

XsdOptional:
	XSDOPTIONAL
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

XsdNillable:
	XSDNILLABLE
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

XsdAttributes:
	XSDATTRS '{' FieldList '}'
	{
		$$ = &ast.XsdAttrs{
			Start:     $1.Start,
			LBrace:    $2.Start,
			FieldList: $3,
			RBrace:    $4.Start,
		}
	}
|	// empty
	{
		$$ = nil
	}

Exception:
	EXCEPTION Identifier '{' FieldList '}' TypeAnnotations
	{
		$$ = ast.Exception{
			Start:       $1.Start,
			Identifier:  *$2,
			LBrace:      $3.Start,
			FieldList:   $4,
			RBrace:      $5.Start,
			Annotations: $6,
		}
	}

Service:
	SERVICE Identifier Extends '{' FunctionList '}' TypeAnnotations
	{
		$$ = ast.Service{
			Start:        $1.Start,
			Identifier:   *$2,
			Extends:      $3,
			LBrace:       $4.Start,
			FunctionList: $5,
			RBrace:       $6.Start,
			Annotations:  $7,
		}
	}

Extends:
	EXTENDS Identifier
	{
		$$ = &ast.Extends{
			Start:      $1.Start,
			Identifier: *$2,
		}
	}
|	// empty
	{
		$$ = nil
	}

FunctionList:
	FunctionList Function
	{
		$$ = append($1, $2)
	}
|	// empty
	{
		$$ = nil
	}

Function:
	Oneway FunctionType Identifier '(' FieldList ')' Throws TypeAnnotations CommaOrSemicolonOptional
	{
		$$ = ast.Function{
			Oneway:       $1,
			FunctionType: $2,
			Identifier:   *$3,
			LChevron:     $4.Start,
			FieldList:    $5,
			RChevron:     $6.Start,
			Throws:       $7,
			Annotations:  $8,
			End:          $9,
		}
	}

Oneway:
	ONEWAY
	{
		$$ = $1
	}
|	// empty
	{
		$$ = nil
	}

FunctionType:
	FieldType
	{
		$$ = $1
	}
|	VOID
	{
		$$ = *$1
	}

Throws:
	THROWS '(' FieldList ')'
	{
		$$ = &ast.Throws{
			Start:     $1.Start,
			LChevron:  $2.Start,
			FieldList: $3,
			RChevron:  $4.Start,
		}
	}
|	// empty
	{
		$$ = nil
	}

%%
