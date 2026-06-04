package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"115Quick_server/internal/service/auth"
	"115Quick_server/internal/service/cloudtask"
	"115Quick_server/internal/service/download"
	"115Quick_server/internal/service/file"
	"115Quick_server/internal/service/video"
	"115Quick_server/internal/store"
	"115Quick_server/internal/version"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tokenMgr    *auth.Manager
	store       *store.Store
	fileSvc     *file.Service
	videoSvc    *video.Service
	cloudSvc    *cloudtask.Service
	downloadEng *download.Engine
}

func NewHandler(
	tokenMgr *auth.Manager,
	store *store.Store,
	fileSvc *file.Service,
	videoSvc *video.Service,
	cloudSvc *cloudtask.Service,
	downloadEng *download.Engine,
) *Handler {
	return &Handler{
		tokenMgr:    tokenMgr,
		store:       store,
		fileSvc:     fileSvc,
		videoSvc:    videoSvc,
		cloudSvc:    cloudSvc,
		downloadEng: downloadEng,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/v1/Download/api")
	{
		// Server info
		api.GET("/getServerInfo", h.getServerInfo)
		api.GET("/getServerLogs", h.getServerLogs)

		// Token management
		api.GET("/getTokenStatus", h.getTokenStatus)
		api.POST("/setToken", h.setToken)

		// Download progress
		api.GET("/getDownloadProgress", h.getDownloadProgress)

		// Cloud files (old path compatible)
		api.GET("/getCloudFiles", h.getCloudFiles)
		api.POST("/createFolder", h.createFolder)
		api.POST("/deleteFile", h.deleteFile)
		api.POST("/renameFile", h.renameFile)
		api.POST("/moveFile", h.moveFile)

		// Cloud tasks (old path compatible)
		api.GET("/getCloudTasks", h.getCloudTasks)
		api.GET("/getTaskHistory", h.getTaskHistory)
		api.POST("/removeDownloadTask", h.removeDownloadTask)
		api.POST("/refreshTasks", h.refreshTasks)
		api.POST("/clearTaskHistory", h.clearTaskHistory)
		api.POST("/clearCompletedTasks", h.clearCompletedTasks)

		// Download link
		api.POST("/addDownloadLink", h.addDownloadLink)
		api.POST("/setDownloadMode", h.setDownloadMode)

		// Default download dir
		api.GET("/getDefaultDownloadDir", h.getDefaultDownloadDir)
		api.POST("/setDefaultDownloadDir", h.setDefaultDownloadDir)

		// Video (new paths)
		api.GET("/video/play", h.getPlayInfo)
		api.GET("/video/subtitles", h.getSubtitles)
		api.GET("/video/history", h.getPlayHistory)
		api.POST("/video/history", h.savePlayHistory)
		api.POST("/video/transcode", h.submitTranscode)
		api.GET("/video/proxy", h.proxyVideo)

		// File download URL
		api.GET("/files/download", h.getDownloadURL)

		// Cloud task management (new paths)
		api.GET("/cloud/tasks", h.getTaskList)
		api.POST("/cloud/addUrls", h.addTaskURLs)
		api.POST("/cloud/addBt", h.addTaskBT)
		api.POST("/cloud/delete", h.deleteTask)
		api.POST("/cloud/clear", h.clearTasks)
		api.GET("/cloud/quota", h.getQuotaInfo)
	}
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

func errorResp(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
}

// ==================== Server Info ====================

func (h *Handler) getServerInfo(c *gin.Context) {
	info := version.GetInfo()
	release, err := version.CheckUpdate()

	result := gin.H{
		"version":    info.Version,
		"git_commit": info.GitCommit,
		"build_time": info.BuildTime,
	}

	if err == nil && release != nil {
		result["latest_version"] = release.TagName
		result["release_url"] = release.HTMLURL
		result["has_update"] = release.TagName != "v"+info.Version
	}

	success(c, result)
}

func (h *Handler) getServerLogs(c *gin.Context) {
	logType := c.DefaultQuery("type", "stdout")
	lines, _ := strconv.Atoi(c.DefaultQuery("lines", "200"))

	log.Printf("getServerLogs called: type=%s, lines=%d", logType, lines)
	success(c, gin.H{
		"logs": []string{"Logs are available in the server console"},
		"type": logType,
	})
}

// ==================== Token Management ====================

func (h *Handler) getTokenStatus(c *gin.Context) {
	status := h.tokenMgr.GetStatus()
	success(c, status)
}

func (h *Handler) setToken(c *gin.Context) {
	var req struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		refreshToken = req.AccessToken
	}

	if err := h.tokenMgr.SetToken("", refreshToken, 7200); err != nil {
		errorResp(c, 500, "set token failed: "+err.Error())
		return
	}

	if err := h.tokenMgr.ForceRefresh(); err != nil {
		errorResp(c, 500, "refresh token failed: "+err.Error())
		return
	}

	success(c, gin.H{"message": "token set successfully"})
}

// ==================== Cloud Files ====================

func (h *Handler) getCloudFiles(c *gin.Context) {
	folderId := c.DefaultQuery("folderId", "0")
	fileType, _ := strconv.Atoi(c.DefaultQuery("fileType", "0"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	cid := folderId
	if cid == "" {
		cid = "0"
	}

	result, err := h.fileSvc.GetFileList(cid, offset, limit, fileType)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) createFolder(c *gin.Context) {
	var req struct {
		ParentId   string `json:"parentId"`
		PID        string `json:"pid"`
		FolderName string `json:"folderName"`
		FileName   string `json:"file_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	pid := req.ParentId
	if pid == "" {
		pid = req.PID
	}
	name := req.FolderName
	if name == "" {
		name = req.FileName
	}

	fileID, err := h.fileSvc.CreateFolder(pid, name)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"file_id": fileID})
}

func (h *Handler) deleteFile(c *gin.Context) {
	var req struct {
		FileIds string `json:"fileIds"`
		FileIDs string `json:"file_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	fileIds := req.FileIds
	if fileIds == "" {
		fileIds = req.FileIDs
	}

	if err := h.fileSvc.DeleteFile(fileIds); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "deleted"})
}

func (h *Handler) renameFile(c *gin.Context) {
	var req struct {
		FileId  string `json:"fileId"`
		FileID  string `json:"file_id"`
		NewName string `json:"newName"`
		FileName string `json:"file_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	fileId := req.FileId
	if fileId == "" {
		fileId = req.FileID
	}
	newName := req.NewName
	if newName == "" {
		newName = req.FileName
	}

	if err := h.fileSvc.RenameFile(fileId, newName); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "renamed"})
}

func (h *Handler) moveFile(c *gin.Context) {
	var req struct {
		FileIds   string `json:"fileIds"`
		FileIDs   string `json:"file_ids"`
		TargetDir string `json:"targetDir"`
		ToCID     string `json:"to_cid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	fileIds := req.FileIds
	if fileIds == "" {
		fileIds = req.FileIDs
	}
	targetDir := req.TargetDir
	if targetDir == "" {
		targetDir = req.ToCID
	}

	if err := h.fileSvc.MoveFile(fileIds, targetDir); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "moved"})
}

func (h *Handler) getDownloadURL(c *gin.Context) {
	pickCode := c.Query("pick_code")
	if pickCode == "" {
		errorResp(c, 400, "pick_code is required")
		return
	}

	url, err := h.fileSvc.GetDownloadURL(pickCode)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"url": url})
}

