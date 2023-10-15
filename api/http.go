package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	svc "heyalley-server/image/service"
)

type HttpService struct {
	ImageStorageService svc.ImageStore
}

func (s *HttpService) StoreImage(w http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Access the photo key - First Approach
	file, h, err := request.FormFile("photo")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.ImageStorageService.StoreImage(context.Background(), file, h)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

type ImageSearchRequest struct {
	SearchKey string `json:"searchkey"`
}

func (s *HttpService) SearchImage(w http.ResponseWriter, r *http.Request) {
	v := ImageSearchRequest{}

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	getImageRes, err := s.ImageStorageService.GetImage(context.Background(), v.SearchKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "image: %v", getImageRes)
}
