package v1

import (
	"encoding/json"
	"fmt"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	Base    = 10
	BitSize = 64
)

func (ur *userRoutes) getAllNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := ur.GetAllNodes(r.Context())
	if err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	renderJSON(w, nodes)
}

func (ur *userRoutes) getNode(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		renderError(w, http.StatusBadRequest, fmt.Errorf("empty node id"))
		return
	}

	id, err := strconv.ParseUint(idStr, Base, BitSize)
	if err != nil {
		renderError(w, http.StatusBadRequest, err)
		return
	}

	node, err := ur.GetNodeWithRelationships(r.Context(), id)
	if err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	if node == nil {
		renderError(w, http.StatusNotFound, fmt.Errorf("node not found"))
		return
	}

	renderJSON(w, node)
}

func (ur *userRoutes) createNode(w http.ResponseWriter, r *http.Request) {
	nodeType := strings.ToLower(r.URL.Query().Get("type"))
	switch nodeType {
	case "user":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			renderError(w, http.StatusBadRequest, err)
		}

		var user models.User
		if err := json.Unmarshal(body, &user); err != nil {
			renderError(w, http.StatusBadRequest, err)
		}

		if err := ur.SaveUser(r.Context(), user); err != nil {
			renderError(w, http.StatusInternalServerError, err)
		}
	case "group":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			renderError(w, http.StatusBadRequest, err)
		}

		var group models.GroupWithSubscribers
		if err := json.Unmarshal(body, &group); err != nil {
			renderError(w, http.StatusBadRequest, err)
		}

		if err := ur.SaveGroup(r.Context(), group); err != nil {
			renderError(w, http.StatusInternalServerError, err)
		}
	default:
		renderError(w, http.StatusBadRequest, fmt.Errorf("empty or invalid type of node: %s, use user or group", nodeType))
	}

	w.WriteHeader(http.StatusCreated)
}

func (ur *userRoutes) deleteNode(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		renderError(w, http.StatusBadRequest, fmt.Errorf("empty node id"))
		return
	}

	id, err := strconv.ParseUint(idStr, Base, BitSize)
	if err != nil {
		renderError(w, http.StatusBadRequest, err)
		return
	}

	if err := ur.DeleteNode(r.Context(), id); err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
