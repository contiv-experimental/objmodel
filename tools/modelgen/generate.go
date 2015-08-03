/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
)

// GenerateGoStructs generates go code from a schema
func (s *Schema) GenerateGoStructs() (string, error) {
	var goStr string

	//  Generate all object definitions
	for _, obj := range s.Objects {
		objStr, err := obj.GenerateGoStructs()
		if err == nil {
			goStr = goStr + objStr
		}
	}

	// Generate a collection definitions to store the objects
	goStr = goStr + fmt.Sprintf("\n\ntype Collections struct {\n")
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("	%ss    map[string]*%s\n", obj.Name, initialCap(obj.Name))
	}
	goStr = goStr + fmt.Sprintf("}\n\n")

	goStr = goStr + fmt.Sprintf("var collections Collections\n\n")

	// Generate callback interface
	goStr = goStr + fmt.Sprintf("type Callbacks interface {\n")
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("	%sCreate(%s *%s) error\n", initialCap(obj.Name), obj.Name, initialCap(obj.Name))
		goStr = goStr + fmt.Sprintf("	%sDelete(%s *%s) error\n", initialCap(obj.Name), obj.Name, initialCap(obj.Name))
	}
	goStr = goStr + fmt.Sprintf("}\n\n")

	goStr = goStr + fmt.Sprintf("var objCallbackHandler Callbacks\n\n")

	// Generate an Init function
	goStr = goStr + fmt.Sprintf("\nfunc Init(handler Callbacks) {\n")
	goStr = goStr + fmt.Sprintf("objCallbackHandler = handler\n\n")
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("	collections.%ss = make(map[string]*%s)\n", obj.Name, initialCap(obj.Name))
	}
	goStr = goStr + fmt.Sprintf("\n")
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("	restore%s()\n", initialCap(obj.Name))
	}

	goStr = goStr + fmt.Sprintf("}\n\n")

	return goStr, nil
}

// GenerateGoHdrs generates go file headers
func (s *Schema) GenerateGoHdrs() string {
	var buf bytes.Buffer

	const hdr = `// {{.Name}}.go
// This file is auto generated by modelgen tool
// Do not edit this file manually

package {{.Name}}

import (
	"errors"
	"net/http"
	"encoding/json"
	"github.com/contiv/objmodel/objdb/modeldb"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
)

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error)

`

	tmpl := template.Must(template.New("hdr").Parse(hdr))
	err := tmpl.Execute(&buf, s)
	if err != nil {
		log.Errorf("Error executing template. Err: %v", err)
		return ""
	}

	return buf.String()
}

