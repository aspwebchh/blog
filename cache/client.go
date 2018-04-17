package cache

import (
	"time"
	"../model"
	"strconv"
    "github.com/go-redis/redis"
	"encoding/json"
	"../common/config_manager"
)

var modelArticle = model.Article{}
var modelArticleCate = model.ArticleCate{}

var client *redis.Client

func InitClientCache()  {
	client = redis.NewClient(&redis.Options{
		Addr:    config_manager.GetConfig().Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetArticleCates() []map[string]string {
	const KEY   =  "articlec_cates"
	var cacheData, err = client.Get(KEY).Result()
	if err != nil {
		var fromDatabase = modelArticleCate.GetEnabledArticleCates()
		var jsonBytes,_ = json.Marshal(fromDatabase)
		var jsonString = string(jsonBytes)
		client.Set(KEY, jsonString,time.Minute * 10)
		return fromDatabase
	} else {
		var data = make([]map[string]string ,0)
		var err = json.Unmarshal([]byte(cacheData), &data)
		if err == nil {
			return data
		} else {
			return nil
		}
	}
}

func GetArticle( articleId int ) map[string]string  {
	const KEY  = "article"
	var sArticleId = strconv.Itoa( articleId )
	var articleCache, err = client.Get( KEY + "-" +sArticleId).Result()
	if err != nil {
		var articleInfo = modelArticle.GetArticle(articleId)
		var jsonBytes ,_ = json.Marshal(articleInfo)
		var jsonString = string(jsonBytes)
		client.Set(KEY +"-"+ sArticleId,jsonString, time.Minute * 10)
		return articleInfo
	} else {
		var data = make(map[string]string)
		var err = json.Unmarshal([]byte(articleCache),&data)
		if err == nil{
			return data
		} else {
			return nil
		}}
}