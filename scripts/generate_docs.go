package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RuleDoc struct {
	Name        string
	Description string
	Fields      []string
}

// AI GENERATED - review required once I have more time for DX
func main() {
	rulesDir := "internal/rules"
	docsDir := "docs"

	files, err := os.ReadDir(rulesDir)
	if err != nil {
		log.Fatalf("failed to read rules directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), "rule_") || strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}

		path := filepath.Join(rulesDir, file.Name())
		doc, err := parseRuleFile(path)
		if err != nil {
			log.Printf("failed to parse %s: %v", path, err)
			continue
		}

		if doc.Name == "" {
			continue
		}

		err = writeMarkdown(docsDir, doc)
		if err != nil {
			log.Printf("failed to write markdown for %s: %v", doc.Name, err)
		}
	}
}

func parseRuleFile(path string) (RuleDoc, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return RuleDoc{}, err
	}

	var doc RuleDoc

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if x.Tok == token.CONST {
				for _, spec := range x.Specs {
					vspec := spec.(*ast.ValueSpec)
					for i, name := range vspec.Names {
						if strings.HasPrefix(name.Name, "ruleName") {
							if len(vspec.Values) > i {
								if lit, ok := vspec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
									doc.Name = strings.Trim(lit.Value, "\"")
								}
							}
						}
						// Search for fields rule<RuleName>FieldName<FieldName>
						// The issue says rule<name-of-rule>FieldName<field-name>
						// and ruleNoDeadEndsFieldNameAllowedNames = "allowed_names"
						// so it's rule + CamelCaseName + FieldName + CamelCaseField
						if strings.Contains(name.Name, "FieldName") {
							if len(vspec.Values) > i {
								if lit, ok := vspec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
									doc.Fields = append(doc.Fields, strings.Trim(lit.Value, "\""))
								}
							}
						}
					}
				}
			}
		case *ast.ValueSpec:
			// Handle descriptions in var rule... = Rule{...}
			for _, name := range x.Names {
				if strings.HasPrefix(name.Name, "rule") && !strings.Contains(name.Name, "Name") && !strings.Contains(name.Name, "FieldName") {
					for _, val := range x.Values {
						if comp, ok := val.(*ast.CompositeLit); ok {
							if typ, ok := comp.Type.(*ast.Ident); ok && typ.Name == "Rule" {
								for _, elt := range comp.Elts {
									if kv, ok := elt.(*ast.KeyValueExpr); ok {
										if kname, ok := kv.Key.(*ast.Ident); ok && kname.Name == "Description" {
											if vlit, ok := kv.Value.(*ast.BasicLit); ok && vlit.Kind == token.STRING {
												doc.Description = strings.Trim(vlit.Value, "\"")
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	return doc, nil
}

func writeMarkdown(docsDir string, doc RuleDoc) error {
	filename := strings.ToLower(doc.Name) + ".md"
	path := filepath.Join(docsDir, filename)

	content := fmt.Sprintf("# %s\n\n", doc.Name)
	content += fmt.Sprintf("## Description\n\n%s\n\n", doc.Description)

	if len(doc.Fields) > 0 {
		content += "## Configuration Fields\n\n"
		for _, field := range doc.Fields {
			content += fmt.Sprintf("- `%s`\n", field)
		}
		content += "\n"
	}

	return os.WriteFile(path, []byte(content), 0644)
}
