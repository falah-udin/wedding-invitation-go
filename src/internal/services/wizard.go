package services

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	WizardSessionName    = "wedding_wizard"
	WizardProjectIDKey   = "wizard_project_id"
	WizardClientIDKey    = "wizard_client_id"
	WizardMaxAge         = 86400 * 3 // 3 hari untuk wizard
)

// WizardGetProjectID — ambil project_id dari session wizard
func WizardGetProjectID(r *http.Request) uint {
	store := GetStore()
	session, err := store.Get(r, WizardSessionName)
	if err != nil {
		return 0
	}

	val, ok := session.Values[WizardProjectIDKey]
	if !ok {
		return 0
	}

	switch v := val.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case int64:
		return uint(v)
	case float64:
		return uint(v)
	}
	return 0
}

// WizardSetProjectID — simpan project_id ke session wizard
func WizardSetProjectID(w http.ResponseWriter, r *http.Request, projectID uint) error {
	store := GetStore()
	session, err := store.Get(r, WizardSessionName)
	if err != nil {
		return err
	}

	session.Values[WizardProjectIDKey] = projectID
	session.Options.MaxAge = WizardMaxAge

	return session.Save(r, w)
}

// WizardGetClientID — ambil client_id dari session wizard
func WizardGetClientID(r *http.Request) uint {
	store := GetStore()
	session, err := store.Get(r, WizardSessionName)
	if err != nil {
		return 0
	}

	val, ok := session.Values[WizardClientIDKey]
	if !ok {
		return 0
	}

	switch v := val.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case int64:
		return uint(v)
	case float64:
		return uint(v)
	}
	return 0
}

// WizardSetClientID — simpan client_id ke session wizard
func WizardSetClientID(w http.ResponseWriter, r *http.Request, clientID uint) error {
	store := GetStore()
	session, err := store.Get(r, WizardSessionName)
	if err != nil {
		return err
	}

	session.Values[WizardClientIDKey] = clientID
	session.Options.MaxAge = WizardMaxAge

	return session.Save(r, w)
}

// WizardClear — hapus session wizard (setelah publish / batal)
func WizardClear(w http.ResponseWriter, r *http.Request) {
	store := GetStore()
	session, err := store.New(r, WizardSessionName)
	if err != nil {
		return
	}

	session.Options.MaxAge = -1
	session.Save(r, w)
}

// suppress unused
var _ sessions.Session
