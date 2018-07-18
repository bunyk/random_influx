package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/docker/docker/pkg/namesgenerator"
)

func ServeRandomMetrics(port string) {
	randometer := CreateRandometer()
	http.HandleFunc("/metrics-app", randometer.HandleRequest)
	http.ListenAndServe(port, nil)
}

type Randometer struct {
	Metrics    map[string]map[string]int
	Platforms  []string
	Affiliates []string
}

func CreateRandometer() *Randometer {
	r := Randometer{}
	r.Platforms = []string{"android", "ios", "web"}
	r.Affiliates = make([]string, 5)
	for i, _ := range r.Affiliates {
		r.Affiliates[i] = namesgenerator.GetRandomName(0)
	}

	r.Metrics = make(map[string]map[string]int)
	for _, p := range r.Platforms {
		r.Metrics[p] = make(map[string]int)
	}

	return &r
}

type Metric struct {
	Affiliate string
	Platform  string
	Received  int
}

func (r Randometer) GetMetrics() []Metric {
	res := make([]Metric, len(r.Affiliates)*len(r.Platforms))
	for i, p := range r.Platforms {
		for j, a := range r.Affiliates {
			metric := Metric{
				Affiliate: a,
				Platform:  p,
				Received:  r.Metrics[p][a],
			}
			res[i*len(r.Affiliates)+j] = metric
			r.Metrics[p][a] += rand.Intn(5)
		}
	}
	return res
}

func (r Randometer) HandleRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		metrics := r.GetMetrics()
		msg, err := json.Marshal(metrics)
		if err != nil {
			errorHandler(w, fmt.Errorf("Metrics error: %s", err.Error), http.StatusInternalServerError)
		} else {
			w.Write(msg)
		}
	} else {
	}
}

func errorHandler(w http.ResponseWriter, err error, code int) bool {
	if err == nil {
		return false
	}
	fmt.Println(err)
	msg, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	w.Write(msg)
	return true
}
