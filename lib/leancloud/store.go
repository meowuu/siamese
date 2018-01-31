package leancloud

import (
	"errors"
	"fmt"
	"strings"
	// "github.com/dghubble/sling"
	// "net/http"
	"github.com/kylelemons/go-gypsy/yaml"
)

type Client struct {
	id        string
	key       string
	token     string
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
	apiversion, _ := (config.Get("version"))
	if len(id) < 8 {
		err = errors.New("AppID配置错误")
		return
	}

	client = &Client{
		id:        id,
		key:       key,
		appdomain: fmt.Sprintf("https://%s.api.lncld.net/%s", strings.ToLower(id[0:10]), apiversion),
	}

	return
}

func (c *Client) Login(name string, password string) {
	// c
}
