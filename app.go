package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/iftachsc/contracts"

	"github.com/gorilla/mux"
	"github.com/iftachsc/vmware"
	_ "github.com/lib/pq"
	"github.com/vmware/govmomi"
)

type App struct {
	Router     *mux.Router
	VimClients map[string]*govmomi.Client
	ctx        context.Context
}

//Initialize this is great func
func (a *App) Initialize(user, password, dbname string) {

	ctx := context.Background()

	a.ctx = ctx
	a.VimClients = make(map[string]*govmomi.Client)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/vms", a.getVms).Methods("GET")
	//a.Router.HandleFunc("/vm", a.registerVm).Methods("POST")
	//a.Router.HandleFunc("/hosts", a.getHosts).Methods("GET")

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

func (a *App) getClientForLocation(locationUuid string) (*govmomi.Client, error) {
	client, exists := a.VimClients[locationUuid]
	var err error
	if !exists {
		//getLocation from location microservice (*contracts.Location)
		//return location.InitizalizeClient() (vimclient)
		location := contracts.Location{ //instead of location microservice
			IaasProvider: VmwareVcenter{
				Host:     "https://vcenter01.mgmt.il-center-1.cloudzone.io",
				User:     "iftachsc",
				Password: "5tw5j;M]HVN0$",
			},
			StoageFilers: nil,
		}

		//TODO: Error handling here
		client, err := location.IaasProvider.InitializeClient(a.ctx)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		fmt.Println(client)
		a.VimClients[locationUuid] = client
	}

	return client, err

}

func (a *App) getVms(w http.ResponseWriter, r *http.Request) {

	// Returns a url.Values, which is a map[string][]string
	params := r.URL.Query()
	locationUuid, ok := params["locationUuid"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Missing locationUuid parameter")
		return
	}

	c, err := a.getClientForLocation(locationUuid[0])

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	vms, err := vmware.GetVM(c, a.ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())

	} else {
		respondWithJSON(w, http.StatusOK, vms)
	}
}

// func (a *App) getHosts(w http.ResponseWriter, r *http.Request) {

// 	hosts, err := vmware.GetEsxHost(a.VimClient, a.ctx)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())

// 	} else {
// 		respondWithJSON(w, http.StatusOK, hosts)
// 	}
// }

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}
