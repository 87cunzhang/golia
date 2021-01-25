package golib

//更新实例化应用
func MiniappTemplateUpdateapp(session string, cients string, appId string, extJson string, templateId string, templateVersion string) string {
	sysParams["session"] = session
	sysParams["method"] = "taobao.miniapp.template.updateapp"
	apiParams := make(map[string]string)
	//要更新的投放端,目前可投放： taobao(淘宝),tmall(天猫)
	apiParams["clients"] = cients
	//应用id，如果应用在对应端上已有的最新版本所使用的模板id,模板version和extjson，与本次入参一致，则认为不需要更新，返回已有的版本。
	apiParams["app_id"] = appId
	//扩展信息。线上版本使用的template_id与传入的template_id一致时，可不填。
	apiParams["ext_json"] = extJson
	//模板id
	apiParams["template_id"] = templateId
	//模板版本
	apiParams["template_version"] = templateVersion
	return ExecuteTaobaoRequest(apiParams)
}

func MiniappTemplateOnlineapp(session string, cients string, appId string, appVersion string, templateId string, templateVersion string) string {
	sysParams["session"] = session
	sysParams["method"] = "taobao.miniapp.template.onlineapp"
	apiParams := make(map[string]string)
	//要更新的投放端,目前可投放： taobao(淘宝),tmall(天猫)
	apiParams["clients"] = cients
	//应用id，如果应用在对应端上已有的最新版本所使用的模板id,模板version和extjson，与本次入参一致，则认为不需要更新，返回已有的版本。
	apiParams["app_id"] = appId
	//待上线的预览版本号
	apiParams["app_version"] = appVersion
	//模板id
	apiParams["template_id"] = templateId
	//模板版本
	apiParams["template_version"] = templateVersion
	return ExecuteTaobaoRequest(apiParams)
}
