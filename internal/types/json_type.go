package types

type RootAuth struct {
	Username string
	Password string
}

// 115AuthInfo 存储115认证信息
type Auth115Info struct {
	AccessToken  string `json:"accessToken"`  // 访问令牌
	RefreshToken string `json:"refreshToken"` // 刷新令牌
	ExpiresIn    int64  `json:"expiresIn"`    // 过期时间（秒）
}

// Search115FolderReq 搜索115文件夹请求
type Search115FolderReq struct {
	SearchValue string `json:"searchValue"` // 搜索关键字
	Limit       int    `json:"limit"`       // 单页记录数
	Offset      int    `json:"offset"`      // 数据显示偏移量
}

// Search115FolderResp 搜索115文件夹响应
type Search115FolderResp struct {
	Count int `json:"count"` // 搜索符合条件的文件(夹)总数
	Data  []struct {
		FileID       string `json:"file_id"`       // 文件ID
		UserID       string `json:"user_id"`       // 用户ID
		Sha1         string `json:"sha1"`          // 文件sha1值
		FileName     string `json:"file_name"`     // 文件名称
		FileSize     string `json:"file_size"`     // 文件大小
		UserPtime    string `json:"user_ptime"`    // 上传时间
		UserUtime    string `json:"user_utime"`    // 更新时间
		PickCode     string `json:"pick_code"`     // 文件提取码
		ParentID     string `json:"parent_id"`     // 父目录ID
		AreaID       string `json:"area_id"`       // 文件的状态
		IsPrivate    int    `json:"is_private"`    // 文件是否隐藏
		FileCategory string `json:"file_category"` // 1：文件；0：文件夹
		Ico          string `json:"ico"`           // 文件后缀
	} `json:"data"` // 搜索结果列表
	State   bool   `json:"state"`   // 状态
	Message string `json:"message"` // 异常信息
	Code    int    `json:"code"`    // 异常码
}

// AddOfflineTaskReq 添加云下载链接任务请求
type AddOfflineTaskReq struct {
	Urls     string `json:"urls"`       // 多个链接url,换行符分隔，支持HTTP(S)、FTP、磁力链和电驴链接
	WpPathID string `json:"wp_path_id"` // 保存目标文件夹id；如不传或传0，则默认保存到根目录下面
}

// AddOfflineTaskResp 添加云下载链接任务响应
type AddOfflineTaskResp struct {
	State   bool   `json:"state"`   // 操作结果状态值，true成功，false失败
	Message string `json:"message"` // 操作返回消息，成功时空值
	Code    int    `json:"code"`    // 操作返回号码，成功时返回0
	Data    []struct {
		State    bool   `json:"state"`     // 链接任务添加状态，成功true；失败false
		Code     int    `json:"code"`      // 链接任务状态码，成功返回0
		Message  string `json:"message"`   // 链接任务状态描述，成功返回空字符串
		InfoHash string `json:"info_hash"` // 链接任务sha1，只有任务成功的时候才会返回
		URL      string `json:"url"`       // 链接任务url
	} `json:"data"` // 数据
}

// GetOfflineTaskListResp 获取云下载任务列表响应
type GetOfflineTaskListResp struct {
	State   bool   `json:"state"`   // 操作结果状态值，true成功，false失败
	Message string `json:"message"` // 操作返回消息，成功时空值
	Code    int    `json:"code"`    // 操作返回号码，成功时返回0
	Data    struct {
		Page      int `json:"page"`       // 当前第几页
		PageCount int `json:"page_count"` // 总页数
		Count     int `json:"count"`      // 总数量
		Tasks     []struct {
			InfoHash     string  `json:"info_hash"`      // 任务sha1
			AddTime      int64   `json:"add_time"`       // 任务添加时间戳
			PercentDone  float64 `json:"percentDone"`    // 任务下载进度
			Size         int64   `json:"size"`           // 任务总大小（字节）
			Peers        int     `json:"peers"`          // 当前连接节点数
			RateDownload int     `json:"rateDownload"`   // 下载速度
			Name         string  `json:"name"`           // 任务名
			LastUpdate   int64   `json:"last_update"`    // 任务最后更新时间戳
			LeftTime     int     `json:"left_time"`      // 剩余时间
			FileID       string  `json:"file_id"`        // 任务源文件（夹）对应文件（夹）id
			DeleteFileID string  `json:"delete_file_id"` // 删除任务需删除源文件（夹）时，对应需传递的文件（夹）id
			Move         int     `json:"move"`           // 移动状态
			Status       int     `json:"status"`         // 任务状态：-1下载失败；0分配中；1下载中；2下载成功
			Err          int     `json:"err"`            // 错误码
			URL          string  `json:"url"`            // 链接任务url
			DelPath      string  `json:"del_path"`       // 删除路径
			WpPathID     string  `json:"wp_path_id"`     // 任务源文件所在父文件夹id
			Def2         int     `json:"def2"`           // 视频清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画
			PlayLong     *int    `json:"play_long"`      // 视频时长
			CanAppeal    int     `json:"can_appeal"`     // 是否可申诉
		} `json:"tasks"` // 云下载任务列表
	} `json:"data"` // 数据
}

