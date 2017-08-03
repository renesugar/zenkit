package wstesthelper

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	genapp "github.com/goadesign/goa/goagen/gen_app"
	"github.com/zenoss/zenkit/generator"
)

// WSTest represents an action for which a WebSocket test helper should be
// generated
type WSTest struct {
	Name           string
	ResourceName   string
	ActionName     string
	ControllerName string
	ContextVarName string
	ContextType    string
	Params         []*genapp.ObjectType
	QueryParams    []*genapp.ObjectType
	Headers        []*genapp.ObjectType
	reservedNames  map[string]bool
}

// Escape escapes given string.
func (t *WSTest) Escape(s string) string {
	if ok := t.reservedNames[s]; ok {
		s = t.Escape("_" + s)
	}
	t.reservedNames[s] = true
	return s
}

func Generate() ([]string, error) {
	var outDir, target, appPkg, ver string

	outFiles := make([]string, 0)

	set := flag.NewFlagSet("resources", flag.PanicOnError)
	set.String("design", "", "")
	set.StringVar(&outDir, "out", "", "")
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&target, "pkg", "test", "")
	set.StringVar(&appPkg, "app", "app", "")
	set.Parse(os.Args[2:])

	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	appPkgPath, err := codegen.PackagePath(filepath.Join(outDir, appPkg))
	if err != nil {
		return nil, fmt.Errorf("invalid app package: %s", err)
	}

	outDir = filepath.Join(outDir, appPkg, "test")

	dataFunc := func(api *design.APIDefinition) (interface{}, error) {
		tests := []*WSTest{}
		if err := api.IterateResources(func(res *design.ResourceDefinition) error {
			if err := res.IterateActions(func(action *design.ActionDefinition) error {
				if err := action.IterateResponses(func(response *design.ResponseDefinition) error {
					if response.Status != 101 { // Only deal with websocket endpoints
						return nil
					}
					for _, route := range action.Routes {
						actionName := codegen.Goify(action.Name, true)
						ctrlName := codegen.Goify(res.Name, true)
						varName := codegen.Goify(action.Name, false)
						path := pathParams(action, route)
						query := queryParams(action)
						header := headers(action, res.Headers)
						tests = append(tests, &WSTest{
							Name:           fmt.Sprintf("%s%sWSTestHelper", actionName, ctrlName),
							ActionName:     actionName,
							ControllerName: fmt.Sprintf("%s.%sController", appPkg, ctrlName),
							ContextVarName: fmt.Sprintf("%sCtx", varName),
							ContextType:    fmt.Sprintf("%s.New%s%sContext", appPkg, actionName, ctrlName),
							Params:         path,
							QueryParams:    query,
							Headers:        header,
							reservedNames:  reservedNames(path, query, header),
						})
					}
					return nil
				}); err != nil {
					return nil
				}
				return nil
			}); err != nil {
				return nil
			}
			return nil
		}); err != nil {
			return nil, err
		}
		if len(tests) == 0 {
			return nil, nil
		}
		return tests, nil
	}
	wsHelperGen := generator.New(target, outDir, "websocket_helpers.go", wsTmpl, "websocket_helpers", dataFunc)
	wsHelperGen.AddImports("context", "net/http", "net/url", "strconv",
		"github.com/goadesign/goa", "github.com/goadesign/goa/goatest",
		"github.com/gorilla/websocket", "github.com/posener/wstest", appPkgPath)
	wsHelperGen.SetFuncs(template.FuncMap{
		"isSlice": isSlice,
	})
	wsHelperFiles, err := wsHelperGen.Generate(design.Design)
	if err != nil {
		return []string{}, err
	}
	outFiles = append(outFiles, wsHelperFiles...)
	return outFiles, nil
}

// pathParams returns the path params for the given action and route.
func pathParams(action *design.ActionDefinition, route *design.RouteDefinition) []*genapp.ObjectType {
	return paramFromNames(action, route.Params())
}

func attToObject(name string, parent, att *design.AttributeDefinition) *genapp.ObjectType {
	obj := &genapp.ObjectType{}
	obj.Label = name
	obj.Name = codegen.Goify(name, false)
	obj.Type = codegen.GoTypeRef(att.Type, nil, 0, false)
	if att.Type.IsPrimitive() && parent.IsPrimitivePointer(name) {
		obj.Pointer = "*"
	}
	return obj
}

// queryParams returns the query string params for the given action.
func queryParams(action *design.ActionDefinition) []*genapp.ObjectType {
	var qparams []string
	if qps := action.QueryParams; qps != nil {
		for pname := range qps.Type.ToObject() {
			qparams = append(qparams, pname)
		}
	}
	sort.Strings(qparams)
	return paramFromNames(action, qparams)
}

func paramFromNames(action *design.ActionDefinition, names []string) (params []*genapp.ObjectType) {
	obj := action.Params.Type.ToObject()
	for _, name := range names {
		params = append(params, attToObject(name, action.Params, obj[name]))
	}
	return
}

