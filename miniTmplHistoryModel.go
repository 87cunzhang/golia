package golib

import (
	"strconv"
	"time"
	"zkds/src/model"
)

const miniTmplHistoryTableName = "txzj_mini_tmpl_history"

type MiniTmplHistory struct {
	Id         int64
	TmplType   int64
	TmplId     string
	OldVersion string
	NewVersion string
	Status     string
	CreateTime int64
	UpdateTime int64
}

//获取最后一条更新记录
func GetLastHistoryRecord() string {
	var newVersion string
	lastRecord, _ := model.DB("member").Table(miniTmplHistoryTableName).QueryValue("select * from " + miniTmplHistoryTableName + " order by id desc limit 1")
	if len(lastRecord) > 0 {
		newVersion := lastRecord[0]["new_version"]
		return string(newVersion)
	} else {
		return newVersion
	}

}

//插入更新版本记录
func AddMiniTmplHistory(r *MiniTmplHistory) int64 {
	r.UpdateTime = time.Now().Unix()
	r.CreateTime = time.Now().Unix()
	model.DB("member").Table(miniTmplHistoryTableName).Insert(r)
	return r.Id
}

//更新批次状态
func UpdateMiniTmplHistoryStatus(Id int64) error {
	_, err := model.DB("member").Table(miniTmplHistoryTableName).SQL("update " + miniTmplHistoryTableName + " set `status` = 1,`update_time` = " + strconv.Itoa(int(time.Now().Unix())) + " where id = '" + strconv.FormatInt(Id, 10) + "'").Execute()
	return err
}
