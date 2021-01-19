package golib

import (
	"strconv"
	"time"
	"zkds/src/model"
)

type MiniTmpl struct {
	AppId  string
	ShopId int64
}

//获取实例化记录表名
func GetTmplTableName(TmplId string) string {
	//截取后四位,总共16位,从12位开始往后截取
	if len(TmplId) != 16 {
		return ""
	}
	return "txzj_mini_tmpl_C" + string([]byte(TmplId)[12:])
}

//更新成功
func UpdateAppSuccess(shopId string, status string, appVersion string, templateVersion string, apiResponse string, historyId string, TmplId string) error {
	_, err := model.DB("member").Table(GetTmplTableName(TmplId)).SQL("update " + GetTmplTableName(TmplId) + " set `status` = '" + status + "',`template_version` = '" + templateVersion + "',`app_version` = '" + appVersion + "',`sub_code` = 'success',`api_response` = '" + apiResponse + "',`history_id` = '" + historyId + "',`update_time` = " + strconv.Itoa(int(time.Now().Unix())) + " where shop_id = '" + shopId + "'").Execute()
	return err
}

//更新失败
func UpdateAppFail(shopId string, status string, subCode string, templateVersion string, apiResponse string, historyId string, TmplId string) error {
	_, err := model.DB("member").Table(GetTmplTableName(TmplId)).SQL("update " + GetTmplTableName(TmplId) + " set `status` = '" + status + "',`template_version` = '" + templateVersion + "',`sub_code` = '" + subCode + "',`api_response` = '" + apiResponse + "',`history_id` = '" + historyId + "',`update_time` = " + strconv.Itoa(int(time.Now().Unix())) + " where shop_id = '" + shopId + "'").Execute()
	return err
}

//上线成功
func OnlineAppSuccess(shopId string, status string, appVersion string, apiResponse string, PreviewUrl string, historyId string, TmplId string) error {
	_, err := model.DB("member").Table(GetTmplTableName(TmplId)).SQL("update " + GetTmplTableName(TmplId) + " set `status` = '" + status + "',`pre_view_url` = '" + PreviewUrl + "',`app_version` = '" + appVersion + "',`sub_code` = 'success',`api_response` = '" + apiResponse + "',`history_id` = '" + historyId + "',`update_time` = " + strconv.Itoa(int(time.Now().Unix())) + " where shop_id = '" + shopId + "'").Execute()
	return err
}

//上线失败
func OnlineAppFail(shopId string, status string, subCode string, apiResponse string, historyId string, TmplId string) error {
	_, err := model.DB("member").Table(GetTmplTableName(TmplId)).SQL("update " + GetTmplTableName(TmplId) + " set `status` = '" + status + "',`sub_code` = '" + subCode + "',`api_response` = '" + apiResponse + "',`history_id` = '" + historyId + "',`update_time` = " + strconv.Itoa(int(time.Now().Unix())) + " where shop_id = '" + shopId + "'").Execute()
	return err
}

//获取要更新的店铺,排除过期的店铺
func GetUpdateShops(offset int64, pageSize int64, TmplId string) []MiniTmpl {
	miniTmplColumns := make([]MiniTmpl, 0)
	offsetString := strconv.FormatInt(offset, 10)
	pageSizeString := strconv.FormatInt(pageSize, 10)
	currTimeStampString := strconv.FormatInt(time.Now().Unix(), 10)
	model.DB("member").Table(GetTmplTableName(TmplId)).Join("INNER", "txzj_shop", GetTmplTableName(TmplId)+".shop_id = txzj_shop.shop_id and txzj_shop.deadline > "+currTimeStampString+" limit "+offsetString+","+pageSizeString).Find(&miniTmplColumns)
	return miniTmplColumns
}
