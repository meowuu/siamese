package leancloud

import (
	"errors"
	"strings"
	"fmt"
	// "github.com/dghubble/sling"
	// "net/http"
	"github.com/kylelemons/go-gypsy/yaml"
)

type Client struct {
	id string
	key string
	token string
	appdomain string
}

func GetClient() (client *Client, err error) {
	config, err := yaml.ReadFile("./lib/leancloud/conf.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	id, _ := (config.Get("AppID"))
	key, _ := (config.Get("AppKey"))
	if len(id) < 8 {
		err = errors.New("AppID配置错误")
		return
	}

	client = &Client{
		id: id,
		key: key,
		appdomain: strings.ToLower(id[0: 10]),
	}

	return
}

func (c *Client) Login(name string, password string) {
	// c
}