// headers builds the template data structure needed to proprely render the code
// for setting the headers for the given action.
func headers(action *design.ActionDefinition, headers *design.AttributeDefinition) []*genapp.ObjectType {
	hds := &design.AttributeDefinition{
		Type: design.Object{},
	}
	if headers != nil {
		hds.Merge(headers)
		hds.Validation = headers.Validation
	}
	if action.Headers != nil {
		hds.Merge(action.Headers)
		hds.Validation = action.Headers.Validation
	}

	if hds == nil {
		return nil
	}
	var headrs []string
	for header := range hds.Type.ToObject() {
		headrs = append(headrs, header)
	}
	sort.Strings(headrs)
	objs := make([]*genapp.ObjectType, len(headrs))
	for i, name := range headrs {
		objs[i] = attToObject(name, hds, hds.Type.ToObject()[name])
		objs[i].Label = http.CanonicalHeaderKey(objs[i].Label)
	}
	return objs
}

func reservedNames(params, queryParams, headers []*genapp.ObjectType) map[string]bool {
	var names = make(map[string]bool)
	for _, param := range params {
		names[param.Name] = true
	}
	for _, param := range queryParams {
		names[param.Name] = true
	}
	for _, header := range headers {
		names[header.Name] = true
	}
	return names
}

func isSlice(typeName string) bool {
	return strings.HasPrefix(typeName, "[]")
}

var convertParamTmpl = `{{ if eq .Type "string" }}		sliceVal := []string{ {{ if .Pointer }}*{{ end }}{{ .Name }}}{{/*
*/}}{{ else if eq .Type "int" }}		sliceVal := []string{strconv.Itoa({{ if .Pointer }}*{{ end }}{{ .Name }})}{{/*
*/}}{{ else if eq .Type "[]string" }}		sliceVal := {{ .Name }}{{/*
*/}}{{ else if (isSlice .Type) }}		sliceVal := make([]string, len({{ .Name }}))
		for i, v := range {{ .Name }} {
			sliceVal[i] = fmt.Sprintf("%v", v)
		}{{/*
*/}}{{ else if eq .Type "time.Time" }}		sliceVal := []string{ {{ if .Pointer }}(*{{ end }}{{ .Name }}{{ if .Pointer }}){{ end }}.Format(time.RFC3339)}{{/*
*/}}{{ else }}		sliceVal := []string{fmt.Sprintf("%v", {{ if .Pointer }}*{{ end }}{{ .Name }})}{{ end }}`

var wsTmpl = `{{ define "convertParam" }}` + convertParamTmpl + `{{ end }}` + `
{{ range $test := . }}
func {{ $test.Name }}(t goatest.TInterface, ctx context.Context, service *goa.Service, ctrl {{ $test.ControllerName}}{{/*
*/}}{{ range $param := $test.Params }}, {{ $param.Name }} {{ $param.Pointer }}{{ $param.Type }}{{ end }}{{/*
*/}}{{ range $param := $test.QueryParams }}, {{ $param.Name }} {{ $param.Pointer }}{{ $param.Type }}{{ end }}{{/*
*/}}{{ range $header := $test.Headers }}, {{ $header.Name }} {{ $header.Pointer }}{{ $header.Type }}{{ end }}{{/*
*/}}) (*websocket.Conn, error) {
	{{ $prms := $test.Escape "prms" }}{{ $prms }} := url.Values{}
{{ range $param := $test.Params }}	{{ $prms }}["{{ $param.Label }}"] = []string{fmt.Sprintf("%v",{{ $param.Name}})}
{{ end }}{{ range $param := $test.QueryParams }}{{ if $param.Pointer }} if {{ $param.Name }} != nil {{ end }} {
{{ template "convertParam" $param }}
		{{ $prms }}[{{ printf "%q" $param.Label }}] = sliceVal
	}
{{ end }}	h := func(rw http.ResponseWriter, req *http.Request) {
		{{ $goaCtx := $test.Escape "goaCtx" }}{{ $goaCtx }} := goa.NewContext(goa.WithAction(ctx, "{{ $test.ResourceName }}Test"), rw, req, {{ $prms }})
		{{ $test.ContextVarName }}, {{ $err := $test.Escape "err" }}{{ $err }} := {{ $test.ContextType }}({{ $goaCtx }}, req, service)
		if {{ $err }} != nil {
			panic("invalid test data " + {{ $err }}.Error()) // bug
		}
		{{ $err }} = ctrl.{{ $test.ActionName}}({{ $test.ContextVarName }})
		if {{ $err }} != nil {
			t.Fatalf("controller returned %%+v", {{ $err }})
		}
	}
	dialer := wstest.NewDialer(http.HandlerFunc(h), nil)
	head := http.Header{}
	head.Set("Origin", "http://localhost/")
	conn, _, err := dialer.Dial("ws://example.com/", head)
	if err != nil {
		t.Fatalf("unable to set up websocket test harness: %%+v", err)
	}
	return conn, nil
}
{{end}}
`
