package golic

const shopTableName = "txzj_shop"

//获取最后一条更新记录
func GetAccessTokenByShopId(shopId string) string {
	var accessTokenMain string
	shopInfo, _ := DB("member").Table(shopTableName).QueryValue("select access_token_main from " + shopTableName + " where shop_id = " + shopId)
	if len(shopInfo) > 0 {
		accessTokenMain := shopInfo[0]["access_token_main"]
		return string(accessTokenMain)
	} else {
		return accessTokenMain
	}

}
