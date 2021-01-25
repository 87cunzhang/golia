package golic

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"zkds/src/confParser"
	"zkds/src/model"
)

var sysParams map[string]string
var secretKey string
var taoBaoUrl = "http://gw.api.taobao.com/router/rest"

func init() {
	conf := confParser.DefaultConf()
	sysParams = make(map[string]string)
	sysParams["app_key"] = conf.String("taobao::appKey")
	sysParams["v"] = conf.String("taobao::apiVersion")
	sysParams["format"] = conf.String("taobao::format")
	sysParams["sign_method"] = conf.String("taobao::signMethod")
	sysParams["partner_id"] = conf.String("taobao::sdkVersion")
	secretKey = conf.String("taobao::secretKey")

}

func logErr(content string, err error) {
	logPath := confParser.DefaultConf().String("errLogPath")
	logData := fmt.Sprintf("%s err: %s, content: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error(), content)
	fileName := logPath + "_" + time.Now().Format("2006-01-02")
	if err := ioutil.WriteFile(fileName, []byte(logData), 0644); err != nil {
		log.Println("write file err:", err)
	}
}

func ExecuteTaobaoRequest(apiParams map[string]string) string {
	datetime := time.Now().Format("2006-01-02 15:04:05")
	sysParams["timestamp"] = datetime
	sign := generateSign(sysParams, apiParams)
	httpUrl := taoBaoUrl + "?"

	for k, v := range sysParams {
		escapeStr := url.QueryEscape(v)
		httpUrl += k + "=" + escapeStr + "&"
	}

	httpUrl += "sign=" + sign
	var postStr string

	for k, v := range apiParams {
		apiStr := url.QueryEscape(v)
		postStr += k + "=" + apiStr + "&"
	}

	resp := post(httpUrl, postStr, "application/x-www-form-urlencoded")
	saveApiLog(resp, apiParams, sysParams)
	return resp
}

//记录淘宝api错误日志
func saveApiLog(resp string, apiParams map[string]string, sysParams map[string]string) {
	errorResponse, _, _, _ := jsonparser.Get([]byte(resp), "error_response")
	if len(errorResponse) > 0 {
		r := new(model.ApiLog)
		r.Api = sysParams["method"]
		r.Code, _ = jsonparser.GetInt(errorResponse, "code")
		r.Msg, _ = jsonparser.GetString(errorResponse, "msg")
		r.Content = resp

		//合并参数,将sysParams合并到apiParams
		for k, v := range sysParams {
			apiParams[k] = v
		}

		paramsJson, _ := json.Marshal(apiParams)
		r.Parameters = string(paramsJson)
		model.AddApiLog(r)
	}
}

//生成淘宝请求签
func generateSign(sysParam map[string]string, apiParam map[string]string) string {
	newMap := make(map[string]string)
	for k, v := range sysParam {
		newMap[k] = v
	}

	for k, v := range apiParam {
		newMap[k] = v
	}

	var newSlice []string
	for k, _ := range newMap {
		newSlice = append(newSlice, k)
	}

	sort.Strings(newSlice)
	var signStr string
	signStr = secretKey
	for _, v := range newSlice {
		signStr += v + newMap[v]
	}

	signStr += secretKey
	return strings.ToUpper(md5V(signStr))
}
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/text
//content:请求返回的内容
func post(url string, postStr string, contentType string) (content string) {
	postStr = strings.TrimRight(postStr, "&")
	resp, err := http.Post(url, contentType, strings.NewReader(postStr))
	if err != nil {
		log.Println("taobao err", err)
		return ""
	}

	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

//TODO::taobao err save
func saveTaobaoErr(data []byte, err error) {
	logPath := confParser.DefaultConf().String("DBerrLogPath")
	logData := fmt.Sprintf("%s dberr: %s,params: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error(), string(data))
	if err := ioutil.WriteFile(logPath, []byte(logData), 0644); err != nil {
		log.Println("write file err:", err)
	}
}
