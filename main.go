package main

import (
    "io/ioutil"
    "encoding/json"
    "github.com/bitly/go-simplejson"
    redis "github.com/go-redis/redis"
    "github.com/apsdehal/go-logger"
    "github.com/gin-gonic/gin"
    "regexp"
    // "strconv"
    "os"
    "flag"
    "fmt"
)


var log *logger.Logger
var options *simplejson.Json


type Info struct {
    Redis_version string `json:"redis_version"`
    Process_id string `json:"process_id"`
    Uptime_in_days string `json:"uptime_in_days"`

    // # Clients
    Connected_clients string `json:"connected_clients"`

    // # Memory
    Used_memory string    `json:"userd_memory"`
    Used_memory_human string    `json:"userd_memory_human"`

    Keyspace_hits string    `json:"keyspace_hits"`
    Keyspace_misses string    `json:"keyspace_misses"`
}





func reg_value(info string, regs string) string {
    reg := regexp.MustCompile(regs)
    val := reg.FindStringSubmatch(info)
    return val[1]
}

func get_info(info string) *Info {
    data := &Info{
        Redis_version: reg_value(info, `redis_version:(\w\S+)`),
        Process_id: reg_value(info, `process_id:(\w+)`),
        Uptime_in_days: reg_value(info, `uptime_in_days:(\w+)`),
        Connected_clients: reg_value(info, `connected_clients:(\w+)`),
        Used_memory: reg_value(info, `used_memory:(\w+)`),
        Used_memory_human: reg_value(info, `used_memory_human:(\w\S+)`),
        Keyspace_hits: reg_value(info, `keyspace_hits:([\d]+)`),
        Keyspace_misses: reg_value(info, `keyspace_misses:([\d]+)`)}
    return data
}


func monitor_redis(addr string, password string) *Info {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password, // no password set
        DB:       0,  // use default DB
    })

    info, err := client.Info().Result()
    if err != nil {
        log.Error(err.Error())
        os.Exit(-1)
    }else{
        log.Info("redis client is ok!")
    }
    client.Close()
    c := get_info(info)
    return c
}

/**
 * [解析配置文件]
 * @type {[type]}
 */

func parse_config(file string){
    bytes, e := ioutil.ReadFile(file)
    if e != nil {
        log.Error(e.Error())
        os.Exit(1)
    }

    opt, err := simplejson.NewJson(bytes)
    if err != nil {
        log.Error(err.Error())
        os.Exit(1)
    }
    options = opt
}


func main(){
    log, _ = logger.New("main", 1, os.Stdout)

    // 解析cli 参数设置
    conf := flag.String("c", "./redis-adm-config.sample", "config file path")
    parse_config(*conf)
    flag.Parse()


    // fmt.Println(*conf)
    // fmt.Println(options.Get("redis").Array())


    // service
    // gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    r.GET("/get/info", func(c *gin.Context){
        cc := monitor_redis("127.0.0.1:6379", "")
        b, _ := json.Marshal(cc)
        callback := c.Query("callback")
        res := fmt.Sprintf("%s(%s)", callback, string(b))
        c.String(200, res)
    })

    r.Run(":5001")
}
