package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	port = ":80"
)

var session *mgo.Session

func main() {

	var err error

	//TODO move to config
	session, err = mgo.Dial("52.168.20.79")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	http.HandleFunc("/search", search)
	http.HandleFunc("/list", list)

	log.Println("Starting server at port", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Panicln("erro", err)

	}
}

func list(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	fmt.Println(page)
	//TODO handle errors
	p, _ := strconv.Atoi(page)
	b, _ := listFromDB(p)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)

}

func search(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	b, _ := searchFromDB(name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func searchFromDB(name string) ([]byte, error) {
	var gh []interface{}
	session.DB("movies").C("list").Find(bson.M{"$text": bson.M{"$search": name, "$diacriticSensitive": true}}).All(&gh)
	return bson.MarshalJSON(gh)

}

func listFromDB(page int) ([]byte, error) {
	skip := page * 10
	var gh []interface{}
	session.DB("movies").C("list").Find(bson.M{}).Skip(skip).Limit(10).All(&gh)
	return bson.MarshalJSON(gh)
}
