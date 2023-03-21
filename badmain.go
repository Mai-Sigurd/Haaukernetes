//https://medium.com/programming-kubernetes/building-stuff-with-the-kubernetes-api-part-4-using-go-b1d0e3c1c899
//https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go

package main

import (
	"log"
	"net/http"
	"os"
	_ "sync"
	_ "sync/atomic"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/middleware/untyped"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatalln("this command requires the swagger spec as argument")
	}
	log.Printf("loading %q as contract for the server", os.Args[1])

	specDoc, err := loads.Spec(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	// our spec doesn't have application/json in the consumes or produces
	// so we need to clear those settings out
	api := untyped.NewAPI(specDoc).WithoutJSONDefaults()

	// register serializers
	mediaType := "application/io.goswagger.examples.todo-list.v1+json"
	api.DefaultConsumes = mediaType
	api.DefaultProduces = mediaType
	api.RegisterConsumer(mediaType, runtime.JSONConsumer())
	api.RegisterProducer(mediaType, runtime.JSONProducer())

	// register the operation handlers
	// Namespaces
	api.RegisterOperation("GET", "/namespace/{id}", getNamespace)
	api.RegisterOperation("POST", "/namespace", postNamespace)
	api.RegisterOperation("DELETE", "/{id}", deleteNamespace)

	// Kali
	api.RegisterOperation("GET", "/kali/{id}", getKali)
	// ID beceause we need the namespace, but maybe thats in the body instead????
	api.RegisterOperation("POST", "/kali", postKali)
	api.RegisterOperation("DELETE", "/kali/{id}", deleteKali)

	// Challenge
	api.RegisterOperation("POST", "/challenge", postChallenge)

	api.RegisterOperation("DELETE", "/challenge/{id}", deleteChallenge)

	// validate the API descriptor, to ensure we don't have any unhandled operations
	if err := api.Validate(); err != nil {
		log.Fatalln(err)
	}

	// construct the application context for this server
	// use the loaded spec document and the api descriptor with the default router
	app := middleware.NewContext(specDoc, api, nil)

	log.Println("serving", specDoc.Spec().Info.Title, "at http://localhost:8000")

	// serve the api with spec and UI
	if err := http.ListenAndServe(":8000", app.APIHandler(nil)); err != nil {
		log.Fatalln(err)
	}
}

var getNamespace = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'getNamespace'")
	log.Printf("%#v\n", params)

	return nil, nil
})

var postNamespace = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'postNamespace'")
	log.Printf("%#v\n", params)

	return nil, nil
})
var deleteNamespace = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'deleteNamespace'")
	log.Printf("%#v\n", params)

	return nil, nil
})

var getKali = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'getKali'")
	log.Printf("%#v\n", params)

	return nil, nil
})

var postKali = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'postKali'")
	log.Printf("%#v\n", params)

	return nil, nil
})

var deleteKali = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'deleteKali'")
	log.Printf("%#v\n", params)

	return nil, nil
})

var postChallenge = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'postChallenge'")
	log.Printf("%#v\n", params)

	return nil, nil
})

var deleteChallenge = runtime.OperationHandlerFunc(func(params interface{}) (interface{}, error) {
	log.Println("received 'deleteChallenge'")
	log.Printf("%#v\n", params)

	return nil, nil
})