// ==================== Cloud Tasks ====================

func (h *Handler) getCloudTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	result, err := h.cloudSvc.GetTaskList(page)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) getTaskList(c *gin.Context) {
	h.getCloudTasks(c)
}

func (h *Handler) getTaskHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	tasks, total, err := h.store.GetTaskHistory(page, pageSize)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{
		"list":  tasks,
		"total": total,
		"page":  page,
	})
}

func (h *Handler) removeDownloadTask(c *gin.Context) {
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	success(c, gin.H{"message": "task removed"})
}

func (h *Handler) refreshTasks(c *gin.Context) {
	success(c, gin.H{"message": "tasks refreshed"})
}

func (h *Handler) clearTaskHistory(c *gin.Context) {
	if err := h.store.ClearTaskHistory(); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "task history cleared"})
}

func (h *Handler) clearCompletedTasks(c *gin.Context) {
	if err := h.store.ClearCompletedTasks(); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "completed tasks cleared"})
}

func (h *Handler) addTaskURLs(c *gin.Context) {
	var req struct {
		URLs     string `json:"urls" binding:"required"`
		WpPathID string `json:"wp_path_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	result, err := h.cloudSvc.AddTaskURLs(req.URLs, req.WpPathID)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) addTaskBT(c *gin.Context) {
	var req struct {
		InfoHash    string `json:"info_hash" binding:"required"`
		Wanted      string `json:"wanted" binding:"required"`
		SavePath    string `json:"save_path" binding:"required"`
		TorrentSHA1 string `json:"torrent_sha1" binding:"required"`
		PickCode    string `json:"pick_code" binding:"required"`
		WpPathID    string `json:"wp_path_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.cloudSvc.AddTaskBT(req.InfoHash, req.Wanted, req.SavePath, req.TorrentSHA1, req.PickCode, req.WpPathID); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "bt task added"})
}

func (h *Handler) deleteTask(c *gin.Context) {
	var req struct {
		InfoHash         string `json:"info_hash" binding:"required"`
		DeleteSourceFile bool   `json:"delete_source_file"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.cloudSvc.DeleteTask(req.InfoHash, req.DeleteSourceFile); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "task deleted"})
}

func (h *Handler) clearTasks(c *gin.Context) {
	var req struct {
		Flag int `json:"flag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.cloudSvc.ClearTasks(req.Flag); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "tasks cleared"})
}

func (h *Handler) getQuotaInfo(c *gin.Context) {
	result, err := h.cloudSvc.GetQuotaInfo()
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

// ==================== Download ====================

func (h *Handler) addDownloadLink(c *gin.Context) {
	var req struct {
		DownloadLink string `json:"downloadLink"`
		URLs         string `json:"urls"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	links := req.DownloadLink
	if links == "" {
		links = req.URLs
	}

	if links == "" {
		errorResp(c, 400, "download link is required")
		return
	}

	result, err := h.cloudSvc.AddTaskURLs(links, "0")
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) setDownloadMode(c *gin.Context) {
	var req struct {
		Mode int `json:"mode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	h.store.SetConfig("download_mode", fmt.Sprintf("%d", req.Mode))
	success(c, gin.H{"mode": req.Mode})
}

func (h *Handler) getDownloadProgress(c *gin.Context) {
	result := h.downloadEng.GetAllProgress()
	success(c, result)
}

func (h *Handler) getDefaultDownloadDir(c *gin.Context) {
	folderId, _ := h.store.GetConfig("default_download_folder_id")
	folderName, _ := h.store.GetConfig("default_download_folder_name")

	success(c, gin.H{
		"folderId":   folderId,
		"folderName": folderName,
	})
}

func (h *Handler) setDefaultDownloadDir(c *gin.Context) {
	var req struct {
		FolderId   string `json:"folderId"`
		FolderName string `json:"folderName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	h.store.SetConfig("default_download_folder_id", req.FolderId)
	h.store.SetConfig("default_download_folder_name", req.FolderName)

	success(c, gin.H{"message": "default download dir set"})
}

// ==================== Video ====================

func (h *Handler) getPlayInfo(c *gin.Context) {
	pickCode := c.Query("pick_code")
	if pickCode == "" {
		errorResp(c, 400, "pick_code is required")
		return
	}

	result, err := h.videoSvc.GetPlayInfo(pickCode)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) getSubtitles(c *gin.Context) {
	pickCode := c.Query("pick_code")
	if pickCode == "" {
		errorResp(c, 400, "pick_code is required")
		return
	}

	result, err := h.videoSvc.GetSubtitles(pickCode)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) getPlayHistory(c *gin.Context) {
	pickCode := c.Query("pick_code")
	if pickCode == "" {
		errorResp(c, 400, "pick_code is required")
		return
	}

	result, err := h.videoSvc.GetPlayHistory(pickCode)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) savePlayHistory(c *gin.Context) {
	var req struct {
		PickCode string `json:"pick_code" binding:"required"`
		Time     int    `json:"time"`
		WatchEnd bool   `json:"watch_end"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.videoSvc.SavePlayHistory(req.PickCode, req.Time, req.WatchEnd); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "saved"})
}

func (h *Handler) submitTranscode(c *gin.Context) {
	var req struct {
		PickCode string `json:"pick_code" binding:"required"`
		Op       string `json:"op" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.videoSvc.SubmitTranscode(req.PickCode, req.Op); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "transcode submitted"})
}

func (h *Handler) proxyVideo(c *gin.Context) {
	pickCode := c.Query("pick_code")
	if pickCode == "" {
		errorResp(c, 400, "pick_code is required")
		return
	}

	definition := c.Query("definition")

	if err := h.videoSvc.ProxyVideoStream(pickCode, definition, c.Writer, c.Request); err != nil {
		if !c.Writer.Written() {
			errorResp(c, 500, err.Error())
		}
	}
}

// ==================== Unused but keep compatibility ====================

func (h *Handler) startRename(c *gin.Context) {
	success(c, gin.H{"message": "rename feature not available"})
}

func (h *Handler) getSMBConfig(c *gin.Context) {
	success(c, gin.H{
		"enabled": false,
		"message": "SMB support has been removed",
	})
}

func (h *Handler) setSMBConfig(c *gin.Context) {
	success(c, gin.H{"message": "SMB support has been removed"})
}

func (h *Handler) testSMBConnection(c *gin.Context) {
	errorResp(c, 400, "SMB support has been removed")
}

func (h *Handler) smbBrowse(c *gin.Context) {
	errorResp(c, 400, "SMB support has been removed")
}

func (h *Handler) smbDownload(c *gin.Context) {
	errorResp(c, 400, "SMB support has been removed")
}
