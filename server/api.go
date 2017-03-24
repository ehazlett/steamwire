package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ehazlett/steamwire/version"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *Server) router() (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/", s.indexHandler)
	r.HandleFunc("/sync", s.syncHandler)
	r.HandleFunc("/apps", s.getHandler).Methods("GET")
	r.HandleFunc("/apps/{appID:.*}", s.addHandler).Methods("POST")
	r.HandleFunc("/apps/{appID:.*}", s.deleteHandler).Methods("DELETE")
	r.HandleFunc("/apps/{appID:.*}/news", s.getNewsHandler).Methods("GET")
	r.HandleFunc("/applist", s.getAppListHandler).Methods("GET")
	r.HandleFunc("/applist/update", s.appListUpdateHandler).Methods("POST")
	r.HandleFunc("/applist/search", s.appListSearchHandler).Methods("POST")

	return r, nil
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version.FullVersion() + "\n"))
}

func (s *Server) syncHandler(w http.ResponseWriter, r *http.Request) {
	s.sync()
}

func (s *Server) addHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID := vars["appID"]

	if appID == "" {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	if err := s.ds.AddApp(appID); err != nil {
		http.Error(w, fmt.Sprintf("error adding app: %s", err), http.StatusInternalServerError)
		return
	}
	logrus.WithFields(logrus.Fields{
		"app": appID,
	}).Info("added app")
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID := vars["appID"]

	if appID == "" {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	if err := s.ds.DeleteApp(appID); err != nil {
		http.Error(w, fmt.Sprintf("error deleting app: %s", err), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"app": appID,
	}).Info("deleted app")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	apps, err := s.ds.GetApps()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting apps: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, fmt.Sprintf("error encoding apps: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getNewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID := vars["appID"]

	if appID == "" {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}
	appNews, err := s.getNews(appID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting app news: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(appNews); err != nil {
		http.Error(w, fmt.Sprintf("error encoding apps: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getAppListHandler(w http.ResponseWriter, r *http.Request) {
	apps, err := s.ds.GetAppList()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting app list: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, fmt.Sprintf("error encoding app list: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) appListUpdateHandler(w http.ResponseWriter, r *http.Request) {
	force := false
	q := r.URL.Query()
	if v, ok := q["force"]; ok {
		switch strings.ToLower(v[0]) {
		case "false", "0":
			// false
		default:
			// everything else is true
			force = true
		}
	}

	if err := s.updateAppList(force); err != nil {
		http.Error(w, fmt.Sprintf("error updating app list: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) appListSearchHandler(w http.ResponseWriter, r *http.Request) {
	q, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("you must POST a query in the body: %s", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	apps, err := s.ds.FindApp(string(q))
	if err != nil {
		http.Error(w, fmt.Sprintf("error searching app list: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, fmt.Sprintf("error encoding app list results: %s", err), http.StatusInternalServerError)
		return
	}
}
