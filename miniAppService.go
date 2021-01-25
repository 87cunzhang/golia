package golib

import (
	"errors"
	"fmt"
	"github.com/87cunzhang/utils"
	"github.com/buger/jsonparser"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

//更新单个店铺小程序
func MiniappUpdateOnline(jsonMsg []byte) error {
	shopId, _ := jsonparser.GetString(jsonMsg, "shop_id")
	appId, _ := jsonparser.GetString(jsonMsg, "app_id")
	newVersion, _ := jsonparser.GetString(jsonMsg, "new_version")
	historyId, _ := jsonparser.GetString(jsonMsg, "history_id")
	tmplId, _ := jsonparser.GetString(jsonMsg, "tmpl_id")
	if len(tmplId) != 16 {
		utils.LogErr("tmplId invalid", errors.New("tmplId 参数错误"))
		return nil
	}
	accessToken := GetSessionByShopId(shopId)
	clients := "taobao,tmall"
	extJson := "{\"ext\":{\"shopId\":" + shopId + ",\"Version\":\"" + newVersion + "\"},\"extEnable\":true}"
	updateResponse := MiniappTemplateUpdateapp(accessToken, clients, appId, extJson, tmplId, newVersion)
	updateData, _, _, _ := jsonparser.Get([]byte(updateResponse), "miniapp_template_updateapp_response")

	if len(updateData) > 0 {
		//更新成功
		appVersion, _ := jsonparser.GetString(updateData, "app_version")
		status := "210"

		if err := updateAppSuccess(shopId, status, appVersion, newVersion, updateResponse, historyId, tmplId); err != nil {
			return err
		}

		//上线小程序
		if err := onlineApp(shopId, accessToken, clients, appId, appVersion, tmplId, newVersion, historyId); err != nil {
			return err
		}

	} else {
		//更新失败
		status := "310"
		errorResponse, _, _, _ := jsonparser.Get([]byte(updateResponse), "error_response")
		subCode, _ := jsonparser.GetString(errorResponse, "sub_code")

		if err := updateAppFail(shopId, status, subCode, newVersion, updateResponse, historyId, tmplId); err != nil {
			return err
		}

	}

	return nil
}

//更新全部店铺
func MiniappUpdateAll(jsonMsg []byte) error {
	//新版本号
	newVersion, _ := jsonparser.GetString(jsonMsg, "new_version")
	tmplId, _ := jsonparser.GetString(jsonMsg, "tmpl_id")
	if len(tmplId) != 16 {
		utils.LogErr("tmplId invalid", errors.New("tmplId 参数错误"))
		return nil
	}
	//生成一条小程序更新记录
	r := new(MiniTmplHistory)
	//老版本号
	r.OldVersion = GetLastHistoryRecord()
	//模板类型:1=商家端应用(移动) 2=商家端应用(PC)
	r.TmplType = 1
	r.TmplId = tmplId
	//模板ID
	r.NewVersion = newVersion
	//生成一条更新记录
	historyId := AddMiniTmplHistory(r)

	var offset, pageSize, pageNum int64
	pageSize = 20
	pageNum = 1

	//分批查询所有要更新的店铺
	for {
		offset = (pageNum - 1) * pageSize
		//店铺列表
		shopList := GetUpdateShops(offset, pageSize, tmplId)
		if len(shopList) == 0 {
			//将状态更新为已全部更新完成
			err := UpdateMiniTmplHistoryStatus(historyId)
			if err != nil {
				utils.LogErr("状态更新失败", errors.New("historyId:"+strconv.FormatInt(historyId, 10)+", 状态更新失败"))
			}
			//退出循环
			break
		}

		//发送消息到更新单个店铺小程序的队列
		for _, v := range shopList {
			body := "{\"data\":{\"shop_id\":\"" + strconv.FormatInt(v.ShopId, 10) + "\",\"app_id\":\"" + v.AppId + "\",\"new_version\":\"" + newVersion + "\",\"tmpl_id\":\"" + tmplId + "\",\"history_id\":\"" + strconv.FormatInt(historyId, 10) + "\"},\"type\":\"miniapp_update_online\"}"
			conf := DefaultConf()
			miniAppExchangeName := conf.DefaultString("AMQP::miniAppExchangeName", "")
			miniAppQueueBindKey := conf.DefaultString("AMQP::miniAppQueueBindKey", "")
			if len(miniAppExchangeName) == 0 || len(miniAppQueueBindKey) == 0 {
				utils.LogErr("miniAppExchangeName and miniAppQueueBindKey not init", errors.New("未配置实例化队列"))
				return nil
			}
			err := publish(miniAppExchangeName, miniAppQueueBindKey, body, true)
			if err != nil {
				utils.LogErr("消息发送失败", err)
			}
		}

		pageNum++
	}

	return nil
}

//更新小程序
func onlineApp(shopId string, accessToken string, clients string, appId string, appVersion string, TmplId string, newVersion string, historyId string) error {
	//上线小程序
	onlineResponse := MiniappTemplateOnlineapp(accessToken, clients, appId, appVersion, TmplId, newVersion)
	onlineData, _, _, _ := jsonparser.Get([]byte(onlineResponse), "miniapp_template_onlineapp_response")

	if len(onlineData) > 0 {
		//上线成功
		onlineStatus := "220"
		preViewUrl, _ := jsonparser.GetString(onlineData, "app_info", "online_url")
		if err := onlineAppSuccess(shopId, onlineStatus, appVersion, onlineResponse, preViewUrl, historyId, TmplId); err != nil {
			return err
		}
	} else {
		//上线失败
		status := "320"
		errorResponse, _, _, _ := jsonparser.Get([]byte(onlineResponse), "error_response")
		subCode, _ := jsonparser.GetString(errorResponse, "sub_code")
		if err := onlineAppFail(shopId, status, subCode, onlineResponse, historyId, TmplId); err != nil {
			return err
		}
	}

	return nil
}

//小程序更新成功
func updateAppSuccess(shopId string, status string, appVersion string, templateVersion string, apiResponse string, historyId string, TmplId string) error {
	err := UpdateAppSuccess(shopId, status, appVersion, templateVersion, apiResponse, historyId, TmplId)
	return err
}

//小程序更新失败
func updateAppFail(shopId string, status string, subCode string, templateVersion string, apiResponse string, historyId string, TmplId string) error {
	err := UpdateAppFail(shopId, status, subCode, templateVersion, apiResponse, historyId, TmplId)
	return err
}

//小程序上线成功
func onlineAppSuccess(shopId string, status string, appVersion string, apiResponse string, preViewUrl string, historyId string, TmplId string) error {
	err := OnlineAppSuccess(shopId, status, appVersion, apiResponse, preViewUrl, historyId, TmplId)
	return err
}

//小程序上线失败
func onlineAppFail(shopId string, status string, subCode string, apiResponse string, historyId string, TmplId string) error {
	err := OnlineAppFail(shopId, status, subCode, apiResponse, historyId, TmplId)
	return err
}

//发布消息
func publish(exchange, routingKey, body string, reliable bool) error {
	conf := DefaultConf()
	user := conf.DefaultString("AMQP::user", "guest")
	pwd := conf.DefaultString("AMQP::password", "guest")
	addr := conf.DefaultString("AMQP::addr", "localhost")
	port := conf.DefaultInt("AMQP::port", 5672)
	amqpURI := "amqp://" + user + ":" + pwd + "@" + addr + ":" + strconv.Itoa(port) + "/"
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		utils.LogErr("rabbitMq 连接失败", err)
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	if reliable {
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms)
	}

	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

//确认消息
func confirmOne(confirms <-chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}

//获取店铺access_token
func GetSessionByShopId(shopId string) string {
	session := GetAccessTokenByShopId(shopId)

	if session == "" {
		utils.LogErr("shopId:"+shopId+", session empty", errors.New("shopId:"+shopId+", session empty"))
		return ""
	}

	return session
}
