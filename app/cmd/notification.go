package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"go.uber.org/zap"
)

type Subscription struct {
	Sub    webpush.Subscription `json:"sub"`
	Active bool                 `json:"active"`
}

func loadSubscriptions() {
	file, err := os.Open("subscriptions.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&subscriptions)
}

func saveSubscriptions() {
	file, err := os.Create("subscriptions.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(&subscriptions)
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var sub webpush.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	zapLogger.Debug("Subscribing", zap.String("id", id))

	subscriptionsMutex.Lock()
	subscriptions[id] = Subscription{Sub: sub, Active: true}
	subscriptionsMutex.Unlock()
	saveSubscriptions()
	w.WriteHeader(http.StatusOK)
}

func UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	subscriptionsMutex.Lock()
	defer subscriptionsMutex.Unlock()

	if subscription, ok := subscriptions[id]; ok {
		subscription.Active = false
		subscriptions[id] = subscription
		saveSubscriptions()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unsubscribed successfully"))
	} else {
		http.Error(w, "Subscription not found", http.StatusNotFound)
	}
}

func PublishHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	var payload map[string]string
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate options in the backend
	options := map[string]interface{}{
		"body":               payload["body"],
		"tag":                "notification: " + time.Now().String(),
		"icon":               "/assets/img/icon/cloud-192.png",
		"badge":              "/assets/img/icon/cloud-192.png",
		"requireInteraction": false,
	}

	message, _ := json.Marshal(map[string]interface{}{
		"title":   payload["title"],
		"options": options,
	})

	id := r.URL.Query().Get("id")
	// zapLogger.Debug("Publishing", zap.String("id", id), zap.String("message", string(message)))
	subscriptionsMutex.Lock()
	defer subscriptionsMutex.Unlock()
	if sub, ok := subscriptions[id]; ok {
		if sub.Active {
			resp, err := webpush.SendNotification(message, &sub.Sub, &webpush.Options{
				Subscriber:      "mailto:example@example.com",
				VAPIDPublicKey:  "BDOpUfHEw7LFRJWhDxF5TW7SR-kiaOY-_6iFrVweY8rfmi9ySzjxSGWbbm-wwriXwAYWVX5808Pb2U2ApYXYKLc",
				VAPIDPrivateKey: "TkyndbWdGc_D3ukx9tbfh5_ElMjRzL0ixQ86JAMtDzI",
			})
			if err != nil {
				zapLogger.Error("Failed to send notification", zap.Error(err))
			} else {
				zapLogger.Info("Notification sent", zap.Int("status", resp.StatusCode))
				body, _ := io.ReadAll(resp.Body)
				zapLogger.Debug("Response Body", zap.String("body", string(body)))
			}
		} else {
			zapLogger.Info("Subscription is inactive", zap.String("id", id))
		}
	} else {
		zapLogger.Info("Subscription not found", zap.String("id", id))
	}
	w.WriteHeader(http.StatusOK)
}

func init() {
	loadSubscriptions()
}
