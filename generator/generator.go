package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/utils"
)

type DataFunc func(api *design.APIDefinition) (interface{}, error)

// Generator defines the behavior of a resulting generated file
type Generator struct {
	target   string
	outDir   string
	outFile  string
	tmplName string
	imports  []string
	funcs    template.FuncMap
	fileTmpl string
	genFiles []string
	dataFunc DataFunc
}

// New creates a new Generator with the given arguments
func New(target, outDir, outFile, fileTmpl, tmplName string, dataFunc DataFunc) *Generator {
	return &Generator{
		target:   target,
		outDir:   outDir,
		outFile:  outFile,
		tmplName: tmplName,
		fileTmpl: fileTmpl,
		dataFunc: dataFunc,
	}
}

func (g *Generator) AddImports(importStr ...string) {
	g.imports = append(g.imports, importStr...)
}

func (g *Generator) SetFuncs(funcs template.FuncMap) {
	g.funcs = funcs
}

// Generate generates the file defined by the Generator
func (g *Generator) Generate(api *design.APIDefinition) (_ []string, err error) {
	if api == nil {
		return nil, fmt.Errorf("missing API definition, make sure design.Design is properly initialized")
	}
	go utils.Catch(nil, func() {
		g.Cleanup()
	})
	defer func() {
		if err != nil {
			g.Cleanup()
		}
	}()
	if err := os.MkdirAll(g.outDir, 0755); err != nil {
		return nil, err
	}
	if err := g.generate(api); err != nil {
		return g.genFiles, err
	}
	return g.genFiles, nil
}

// Cleanup cleans stuff up
func (g *Generator) Cleanup() {
	if len(g.genFiles) == 0 {
		return
	}
	g.genFiles = nil
}

func (g *Generator) generate(api *design.APIDefinition) error {
	// Make the dir if it doesn't exist
	if err := os.MkdirAll(g.outDir, 0755); err != nil {
		return err
	}

	// Remove the old file, if it exists
	outFile := filepath.Join(g.outDir, g.outFile)
	os.Remove(outFile)

	// Create the data for the template now, so we avoid writing the file to
	// disk if there's nothing in it
	var data interface{}
	if g.dataFunc != nil {
		data, err := g.dataFunc(api)
		if err != nil {
			return err
		}
		if data == nil {
			return nil
		}
	} else {
		data = map[string]interface{}{
			"API": api,
		}
	}

	// Create the new source file
	g.genFiles = append(g.genFiles, outFile)
	file, err := codegen.SourceFileFor(outFile)
	if err != nil {
		return err
	}

	// Add imports
	imports := make([]*codegen.ImportSpec, 0)
	for _, importStr := range g.imports {
		imports = append(imports, codegen.SimpleImport(importStr))
	}
	file.WriteHeader("", g.target, imports)

	if err = file.ExecuteTemplate(g.tmplName, g.fileTmpl, g.funcs, data); err != nil {
		return err
	}
	return file.FormatCode()
}
