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
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

// Simple test to parse json schema
func TestParseJsonSchema(t *testing.T) {
	inputStr, err := ioutil.ReadFile("./test_input.json")
	if err != nil {
		t.Fatalf("Could not read expected output file ./test_input.json")
	}

	// Parse the input json string
	schema, err := ParseSchema([]byte(inputStr))
	if err != nil {
		t.Fatalf("Error parsing json schema. Err: %v", err)
	}

	log.Printf("Parsed json schema: %+v", schema)

	// Generate the code
	goStr, err := schema.GenerateGo()
	if err != nil {
		t.Fatalf("Error generating go code. Err: %v", err)
	}

	// Write the output
	log.Debugf("Generated go code: \n\n%s", goStr)
	gotFile, _ := os.Create("./test_got.txt")
	fmt.Fprintln(gotFile, goStr)

	// Read the expected output file
	b, err := ioutil.ReadFile("./test_exp.txt")
	if err != nil {
		t.Fatalf("Could not read expected output file ./test_exp.txt")
	}

	// Make sure every line in expected output is present in the gotten output
	exp_lines := strings.Split(string(b), "\n")
	for _, line := range exp_lines {
		if !strings.Contains(goStr, line) {
			t.Fatalf("Generated code does not match expected output")
		}
	}
}
