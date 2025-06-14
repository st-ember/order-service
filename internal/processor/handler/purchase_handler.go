package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/st-ember/ecommerceprocessor/internal/cache"
	"github.com/st-ember/ecommerceprocessor/internal/enum"
	"github.com/st-ember/ecommerceprocessor/internal/processor/redis"
	"github.com/st-ember/ecommerceprocessor/internal/processor/storage"
	"github.com/st-ember/ecommerceprocessor/internal/request"
)

func CreatePurchase(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot access request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var purchase request.Purchase
	err = json.Unmarshal(body, &purchase)
	if err != nil {
		http.Error(w, "Request body faulty", http.StatusBadRequest)
		return
	}

	data := cache.Purchase{
		Timestamp: time.Now().UnixNano(),
		Body: request.Purchase{
			Product:  purchase.Product,
			Customer: purchase.Customer,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Cannot serialize request", http.StatusInternalServerError)
		return
	}

	if err := redis.LPush("purchase_que", jsonData); err != nil {
		http.Error(w, "Cannot enque request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Request enqueued successfully"))
}

func PurchaseStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Cannot parse id to uuid", http.StatusBadRequest)
		return
	}

	var status *enum.PurchaseStatus
	status, err = storage.PurchaseStatus(id)
	if err != nil {
		http.Error(w, "Cannot retrieve purchase status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status.String()))
}

func DeletePurchase(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Cannot parse id to uuid", http.StatusBadRequest)
		return
	}

	err = storage.DeletePurchase(id)
	if err != nil {
		http.Error(w, "Cannot delete purchase item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted purchase item"))
}
