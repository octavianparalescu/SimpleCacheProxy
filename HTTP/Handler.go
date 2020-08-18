package HTTP

import (
	"fmt"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
)

func HandlerFactory(globalRedisClient *redis.Client, source string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := Response{}
		path := r.URL.Path
		fmt.Println("------------------------------")
		fmt.Println("Path is " + path)
		if globalRedisClient != nil {
			pathEncoded := EncodePath(path)

			fmt.Println("Md5 is " + pathEncoded)

			// See if a cache hit
			cacheKey := "page_" + pathEncoded

			// [CACHE] Get
			cacheEntry, errB := globalRedisClient.Get(cacheKey).Result()

			if errB != nil {
				fmt.Println("Not a hit")

				response = downloadPage(source, path, response)

				// [CACHE] Save
				defer globalRedisClient.Set(cacheKey, EncodeResponse(response), 0)
			} else {
				fmt.Println("A hit")

				response = DecodeResponse(cacheEntry)
			}
		} else {
			response = downloadPage(source, path, response)
		}

		// Show page
		sendResponse(w, response)

		fmt.Println("------------------------------")

	}
}

func downloadPage(source string, path string, response Response) Response {
	resp, err := http.Get(source + path)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	body = GetProperBODY(path, body)
	headers := GetProperHeaders(resp.Header)

	response = Response{Headers: headers, Body: body}
	return response
}

func sendResponse(w http.ResponseWriter, response Response) {
	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}
	w.Write(response.Body)
}
