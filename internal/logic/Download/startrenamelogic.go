package Download

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"115Quick_server/internal/svc"
	"115Quick_server/internal/types"
	"115Quick_server/internal/utils/request115"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartReNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重命名
func NewStartReNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartReNameLogic {
	return &StartReNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartReNameLogic) StartReName(req *types.StartReNameReq) (resp *types.NilResp, err error) {
	resp = &types.NilResp{}

	// 使用默认正则表达式，如果请求中没有指定
	filePattern := req.FilePattern
	if filePattern == "" {
		filePattern = `(?i)[a-z]+-\d+` // 默认模式匹配类似 "snis-1234" 的格式
	}

	// 设置起始文件夹ID
	startFolderID := req.FolderID
	if startFolderID == "" {
		startFolderID = "0" // 根目录ID
	}

	// 根据模式选择不同的处理方式
	if req.Mode == 0 {
		// 全量重命名模式 - 递归处理所有文件夹
		err = l.processFolder(startFolderID, filePattern, req.DeleteOthers, req.RenameToMatch)
	} else {
		// 增量重命名模式 - 只处理指定文件夹
		// 检查文件夹名是否匹配特定模式
		folderInfo, err := request115.GetFolderInfo(startFolderID)
		if err != nil {
			l.Error(fmt.Errorf("获取文件夹信息失败: %v", err))
			resp.Success = false
			resp.Message = err.Error()
			return resp, nil
		}

		// 编译正则表达式
		pattern, err := regexp.Compile(filePattern)
		if err != nil {
			l.Error(fmt.Errorf("无效的正则表达式: %v", err))
			resp.Success = false
			resp.Message = "无效的正则表达式: " + err.Error()
			return resp, nil
		}

		if pattern.MatchString(folderInfo.FileName) {
			// 如果文件夹名符合特定模式，处理该文件夹内的文件
			targetName := pattern.FindString(folderInfo.FileName)
			err = l.processSpecialFolder(startFolderID, targetName, req.DeleteOthers, req.RenameToMatch)
			if err != nil {
				l.Error(fmt.Errorf("处理特定文件夹失败: %v", err))
				resp.Success = false
				resp.Message = err.Error()
			}
		} else {
			l.Infof("指定文件夹 %s 不符合正则表达式格式 %s", folderInfo.FileName, filePattern)
		}
	}

	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
	}

	return resp, nil
}

// 递归处理文件夹
func (l *StartReNameLogic) processFolder(folderID string, filePattern string, deleteOthers bool, renameToMatch bool) error {
	// 获取文件夹内容
	fileList, err := request115.GetFileList(folderID, 1, 1000) // 一次获取1000个文件/文件夹
	if err != nil {
		return fmt.Errorf("获取文件夹 %s 内容失败: %v", folderID, err)
	}

	// 检查当前文件夹名称是否符合特定模式
	folderInfo, err := request115.GetFolderInfo(folderID)
	if err != nil {
		return fmt.Errorf("获取文件夹 %s 详情失败: %v", folderID, err)
	}

	// 编译正则表达式
	pattern, err := regexp.Compile(filePattern)
	if err != nil {
		return fmt.Errorf("无效的正则表达式 %s: %v", filePattern, err)
	}

	// 如果当前文件夹名符合特定模式，处理该文件夹内的文件
	if pattern.MatchString(folderInfo.FileName) {
		l.Logger.Infof("发现特定格式文件夹: %s (ID: %s)", folderInfo.FileName, folderID)
		// 提取特定格式字符串 (如 snis-1234)
		targetName := pattern.FindString(folderInfo.FileName)
		err = l.processSpecialFolder(folderID, targetName, deleteOthers, renameToMatch)
		if err != nil {
			return fmt.Errorf("处理特定文件夹 %s 失败: %v", folderInfo.FileName, err)
		}
	}

	// 递归处理所有子文件夹
	for _, file := range fileList.Data {
		// 只处理文件夹
		if file.FileCategory == "0" {
			err = l.processFolder(file.FileID, filePattern, deleteOthers, renameToMatch)
			if err != nil {
				l.Error(err)
				// 继续处理其他文件夹，不中断整个过程
			}
		}
	}

	return nil
}

// 处理特定格式文件夹 (包含类似 snis-1234 格式的文件夹)
func (l *StartReNameLogic) processSpecialFolder(folderID string, targetName string, deleteOthers bool, renameToMatch bool) error {
	// 获取文件夹中的所有文件
	fileList, err := request115.GetFileList(folderID, 1, 1000)
	if err != nil {
		return err
	}

	// 筛选出只有文件（不包括文件夹）
	var files []request115.File
	for _, file := range fileList.Data {
		if file.FileCategory == "1" { // 只处理文件，不处理文件夹
			files = append(files, file)
		}
	}

	if len(files) == 0 {
		return nil // 没有文件需要处理
	}

	// 按文件大小排序
	sort.Slice(files, func(i, j int) bool {
		// 降序排列，最大的文件排在前面
		return files[i].FileSize > files[j].FileSize
	})

	// 获取最大的文件
	largestFile := files[0]
	l.Logger.Infof("最大文件: %s (大小: %s)", largestFile.FileName, largestFile.FileSize)

	// 检查最大文件的名称是否需要重命名
	shouldRename := false

	if renameToMatch {
		// 编译正则表达式匹配类似 snis-1234 的格式
		pattern := regexp.MustCompile(`(?i)[a-z]+-\d+`)

		// 如果文件名中没有包含特定格式 或者 文件名混乱，需要重命名
		if !pattern.MatchString(largestFile.FileName) {
			shouldRename = true
		} else {
			// 如果文件名中已包含特定格式但与文件夹格式不一致，也需要重命名
			existingFormat := pattern.FindString(largestFile.FileName)
			if !strings.EqualFold(existingFormat, targetName) {
				shouldRename = true
			}
		}

		// 如果需要重命名最大文件
		if shouldRename {
			// 获取文件扩展名
			fileParts := strings.Split(largestFile.FileName, ".")
			var extension string
			if len(fileParts) > 1 {
				extension = "." + fileParts[len(fileParts)-1]
			}

			// 新文件名: 特定格式 + 扩展名
			newFileName := targetName + extension

			// 重命名文件
			err = request115.RenameFile(largestFile.FileID, newFileName)
			if err != nil {
				return fmt.Errorf("重命名文件 %s 失败: %v", largestFile.FileName, err)
			}

			l.Logger.Infof("文件已重命名: %s -> %s", largestFile.FileName, newFileName)
		}
	}

	// 删除其他所有文件（保留最大的文件）
	if deleteOthers && len(files) > 1 {
		l.Logger.Infof("准备删除 %d 个较小的文件", len(files)-1)
		for i := 1; i < len(files); i++ {
			err = request115.DeleteFile(files[i].FileID)
			if err != nil {
				l.Logger.Errorf("删除文件 %s 失败: %v", files[i].FileName, err)
				// 继续删除其他文件，不中断整个过程
			} else {
				l.Logger.Infof("已删除文件: %s", files[i].FileName)
			}
		}
	}

	return nil
}
