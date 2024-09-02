package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CodePayload struct {
	Code      string `json:"code"`
	Directory string `json:"directory"`
	FileName  string `json:"file_name"`
}

var BaseCodeFolder = "codes/"

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/codes", func(w http.ResponseWriter, r *http.Request) {
		uploadHandler(w, r)
	})

	router.HandleFunc("GET /api/v1/codes/{directory_type}/{directory}/{filename}", func(w http.ResponseWriter, r *http.Request) {
		downloadHandler(w, r)
	})

	router.HandleFunc("POST /api/v1/run-test-python", func(w http.ResponseWriter, r *http.Request) {
		runTestPython(w, r)
	})

	router.HandleFunc("POST /api/v1/run-test-golang", func(w http.ResponseWriter, r *http.Request) {
		runTestGolang(w, r)
	})

	defaultMiddlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		CorsMiddleware,
	)

	server := http.Server{
		Addr:    ":8083",
		Handler: defaultMiddlewareChain(router),
	}

	fmt.Println("Server started on port 8083")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error on start server:", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var payload CodePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error on decode JSON: " + err.Error()})
		return
	}

	if strings.Contains(payload.Code, ",") {
		parts := strings.Split(payload.Code, ",")
		if len(parts) == 2 {
			payload.Code = parts[1]
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Base64 in invalid format"})
			return
		}
	}

	// codeData := payload.Code
	// codeData, err := base64.StdEncoding.DecodeString(payload.Code)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": "Error on decode base64: " + err.Error()})
	// 	return
	// }

	if err := os.MkdirAll(BaseCodeFolder+payload.Directory, os.ModePerm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error on create folder: " + err.Error()})
		return
	}

	filePath := filepath.Join(BaseCodeFolder, payload.Directory, payload.FileName)

	if err := os.WriteFile(filePath, []byte(payload.Code), 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error on save code: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Code received successfully"))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	directory := r.PathValue("directory")
	fileName := r.PathValue("filename")
	directoryType := r.PathValue("directory_type")

	if directory == "" || fileName == "" || directoryType == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Parameters 'directory_type', 'directory', and 'filename' are mandatory"})
		return
	}

	filePath := filepath.Join(BaseCodeFolder, directoryType, directory, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "File not found"})
		return
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error on reading file: " + err.Error()})
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fileData)
}

type RequestData struct {
	Text string `json:"text"`
}

type ResponseData struct {
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
}

func runTestPython(w http.ResponseWriter, r *http.Request) {
	var data RequestData

	// Decodifica o JSON recebido
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Salva o código em um arquivo .py
	fileName := "user_code.py"
	err = ioutil.WriteFile(fileName, []byte(data.Text), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	// Executa o arquivo .py usando pytest
	cmd := exec.Command("pytest", fileName)
	output, err := cmd.CombinedOutput()

	var response ResponseData
	if err != nil {
		response = ResponseData{
			Success: false,
			Error:   string(output),
		}
	} else {
		response = ResponseData{
			Success: true,
			Output:  string(output),
		}
	}

	// Remove o arquivo .py após a execução
	os.Remove(fileName)

	// Retorna o resultado como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func runTestGolang(w http.ResponseWriter, r *http.Request) {
	var data RequestData

	// Decodifica o JSON recebido
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Salva o código em um arquivo .go
	fileName := "user_code.go"
	err = ioutil.WriteFile(fileName, []byte(data.Text), 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	// Executa o arquivo .go com o comando `go run`
	cmd := exec.Command("go", "run", fileName)
	output, err := cmd.CombinedOutput()

	var response ResponseData
	if err != nil {
		response = ResponseData{
			Success: false,
			Error:   string(output),
		}
	} else {
		response = ResponseData{
			Success: true,
			Output:  string(output),
		}
	}

	// Remove o arquivo .go após a execução
	os.Remove(fileName)

	// Retorna o resultado como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next.ServeHTTP
	}
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func CorsMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Se for uma requisição OPTIONS, responda com status 200 e termine a execução
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		//	only return json
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Chame o próximo handler
		next.ServeHTTP(w, r)
	})
}
