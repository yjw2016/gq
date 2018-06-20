package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"time"
	etcd_client "github.com/coreos/etcd/clientv3"
	"encoding/json"
	"fmt"
	"context"
)
var (
	redisPool *redis.Pool
	etcdClient *etcd_client.Client
)
func initRedis() (err error) {
	redisPool = &redis.Pool{
		// Dial is an application supplied function for creating and configuring a
		// connection.
		//
		// The connection returned from Dial must not be in a special state
		// (subscribed to pubsub channel, transaction started, ...).
		Dial: func() (redis.Conn, error) {
			logs.Debug("redis adrr is ",secKillConf.redisConf.redisAddr)
			return redis.Dial("tcp", secKillConf.redisConf.redisAddr)
		},

		// Maximum number of idle connections in the pool.
		MaxIdle : secKillConf.redisConf.redisMaxIdle,

		// Maximum number of connections allocated by the pool at a given time.
		// When zero, there is no limit on the number of connections in the pool.
		MaxActive :secKillConf.redisConf.redisMaxActive,

		// Close connections after remaining idle for this duration. If the value
		// is zero, then idle connections are not closed. Applications should set
		// the timeout to a value less than the server's timeout.
		IdleTimeout : time.Duration(secKillConf.redisConf.redisIdleTimeout)*time.Second,
	}

	conn := redisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis error")
	}
	return
}

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{secKillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(secKillConf.EtcdConf.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdClient = cli
	return
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = secKillConf.logConf.logPath
	config["level"] = convertLogLevel(secKillConf.logConf.logLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
func loadSecConf() (err error) {

	logs.Debug("secKillConf.EtcdConf.EtcdSecProductKey is :",secKillConf.EtcdConf.EtcdSecProductKey)
	resp, err := etcdClient.Get(context.Background(), secKillConf.EtcdConf.EtcdSecProductKey)
	if err != nil {
		logs.Error("get [%s] from etcd failed, err:%v", secKillConf.EtcdConf.EtcdSecProductKey, err)
		return
	}

	var secProductInfo []SecProductInfoConf
	for k,v := range resp.Kvs{
		logs.Debug("key[%v] value[%v]",k,v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("Unmarsha sec product info failed ,err:%v",err)
			return
		}

		logs.Debug("sec info conf is [%v]",secProductInfo)

	}

	secKillConf.SecProductInfoConf = secProductInfo
	//var secProductInfo []service.SecProductInfoConf
	//for k, v := range resp.Kvs {
	//	logs.Debug("key[%v] valud[%v]", k, v)
	//	err = json.Unmarshal(v.Value, &secProductInfo)
	//	if err != nil {
	//		logs.Error("Unmarshal sec product info failed, err:%v", err)
	//		return
	//	}
	//
	//	logs.Debug("sec info conf is [%v]", secProductInfo)
	//}
	//
	//updateSecProductInfo(secProductInfo)
	return
}
func initSec() (err error) {

	err = initLogger()
	if err != nil {
		logs.Error("init logger failed,err:%v",err)
		return
	}
	err = initRedis()
	if err != nil {
		logs.Error("init redis failed,err:%v",err)
		return
	}

	err = initEtcd()
	if err != nil {
		logs.Error("init etcd failed,err:%v",err)
		return
	}

	err = loadSecConf()
	if err != nil {
		return
	}

	logs.Info("init sec succ")
	return
}