func (s *Schema) GenerateGoFuncs() string {
	var buf bytes.Buffer
	var goStr string

	routeFunc := `
// Simple Wrapper for http handlers
func makeHttpHandler(handlerFunc HttpApiFunc) http.HandlerFunc {
	// Create a closure and return an anonymous function
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the handler
		resp, err := handlerFunc(w, r, mux.Vars(r))
		if err != nil {
			// Log error
			log.Errorf("Handler for %s %s returned error: %s", r.Method, r.URL, err)

			// Send HTTP response
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// Send HTTP response as Json
			err = writeJSON(w, http.StatusOK, resp)
			if err != nil {
				log.Errorf("Error generating json. Err: %v", err)
			}
		}
	}
}

// writeJSON: writes the value v to the http response stream as json with standard
// json encoding.
func writeJSON(w http.ResponseWriter, code int, v interface{}) error {
	// Set content type as json
	w.Header().Set("Content-Type", "application/json")

	// write the HTTP status code
	w.WriteHeader(code)

	// Write the Json output
	return json.NewEncoder(w).Encode(v)
}

// Add all routes for REST handlers
func AddRoutes(router *mux.Router) {
	var route, listRoute string
`
	funcMap := template.FuncMap{
		"initialCap": initialCap,
	}
	routeTmpl := `
	// Register {{.}}
	route = "/api/{{.}}s/{key}/"
	listRoute = "/api/{{.}}s/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpList{{initialCap .}}s))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGet{{initialCap .}}))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreate{{initialCap .}}))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreate{{initialCap .}}))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDelete{{initialCap .}}))
`
	// Output the functions and routes
	rfTmpl, _ := template.New("routeTmpl").Parse(routeFunc)
	rfTmpl.Execute(&buf, "")
	goStr = goStr + buf.String()

	// add a path for each object
	for _, obj := range s.Objects {
		var buf bytes.Buffer

		// Create a template, add the function map, and parse the text.
		tmpl, err := template.New("routeTmpl").Funcs(funcMap).Parse(routeTmpl)
		if err != nil {
			log.Fatalf("parsing: %s", err)
		}

		// Run the template.
		err = tmpl.Execute(&buf, obj.Name)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		goStr = goStr + buf.String()
	}

	goStr = goStr + fmt.Sprintf("\n}\n")

	// template for handler functions
	handlerFuncs := `
// LIST REST call
func httpList{{initialCap .}}s(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpList{{initialCap .}}s: %+v", vars)

	list := make([]*{{initialCap .}}, 0)
	for _, obj := range collections.{{.}}s {
		list = append(list, obj)
	}

	// Return the list
	return list, nil
}

// GET REST call
func httpGet{{initialCap .}}(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGet{{initialCap .}}: %+v", vars)

	key := vars["key"]

	obj := collections.{{.}}s[key]
	if obj == nil {
		log.Errorf("{{.}} %s not found", key)
		return nil, errors.New("{{.}} not found")
	}

	// Return the obj
	return obj, nil
}

// CREATE REST call
func httpCreate{{initialCap .}}(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpGet{{initialCap .}}: %+v", vars)

	var obj {{initialCap .}}
	key := vars["key"]

	// Get object from the request
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		log.Errorf("Error decoding {{.}} create request. Err %v", err)
		return nil, err
	}

	// set the key
	obj.Key = key

	// Create the object
	err = Create{{initialCap .}}(&obj)
	if err != nil {
		log.Errorf("Create{{initialCap .}} error for: %+v. Err: %v", obj, err)
		return nil, err
	}

	// Return the obj
	return obj, nil
}

// DELETE rest call
func httpDelete{{initialCap .}}(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error) {
	log.Debugf("Received httpDelete{{initialCap .}}: %+v", vars)

	key := vars["key"]

	// Delete the object
	err := Delete{{initialCap .}}(key)
	if err != nil {
		log.Errorf("Delete{{initialCap .}} error for: %s. Err: %v", key, err)
		return nil, err
	}

	// Return the obj
	return key, nil
}

// Create a {{.}} object
func Create{{initialCap .}}(obj *{{initialCap .}}) error {
	// save it in cache
	collections.{{.}}s[obj.Key] = obj

	// Perform callback
	err := objCallbackHandler.{{initialCap .}}Create(obj)
	if err != nil {
		log.Errorf("{{initialCap .}}Create retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// Write it to modeldb
	err = obj.Write()
	if err != nil {
		log.Errorf("Error saving {{.}} %s to db. Err: %v", obj.Key, err)
		return err
	}

	return nil
}

// Return a pointer to {{.}} from collection
func Find{{initialCap .}}(key string) *{{initialCap .}} {
	obj := collections.{{.}}s[key]
	if obj == nil {
		log.Errorf("{{.}} %s not found", key)
		return nil
	}

	return obj
}

// Delete a {{.}} object
func Delete{{initialCap .}}(key string) error {
	obj := collections.{{.}}s[key]
	if obj == nil {
		log.Errorf("{{.}} %s not found", key)
		return errors.New("{{.}} not found")
	}

	// set the key
	obj.Key = key

	// Perform callback
	err := objCallbackHandler.{{initialCap .}}Delete(obj)
	if err != nil {
		log.Errorf("{{initialCap .}}Delete retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// delete it from modeldb
	err = obj.Delete()
	if err != nil {
		log.Errorf("Error deleting {{.}} %s. Err: %v", obj.Key, err)
	}

	// delete it from cache
	delete(collections.{{.}}s, key)

	return nil
}

func (self *{{initialCap .}}) GetType() string {
	return "{{.}}"
}

func (self *{{initialCap .}}) GetKey() string {
	return self.Key
}

func (self *{{initialCap .}}) Read() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to read {{.}} object")
		return errors.New("Empty key")
	}

	return modeldb.ReadObj("{{.}}", self.Key, self)
}

func (self *{{initialCap .}}) Write() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Write {{.}} object")
		return errors.New("Empty key")
	}

	return modeldb.WriteObj("{{.}}", self.Key, self)
}

func (self *{{initialCap .}}) Delete() error {
	if self.Key == "" {
		log.Errorf("Empty key while trying to Delete {{.}} object")
		return errors.New("Empty key")
	}

	return modeldb.DeleteObj("{{.}}", self.Key)
}

func restore{{initialCap .}}() error {
	strList, err := modeldb.ReadAllObj("{{.}}")
	if err != nil {
		log.Errorf("Error reading {{.}} list. Err: %v", err)
	}

	for _, objStr := range strList {
		// Parse the json model
		var {{.}} {{initialCap .}}
		err = json.Unmarshal([]byte(objStr), &{{.}})
		if err != nil {
			log.Errorf("Error parsing object %s, Err %v", objStr, err)
			return err
		}

		// add it to the collection
		collections.{{.}}s[{{.}}.Key] = &{{.}}
	}

	return nil
}
`
	// Generate REST handlers for each object
	for _, obj := range s.Objects {
		var buf bytes.Buffer
		// Create a template, add the function map, and parse the text.
		tmpl, err := template.New("routeTmpl").Funcs(funcMap).Parse(handlerFuncs)
		if err != nil {
			log.Fatalf("parsing: %s", err)
		}

		// Run the template.
		err = tmpl.Execute(&buf, obj.Name)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		goStr = goStr + buf.String()
	}

	return goStr
}

