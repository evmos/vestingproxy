package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	var err error
	port := int64(8080)
	if len(os.Args) == 2 {
		port, err = strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatal("invalid port")
		}
	}

	targetURL, err := url.Parse("https://proxy.evmos.org/cosmos")
	if err != nil {
		log.Fatal("Error parsing target URL:", err)
	}
	cosmosProxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.RequestURI()
		if strings.Contains(url, "/v2/vesting/") {
			address := strings.Replace(url, "/v2/vesting/", "", 1)
			fmt.Println(address)
			getURL := fmt.Sprintf("https://proxy.evmos.org/evmos/vesting/v2/balances/%s", address)
			resp, err := http.Get(getURL)
			if err != nil {
				http.Error(w, "Failed to fetch from server", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			// Copy the response from the special URL to the client
			w.WriteHeader(resp.StatusCode)
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				log.Println("Error copying response body:", err)
			}
			return
		}

		cosmosProxy.ServeHTTP(w, r)
	},
	)

	// start the server on port 8080
	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		// log the error if server fails to start
		fmt.Println("Error starting server:", err)
	}
}
