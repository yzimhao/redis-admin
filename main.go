package main

import (
    "crypto/md5"
    "encoding/hex"
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
    Addr string                         `json:"addr"`
    Id string                           `json:"id"`
    Redis_version string                `json:"redis_version"`
    Os string                           `json:"os"`
    Process_id string                   `json:"process_id"`
    Tcp_port string                     `json:"tcp_port"`
    Uptime_in_seconds string            `json:"uptime_in_seconds"`
    Uptime_in_days string               `json:"uptime_in_days"`

    Config_file string                  `json:"config_file"`
    // # Clients
    Connected_clients string            `json:"connected_clients"`

    // # Memory
    Used_memory string                  `json:"used_memory"`
    Used_memory_human string            `json:"used_memory_human"`

    Keyspace_hits string                `json:"keyspace_hits"`
    Keyspace_misses string              `json:"keyspace_misses"`
    Expired_keys string                 `json:"expired_keys"`

    Used_cpu_sys string                 `json:"used_cpu_sys"`
    Used_cpu_user string                `json:"used_cpu_user"`
}





func reg_value(info string, regs string) string {
    reg := regexp.MustCompile(regs)
    val := reg.FindStringSubmatch(info)
    return val[1]
}

func get_info(addr, info string) *Info {
    hexid := md5.Sum([]byte(addr))
    id := hex.EncodeToString(hexid[:])

    data := &Info{
        Addr: addr,
        Id: id[8:16],
        Redis_version: reg_value(info, `redis_version:(\w\S+)`),
        Os: reg_value(info, `os:(\w\S+)`),
        Process_id: reg_value(info, `process_id:(\w+)`),
        Tcp_port: reg_value(info, `tcp_port:(\w+)`),

        Uptime_in_seconds: reg_value(info, `uptime_in_seconds:(\w+)`),
        Uptime_in_days: reg_value(info, `uptime_in_days:(\w+)`),
        // Config_file: reg_value(info, `config_file:(\w+)`),
        Connected_clients: reg_value(info, `connected_clients:(\w+)`),
        Used_memory: reg_value(info, `used_memory:(\w+)`),
        Used_memory_human: reg_value(info, `used_memory_human:(\w\S+)`),
        Expired_keys: reg_value(info, `expired_keys:([\d]+)`),

        
        Used_cpu_sys: reg_value(info, `used_cpu_sys:([\d]+)`),
        Used_cpu_user: reg_value(info, `used_cpu_user:([\d]+)`),

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

    c := get_info(addr, info)
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

    // service
    // gin.SetMode(gin.ReleaseMode)

    // gin.Dir("/data/wwwroot/redis-admin", true)
    r := gin.Default()
    r.Static("/webui", "./webui/")
    r.GET("/get/info", func(c *gin.Context){
        addr := c.Query("addr")
        password :=  "" //c.Query("password")
        callback := c.Query("callback")


        cc := monitor_redis(addr, password)
        b, _ := json.Marshal(cc)

        res := fmt.Sprintf("%s(%s)", callback, string(b))
        c.String(200, res)
    })

    r.Run(":5001")
}