// FileInfoResp 获取文件(夹)详情响应
type FileInfoResp struct {
	State   bool   `json:"state"`   // 操作结果状态值，true成功，false失败
	Message string `json:"message"` // 操作返回消息，成功时空值
	Code    int    `json:"code"`    // 操作返回号码，成功时返回0
	Data    struct {
		Count        int64  `json:"count"`          // 包含文件总数量
		Size         string `json:"size"`           // 文件(夹)总大小
		FolderCount  int    `json:"folder_count"`   // 包含文件夹总数量
		PlayLong     int    `json:"play_long"`      // 视频时长；-1：正在统计，其他数值为视频时长的数值(单位秒)
		ShowPlayLong int    `json:"show_play_long"` // 是否开启展示视频时长
		Ptime        string `json:"ptime"`          // 上传时间
		Utime        string `json:"utime"`          // 修改时间
		FileName     string `json:"file_name"`      // 文件名
		PickCode     string `json:"pick_code"`      // 文件提取码
		Sha1         string `json:"sha1"`           // sh1值
		FileID       string `json:"file_id"`        // 文件(夹)ID
		IsMark       string `json:"is_mark"`        // 是否星标
		OpenTime     int64  `json:"open_time"`      // 文件(夹)最近打开时间
		FileCategory string `json:"file_category"`  // 文件属性；1；文件；0：文件夹
		Paths        []struct {
			FileID   interface{} `json:"file_id"`   // 父目录ID，可能是字符串或数字
			FileName string      `json:"file_name"` // 父目录名称
			Iss      string      `json:"iss"`       // 路径标识
		} `json:"paths"` // 文件(夹)所在的路径
	} `json:"data"` // 数据
}

// GetFileDownloadUrlResp 获取文件下载地址响应
type GetFileDownloadUrlResp struct {
	State   bool   `json:"state"`   // 操作结果状态值，true成功，false失败
	Message string `json:"message"` // 操作返回消息，成功时空值
	Code    int    `json:"code"`    // 操作返回号码，成功时返回0
	Data    map[string]struct {
		FileName string `json:"file_name"` // 文件名
		FileSize int64  `json:"file_size"` // 文件大小
		PickCode string `json:"pick_code"` // 文件提取码
		Sha1     string `json:"sha1"`      // 文件sha1值
		URL      struct {
			URL string `json:"url"` // 文件下载地址
		} `json:"url"` // 文件下载地址信息
	} `json:"data"` // 数据
}

// GetFolderFilesResp 获取文件夹内文件列表响应
type GetFolderFilesResp struct {
	State   bool   `json:"state"`   // 操作结果状态值，true成功，false失败
	Message string `json:"message"` // 操作返回消息，成功时空值
	Code    int    `json:"code"`    // 操作返回号码，成功时返回0
	Data    []struct {
		FID   string `json:"fid"`   // 文件ID
		AID   string `json:"aid"`   // 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除
		PID   string `json:"pid"`   // 父目录ID
		FC    string `json:"fc"`    // 文件分类。0 文件夹，1 文件
		FN    string `json:"fn"`    // 文件(夹)名称
		FCO   string `json:"fco"`   // 文件夹封面
		ISM   string `json:"ism"`   // 是否星标，1：星标
		ISP   int    `json:"isp"`   // 是否加密；1：加密
		PC    string `json:"pc"`    // 文件提取码
		FS    int64  `json:"fs"`    // 文件大小
		UPT   int    `json:"upt"`   // 修改时间
		UET   int    `json:"uet"`   // 修改时间
		UPPT  int    `json:"uppt"`  // 上传时间
		CM    int    `json:"cm"`    // 未知字段
		FDESC string `json:"fdesc"` // 文件备注
		ISPL  int    `json:"ispl"`  // 是否统计文件夹下视频时长开关
		ICO   string `json:"ico"`   // 文件后缀名
	} `json:"data"` // 数据
}
