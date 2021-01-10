

#### Simple test api for collection payment information with product ID

##### User story
User chooses one of the avialable membership plan on onboarding page. After that, the user is shown a checkout form with the ability to pay with credit card, PayPal, Apple Pay, Google Pay. If we cannot fully prepare the checkout form (there is an error with at least one payment provider), then we show links to the application in the store.

The goal is to write an API method that accepts the product ID, "goes" to payment gateways (PayPal, Apple Pay, Google Pay, Stripe ...) via https and receives the urls of the payment buttons from them.

Technical requirements:
*  The application should not be written using any framework,
  however different packages can be installed
*  The result of the request must be presented in JSON format
*  The query result must contain the X-Response-Time and X-Server headers
Name. X-Response-Time - time spent processing the request in
microseconds, X-Server-Name is the server name (hostname)

##### Usage:
```
git clone 
go build -o api cmd/main.go
./server
```
##### Request example:
curl -X GET http://localhost:8090/ -d '{"id":42}' -v
