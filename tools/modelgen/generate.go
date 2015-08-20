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
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/objmodel/tools/modelgen/generators"
	"github.com/contiv/objmodel/tools/modelgen/texthelpers"
)

var validPropertyTypes = []string{
	"string",
	"bool",
	"array",
	"number",
	"int",
}

// GenerateGo generates go code for the schema
func (s *Schema) GenerateGo() (string, error) {
	// Generate file headers
	outStr, err := s.GenerateGoHdrs()
	if err != nil {
		return "", err
	}

	// Generate structs
	structStr, err := s.GenerateGoStructs()
	if err != nil {
		log.Errorf("Error generating go structs. Err: %v", err)
		return "", err
	}

	// Merge the header and struct
	outStr = outStr + structStr

	// Merge rest handler
	str, err := s.GenerateGoFuncs()
	if err != nil {
		return "", err
	}

	return outStr + str, nil
}

// GenerateGoStructs generates go code from a schema
func (s *Schema) GenerateGoStructs() (string, error) {
	var goStr string

	//  Generate all object definitions
	for _, obj := range s.Objects {
		objStr, err := obj.GenerateGoStructs()
		if err != nil {
			return "", err
		}

		goStr += objStr
	}

	for _, name := range []string{"gostructs", "callbacks", "init", "register"} {
		str, err := generators.RunTemplate(name, s)
		if err != nil {
			return "", err
		}

		goStr += str
	}

	return goStr, nil
}

// GenerateGoHdrs generates go file headers
func (s *Schema) GenerateGoHdrs() (string, error) {
	return generators.RunTemplate("hdr", s)
}

func (s *Schema) GenerateGoFuncs() (string, error) {
	// Output the functions and routes
	return generators.RunTemplate("routeFunc", s)
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
		if err != nil {
			return "", err
		}
		goStr += propStr
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

func (prop *Property) GenerateGoStructs() (string, error) {
	var found bool

	for _, myType := range validPropertyTypes {
		if myType == prop.Type {
			found = true
		}
	}

	if !found {
		return "", errors.New("Unknown Property")
	}

	return generators.RunTemplate("propstruct", prop)
}
