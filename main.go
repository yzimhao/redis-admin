package main

import (
    // io "io/ioutil"
    "encoding/json"
    // "github.com/bitly/go-simplejson"
    redis "github.com/go-redis/redis"
    "github.com/apsdehal/go-logger"
    "regexp"
    // "strconv"
    "os"

    "fmt"
)

type Info struct {

    // # Clients
    // connected_clients:1

    // # Memory
    Used_memory string    `json:"userd_memory"`
    Used_memory_human string    `json:"userd_memory_human"`
    // used_memory_rss:3600384
    // used_memory_peak:815960
    // used_memory_peak_human:796.84K
    // used_memory_lua:36864

    // # Persistence
    // loading:0
    // rdb_changes_since_last_save:20
    // rdb_bgsave_in_progress:0
    // rdb_last_save_time:1491979572
    // rdb_last_bgsave_status:ok
    // rdb_last_bgsave_time_sec:-1
    // rdb_current_bgsave_time_sec:-1
    //
    //
    // // # Stats
    //
    // expired_keys:1
    // evicted_keys:0
    Keyspace_hits string    `json:"keyspace_hits"`
    Keyspace_misses string    `json:"keyspace_misses"`
    //
    //
    // // # CPU
    // used_cpu_sys:318.30
    // used_cpu_user:117.02
}



func reg_value(info string, regs string) string {
    reg := regexp.MustCompile(regs)
    val := reg.FindStringSubmatch(info)
    return val[1]
}

func get_info(info string) *Info {
    khit := reg_value(info, `keyspace_hits:([\d]+)`)
    kms := reg_value(info, `keyspace_misses:([\d]+)`)
    um := reg_value(info, `used_memory:(\w+)`)
    umh := reg_value(info, `used_memory_human:(\w\S+)`)

    data := &Info{
        Used_memory: um,
        Used_memory_human: umh,
        Keyspace_hits: khit,
        Keyspace_misses: kms}
    return data
}


func main(){
    log, err := logger.New("my", 1, os.Stdout)

    client := redis.NewClient(&redis.Options{
        Addr:     "127.0.0.1:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })


    info, err := client.Info().Result()
    if err != nil {
        log.Error(err.Error())
        os.Exit(-1)
    }


    c := get_info(info)
    b, _ := json.Marshal(c)
    fmt.Println(c)
    fmt.Println(string(b))
}
