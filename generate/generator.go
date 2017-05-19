package generator

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/utils"
)

type Generator struct {
	genfiles   []string
	outDir     string
	target     string
	appPkg     string
	appPkgPath string
}

func Generate() ([]string, error) {
	var outDir, target, appPkg, ver string

	set := flag.NewFlagSet("resources", flag.PanicOnError)
	set.String("design", "", "")
	set.StringVar(&outDir, "out", "", "")
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&target, "pkg", "resources", "")
	set.StringVar(&appPkg, "app", "app", "")
	set.Parse(os.Args[2:])

	fmt.Println(outDir, ver)

	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	appPkgPath, err := codegen.PackagePath(filepath.Join(outDir, appPkg))
	if err != nil {
		return nil, fmt.Errorf("invalid app package: %s", err)
	}

	g := &Generator{
		outDir:     outDir,
		target:     target,
		appPkg:     appPkg,
		appPkgPath: appPkgPath,
	}

	return g.Generate(design.Design)
}

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
	if err := g.generateControllerRegistration(api); err != nil {
		return g.genfiles, err
	}
	return g.genfiles, nil
}

func (g *Generator) Cleanup() {
	if len(g.genfiles) == 0 {
		return
	}
	g.genfiles = nil
}

// tempCount is the counter used to create unique temporary variable names.
var tempCount int

// tempvar generates a unique temp var name.
func tempvar() string {
	tempCount++
	if tempCount == 1 {
		return "c"
	}
	return fmt.Sprintf("c%d", tempCount)
}

func (g *Generator) generateControllerRegistration(api *design.APIDefinition) error {
	if err := os.MkdirAll(g.outDir, 0755); err != nil {
		return err
	}
	outFile := filepath.Join(g.outDir, "controller_reg.go")
	os.Remove(outFile)
	g.genfiles = append(g.genfiles, outFile)
	file, err := codegen.SourceFileFor(outFile)
	if err != nil {
		return err
	}
	imports := []*codegen.ImportSpec{
		codegen.SimpleImport("github.com/goadesign/goa"),
		codegen.SimpleImport(g.appPkgPath),
	}
	file.WriteHeader("", g.target, imports)
	data := map[string]interface{}{
		"API": api,
	}
	funcs := template.FuncMap{
		"tempvar":   tempvar,
		"targetPkg": func() string { return g.appPkg },
	}
	if err = file.ExecuteTemplate("controller_reg", fileT, funcs, data); err != nil {
		return err
	}
	return file.FormatCode()
}

const fileT = `
func MountAllControllers(service *goa.Service) {
{{ $api := .API }}
{{ range $name, $res := $api.Resources }}{{ $name := goify $res.Name true }} // Mount "{{$res.Name}}" controller
	{{ $tmp := tempvar }}{{ $tmp }} := New{{ $name }}Controller(service)
	{{ targetPkg }}.Mount{{ $name }}Controller(service, {{ $tmp }})
{{ end }}
}
`
