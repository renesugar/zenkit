package generator

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
)

func Generate() ([]string, error) {
	var (
		ver    string
		outDir string
	)
	set := flag.NewFlagSet("resources", flag.PanicOnError)
	set.String("design", "", "")
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&outDir, "out", "", "")
	set.Parse(os.Args[2:])

	fmt.Println(outDir, ver)

	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}
	fmt.Println(design.Design)

	return WriteNames(design.Design, outDir)
}

func WriteNames(api *design.APIDefinition, outDir string) ([]string, error) {
	names := make([]string, len(api.Resources))
	i := 0
	api.IterateResources(func(res *design.ResourceDefinition) error {
		names[i] = res.Name
		i++
		return nil
	})

	content := strings.Join(names, "\n")
	outputFile := filepath.Join(outDir, "names.txt")
	if err := ioutil.WriteFile(outputFile, []byte(content), 0644); err != nil {
		return nil, err
	}
	return []string{outputFile}, nil
}
