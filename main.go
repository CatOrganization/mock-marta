package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	responsesDir = "./responses/"

	realTimeTrainsPath         = "/RealtimeTrain/RestServiceNextTrain/GetRealtimeArrivals"
	realTimeTrainsJSONFileName = "real_time_trains.json"

	realTimeAllBusesPath         = "/BRDRestService/RestBusRealTimeService/GetAllBus"
	realTimeAllBusesJSONFileName = "real_time_all_buses.json"

	realTimeBusPath         = "/BRDRestService/RestBusRealTimeService/GetBusByRoute/{route}"
	realTimeBusJSONFileName = "real_time_bus.json"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(responsesDir))))

	r.HandleFunc(realTimeTrainsPath, realTimeTrains)
	r.HandleFunc(realTimeAllBusesPath, realTimeAllBuses)
	r.HandleFunc(realTimeBusPath, realTimeBus)

	return r
}

func realTimeTrains(w http.ResponseWriter, r *http.Request) {
	proxyToStaticFile(realTimeTrainsJSONFileName, w, r)
}

func realTimeAllBuses(w http.ResponseWriter, r *http.Request) {
	proxyToStaticFile(realTimeAllBusesJSONFileName, w, r)
}

func realTimeBus(w http.ResponseWriter, r *http.Request) {
	proxyToStaticFile(realTimeBusJSONFileName, w, r)
}

func proxyToStaticFile(fileName string, w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://" + r.Host + "/static/" + fileName)
	if err != nil {
		log.Errorf("failed to proxy static file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	err = resp.Write(w)
	if err != nil {
		log.Errorf("failed to write response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	r := routes()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(srv.ListenAndServe())
}
