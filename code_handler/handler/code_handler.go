package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	codeHandlerKey         = os.Getenv("CODE_HANDLER_KEY")
	codeHandlerRemoveFiles = bool(os.Getenv("CODE_HANDLER_REMOVE_FILES") == "true")
)

type CodePayload struct {
	Code      string `json:"code"`
	Directory string `json:"directory"`
	FileName  string `json:"file_name"`
}

type RequestData struct {
	Text string `json:"text"`
}

type ResponseData struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

type UnitTestHandler struct {
	SourceFromUser  string `json:"source_from_user"`
	SourceUnitTeste string `json:"source_unit_teste"`
	Lista           string `json:"lista"`
	Exercicio       string `json:"exercicio"`
	Usuario         string `json:"usuario"`
}

type UnitTestResponse struct {
	OutPutFromUser  ResponseData `json:"out_put_from_user"`
	OutPutUnitTeste ResponseData `json:"out_put_unit_teste"`
}

var BaseCodeFolder = "/codes/"

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/run-test-python", func(w http.ResponseWriter, r *http.Request) {
		runTestPython(w, r)
	})

	router.HandleFunc("POST /api/v1/run-test-golang", func(w http.ResponseWriter, r *http.Request) {
		runTestGolang(w, r)
	})

	defaultMiddlewareChain := MiddlewareChain(
		CheckCodeHandlerKeyMiddleware,
		RequestLoggerMiddleware,
		CorsMiddleware,
	)

	server := http.Server{
		Addr:    ":8082",
		Handler: defaultMiddlewareChain(router),
	}

	log.Println("Server started on port 8082")
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error on start server:", err)
	}
}

func decodeJsonAndMakeDir(r *http.Request) (UnitTestHandler, string, error) {
	var data UnitTestHandler

	// Decodifica o JSON recebido
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		return data, "", errors.New("Requisição com corpo inválido.")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Erro ao pegar o caminho de execução atual:", err)
		return data, "", errors.New("Erro ao pegar o caminho de execução atual")
	}

	folderHandler := filepath.Join(currentDir, BaseCodeFolder, data.Lista, data.Exercicio, data.Usuario)
	if err := os.MkdirAll(folderHandler, os.ModePerm); err != nil {
		return data, "", errors.New("Error ao criar o diretório: " + err.Error())
	}

	return data, folderHandler, nil
}

func decodeBase64AndQueryAndWriteFile(sourceFromUser string, filePathUser string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(sourceFromUser)
	if err != nil {
		return errors.New("erro ao decodificar Base64. File: " + filePathUser)
	}

	sourceFromUserBytes, err := url.QueryUnescape(string(decodedBytes))
	if err != nil {
		return errors.New("erro ao decodificar URI component. File: " + filePathUser)
	}

	err = os.WriteFile(filePathUser, []byte(sourceFromUserBytes), 0644)
	if err != nil {
		return errors.New("Falha ao escrever o arquivo na criação do código do usuário. File: " + filePathUser)
	}
	return nil
}

