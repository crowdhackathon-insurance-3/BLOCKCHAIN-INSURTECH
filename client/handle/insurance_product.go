package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
	"io"
	"github.com/gorilla/mux"
	"blockchain_at_insurtech/client/client"
	"blockchain_at_insurtech/client/model"
)

func CreateInsuranceProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateInsuranceProductHandler, received request")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Printf("CreateInsuranceProductHandler, error: %+v\n", err)
		panic(err) // TODO remove panic
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("CreateInsuranceProductHandler, error: %+v\n", err)
		panic(err) // TODO remove panic
	}
	//var returnStatus int
	var msgResponse MsgResponse
	var insuranceProductTypeReq model.InsuranceProductType
	if err := json.Unmarshal(body, &insuranceProductTypeReq); err != nil {
		log.Printf("CreateInsuranceProductHandler, error: %+v\n", err)
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		msgResponse.Msg = "error"
		err := json.NewEncoder(w).Encode(msgResponse)
		if err != nil {
			log.Printf("error: %+v\n", err)
			panic(err)
		}
	} else {
		log.Printf("CreateInsuranceProductHandler, creating insurance product: %s", insuranceProductTypeReq.Id)
		// process request
		responsePayload, err := client.CreateInsuranceProduct(
			insuranceProductTypeReq,
			"",
		)
		if err != nil {
			log.Printf("CreateInsuranceProductHandler, error: %+v", err)
			msgResponse.Msg = "error"
			// respond
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotAcceptable)
			err := json.NewEncoder(w).Encode(msgResponse)
			if err != nil {
				log.Printf("error: %+v\n", err)
				panic(err)
			}
		} else {
			log.Printf("CreateInsuranceProductHandler, success")
			// respond
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(responsePayload))
		}
	}
}

func GetInsuranceProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInsuranceProductHandler, received request")
	vars := mux.Vars(r)
	id := vars["id"]

	var msgResponse MsgResponse
	responsePayload, err := client.RetrieveInsuranceProduct(
		id,
		"",
	)
	if err != nil {
		log.Printf("GetInsuranceProductHandler, error: %+v", err)
		// respond
		msgResponse.Msg = "error"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(msgResponse); err != nil {
			log.Printf("error: %+v\n", err)
			panic(err)
		}
	} else {
		log.Printf("GetInsuranceProductHandler, success")
		msgResponse.Msg = responsePayload
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responsePayload))
	}
}
func GetInsuranceProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInsuranceProductsHandler, received request")

	var msgResponse MsgResponse
	responsePayload, err := client.RetrieveInsuranceProducts(
		"",
	)
	if err != nil {
		log.Printf("GetInsuranceProductsHandler, error: %+v", err)
		msgResponse.Msg = "error"
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(msgResponse); err != nil {
			log.Printf("error: %+v\n", err)
			panic(err)
		}
	} else {
		log.Printf("GetInsuranceProductsHandler, success")
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responsePayload))
	}
}