func (obj *Object) GenerateGoStructs() (string, error) {
	var goStr string

	objName := initialCap(obj.Name)
	goStr = goStr + fmt.Sprintf("type %s struct {\n", objName)

	// every object has a key
	goStr = goStr + fmt.Sprintf("	Key		string		`json:\"key,omitempty\"`\n")

	// Walk each property and generate code for it
	for _, prop := range obj.Properties {
		propStr, err := prop.GenerateGoStructs()
		if err == nil {
			goStr = goStr + propStr
		}
	}

	// add link-sets
	if len(obj.LinkSets) > 0 {
		goStr = goStr + fmt.Sprintf("	LinkSets	%sLinkSets		`json:\"link-sets,omitempty\"`\n", objName)
	}

	// add links
	if len(obj.Links) > 0 {
		goStr = goStr + fmt.Sprintf("	Links	%sLinks		`json:\"links,omitempty\"`\n", objName)
	}

	goStr = goStr + fmt.Sprintf("}\n\n")

	// define object's linkset
	if len(obj.LinkSets) > 0 {
		goStr = goStr + fmt.Sprintf("type %sLinkSets struct {\n", objName)
		for lsName := range obj.LinkSets {
			goStr = goStr + fmt.Sprintf("	%s	map[string]modeldb.Link		`json:\"%s,omitempty\"`\n", initialCap(lsName), lsName)
		}
		goStr = goStr + fmt.Sprintf("}\n\n")
	}
	/*
		// Define each link-sets
		for _, linkSet := range obj.LinkSets {
			subStr, err := linkSet.GenerateGoStructs()
			if err == nil {
				goStr = goStr + subStr
			}
		}
	*/
	// Define object's links
	if len(obj.Links) > 0 {
		goStr = goStr + fmt.Sprintf("type %sLinks struct {\n", objName)
		for lName := range obj.Links {
			goStr = goStr + fmt.Sprintf("	%s	modeldb.Link		`json:\"%s,omitempty\"`\n", initialCap(lName), lName)
		}
		goStr = goStr + fmt.Sprintf("}\n\n")
	}
	/*
		// define each link
		for _, link := range obj.Links {
			subStr, err := link.GenerateGoStructs()
			if err == nil {
				goStr = goStr + subStr
			}
		}
	*/

	return goStr, nil
}