func runTestPython(w http.ResponseWriter, r *http.Request) {
	data, folderHandler, err := decodeJsonAndMakeDir(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Salva o código em um arquivo .py
	fileNameUser := "user_code.py"
	filePathUser := filepath.Join(folderHandler, fileNameUser)

	err = decodeBase64AndQueryAndWriteFile(data.SourceFromUser, filePathUser)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	cmd := exec.Command("python", filePathUser)
	output, err := cmd.CombinedOutput()

	var response UnitTestResponse
	if err != nil {
		response.OutPutFromUser.Success = false
		response.OutPutFromUser.Error = string(output) + " " + err.Error()
		response.OutPutFromUser.Output = ""
	} else {
		response.OutPutFromUser.Success = true
		response.OutPutFromUser.Error = ""
		response.OutPutFromUser.Output = string(output)
	}

	// syntax of file it's ok, now let's start a unit test
	if err == nil {
		// Salva o código em um arquivo .py
		fileNameUnitTest := "user_code_test.py"
		filePathUnitTest := filepath.Join(folderHandler, fileNameUnitTest)

		err = decodeBase64AndQueryAndWriteFile(data.SourceUnitTeste, filePathUnitTest)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Executa o arquivo .py com o comando `python -m pytest`
		cmd := exec.Command("python", "-m", "pytest", filePathUnitTest)
		cmd.Dir = folderHandler
		output, err := cmd.CombinedOutput()

		// response ResponseData
		if err != nil {
			log.Println("After python test run: ", err)
			response.OutPutUnitTeste.Success = false
			response.OutPutUnitTeste.Error = string(output) + " " + err.Error()
			response.OutPutUnitTeste.Output = ""
		} else {
			response.OutPutUnitTeste.Success = true
			response.OutPutUnitTeste.Error = ""
			response.OutPutUnitTeste.Output = string(output)
		}

		if codeHandlerRemoveFiles {
			// Remove o arquivo .py após a execução
			os.Remove(fileNameUnitTest)
		}
	}

	if codeHandlerRemoveFiles {
		// Remove o arquivo .go após a execução
		os.Remove(filePathUser)
	}

	if !response.OutPutFromUser.Success || !response.OutPutUnitTeste.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Retorna o resultado como JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func runTestGolang(w http.ResponseWriter, r *http.Request) {
	data, folderHandler, err := decodeJsonAndMakeDir(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	fileNameGoMod := "go.mod"
	filePathGoMod := filepath.Join(folderHandler, fileNameGoMod)
	err = os.WriteFile(filePathGoMod, []byte(`module tcc_ead/codes

go 1.23.4
	`), 0644)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Falha ao escrever o arquivo na criação do go.mod"})
		return
	}

	// Salva o código em um arquivo .go
	fileNameUser := "code.go"
	filePathUser := filepath.Join(folderHandler, fileNameUser)

	err = decodeBase64AndQueryAndWriteFile(data.SourceFromUser, filePathUser)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Executa o arquivo .go com o comando `go run`
	cmd := exec.Command("go", "run", filePathUser)

	output, err := cmd.CombinedOutput()

	var response UnitTestResponse
	if err != nil {
		response.OutPutFromUser.Success = false
		response.OutPutFromUser.Error = string(output) + " " + err.Error()
		response.OutPutFromUser.Output = ""
	} else {
		response.OutPutFromUser.Success = true
		response.OutPutFromUser.Error = ""
		response.OutPutFromUser.Output = string(output)
	}

	// syntax of file it's ok, now let's start a unit test
	if err == nil {
		// Salva o código em um arquivo .go
		fileNameUnitTest := "code_test.go"
		filePathUnitTest := filepath.Join(folderHandler, fileNameUnitTest)

		err = decodeBase64AndQueryAndWriteFile(data.SourceUnitTeste, filePathUnitTest)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Executa o arquivo .go com o comando `go test`
		// cmd := exec.Command("go", "test", filePathUnitTest)
		cmd := exec.Command("go", "test", "-v")
		cmd.Dir = folderHandler
		output, err := cmd.CombinedOutput()

		// response ResponseData
		if err != nil {
			log.Println("After go test run: ", err)
			response.OutPutUnitTeste.Success = false
			response.OutPutUnitTeste.Error = string(output) + " " + err.Error()
			response.OutPutUnitTeste.Output = ""
		} else {
			response.OutPutUnitTeste.Success = true
			response.OutPutUnitTeste.Error = ""
			response.OutPutUnitTeste.Output = string(output)
		}

		if codeHandlerRemoveFiles {
			// Remove o arquivo .go após a execução
			os.Remove(fileNameUnitTest)
		}
	}

	if codeHandlerRemoveFiles {
		// Remove o arquivo .go após a execução
		os.Remove(filePathUser)

		// Remove o arquivo .mod
		os.Remove(filePathGoMod)
	}

	if !response.OutPutFromUser.Success || !response.OutPutUnitTeste.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Retorna o resultado como JSON
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

func CheckCodeHandlerKeyMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Auth-Token")
		if apiKey != codeHandlerKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Access unauthorized to code handler"})
			return
		}
		next.ServeHTTP(w, r)
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
