package generators

var templates = map[string]string{
	"handlerFuncs": `
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
	// Validate parameters
	err := Validate{{initialCap .}}(obj)
	if err != nil {
		log.Errorf("Validate{{initialCap .}} retruned error for: %+v. Err: %v", obj, err)
		return err
	}

	// Check if we handle this object
	if objCallbackHandler.{{initialCap .}}Cb == nil {
		log.Errorf("No callback registered for {{.}} object")
		return errors.New("Invalid object type")
	}

	// Check if object already exists
	if collections.{{.}}s[obj.Key] != nil {
		// Perform Update callback
		err = objCallbackHandler.{{initialCap .}}Cb.{{initialCap .}}Update(collections.{{.}}s[obj.Key], obj)
		if err != nil {
			log.Errorf("{{initialCap .}}Update retruned error for: %+v. Err: %v", obj, err)
			return err
		}
	} else {
		// save it in cache
		collections.{{.}}s[obj.Key] = obj

		// Perform Create callback
		err = objCallbackHandler.{{initialCap .}}Cb.{{initialCap .}}Create(obj)
		if err != nil {
			log.Errorf("{{initialCap .}}Create retruned error for: %+v. Err: %v", obj, err)
			delete(collections.{{.}}s, obj.Key)
			return err
		}
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

	// Check if we handle this object
	if objCallbackHandler.{{initialCap .}}Cb == nil {
		log.Errorf("No callback registered for {{.}} object")
		return errors.New("Invalid object type")
	}

	// Perform callback
	err := objCallbackHandler.{{initialCap .}}Cb.{{initialCap .}}Delete(obj)
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
  `,
	"hdr": `
// {{.Name}}.go
// This file is auto generated by modelgen tool
// Do not edit this file manually

package {{.Name}}

import (
	"errors"
	"regexp"
	"net/http"
	"encoding/json"
	"github.com/contiv/objmodel/objdb/modeldb"
	"github.com/gorilla/mux"
	log "github.com/Sirupsen/logrus"
)

type HttpApiFunc func(w http.ResponseWriter, r *http.Request, vars map[string]string) (interface{}, error)
  `,
	"routeFunc": `
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
  `,
	"routeTmpl": `
	// Register {{.}}
	route = "/api/{{.}}s/{key}/"
	listRoute = "/api/{{.}}s/"
	log.Infof("Registering %s", route)
	router.Path(listRoute).Methods("GET").HandlerFunc(makeHttpHandler(httpList{{initialCap .}}s))
	router.Path(route).Methods("GET").HandlerFunc(makeHttpHandler(httpGet{{initialCap .}}))
	router.Path(route).Methods("POST").HandlerFunc(makeHttpHandler(httpCreate{{initialCap .}}))
	router.Path(route).Methods("PUT").HandlerFunc(makeHttpHandler(httpCreate{{initialCap .}}))
	router.Path(route).Methods("DELETE").HandlerFunc(makeHttpHandler(httpDelete{{initialCap .}}))

  `,
	"validateFunc": `
// Validate a {{.Name}} object
func Validate{{initialCap .Name}}(obj *{{initialCap .Name}}) error {
	// Validate key is correct
	keyStr := {{range $index, $element := .Key}}{{if eq 0 $index }}obj.{{initialCap .}} {{else}}+ ":" + obj.{{initialCap .}} {{end}}{{end}}
	if obj.Key != keyStr {
		log.Errorf("Expecting {{initialCap .Name}} Key: %s. Got: %s", keyStr, obj.Key)
		return errors.New("Invalid Key")
	}

	// Validate each field
	{{range $element := .Properties}}{{if eq $element.Type "int"}}{{if ne $element.Default ""}}
	if obj.{{initialCap $element.Name}} == 0 {
		obj.{{initialCap $element.Name}} = {{$element.Default}}
	}
{{end}} {{if ne $element.Min 0.0}}
	if obj.{{initialCap $element.Name}} < {{$element.Min}} {
		return errors.New("{{$element.Name}} Value Out of bound")
	}
{{end}} {{if ne $element.Max 0.0}}
	if obj.{{initialCap $element.Name}} > {{$element.Max}} {
		return errors.New("{{$element.Name}} Value Out of bound")
	}
{{end}} {{else if eq $element.Type "number"}} {{if ne $element.Default ""}}
	if obj.{{initialCap $element.Name}} == 0 {
		obj.{{$element.Name}} = {{$element.Default}}
	}
{{end}} {{if ne $element.Min 0.0}}
	if obj.{{initialCap $element.Name}} < {{$element.Min}} {
		return errors.New("{{$element.Name}} Value Out of bound")
	}
{{end}} {{if ne $element.Max 0.0}}
	if obj.{{initialCap $element.Name}} > {{$element.Max}} {
		return errors.New("{{$element.Name}} Value Out of bound")
	}
{{end}} {{else if eq $element.Type "bool"}} {{if ne $element.Default ""}}
	if obj.{{initialCap $element.Name}} == false {
		obj.{{initialCap $element.Name}} = {{$element.Default}}
	}
{{end}} {{else if eq $element.Type "string"}} {{if ne $element.Default ""}}
	if obj.{{initialCap $element.Name}} == "" {
		obj.{{initialCap $element.Name}} = {{$element.Default}}
	}
{{end}} {{if ne $element.Length 0}}
	if len(obj.{{initialCap $element.Name}}) > {{$element.Length}} {
		return errors.New("{{$element.Name}} string too long")
	}
{{end}} {{if ne $element.Format ""}}
	{{$element.Name}}Match := regexp.MustCompile("{{$element.Format}}")
	if {{$element.Name}}Match.MatchString(obj.{{initialCap $element.Name}}) == false {
		return errors.New("{{$element.Name}} string invalid format")
	}
{{end}} {{end}} {{end}}

	return nil
}

  `,
}
