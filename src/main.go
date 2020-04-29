package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jugalw13/academicdept-go-node/blockchain"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("verify-AcademicDept -student USN \tAcademicDept Verifies Student's data")
}

func verificationByAcademicDept(name string, company string) bool {
	fmt.Println("\nStarting Verification by Academic Dept\n")
	flag := blockchain.AcademicDeptVerification(name, company)
	if flag == true {
		fmt.Println("Verification By Academic Dept Successfully completed!")
		return true
	} else {
		fmt.Println("Verification By Academic Dept Failed!")
		return false
	}
}

func callverificationByAcademicDept(w http.ResponseWriter, r *http.Request) {
	type jsonBody struct {
		Name    string `json:"name"`
		Company string `json:"company"`
	}
	decoder := json.NewDecoder(r.Body)
	var b jsonBody
	if err := decoder.Decode(&b); err != nil {
		log.Fatal(err)
	}

	message := ""
	flag := verificationByAcademicDept(b.Name, b.Company)
	if flag == true {
		message = "Verification By Academic Dept Successfully completed!"
	} else {
		message = "Verification By Academic Dept Failed!"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(message))

	fmt.Println("\n\nSending Notification to Placement Dept for Verification\n\n")
	callPlacementDeptVerification(b.Name, b.Company)
}

func callPlacementDeptVerification(name string, company string) {
	reqBody, err := json.Marshal(map[string]string{
		"name":    name,
		"company": company,
	})
	if err != nil {
		print(err)
	}
	resp, err := http.Post("http://localhost:5004/verify-PlacementDept",
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))
}

func callprintUsage(w http.ResponseWriter, r *http.Request) {

	printUsage()

	w.Header().Set("Content-Type", "application/json")
	message := "Printed Usage!!"
	w.Write([]byte(message))
}

func main() {
	port := "5003"
	http.HandleFunc("/verify-AcademicDept", callverificationByAcademicDept)
	http.HandleFunc("/usage", callprintUsage)
	fmt.Printf("Server listening on localhost:%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}