package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"io/ioutil"
	"time"
	"fmt"
	"strings"
)

type Recipe struct {
	Name			string `json:"name"`
	Ingredients		string `json:"ingredients"`
	Instructions	string `json:"instructions"`
	ID				string `json:"id"`
}

type recipesHandlers struct {
	//concurrent access handling
	sync.Mutex 

	// map of all recipes 
	store map[string] Recipe
}

// Handles general operations
func (h *recipesHandlers) recipes (w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case "GET":
		h.get(w,r)
		return
	case "POST":
		h.post(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

// Handles individual operations
func (h *recipesHandlers) singleRecipe (w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case "GET":
		h.getRecipe(w,r)
		return
	case "UPDATE":
		h.updateRecipe(w,r)
		return
	case "DELETE":
		h.deleteRecipe(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

// Get all recipes 
func (h *recipesHandlers) get(w http.ResponseWriter, r *http.Request){
	recipes := make([]Recipe,len(h.store))
	
	// handle concurrent access by locking
	h.Lock()
	i:= 0
	for _, recipe := range h.store{
		recipes[i] = recipe
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(recipes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/server")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Add a recipe
func (h *recipesHandlers) post(w http.ResponseWriter, r *http.Request){
	bodyBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var recipe Recipe
	err = json.Unmarshal(bodyBytes, &recipe)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	recipe.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[recipe.ID] = recipe
	defer h.Unlock()
}

// Get a specific recipe with ID
func (h *recipesHandlers) getRecipe(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	recipe, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Delete a specific recipe with ID
func (h *recipesHandlers) deleteRecipe(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	_, ok := h.store[parts[2]]
	h.Unlock()
	fmt.Sprintf("%s", h.store[parts[2]])
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Lock()
	delete(h.store,parts[2])
	h.Unlock()

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Update a specific recipe but doesn't allow one field manipulation
// TODO: to improve, we could only give the change the the rest of JSON should remain the same
func (h *recipesHandlers) updateRecipe(w http.ResponseWriter, r *http.Request){
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var recipe Recipe
	err = json.Unmarshal(bodyBytes, &recipe)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	h.Lock()
	_, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	recipe.ID = parts[2]
	h.Lock()
	h.store[parts[2]] = recipe
	defer h.Unlock()
}

func newRecipesHandlers() *recipesHandlers{
	return &recipesHandlers{
		store: map[string]Recipe{
			"id1": Recipe{
				Name: "Honey Garlic Glazed Salmon",
				Instructions: "whisk, heat,season, cook, flip, add, garnish, serve",
				Ingredients: "honey, soy sauce, lemon juice, red pepper, olive oil, salmon, salt, black pepper, garlic , lemon",
				ID: "id1",
			},
		},
	}
}

func main() {
	recipesHandlers := newRecipesHandlers()
	http.HandleFunc("/recipes", recipesHandlers.recipes)
	http.HandleFunc("/recipes/", recipesHandlers.singleRecipe)
	err:= http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}