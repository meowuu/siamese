package leancloud

import (
	"fmt"
	// "github.com/dghubble/sling"
	// "net/http"
	"github.com/kylelemons/go-gypsy/yaml"
)

type Client struct {
	id string
	key string
}

func GetClient() (client *Client) {
	config, err := yaml.ReadFile("./lib/leancloud/conf.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	id, _ := (config.Get("AppID"))
	key, _ := (config.Get("AppKey"))

	client = &Client{
		id: id,
		key: key,
	}

	return
}
