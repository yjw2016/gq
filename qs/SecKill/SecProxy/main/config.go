package main

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strings"
)

var (
	secKillConf = &SecSkillConf{}
)
type RedisConf struct {
	redisAddr string
	redisMaxIdle int
	redisMaxActive int
	redisIdleTimeout int
}
type EtcdConf struct {
	EtcdAddr string
	Timeout int
	EtcdSecKey string
	EtcdSecProductKey string
	EtcdSecKeyPrefix string
}
type logConf struct {
	logLevel string
	logPath string
}
type SecSkillConf struct {
	redisConf RedisConf
	EtcdConf EtcdConf
	logConf logConf
	SecProductInfoConf []SecProductInfoConf

}
type SecProductInfoConf struct {
	ProductId int
	StartTime int
	EndTime int
	Status int
	Total int
	Left int
}
func initConfig() (err error)  {

	err = redisConfig()
	if err != nil {
		err = fmt.Errorf("redis config failed,err:%v",err)
	}
	err = etcdConfig()
	if err != nil {
		err = fmt.Errorf("etcd config failed,err:%v",err)
	}

	err = logConfig()
	if err != nil {
		err = fmt.Errorf("log config failed,err:%v",err)
	}
	return
}



func redisConfig() (err error) {
	redisAddr := beego.AppConfig.String("redis_addr")
	redisMaxIdle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_idle error:%v", err)
		return
	}

	redisMaxActive, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_active error:%v", err)
		return
	}

	redisIdleTimeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_idle_timeout error:%v", err)
		return
	}


	secKillConf.redisConf.redisMaxIdle = redisMaxIdle
	secKillConf.redisConf.redisMaxActive = redisMaxActive
	secKillConf.redisConf.redisIdleTimeout = redisIdleTimeout
	if len(redisAddr)==0{
		err = fmt.Errorf("init config failed,redis[%s]",redisAddr)
		return
	}
	secKillConf.redisConf.redisAddr = redisAddr
	logs.Debug("read config succ,redis:%v",redisAddr)
	return
}

func etcdConfig() (err error) {
	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	etcdAddr := beego.AppConfig.String("etcd_addr")

	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}
	secKillConf.EtcdConf.Timeout = etcdTimeout
	secKillConf.EtcdConf.EtcdAddr = etcdAddr
	if len(etcdAddr)==0 {
		err = fmt.Errorf("init config failed,etcd[%s] config is null",etcdAddr)
		return
	}
	if err != nil {
		err = fmt.Errorf("init config failed, read etcd_timeout error:%v", err)
		return
	}

	secKillConf.EtcdConf.Timeout = etcdTimeout
	secKillConf.EtcdConf.EtcdSecKeyPrefix = beego.AppConfig.String("etcd_sec_key_prefix")
	if len(secKillConf.EtcdConf.EtcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("init config failed, read etcd_sec_key error:%v", err)
		return
	}

	productKey := beego.AppConfig.String("etcd_product_key")
	if len(productKey) == 0 {
		err = fmt.Errorf("init config failed, read etcd_product_key error:%v", err)
		return
	}

	if strings.HasSuffix(secKillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		secKillConf.EtcdConf.EtcdSecKeyPrefix = secKillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}

	secKillConf.EtcdConf.EtcdSecProductKey = fmt.Sprintf("%s%s", secKillConf.EtcdConf.EtcdSecKeyPrefix, productKey)

	logs.Debug("read config succ,etcd:%v",etcdAddr)

	return
}

func logConfig() (err error) {
	logPath := beego.AppConfig.String("log_path")
	logLevel := beego.AppConfig.String("log_level")

	secKillConf.logConf.logLevel=logLevel
	secKillConf.logConf.logPath = logPath
	return

}

