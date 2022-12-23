package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	Db      *gorm.DB
	MyRedis *redis.Client
)

func InitConfig() {
	cmd := exec.Command("swag", "init")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	_ = cmd.Start()
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("config app initialized.")
	//fmt.Println("config mysql", viper.Get("mysql"))
}

func InitMysql() {
	var newLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second, // 慢sql的阈值
		LogLevel:      logger.Info,
		Colorful:      true,
	})
	Db, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
	fmt.Println("config mysql initialized.")
}

func InitRedis() {
	MyRedis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
	//var pong, err = MyRedis.Ping().Result()
	//if err != nil {
	//	fmt.Println("initialize redis err: ", err)
	//} else {
	//	fmt.Println("redis initialized: ", pong)
	//}
}

const (
	PublishKey = "websocket"
)

// 消息发布
func Publish(ctx context.Context, channel string, msg string) error {
	fmt.Println("消息订阅", msg)
	var err error
	err = MyRedis.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// 消息订阅
func Subscribe(ctx context.Context, channel string) (string, error) {
	var sub = MyRedis.Subscribe(ctx, channel)
	var msg, err = sub.ReceiveMessage(ctx)
	fmt.Println("消息发布", msg.Payload)
	return msg.Payload, err
}
