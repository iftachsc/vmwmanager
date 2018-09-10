package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/iftachsc/contracts"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router        *mux.Router
	IaasEndpoints map[string]contracts.IaasEndpoint
	ctx           context.Context
}

//Initialize this is great func
func (a *App) Initialize(user, password, dbname string) {

	a.ctx = context.Background()
	a.IaasEndpoints = make(map[string]contracts.IaasEndpoint)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/vms", a.getVms).Methods("GET")
	//a.Router.HandleFunc("/vm", a.registerVm).Methods("POST")
	a.Router.HandleFunc("/scsi_luns", a.getScsiLunDisks).Methods("GET")
	a.Router.HandleFunc("/hosts", a.getHosts).Methods("GET")
	// a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	// a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	// a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

//Run is
func (a *App) Run(addr string) {
	println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getIaasEndpointForLocation(locationUuid string) (contracts.IaasEndpoint, error) {

	iaas, exists := a.IaasEndpoints[locationUuid]
	if !exists {
		//getLocation from location microservice (*contracts.Location)
		//return location.InitizalizeClient() (vimclient)
		location := contracts.Location{ //instead of location microservice
			IaasProvider: &VmwareVcenter{
				Host:     "https://vcenter01.mgmt.il-center-1.cloudzone.io",
				User:     "iftachsc",
				Password: "5tw5j;M]HVN0$",
				Client:   nil,
			},
			StorageFilers: nil,
		}

		//TODO: Error handling here
		err := location.IaasProvider.InitializeClient(a.ctx)
		fmt.Println("---->>>", location.IaasProvider.(*VmwareVcenter).Client)

		if err != nil {
			return nil, err
		}
		a.IaasEndpoints[locationUuid] = location.IaasProvider

		return location.IaasProvider, nil
	} else {
		return iaas, nil
	}
}

func (a *App) getVms(w http.ResponseWriter, r *http.Request) {

	// Returns a url.Values, which is a map[string][]string
	params := r.URL.Query()
	locationUUID, ok := params["locationUuid"]

	if !ok {
		respondWithError(w, http.StatusBadRequest, "Missing locationUuid parameter")
		return
	}
	iaas, err := a.getIaasEndpointForLocation(locationUUID[0])
	//fmt.Println("---->>>", iaas.(*VmwareVcenter).Client)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	vms, err := iaas.GetVMs(a.ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())

	} else {
		respondWithJSON(w, http.StatusOK, vms)
	}
}

type VmObject interface{}

func (a *App) getHosts(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	locationUUID, ok := params["locationUuid"]

	if !ok {
		respondWithError(w, http.StatusBadRequest, "Missing locationUuid parameter")
		return
	}

	iaas, err := a.getIaasEndpointForLocation(locationUUID[0])

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	hosts, err := iaas.GetHost(a.ctx)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())

	} else {
		respondWithJSON(w, http.StatusOK, hosts)
	}
}

func (a *App) getScsiLunDisks(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	locationUUID, ok := params["locationUuid"]

	if !ok {
		respondWithError(w, http.StatusBadRequest, "Missing locationUuid parameter")
		return
	}

	iaas, err := a.getIaasEndpointForLocation(locationUUID[0])

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	hosts, err := iaas.GetScsiLunDisks(a.ctx)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())

	} else {
		respondWithJSON(w, http.StatusOK, hosts)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}
