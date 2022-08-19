package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Operation struct {
	Value1 float64 `json:"value1"`
	Value2 float64 `json:"value2"`
	Result float64 `json:"result"`
	Operation string `json:"operation"`
}

type ResponseBody struct {
	
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

const DefaultPort int = 8080

var history []Operation

func createJsonResponseWithData(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ResponseBody{message, data})
}


func createSimpleJsonResponse(w http.ResponseWriter, message string, statusCode int) {
	createJsonResponseWithData(w, message, nil, statusCode)
}

func calculate(value1, value2 float64, operation string) (result float64, err error) {
	switch operation {
	case "sum":
		result = value1 + value2
	case "sub":
		result = value1 - value2
	case "mul":
		result = value1 * value2
	case "div":
		if value2 == 0.0 {
			err = errors.New("Cannot divide by zero")
			return
		}
		result = value1 / value2
	default:
		err = errors.New(operation + " is not a valid operation")
		return
	}
	return
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		createSimpleJsonResponse(w, "Missing path parameters in /calc/{operation}/{value1}/{value2}", http.StatusBadRequest)
		return
	}

	operation := pathParts[2]
	value1, conversionErr1 := strconv.ParseFloat(pathParts[3], 64)
	if conversionErr1 != nil {
		createSimpleJsonResponse(w, conversionErr1.Error(), http.StatusBadRequest)
		return
	}
	value2, conversionErr2 := strconv.ParseFloat(pathParts[4], 64)
	if conversionErr2 != nil {
		createSimpleJsonResponse(w, conversionErr2.Error(), http.StatusBadRequest)
		return
	}

	result, calculatorError := calculate(value1, value2, operation)
	if calculatorError != nil {
		createSimpleJsonResponse(w, calculatorError.Error(), http.StatusBadRequest)
		return
	}

	operationInfo := Operation{value1, value2, result, operation}
	history = append(history, operationInfo)

	createJsonResponseWithData(w, "Operation performed successfully", operationInfo, 200)
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	if len(history) == 0 {
		createSimpleJsonResponse(w, "Empty history", 200)
		return
	}
	createJsonResponseWithData(w, "Operations history", history, 200)
}

func getServerPort() (port int) {
	value := os.Getenv("SERVER_PORT")
	if value == "" {
		return DefaultPort
	}
	portValue, err := strconv.Atoi(value)
	if err != nil {
		return DefaultPort
	}
	return portValue
}

func main() {
	fmt.Println("Server running")
	http.HandleFunc("/calc/", calcHandler)
	http.HandleFunc("/calc/history", historyHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", getServerPort()), nil)
}
