// golang-rest-mongo project main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
)

// error response contains everything we need to use http.Error
type handlerError struct {
	Error   error
	Message string
	Code    int
}

//Mongo collection & session
var (
	session    *mgo.Session
	collection *mgo.Collection
)

// a custom type that we can use for handling errors and formatting responses
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)

// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

//get all of the studies from the data repository
func listStudies(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	log.Printf("Call to studies list")

	studies := make([]study, 0)
	log.Printf("Call to mongo to get studies")
	iter := collection.Find(nil).Iter()
	result := study{}
	for iter.Next(&result) {
		log.Printf(result.StudyName)
		studies = append(studies, result)
	}
	log.Printf("Number of studies retrieved %d", len(studies))
	return studies, nil
}

//add a study to the data repository
func addOrUpdateStudy(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	payload, e := parseStudyRequest(r)
	if e != nil {
		return nil, e
	}

	//If no id then set one
	if payload.Id == "" {
		payload.Id = bson.NewObjectId()
	}

	//perform the upsert
	_, err := collection.UpsertId(payload.Id, payload)
	if err != nil {
		return nil, &handlerError{err, "Could not upsert study", http.StatusBadRequest}
	} else {
		log.Printf("Inserted new study %s with name %s", payload.Id, payload.StudyName)
	}

	//log.Printf("Records Updated: %d  MongoId: %d", info.Updated, info.UpsertedId)

	// we return the study we just made so the client can see the ID if they want
	return payload, nil
}

//get a study based on id
func getStudy(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	// mux.Vars grabs variables from the path
	id := mux.Vars(r)["id"]
	log.Printf("Trying to find study: " + id)

	s := study{}
	e := collection.FindId(bson.ObjectIdHex(id)).One(&s)

	if e != nil {
		return nil, &handlerError{nil, "Could not find study " + id, http.StatusNotFound}
	}

	return s, nil
}

//remove the study from the database
func removeStudy(w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	// mux.Vars grabs variables from the path
	vars := mux.Vars(r)
	id := vars["id"]

	// Remove it from database
	err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Printf("Could not find study %s to delete", id)
	}

	return make(map[string]string), nil
}

func parseStudyRequest(r *http.Request) (study, *handlerError) {
	// the study payload is in the request body
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return study{}, &handlerError{e, "Could not read request", http.StatusBadRequest}
	}

	// turn the request body (JSON) into a study object
	var payload study
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return study{}, &handlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	return payload, nil
}

func main() {

	// command line flags
	port := flag.Int("port", 9000, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()
	router.Handle("/", http.RedirectHandler("/static/", 302))
	router.Handle("/studies", handler(listStudies)).Methods("GET")
	router.Handle("/studies", handler(addOrUpdateStudy)).Methods("POST")
	router.Handle("/studies/{id}", handler(getStudy)).Methods("GET")
	router.Handle("/studies/{id}", handler(addOrUpdateStudy)).Methods("POST")
	router.Handle("/studies/{id}", handler(removeStudy)).Methods("DELETE")

	//Base routes
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileHandler))
	http.Handle("/", router)

	//MongoDB connection
	log.Println("Starting mongo db session")
	var mongoErr error
	session, mongoErr = mgo.Dial("localhost")
	if mongoErr != nil {
		panic(mongoErr)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("test").C("studies")

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())

}
