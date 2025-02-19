package webhook

import (
	"fmt"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"tarscontroller/controller"
	"tarscontroller/webhook/conversion"
	"tarscontroller/webhook/mutating"
	"tarscontroller/webhook/validating"
	"time"
)

const CertFile = "/etc/tarscontroller-cert/tls.crt"
const CertKey = "/etc/tarscontroller-cert/tls.key"

type Webhook struct {
	mutating   *mutating.Mutating
	validating *validating.Validating
	conversion *conversion.Conversion
}

func New(clients *controller.Clients, informers *controller.Informers) *Webhook {
	webhook := &Webhook{
		conversion: conversion.New(),
		mutating:   mutating.New(clients, informers),
		validating: validating.New(clients, informers),
	}
	return webhook
}

func (h *Webhook) Start(stopCh chan struct{}) {
	go wait.Until(func() {
		validatingFunc := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Connection", "keep-alive")
			h.validating.Handle(w, r)
		}

		mutatingFunc := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Connection", "keep-alive")
			h.mutating.Handle(w, r)
		}

		conversionFunc := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Add("Connection", "keep-alive")
			h.conversion.Handle(w, r)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/validating", validatingFunc)
		mux.HandleFunc("/mutating", mutatingFunc)
		mux.HandleFunc("/conversion", conversionFunc)

		srv := &http.Server{
			Addr:              ":443",
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      12 * time.Second,
		}
		// ListenAndServe always returns a non-nil error. After Shutdown or Close,
		// the returned error is ErrServerClosed.
		err := srv.ListenAndServeTLS(CertFile, CertKey)
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf("will exist because : %s \n", err.Error()))
			close(stopCh)
		}
	}, time.Second, stopCh)
}
