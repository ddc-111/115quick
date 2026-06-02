package Download

import (
	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddDownloadLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加下载链接
func NewAddDownloadLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDownloadLinkLogic {
	return &AddDownloadLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddDownloadLinkLogic) AddDownloadLink(req *types.AddDownloadLinkReq) (resp *types.NilResp, err error) {
	// 获取115访问令牌
	// if l.svcCtx.Auth.IsExpired() {
	// 	l.Logger.Errorf("权限已经过期")
	// 	return nil, fmt.Errorf("权限已经过期")
	// }
	//token, err := l.svcCtx.Auth115.GetAccessToken()
	//if err != nil {
	//	l.Logger.Errorf("获取115访问令牌失败: %v", err)
	//	return nil, fmt.Errorf("获取115访问令牌失败: %v", err)
	//}
	//var WpPathId string
	//if err := l.svcCtx.GetSecureData(svc.SecureKey115FolderID, &WpPathId); err != nil {
	//	l.Logger.Errorf("获取115文件夹ID失败: %v", err)
	//	return nil, fmt.Errorf("获取115文件夹ID失败: %v", err)
	//}
	//
	//// 构建表单数据
	//formData := url.Values{}
	//formData.Set("urls", req.DownloadLink)
	//formData.Set("wp_path_id", WpPathId)
	//
	//// 创建HTTP请求
	//apiURL := "https://proapi.115.com/open/offline/add_task_urls"
	//request, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	//if err != nil {
	//	l.Logger.Errorf("创建HTTP请求失败: %v", err)
	//	return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	//}
	//
	//// 设置请求头
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//request.Header.Set("Authorization", "Bearer "+token)
	//
	//// 发送请求
	//client := &http.Client{}
	//response, err := client.Do(request)
	//if err != nil {
	//	l.Logger.Errorf("发送HTTP请求失败: %v", err)
	//	return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	//}
	//defer response.Body.Close()
	//
	//// 读取响应体
	//body, err := io.ReadAll(response.Body)
	//if err != nil {
	//	l.Logger.Errorf("读取响应体失败: %v", err)
	//	return nil, fmt.Errorf("读取响应体失败: %v", err)
	//}
	//
	//// 解析响应
	//var respData types.AddOfflineTaskResp
	//if err := json.Unmarshal(body, &respData); err != nil {
	//	l.Logger.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	//	return nil, fmt.Errorf("解析响应失败: %v", err)
	//}
	//
	//// 检查响应状态
	//if !respData.State {
	//	l.svcCtx.Auth115.RefreshToken()
	//	if respData.Code == 10008 {
	//		// 添加一个重复的云下载链接 但是还进行下载
	//		l.Logger.Infof("添加一个重复的云下载链接 但是还进行下载: %s", req.DownloadLink)
	//		l.svcCtx.Auth115.AddFileInfo("", req.DownloadLink)
	//
	//		l.Logger.Info("启动定时器，60秒后开始轮询文件 如果秒链接完成则直接开始下载")
	//		go func() {
	//			time.Sleep(60 * time.Second)
	//			l.Logger.Info("定时器触发，开始轮询文件")
	//			l.svcCtx.Auth115.PollFiles()
	//		}()
	//
	//		return &types.NilResp{}, nil
	//	}
	//	l.Logger.Errorf("添加下载任务失败: %s, 错误码: %d", respData.Message, respData.Code)
	//	return nil, fmt.Errorf("添加下载任务失败: %s", respData.Message)
	//}
	//
	//// 打印响应结果
	//l.Logger.Infof("添加下载任务成功: %+v", respData)
	//
	//// 检查每个链接的添加状态并保存文件信息
	//hasSuccess := false
	//for i, item := range respData.Data {
	//	if !item.State {
	//		if item.Code == 10008 {
	//			l.svcCtx.Auth115.AddFileInfo(item.InfoHash, item.URL)
	//			hasSuccess = true
	//			//添加了一个重复的云下载链接 但是还进行下载
	//			l.Logger.Infof("链接 %d 重复添加成功: URL: %s", i+1, item.URL)
	//		} else {
	//			l.Logger.Infof("链接 %d 添加失败: %s, 错误码: %d, URL: %s", i+1, item.Message, item.Code, item.URL)
	//		}
	//
	//	} else {
	//		l.Logger.Infof("链接 %d 添加成功: URL: %s, InfoHash: %s", i+1, item.URL, item.InfoHash)
	//		// 保存文件信息到 Auth115Manager
	//		l.svcCtx.Auth115.AddFileInfo(item.InfoHash, item.URL)
	//		hasSuccess = true
	//	}
	//}
	//
	//// 如果有成功添加的链接，启动定时器在10秒后调用pollFiles
	//if hasSuccess {
	//	l.Logger.Info("启动定时器，60秒后开始轮询文件 如果秒链接完成则直接开始下载")
	//	go func() {
	//		time.Sleep(60 * time.Second)
	//		l.Logger.Info("定时器触发，开始轮询文件")
	//		l.svcCtx.Auth115.PollFiles()
	//	}()
	//}

	l.svcCtx.AddDownloadLinkMessage(req.DownloadLink)
	return &types.NilResp{}, nil
}
