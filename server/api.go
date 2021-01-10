package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

// API - main instance of application
type API struct {
	Port            string
	PaymentMethods  []string
	DB              []Product
	PaymentLinks    map[string]string
	DownloadAppLink string
	Host            string
}

// New - creates new api instance
func New(port int) *API {
	var err error
	api := new(API)
	portStr := strconv.Itoa(port)
	api.Port = portStr

	api.PaymentMethods = []string{ // payment methods option
		"card",   // valid
		"alipay", // valid
		//"p24", // will cause error (invalid_request_error, currency not supported), uncomment for receiving download link
	}
	api.PaymentLinks = make(map[string]string)
	api.Host, err = os.Hostname()
	if err != nil {
		api.Host = "localhost"
	}
	api.DownloadAppLink = "https://applestore/download?id=1000" // link that we using in case we have error from payment server
	api.fillDatabase()
	return api
}

func (api *API) getPaymentUrls(url string, wg *sync.WaitGroup, errorChan chan error, index int) {
	defer wg.Done()

	stripeKey := os.Getenv("StripeKey") // get secret key from env in production
	stripe.Key = stripeKey
	stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc" // use test value for development

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(api.DB[index].Amount),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			url,
		}),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		errorChan <- err
	}

	api.PaymentLinks[url] = pi.ClientSecret
}

func (api *API) paymentHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	if r.Method != http.MethodGet {
		api.handleError(w, errors.New("API accepts only GET method"), 405, countRequestTime(startTime))
		return
	}
	requestBody := make(map[string]int)
	err := json.NewDecoder(r.Body).Decode(&requestBody) // get product id from request
	if err != nil {
		api.handleError(w, err, http.StatusBadRequest, countRequestTime(startTime))
		return
	}
	index, find := api.findProduct(requestBody["id"])
	if !find {
		api.handleError(w, errors.New("Product with this ID doesn't exist"), http.StatusBadRequest, countRequestTime(startTime))
		return
	}

	errorChan := make(chan error)
	var wg sync.WaitGroup

	for _, url := range api.PaymentMethods {
		wg.Add(1)
		go api.getPaymentUrls(url, &wg, errorChan, index)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		if err != nil {
			downloadLink := make(map[string]string)
			downloadLink["downloadAppLink"] = api.DownloadAppLink
			api.handleSuccess(w, downloadLink, countRequestTime(startTime))
			return
		}
	}

	api.handleSuccess(w, api.PaymentLinks, countRequestTime(startTime))
}

func (api *API) findProduct(id int) (int, bool) {
	for i, product := range api.DB { // grab product from database
		if product.ID == id {
			return i, true
		}
	}
	return -1, false
}

func countRequestTime(startTime time.Time) string {
	elapsed := time.Since(startTime)
	res := strconv.FormatInt(elapsed.Microseconds(), 10)

	return res
}

func (api *API) handleSuccess(w http.ResponseWriter, output interface{}, time string) {
	w.Header().Set("X-Server-Name", api.Host)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Response-Time", time)

	json.NewEncoder(w).Encode(output)
}

func (api *API) handleError(w http.ResponseWriter, err error, statusCode int, time string) {
	w.Header().Set("X-Server-Name", "Payment API Server")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Response-Time", time)
	w.WriteHeader(statusCode)

	finalResult := make(map[string]interface{}, 0)
	finalResult["error"] = err.Error()
	json.NewEncoder(w).Encode(finalResult)
}

// Start - start to serve
func (api *API) Start() {
	mux := http.NewServeMux()

	handler := http.HandlerFunc(api.paymentHandler)
	mux.Handle("/", handler)
	//http.HandleFunc("/", api.paymentHandler)

	fmt.Printf("Starting server at port %s\n", api.Port)
	go func() {
		for {
			log.Fatal(http.ListenAndServe(":"+api.Port, mux))
		}
	}()
}

// Configure http.Client for get request go get payment links

// tr := &http.Transport{
// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// }
// var netClient = &http.Client{
// 	Timeout:   time.Second * 3,
// 	Transport: tr,
// }

// resp, err := netClient.Get(url)
// fmt.Println("Error is ----- ", err)
// if err != nil {
// 	errorChan <- err
// }
// //data, _ := ioutil.ReadAll(resp.Body)
// resp.Body.Close()
