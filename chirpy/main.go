package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
}

//  --------- chirp validation -----------

// shape of the incoming request body
type chirpRequest struct {
	Body string `json:"body"`
}

//shape of an error response
type errorResponse struct {
	Error string `json:"error"`
}

// shape of a success response
// type validResponse struct {
// 	Valid bool `json:"valid"`
// }
//clean response type
type rs struct {
	CleanedBody string `json:"cleaned_body"`
}
// helper: writes a Json error response with the given status code
func respondWithError(w http.ResponseWriter,code int, msg string) {
	respondWithJSON(w, code, errorResponse{Error: msg})
}

//helper: writes any JSON payload with the given status code
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

// 1. accepts a POST request with a JSON body {"body": "..."}
// 2. decodes the request body into a chirpRequest struct
// 3. validates the length is <= 140 characters
// 4. responds with {"valid": true} on success, or {"error": "..."} on failure
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	//decode
	decoder := json.NewDecoder(r.Body)
	reqBody := chirpRequest{}

	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("error decoding request body: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(reqBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleaned := cleanBody(reqBody.Body)
	respondWithJSON(w, 200, rs{CleanedBody:cleaned})
}

// 1. split the whole sentence into individual words (by spaces)
// 2. go through each word one at a time
// 3. for each word, check if its LOWERCASE version matches one of the 3 banned words
// 4. if it matches, replace that word with ****
//    if it doesn't, leave it untouched (keep original casing!)
// 5. join all the words back together with spaces

func cleanBody(body string) string{
	//split words by space
	splitWords := strings.Split(body, " ")
	fmt.Println(splitWords)

	badWordSlice := []string{"kerfuffle","sharbert","fornax"}
	//loop to compare words 
	for i, word := range splitWords{
		for _, badWord := range badWordSlice{
			//compare words
			if strings.ToLower(word) == badWord {
				splitWords[i] = "****"
			}
		}
	}
	return strings.Join(splitWords, " ")
}
func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir("."))),
	))

	mux.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/validate_chirp/", handlerValidateChirp)
	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}