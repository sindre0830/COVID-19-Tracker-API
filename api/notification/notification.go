package notification

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Initialize signature (via init())
var SignatureKey = "X-SIGNATURE"
//var Mac hash.Hash
var Secret []byte

var webhooks []Notification

type Notification struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

func MethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodPost:
			// Expects incoming body in terms of WebhookRegistration struct
			webhook := Notification{}
			err := json.NewDecoder(r.Body).Decode(&webhook)
			if err != nil {
				http.Error(w, "Something went wrong: " + err.Error(), http.StatusBadRequest)
			}
			webhooks = append(webhooks, webhook)
			// Note: Approach does not guarantee persistence or permanence of resource id (for CRUD)
			//fmt.Fprintln(w, len(webhooks)-1)
			fmt.Println("Webhook " + webhook.URL + " has been registered.")
			http.Error(w, strconv.Itoa(len(webhooks)-1), http.StatusCreated)
		case http.MethodGet:
			// For now just return all webhooks, don't respond to specific resource requests
			err := json.NewEncoder(w).Encode(webhooks)
			if err != nil {
				http.Error(w, "Something went wrong: " + err.Error(), http.StatusInternalServerError)
			}
		case http.MethodDelete:
			//
		default: http.Error(w, "Invalid method " + r.Method, http.StatusBadRequest)
	}
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("Received POST request...")
		for _, v := range webhooks {
			go CallUrl(v.URL, "Trigger event: Call to service endpoint with method " + v.Trigger)
		}
	default:
		http.Error(w, "Invalid method "+r.Method, http.StatusBadRequest)
	}
}

func CallUrl(url string, content string) {
	fmt.Println("Attempting invocation of url " + url + " ...")
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(content)))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Hash content
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(content))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Convert to string & add to header
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in HTTP request: " + err.Error())
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Something is wrong with invocation response: " + err.Error())
	}

	fmt.Println("Webhook invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}
