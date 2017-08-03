package controllerreg

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/zenoss/zenkit/generator"
)

const ctrlRegTmpl = `
func MountAllControllers(service *goa.Service) {
{{ $api := .API }}
{{ range $name, $res := $api.Resources }}{{ $name := goify $res.Name true }} // Mount {{$res.Name}} controller
	{{ $tmp := tempvar }}{{ $tmp }} := New{{ $name }}Controller(service)
	{{ targetPkg }}.Mount{{ $name }}Controller(service, {{ $tmp }})
{{ end }}
}
`

const ctrlRegTestTmpl = `
var _ = Describe("ControllerReg", func() {
	var (
		svc = goa.New("controller-test")
	)
	Context("when mounting all controllers", func() {
		MountAllControllers(svc)
		{{ range $name, $res := $api.Resources }}{{ $name := goify $res.Name true }}
			It("should mount the {{$res.Name}} controller", func() {
				// Put your logic here to test that the controller is mounted
			})
		{{ end }}
	})
})
`

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

func Generate() ([]string, error) {
	var outDir, target, appPkg, ver string

	outFiles := make([]string, 0)

	set := flag.NewFlagSet("resources", flag.PanicOnError)
	set.String("design", "", "")
	set.StringVar(&outDir, "out", "", "")
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&target, "pkg", "resources", "")
	set.StringVar(&appPkg, "app", "app", "")
	set.Parse(os.Args[2:])

	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	appPkgPath, err := codegen.PackagePath(filepath.Join(outDir, appPkg))
	if err != nil {
		return nil, fmt.Errorf("invalid app package: %s", err)
	}

	ctrlRegGen := generator.New(target, outDir, "controller_reg.go", ctrlRegTmpl, "controller_reg")
	ctrlRegGen.AddImports("github.com/goadesign/goa", appPkgPath)
	ctrlRegGen.SetFunc(template.FuncMap{
		"tempvar":   tempvar,
		"targetPkg": func() string { return g.appPkg },
	})
	ctrlRegGen.Generate(design.Design)

	ctrlRegTestGen := generator.New(target+"_test", outDir, "controller_reg_test.go", ctrlRegTmpl, "controller_reg_test")
	ctrlRegGen.AddImports("github.com/goadesign/goa",
		". github.com/onsi/ginkgo",
		". github.com/onsi/gomega",
		appPkgPath,
	)
	ctrlRegGen.Generate(design.Design)

	gTargetTest := &generate.Generator{
		outDir:       outDir,
		outFile:      "controller_reg_test.go",
		target:       target,
		templateName: "controller_reg_test",
		imports: {
			"github.com/goadesign/goa",
			". github.com/onsi/ginkgo",
			". " + gTarget.outDir + "/" + gTarget.target,
		},
		funcs: template.FuncMap{
			"tempvar":   tempvar,
			"targetPkg": func() string { return g.appPkg },
		},
		fileTemplate: ctrlRegTempl,
	}

	return outFiles
}
