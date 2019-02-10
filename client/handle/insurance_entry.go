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

func CreateInsuranceEntryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateInsuranceEntryHandler, received request")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Printf("CreateInsuranceEntryHandler, error: %+v\n", err)
		panic(err) // TODO remove panic
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("CreateInsuranceEntryHandler, error: %+v\n", err)
		panic(err) // TODO remove panic
	}
	//var returnStatus int
	var msgResponse MsgResponse
	var insuranceEntryTypeReq model.InsuranceEntryType
	if err := json.Unmarshal(body, &insuranceEntryTypeReq); err != nil {
		log.Printf("CreateInsuranceEntryHandler, error: %+v\n", err)
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
		log.Printf("CreateInsuranceEntryHandler, creating insurance entry: %s", insuranceEntryTypeReq.Id)
		// process request
		responsePayload, err := client.CreateInsuranceEntry(
			insuranceEntryTypeReq,
			"",
		)
		if err != nil {
			log.Printf("CreateInsuranceEntryHandler, error: %+v", err)
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
			log.Printf("CreateInsuranceEntryHandler, success")
			// respond
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(responsePayload))
		}
	}
}

func GetInsuranceEntryHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInsuranceEntryHandler, received request")
	vars := mux.Vars(r)
	id := vars["id"]

	//var returnStatus int
	var msgResponse MsgResponse
	responsePayload, err := client.RetrieveInsuranceEntry(
		id,
		"",
	)
	if err != nil {
		log.Printf("GetInsuranceEntryHandler, error: %+v", err)
		// respond
		msgResponse.Msg = "error"
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(msgResponse); err != nil {
			log.Printf("error: %+v\n", err)
			panic(err)
		}
	} else {
		log.Printf("GetInsuranceEntryHandler, success")
		msgResponse.Msg = responsePayload
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responsePayload))

	}
}

func GetInsuranceEntriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetInsuranceEntriesHandler, received request")

	var msgResponse MsgResponse
	responsePayload, err := client.RetrieveInsuranceEntries(
		"",
	)
	if err != nil {
		log.Printf("GetInsuranceEntriesHandler, error: %+v", err)
		msgResponse.Msg = "error"
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(msgResponse); err != nil {
			log.Printf("error: %+v\n", err)
			panic(err)
		}
	} else {
		log.Printf("GetInsuranceEntriesHandler, success")
		// respond
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responsePayload))
	}
}
