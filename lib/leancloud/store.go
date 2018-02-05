package leancloud

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/dghubble/sling"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/tidwall/gjson"
)

type Client struct {
	id        string
	key       string
	token     string
	appdomain string
}

type Book struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type SearchQuery struct {
	Where string `url:"where"`
}

type Section struct {
	Name   string   `json:"name"`
	ID     int      `json:"id"`
	Images []string `json:"images"`
	Url    string   `json:"url"`
	BookID string   `json:"bookid"`
	Index  int      `json:"index"`
}

// Wait Group
var wg sync.WaitGroup

// GetClient return leancloud client
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
		appdomain: fmt.Sprintf("https://%s.api.lncld.net/%s/", strings.ToLower(id[0:10]), apiversion),
	}

	return
}

func addHeader(req *http.Request, c *Client) {
	req.Header.Add("X-LC-Id", c.id)
	req.Header.Add("X-LC-Key", c.key)
}

func (c *Client) baseRequest() *sling.Sling {
	return sling.New().Base(c.appdomain)
}

// Login user
func (c *Client) Login(name string, password string) (state bool, err error) {
	client := &http.Client{}

	type UserInfo struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}

	req, _ := sling.New().Base(c.appdomain).Post("login").BodyJSON(&UserInfo{
		Name:     name,
		Password: password,
	}).Request()

	addHeader(req, c)

	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))

	if gjson.Get(string(data), "code").Int() != 200 {
		state = false
		err = errors.New(gjson.Get(string(data), "error").String())
		return
	}

	state = true
	c.token = gjson.Get(string(data), "sessionToken").String()
	return
}

func (c *Client) SaveBook(book Book) string {
	client := &http.Client{}

	req, _ := c.baseRequest().Post("classes/book").BodyJSON(book).Request()

	addHeader(req, c)

	res, _ := client.Do(req)
	data, _ := ioutil.ReadAll(res.Body)

	return gjson.Get(string(data), "objectId").String()
}

func (c *Client) GetBookByName(name string) (objectid string) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	client := &http.Client{}

	req, _ := c.baseRequest().Get("classes/book").QueryStruct(SearchQuery{
		Where: fmt.Sprintf(`{"name":"%s"}`, name),
	}).Request()

	addHeader(req, c)

	res, _ := client.Do(req)

	data, _ := ioutil.ReadAll(res.Body)

	objectid = gjson.Get(string(data), "results").Array()[0].Get("objectId").String()
	return
}

func (c *Client) GetSectionByName(name string) (objectid string) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	client := &http.Client{}

	req, _ := c.baseRequest().Get("classes/section").QueryStruct(SearchQuery{
		Where: fmt.Sprintf(`{"name":"%s"}`, name),
	}).Request()

	addHeader(req, c)

	res, _ := client.Do(req)

	data, _ := ioutil.ReadAll(res.Body)

	objectid = gjson.Get(string(data), "results").Array()[0].Get("objectId").String()
	return
}

func (c *Client) SaveSection(section Section) {
	defer wg.Done()

	if c.GetSectionByName(section.Name) != "" {
		return
	}

	client := &http.Client{}

	req, _ := c.baseRequest().Post("classes/section").BodyJSON(section).Request()

	addHeader(req, c)
	res, _ := client.Do(req)
	if res.StatusCode == 200 || res.StatusCode == 201 {
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("%s保存失败, 结果%s\n", section.Name, string(data))
	}
}

// Save data
func (c *Client) Save(book Book, sections []Section) {
	fmt.Println("开始保存")
	objectId := c.GetBookByName(book.Name)
	if objectId == "" {
		objectId = c.SaveBook(book)
	}

	gornums := 0
	for _, section := range sections {
		gornums++
		wg.Add(1)

		section.BookID = objectId
		go c.SaveSection(section)

		if gornums == 6 {
			gornums = 0
			wg.Wait()
		}
	}

	wg.Wait()
	fmt.Println("保存结束 🐳")
}
