package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// /api/address
func addressRootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllAddress(w, r)
	case http.MethodPost:
		createAddress(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// /api/address/{id}
func addressByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAddressByID(w, r, id)
	case http.MethodPut:
		updateAddress(w, r, id)
	case http.MethodDelete:
		deleteAddress(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// CRUD
func getAllAddress(w http.ResponseWriter, r *http.Request) {
	var list []AddressBook
	DB.Find(&list)
	respondJSON(w, http.StatusOK, "success", list)
}

func getAddressByID(w http.ResponseWriter, r *http.Request, id int) {
	var item AddressBook
	if err := DB.First(&item, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Not found")
		return
	}
	respondJSON(w, http.StatusOK, "success", item)
}

func createAddress(w http.ResponseWriter, r *http.Request) {
	var input AddressBook
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	DB.Create(&input)
	respondJSON(w, http.StatusCreated, "created", input)
}

func updateAddress(w http.ResponseWriter, r *http.Request, id int) {
	var item AddressBook
	if err := DB.First(&item, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Not found")
		return
	}

	var input AddressBook
	json.NewDecoder(r.Body).Decode(&input)

	item.Firstname = input.Firstname
	item.Lastname = input.Lastname
	item.Code = input.Code
	item.Phone = input.Phone

	DB.Save(&item)
	respondJSON(w, http.StatusOK, "updated", item)
}

func deleteAddress(w http.ResponseWriter, r *http.Request, id int) {
	if err := DB.Delete(&AddressBook{}, id).Error; err != nil {
		respondError(w, http.StatusNotFound, "Not found")
		return
	}
	respondJSON(w, http.StatusOK, "deleted", nil)
}

// Utils
func extractID(path string) (int, error) {
	parts := strings.Split(path, "/")
	return strconv.Atoi(parts[len(parts)-1])
}
