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

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/objmodel/tools/modelgen/generators"
	"github.com/contiv/objmodel/tools/modelgen/texthelpers"
)

// GenerateGo generates go code for the schema
func (s *Schema) GenerateGo() (string, error) {
	// Generate file headers
	outStr := s.GenerateGoHdrs()

	// Generate structs
	structStr, err := s.GenerateGoStructs()
	if err != nil {
		log.Errorf("Error generating go structs. Err: %v", err)
		return "", err
	}

	// Merge the header and struct
	outStr = outStr + structStr

	// Merge rest handler
	outStr = outStr + s.GenerateGoFuncs()

	return outStr, nil
}

// GenerateGoStructs generates go code from a schema
func (s *Schema) GenerateGoStructs() (string, error) {
	var goStr string

	//  Generate all object definitions
	for _, obj := range s.Objects {
		objStr, err := obj.GenerateGoStructs()
		if err == nil {
			goStr += objStr
		}
	}

	buf := new(bytes.Buffer)

	tmpl := generators.GetTemplate("gostructs")
	if err := tmpl.Execute(buf, s); err != nil {
		return "", err
	}

	goStr += buf.String()

	tmpl = generators.GetTemplate("callbacks")
	if err := tmpl.Execute(buf, s); err != nil {
		return "", err
	}

	goStr += buf.String()

	// generate callback handler

	// Generate an Init function
	goStr = goStr + fmt.Sprintf("\nfunc Init() {\n")
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("	collections.%ss = make(map[string]*%s)\n", obj.Name, texthelpers.InitialCap(obj.Name))
	}
	goStr = goStr + fmt.Sprintf("\n")
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("	restore%s()\n", texthelpers.InitialCap(obj.Name))
	}

	goStr = goStr + fmt.Sprintf("}\n\n")

	// Generate callback register functions
	for _, obj := range s.Objects {
		goStr = goStr + fmt.Sprintf("func Register%sCallbacks(handler %sCallbacks) {\n", texthelpers.InitialCap(obj.Name), texthelpers.InitialCap(obj.Name))
		goStr = goStr + fmt.Sprintf("	objCallbackHandler.%sCb = handler\n", texthelpers.InitialCap(obj.Name))
		goStr = goStr + fmt.Sprintf("}\n\n")
	}
	return goStr, nil
}

// GenerateGoHdrs generates go file headers
func (s *Schema) GenerateGoHdrs() string {
	var buf bytes.Buffer

	tmpl := generators.GetTemplate("hdr")
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

	// Output the functions and routes
	rfTmpl := generators.GetTemplate("routeFunc")
	rfTmpl.Execute(&buf, "")
	goStr = goStr + buf.String()

	// add a path for each object
	for _, obj := range s.Objects {
		var buf bytes.Buffer

		// Create a template, add the function map, and parse the text.
		tmpl := generators.GetTemplate("routeTmpl")

		// Run the template.
		if err := tmpl.Execute(&buf, obj.Name); err != nil {
			log.Fatalf("execution: %s", err)
		}

		goStr = goStr + buf.String()
	}

	goStr = goStr + fmt.Sprintf("\n}\n")

	// Generate REST handlers for each object
	for _, obj := range s.Objects {
		var buf bytes.Buffer
		// Create a template, add the function map, and parse the text.
		tmpl := generators.GetTemplate("handlerFuncs")

		// Run the template.
		if err := tmpl.Execute(&buf, obj.Name); err != nil {
			log.Fatalf("execution: %s", err)
		}

		goStr = goStr + buf.String()

		//  Generate object validators
		objStr, err := obj.GenerateValidate()
		if err == nil {
			goStr = goStr + objStr
		}
	}

	return goStr
}

func (obj *Object) GenerateGoStructs() (string, error) {
	var goStr string

	objName := texthelpers.InitialCap(obj.Name)
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
			goStr = goStr + fmt.Sprintf("	%s	map[string]modeldb.Link		`json:\"%s,omitempty\"`\n", texthelpers.InitialCap(lsName), lsName)
		}
		goStr = goStr + fmt.Sprintf("}\n\n")
	}

	// Define object's links
	if len(obj.Links) > 0 {
		goStr = goStr + fmt.Sprintf("type %sLinks struct {\n", objName)
		for lName := range obj.Links {
			goStr = goStr + fmt.Sprintf("	%s	modeldb.Link		`json:\"%s,omitempty\"`\n", texthelpers.InitialCap(lName), lName)
		}
		goStr = goStr + fmt.Sprintf("}\n\n")
	}

	return goStr, nil
}

func (obj *Object) GenerateValidate() (string, error) {
	var goStr string

	var buf bytes.Buffer
	// Create a template, add the function map, and parse the text.
	tmpl := generators.GetTemplate("validateFunc")

	// Run the template.
	if err := tmpl.Execute(&buf, obj); err != nil {
		log.Fatalf("execution: %s", err)
	}

	goStr = goStr + buf.String()

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

	goStr = fmt.Sprintf("	%s	", texthelpers.InitialCap(prop.Name))
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
