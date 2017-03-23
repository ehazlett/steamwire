package server

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	return r, nil
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version.FullVersion() + "\n"))
}

func (s *Server) syncHandler(w http.ResponseWriter, r *http.Request) {
	s.Sync()
}

func (s *Server) addHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID := vars["appID"]

	if appID == "" {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	if err := s.AddApp(appID); err != nil {
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

	if err := s.DeleteApp(appID); err != nil {
		http.Error(w, fmt.Sprintf("error deleting app: %s", err), http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"app": appID,
	}).Info("deleted app")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	apps, err := s.GetApps()
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
	appNews, err := s.GetNews(appID)
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