func (ls *LinkSet) GenerateGoStructs() (string, error) {
	var goStr string

	goStr = goStr + fmt.Sprintf("type %sLinkSet struct {\n", ls.Name)
	goStr = goStr + fmt.Sprintf("	Type	string		`json:\"type,omitempty\"`\n")
	goStr = goStr + fmt.Sprintf("	Key		string		`json:\"key,omitempty\"`\n")
	goStr = goStr + fmt.Sprintf("	%s		*%s			`json:\"-\"`\n", ls.Ref, initialCap(ls.Ref))
	goStr = goStr + fmt.Sprintf("}\n\n")

	return goStr, nil
}

func (link *Link) GenerateGoStructs() (string, error) {
	var goStr string

	goStr = goStr + fmt.Sprintf("type %sLink struct {\n", link.Name)
	goStr = goStr + fmt.Sprintf("	Type	string		`json:\"type,omitempty\"`\n")
	goStr = goStr + fmt.Sprintf("	Key		string		`json:\"key,omitempty\"`\n")
	goStr = goStr + fmt.Sprintf("	%s		*%s		`json:\"-\"`\n", link.Ref, initialCap(link.Ref))
	goStr = goStr + fmt.Sprintf("}\n\n")

	return goStr, nil
}

func xlatePropType(propType string) string {
	var goStr string
	switch propType {
	case "string":
		goStr = goStr + fmt.Sprintf("string")
	case "number":
		goStr = goStr + fmt.Sprintf("float64")
	case "int":
		goStr = goStr + fmt.Sprintf("int64")
	case "bool":
		goStr = goStr + fmt.Sprintf("bool")
	default:
		return ""
	}

	return goStr
}

func (prop *Property) GenerateGoStructs() (string, error) {
	var goStr string

	goStr = fmt.Sprintf("	%s	", initialCap(prop.Name))
	switch prop.Type {
	case "string":
		fallthrough
	case "number":
		fallthrough
	case "int":
		fallthrough
	case "bool":
		subStr := xlatePropType(prop.Type)
		goStr = goStr + fmt.Sprintf("%s		`json:\"%s,omitempty\"`\n", subStr, prop.Name)
	case "array":
		subStr := xlatePropType(prop.Items)
		if subStr == "" {
			return "", errors.New("Unknown array items")
		}

		goStr = goStr + fmt.Sprintf("[]%s		`json:\"%s,omitempty\"`\n", subStr, prop.Name)
	default:
		return "", errors.New("Unknown Property")
	}

	return goStr, nil
}

/*********************** Helper funcs *************************/
var (
	newlines  = regexp.MustCompile(`(?m:\s*$)`)
	acronyms  = regexp.MustCompile(`(Url|Http|Id|Io|Uuid|Api|Uri|Ssl|Cname|Oauth|Otp)$`)
	camelcase = regexp.MustCompile(`(?m)[-.$/:_{}\s]`)
)

func initialCap(ident string) string {
	if ident == "" {
		panic("blank identifier")
	}
	return depunct(ident, true)
}

func initialLow(ident string) string {
	if ident == "" {
		panic("blank identifier")
	}
	return depunct(ident, false)
}

func depunct(ident string, initialCap bool) string {
	matches := camelcase.Split(ident, -1)
	for i, m := range matches {
		if initialCap || i > 0 {
			m = capFirst(m)
		}
		matches[i] = acronyms.ReplaceAllStringFunc(m, func(c string) string {
			if len(c) > 4 {
				return strings.ToUpper(c[:2]) + c[2:]
			}
			return strings.ToUpper(c)
		})
	}
	return strings.Join(matches, "")
}

func capFirst(ident string) string {
	r, n := utf8.DecodeRuneInString(ident)
	return string(unicode.ToUpper(r)) + ident[n:]
}
