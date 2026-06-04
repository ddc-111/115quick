package handler

import (
	"net/http"
	"strconv"

	"115Quick_server/internal/service/auth"
	"115Quick_server/internal/service/cloudtask"
	"115Quick_server/internal/service/download"
	"115Quick_server/internal/service/file"
	"115Quick_server/internal/service/video"
	"115Quick_server/internal/version"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tokenMgr    *auth.Manager
	fileSvc     *file.Service
	videoSvc    *video.Service
	cloudSvc    *cloudtask.Service
	downloadEng *download.Engine
}

func NewHandler(
	tokenMgr *auth.Manager,
	fileSvc *file.Service,
	videoSvc *video.Service,
	cloudSvc *cloudtask.Service,
	downloadEng *download.Engine,
) *Handler {
	return &Handler{
		tokenMgr:    tokenMgr,
		fileSvc:     fileSvc,
		videoSvc:    videoSvc,
		cloudSvc:    cloudSvc,
		downloadEng: downloadEng,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/v1/Download/api")
	{
		api.GET("/getServerInfo", h.getServerInfo)
		api.GET("/getTokenStatus", h.getTokenStatus)
		api.POST("/setToken", h.setToken)

		files := api.Group("/files")
		{
			files.GET("/list", h.getFileList)
			files.GET("/info", h.getFileInfo)
			files.GET("/search", h.searchFiles)
			files.POST("/createFolder", h.createFolder)
			files.POST("/delete", h.deleteFile)
			files.POST("/rename", h.renameFile)
			files.POST("/move", h.moveFile)
			files.GET("/download", h.getDownloadURL)
		}

		video := api.Group("/video")
		{
			video.GET("/play", h.getPlayInfo)
			video.GET("/subtitles", h.getSubtitles)
			video.GET("/history", h.getPlayHistory)
			video.POST("/history", h.savePlayHistory)
			video.POST("/transcode", h.submitTranscode)
			video.GET("/proxy", h.proxyVideo)
		}

		cloud := api.Group("/cloud")
		{
			cloud.GET("/tasks", h.getTaskList)
			cloud.POST("/addUrls", h.addTaskURLs)
			cloud.POST("/addBt", h.addTaskBT)
			cloud.POST("/delete", h.deleteTask)
			cloud.POST("/clear", h.clearTasks)
			cloud.GET("/quota", h.getQuotaInfo)
		}

		download := api.Group("/download")
		{
			download.GET("/progress", h.getDownloadProgress)
			download.POST("/start", h.startDownload)
		}
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

func (h *Handler) getTokenStatus(c *gin.Context) {
	status := h.tokenMgr.GetStatus()
	success(c, status)
}

func (h *Handler) setToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.tokenMgr.SetToken("", req.RefreshToken, 7200); err != nil {
		errorResp(c, 500, "set token failed: "+err.Error())
		return
	}

	if err := h.tokenMgr.ForceRefresh(); err != nil {
		errorResp(c, 500, "refresh token failed: "+err.Error())
		return
	}

	success(c, gin.H{"message": "token set successfully"})
}

func (h *Handler) getFileList(c *gin.Context) {
	cid := c.Query("cid")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	fileType, _ := strconv.Atoi(c.Query("type"))

	result, err := h.fileSvc.GetFileList(cid, offset, limit, fileType)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) getFileInfo(c *gin.Context) {
	fileID := c.Query("file_id")
	if fileID == "" {
		errorResp(c, 400, "file_id is required")
		return
	}

	result, err := h.fileSvc.GetFileInfo(fileID)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) searchFiles(c *gin.Context) {
	keyword := c.Query("keyword")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.fileSvc.SearchFiles(keyword, offset, limit)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
}

func (h *Handler) createFolder(c *gin.Context) {
	var req struct {
		PID      string `json:"pid" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	fileID, err := h.fileSvc.CreateFolder(req.PID, req.FileName)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"file_id": fileID})
}

func (h *Handler) deleteFile(c *gin.Context) {
	var req struct {
		FileIDs string `json:"file_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.fileSvc.DeleteFile(req.FileIDs); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "deleted"})
}

func (h *Handler) renameFile(c *gin.Context) {
	var req struct {
		FileID   string `json:"file_id" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.fileSvc.RenameFile(req.FileID, req.FileName); err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, gin.H{"message": "renamed"})
}

func (h *Handler) moveFile(c *gin.Context) {
	var req struct {
		FileIDs string `json:"file_ids" binding:"required"`
		ToCID   string `json:"to_cid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	if err := h.fileSvc.MoveFile(req.FileIDs, req.ToCID); err != nil {
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
		errorResp(c, 500, err.Error())
	}
}

func (h *Handler) getTaskList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	result, err := h.cloudSvc.GetTaskList(page)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, result)
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
		InfoHash        string `json:"info_hash" binding:"required"`
		DeleteSourceFile bool  `json:"delete_source_file"`
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

func (h *Handler) getDownloadProgress(c *gin.Context) {
	result := h.downloadEng.GetAllProgress()
	success(c, result)
}

func (h *Handler) startDownload(c *gin.Context) {
	var req struct {
		URL      string `json:"url" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
		FileSize int64  `json:"file_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp(c, 400, "invalid request: "+err.Error())
		return
	}

	go func() {
		h.downloadEng.Download(c.Request.Context(), req.URL, req.FileName, req.FileSize)
	}()

	success(c, gin.H{"message": "download started"})
}
