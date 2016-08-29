package main

import (
	"net/http"
	"gopkg.in/redis.v4"
	"io/ioutil"
	"github.com/OctavianParalescu/SimpleCacheProxy/Tools"
	"github.com/NYTimes/gziphandler"
)

func redisConnect() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0, // use default DB
	})

	return redisClient
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathEncoded := Tools.EncodePath(path)

	//fmt.Println("------------------------------")
	//fmt.Println("Path is " + path)
	//fmt.Println("Md5 is " + pathEncoded)

	// See if a cache hit
	cacheKey := "page_" + pathEncoded

	// [CACHE] Get
	cacheEntry, errB := globalRedisClient.Get(cacheKey).Result()

	response := Tools.HTTPResponse{}
	if (errB != nil) {
		//fmt.Println("Not a hit");

		resp, err := http.Get("http://develop.drewberry.co.uk" + path)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		body = Tools.GetProperBODY(path, body)
		headers := Tools.GetProperHeaders(resp.Header)

		response = Tools.HTTPResponse{Headers: headers, Body: body}

		// [CACHE] Save
		defer globalRedisClient.Set(cacheKey, Tools.EncodeResponse(response), 0)
	} else {
		//fmt.Println("A hit");

		response = Tools.DecodeResponse(cacheEntry)
	}

	// Show page
	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}
	w.Write(response.Body)

	//fmt.Println("------------------------------")

}

var globalRedisClient *redis.Client

func main() {
	globalRedisClient = redisConnect()

	withGz := gziphandler.GzipHandler(http.HandlerFunc(handler))

	http.Handle("/", withGz)
	http.ListenAndServe(":8040", nil)
}