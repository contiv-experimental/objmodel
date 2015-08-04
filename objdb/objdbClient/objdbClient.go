package objdbClient

import (
	log "github.com/Sirupsen/logrus"
	"github.com/contiv/objmodel/objdb"
	"github.com/contiv/objmodel/objdb/plugins"
)

// Create a new conf store
func NewClient() objdb.ObjdbApi {
	defaultConfStore := "etcd"

	// Init all plugins
	plugins.Init()

	// Get the plugin
	plugin := objdb.GetPlugin(defaultConfStore)

	// Initialize the objdb client
	err := plugin.Init([]string{})
	if err != nil {
		log.Errorf("Error initializing confstore plugin. Err: %v", err)
		log.Fatal("Error initializing confstore plugin")
	}

	return plugin
}