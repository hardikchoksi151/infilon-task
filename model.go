package main

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Phone struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	PersonID int    `json:"person_id"`
}

type Address struct {
	ID      int    `json:"id"`
	City    string `json:"city"`
	State   string `json:"state"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	ZipCode string `json:"zip_code"`
}

type AddressJoin struct {
	ID        int `json:"id"`
	PersonID  int `json:"person_id"`
	AddressID int `json:"address_id"`
}

type PersonInfo struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code"`
}

type PersonCreate struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	Street1     string `json:"street1" binding:"required"`
	Street2     string `json:"street2"`
	ZipCode     string `json:"zip_code" binding:"required"`
}
