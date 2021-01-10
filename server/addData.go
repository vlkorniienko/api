package server

// Product - products information
type Product struct {
	ID     int
	Name   string
	Amount int64
}

func (api *API) fillDatabase() {
	product1 := Product{ID:1, Name: "phone", Amount: 1000}
	product2 := Product{ID:42, Name: "laptop", Amount: 1500}

	api.DB = append(api.DB, product1, product2)
}