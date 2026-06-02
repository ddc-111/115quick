package svc

import (
	"115Quick_server/internal/config"
	"115Quick_server/internal/types"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config                  config.Config
	Auth115                 *Auth115Manager
	Store                   *Store
	Mode                    int64
	DownloadLinkMessageChan chan *DownloadLinkMessage
	folderID                string
}

type DownloadLinkMessage struct {
	DownloadLink string
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx := &ServiceContext{
		Config:                  c,
		Mode:                    0,
		DownloadLinkMessageChan: make(chan *DownloadLinkMessage, 1000000),
	}

	store, err := NewStore(c.DBPath)
	if err != nil {
		logx.Errorf("初始化数据库失败: %v", err)
		panic(err)
	}
	ctx.Store = store

	ctx.Auth115 = NewAuth115Manager(ctx)

	return ctx
}

func (s *ServiceContext) SetFolderID(folderID string) {
	s.folderID = folderID
}

func (s *ServiceContext) GetFolderID() string {
	return s.folderID
}

func (s *ServiceContext) AddDownloadLinkMessage(DownloadLink string) {
	go func() {
		select {
		case s.DownloadLinkMessageChan <- &DownloadLinkMessage{DownloadLink: DownloadLink}:
		default:
			log.Println("broadcast queue full, message dropped")
		}
	}()
}

func (s *ServiceContext) StarDownloadLinkWorker() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for msg := range s.DownloadLinkMessageChan {
			<-ticker.C

			randomDelay := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(randomDelay)

			s.addDownloadLink(msg.DownloadLink)
		}
	}()
}

func (s *ServiceContext) addDownloadLink(DownloadLink string) {
	defer s.Auth115.DoOnceWithCooldownPoll()
	token, err := s.Auth115.GetAccessToken()
	if err != nil {
		fmt.Errorf("获取115访问令牌失败: %v", err)
	}
	WpPathId := s.GetFolderID()

	formData := url.Values{}
	formData.Set("urls", DownloadLink)
	formData.Set("wp_path_id", WpPathId)

	apiURL := "https://proapi.115.com/open/offline/add_task_urls"
	request, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("读取响应体失败: %v", err)
	}

	var respData types.AddOfflineTaskResp
	if err := json.Unmarshal(body, &respData); err != nil {
		fmt.Errorf("解析响应失败: %v", err)
	}

	if !respData.State {
		if respData.Code == 10008 {
			s.Auth115.AddFileInfo("", DownloadLink)
			return
		}
		fmt.Errorf("添加下载任务失败: %s", respData.Message)
	}

	fmt.Printf("添加下载任务成功: %+v", respData)

	for i, item := range respData.Data {
		if !item.State {
			if item.Code == 10008 {
				s.Auth115.AddFileInfo(item.InfoHash, item.URL)
				fmt.Printf("链接 %d 重复添加成功: URL: %s", i+1, item.URL)
			} else {
				fmt.Printf("链接 %d 添加失败: %s, 错误码: %d, URL: %s", i+1, item.Message, item.Code, item.URL)
			}
		} else {
			fmt.Printf("链接 %d 添加成功: URL: %s, InfoHash: %s", i+1, item.URL, item.InfoHash)
			s.Auth115.AddFileInfo(item.InfoHash, item.URL)
		}
	}

	return
}
