package scanner

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	packageName  = "babel"
	singularFunc = "Tr"
	pluralFunc   = "Trn"
)

type references []string

type translation interface {
	Serialize() string
}

type singularized struct {
	text string
}

func (s singularized) Serialize() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("msgid %s\n", s.text))
	buffer.WriteString(fmt.Sprintf("msgstr %s\n\n", s.text))
	return buffer.String()
}

type pluralized struct {
	singular string
	plural   string
}

func (s pluralized) Serialize() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("msgid %s\n", s.singular))
	buffer.WriteString(fmt.Sprintf("msgid_plural %s\n", s.plural))
	buffer.WriteString(fmt.Sprintf("msgstr[0] %s\n", s.singular))
	buffer.WriteString(fmt.Sprintf("msgstr[1] %s\n\n", s.plural))
	return buffer.String()
}

type scanner struct {
	translations map[translation]references
	fileset      *token.FileSet
}

func NewFileScanner() *scanner {
	return &scanner{translations: make(map[translation]references), fileset: token.NewFileSet()}
}

func (s *scanner) Scan(filename string) error {
	file, err := parser.ParseFile(s.fileset, filename, nil, 0)
	if err != nil {
		return err
	}
	ast.Walk(s, file)
	return nil
}

func (s *scanner) Save(filename string) error {
	var buffer bytes.Buffer

	for translation, references := range s.translations {
		for _, reference := range references {
			buffer.WriteString(fmt.Sprintf("#: %s\n", reference))
		}
		buffer.WriteString(translation.Serialize())
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return errors.Wrap(err, "Error trying to create the output dir")
	}
	if output, err := os.Create(filename); err != nil {
		return errors.Wrap(err, "Error trying to create the output file")
	} else {
		defer output.Close()
		io.WriteString(output, buffer.String())
	}
	return nil
}

func (s *scanner) Visit(node ast.Node) ast.Visitor {
	ok := false

	// search for function calls
	var call *ast.CallExpr
	if call, ok = node.(*ast.CallExpr); !ok {
		return s
	}

	// only functions with a selector (foo.bar(), nothing.Do(), fmt.Println("hello world"))
	var fun *ast.SelectorExpr
	if fun, ok = call.Fun.(*ast.SelectorExpr); !ok {
		return s
	}

	// check the selector identifier
	var pkg *ast.Ident
	if pkg, ok = fun.X.(*ast.Ident); !ok {
		return s
	}

	// check if the call is babel.T(locale, "some key", ...)
	if pkg.Name == packageName {
		if fun.Sel.Name == singularFunc {
			if literal, ok := call.Args[1].(*ast.BasicLit); ok && literal.Kind == token.STRING {
				s.addTranslation(singularized{literal.Value}, node)
			}
		} else if fun.Sel.Name == pluralFunc {
			if singular, ok := call.Args[2].(*ast.BasicLit); ok && singular.Kind == token.STRING {
				if plural, ok := call.Args[3].(*ast.BasicLit); ok && plural.Kind == token.STRING {
					s.addTranslation(pluralized{singular.Value, plural.Value}, node)
				}
			}
		}
	}

	return s
}

func (s *scanner) addTranslation(t translation, node ast.Node) {
	s.translations[t] = append(s.translations[t], s.fileset.Position(node.Pos()).String())
}
