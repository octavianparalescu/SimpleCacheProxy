package main

import (
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/OctavianParalescu/SimpleCacheProxy/Config"
	"github.com/OctavianParalescu/SimpleCacheProxy/HTTP"
	"gopkg.in/redis.v5"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ProgramConfig struct {
	Source              string `yaml:"source"`
	DoExposeSecure      bool   `yaml:"doExposeSecure"`
	CertificateLocation string `yaml:"certificateLocation"`
	ExposedPort         int    `yaml:"exposedPort"`
	DoExposeGzip        bool   `yaml:"doExposeGzip"`
	DoCache             bool   `yaml:"doCache"`
	CacheType           string `yaml:"cacheType"`
	RedisCache          struct {
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db,omitempty"`
	} `yaml:"redisCache"`
}

var globalRedisClient *redis.Client

func main() {
	dat, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))
	programCfg := ProgramConfig{}
	err = yaml.Unmarshal(dat, &programCfg)
	if err != nil {
		panic("config.yaml cannot be parsed: " + err.Error())
	}
	fmt.Println(programCfg)

	if programCfg.DoCache {
		if programCfg.CacheType == "redis" {
			if programCfg.RedisCache.Hostname == "" {
				programCfg.RedisCache.Hostname = "localhost"
			}
			if programCfg.RedisCache.Port == 0 {
				programCfg.RedisCache.Port = 6379
			}
			addr := programCfg.RedisCache.Hostname + ":" + strconv.Itoa(programCfg.RedisCache.Port)
			globalRedisClient = Config.RedisConnect(&redis.Options{
				Addr:     addr,
				Password: programCfg.RedisCache.Password, // no password set
				DB:       programCfg.RedisCache.DB,       // use default DB
			})
		}
	}

	var withGz http.Handler
	if programCfg.DoExposeGzip {
		withGz = gziphandler.GzipHandler(http.HandlerFunc(HTTP.HandlerFactory(globalRedisClient, programCfg.Source)))
	}

	fmt.Println("Starting the handler")
	if programCfg.DoExposeGzip {
		http.Handle("/", withGz)
	} else {
		http.Handle("/", http.HandlerFunc(HTTP.HandlerFactory(globalRedisClient, programCfg.Source)))
	}
	if programCfg.ExposedPort == 0 {
		programCfg.ExposedPort = 80
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(programCfg.ExposedPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
