# 115开放平台 API 文档

> 来源: https://www.yuque.com/115yun/open

---

## 简介

### 概述

**一、“115生活”简介**
“115生活”是一款面向个人用户的数字生活平台，提供海量数据的安全存储，多端同步与快速访问。用户不仅可以便捷地管理和使用各类数字资源，还能尽享多维社交、生活服务等多元化功能，畅享高品质的数字生活。
**二、115生活开放平台能力说明**
115生活开放平台功能提供“115生活”数据存储、同步、管理等功能的API服务，开发者通过对接 API 的方式，可以方便快捷地将“115生活”云存储能力集成到开发者的应用中。
目前，已开放的能力包括：
**•用户管理能力：**用户授权与信息查询等；
**•文件管理能力：**获取文件列表，查看文件详情，文件上传、下载、搜索、移动、删除等；
**•视频管理能力：**视频文件的在线转码与播放等；
**•云下载服务：**获取云下载任务列表、配额信息，下载任务添加、删除等；
**•商业价值转化：**推出“推广产品得收益”计划，开发者可基于用户实际购买的产品获取相应推广收益。

--

---

### 更新记录

| 更新模块 | 更新内容 | 更新时间 |
| --- | --- | --- |
| 增值服务产品 | 新增增值服务产品，获得推广收益 | 2025年4月17日周四 |
| 视频播放、云下载 | 新增视频播放、云下载接口 | 2025年4月3日周四 |
| 接入授权 | 支持H5账号密码/短信授权 | 2025年4月2日周三 |
| 基本框架 | 新增开放平台文档接口 | 2025年1月22日周三 |

---



### 获取产品列表地址

### 基本信息
**Path**：GET 域名 + /open/vip/qr_url
**Method：** GET
**接口描述：**获取开放平台产品列表链接地址

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

#### **Query**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| default_product_id | 1 | 否 | 打开产品列表默认选中的产品对应的产品id，如果没有则使用默认的产品顺序。月费：5; 年费：1； 尝鲜1天：101 长期VIP(长期): 24072401 |
| open_device | ​ | 是 | 设备号 |
| hide_title | 0 | 否 | 是否隐藏购买页推荐人信息0：不隐藏（默认）1：隐藏 |

### 返回数据

| 名称 | 类型 | 备注 |
| --- | --- | --- |
| state | boolean | 接口状态；true:正常；false：异常 |
| message | string | 异常信息 |
| code | string | 异常码 |
| data | object | 数据 |
| └─ qrcode_url | string | 开放平台产品列表地址 |

### 返回示例

## 3.场景触发与收益
接入成功后，以下场景将自动触发：
• 用户对应功能使用权限不足时，引导用户升级至更高的VIP类型；
• 用户云存储空间不足时，引导用户购买VIP服务或单独购买空间扩容套餐；
• 用户云下载配额不足时，引导用户升级至更高的VIP类型或单独购买云下载配额。
用户完成购买后，开发者即可获得对应的推广收益。



## API 接口参考

### 授权管理

#### 手机扫码授权PKCE模式

此模式适用于无后端服务的第三方客户端，使用 OAuth 2.0 + PKCE 模式授权，无需提供 AppSecret。

## 1.获取设备码和二维码内容
使用接口返回的 data.qrcode 作为二维码的内容，在第三方客户端展示二维码，用于给115客户端扫码授权。

### 基本信息
**接口名称：**
设备码方式授权

**接口路径：**
POST https://passportapi.115.com/open/authDeviceCode

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| Content-Type | application/x-www-form-urlencoded | 是 |  |  |

**Body**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| client_id | text | 是 |  | APP ID |
| code_challenge | text | 是 | THHodGWg-FZfv8XYz7QArNGIK_aVomSHPldlSOTUtkw | PKCE 相关参数，计算如下：$code_verifier = <43~128为随机字符串>;$code_challenge = url_safe(base64_encode(sha256($code_verifier)));注意 hash 的结果是二进制格式 |
| code_challenge_method | ​ | 是 | sha256 | 计算 code_challenge 的 hash算法，支持 md5, sha1, sha256 |

​

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | number | 非必须 |  |  |  |
| code | number | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| data | object | 非必须 |  |  |  |
| ├─ uid | string | 非必须 |  | 设备码 |  |
| ├─ time | number | 非必须 |  | 校验用的时间戳，轮询设备码状态用到 |  |
| ├─ qrcode | string | 非必须 |  | 二维码内容，第三方客户端需要根据此内容生成设备二维码，提供给115客户端扫码 |  |
| ├─ sign | string | 非必须 |  | 校验用的签名，轮询设备码状态用到 |  |
| error | string | 非必须 |  |  |  |
| errno | number | 非必须 |  |  |  |

​

## 2.轮询二维码状态
此为长轮询接口。注意：当二维码状态没有更新时，此接口不会立即响应，直到接口超时或者二维码状态有更新。

### 基本信息
**接口名称：**
轮询二维码状态

**接口路径：**
GET https://qrcodeapi.115.com/get/status/

### 请求参数
**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| uid | 是 |  | 二维码ID/设备码，从 /open/authDeviceCode 接口 data.uid 获取 |
| time | 是 |  | 校验参数，从 /open/authDeviceCode 接口 data.time 获取 |
| sign | 是 |  | 校验签名，从 /open/authDeviceCode 接口 data.sign 获取 |

​

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | number | 非必须 |  | 0.二维码无效，结束轮询；1.继续轮询； |  |
| code | number | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| data | object | 非必须 |  | 115客户端扫码或者输入设备码后才有值 |  |
| ├─ msg | string | 非必须 |  | 操作提示 |  |
| ├─ status | number | 非必须 |  | 二维码状态；1.扫码成功，等待确认；2.确认登录/授权，结束轮询; |  |
| ├─ version | string | 非必须 |  |  |  |

​

## 3.获取 access_token

### 基本信息
**接口名称：**
用设备码换取 access_token
​

**接口路径：**
POST https://passportapi.115.com/open/deviceCodeToToken

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| Content-Type | application/x-www-form-urlencoded | 是 |  |  |

**Body**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| uid | text | 是 |  | 二维码ID/设备码 |
| code_verifier | text | 是 | IGKN6CJanWxCDPDhHZJrhswQdlcPBGLqExkhyujysXaQ4fJKBk_6dlPJo47s | 上一步计算 code_challenge 的原值 code_verifier |

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | number | 非必须 |  |  |  |
| code | number | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| data | object | 非必须 |  |  |  |
| ├─ access_token | string | 非必须 |  | 用于访问资源接口的凭证 |  |
| ├─ refresh_token | string | 非必须 |  | 用于刷新 access_token，有效期1年 |  |
| ├─ expires_in | number | 非必须 |  | access_token 有效期，单位秒 | mock: 7200 |
| error | string | 非必须 |  |  |  |
| errno | number | 非必须 |  |  |  |

## 更新记录

| 更新时间 | 更新内容 |
| --- | --- |
| 2025年4月7日周一 | 接口 /open/authDeviceCode code_challenge 参数兼容调整，兼容 url safe |

---

#### 授权码模式

本模式建议开发者服务端参与授权

​

基础域名：https://passportapi.115.com

## 授权码方式请求授权

### 基本信息
**Path：** /open/authorize
**Method：** GET
**接口描述：**
用户未登录情况下，会重定向到登录页面。
用户在已登录情况下，会自动完成授权并重定向到 redirect_uri 指定的地址

### 请求参数
**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| client_id | 是 |  | APP ID |
| redirect_uri | 是 | https://foo.com/bar | 授权成功后重定向到指定的地址(记得要urlencode处理一下)并附上授权码 code，如果本接口有传 state 参数，会附带到重定向地址去。重定向地址的域名，需要先到 https://open.115.com/ 应用管理应用域名设置 |
| response_type | 是 | code | 授权模式，固定 为code, 表示授权码模式 |
| state | 否 | 123456 | 随机值，会通过 redirect_uri 原样返回，防止CSRF攻击，强烈建议开发者带上并在请求 换取access_token 接口之前验证 state 是否一致 |

### 返回数据（接口调用失败时候返回）

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | number | 非必须 |  | 0.失败；1.成功 |  |
| code | number | 非必须 |  |  |  |
| errno | number | 非必须 |  |  |  |
| data | object | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| error | string | 非必须 |  |  |  |

## 用授权码换取 access_token

### 基本信息
**Path：** /open/authCodeToToken
**Method：** POST
**接口描述：**
强烈建议在**服务端**请求本接口，防止 APP Secret 泄露！

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| Content-Type | application/x-www-form-urlencoded | 是 |  |  |

 **Body**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| client_id | text | 是 |  | APP ID |
| client_secret | text | 是 |  | APP Secret |
| code | text | 是 |  | 授权码，/open/authCodeToToken 重定向地址里面 |
| redirect_uri | text | 是 | https://foo.com?state=123456 | 与 /open/authCodeToToken 传的 redirect_uri 一致，防 MITM, CSRF |
| grant_type | text | 是 | authorization_code | 授权类型，固定 为 authorization_code, 表示授权码类型 |

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | number | 非必须 |  | 0.失败；1.成功 |  |
| code | number | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| data | object | 非必须 |  |  |  |
| ├─ access_token | string | 非必须 |  | 用于访问资源接口的凭证 |  |
| ├─ refresh_token | string | 非必须 |  | 用于刷新 access_token，有效期1年 |  |
| ├─ expires_in | number | 非必须 |  | access_token 有效期，单位秒 | mock: 7200 |
| error | string | 非必须 |  |  |  |
| errno | number | 非必须 |  |  |  |

---

#### 刷新 access_token

## 刷新 access_token

### 基本信息
**接口名称：**
刷新 access_token

**接口路径：**
POST https://passportapi.115.com/open/refreshToken

### 备注
请勿频繁刷新，否则列入频控。

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| Content-Type | application/x-www-form-urlencoded | 是 |  |  |

**Body**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| refresh_token | text | 是 |  | 刷新的凭证 |

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | number | 非必须 |  |  |  |
| code | number | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| data | object | 非必须 |  |  |  |
| ├─ access_token | string | 非必须 |  | 新的 access_token，同时刷新有效期 |  |
| ├─ refresh_token | string | 非必须 |  | 新的 refresh_token，有效期不延长不改变 |  |
| ├─ expires_in | number | 非必须 |  | access_token 有效期，单位秒 | mock: 2592000 |
| error | string | 非必须 |  |  |  |
| errno | number | 非必须 |  |  |  |

---

### 用户管理

#### 用户信息

API：GET 域名 + /open/user/info
说明：获取用户空间和vip信息
**请求参数**
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**返回数据**

| state | boolean | 接口状态；true正常；false：异常 |
| --- | --- | --- |
| message | string | 异常信息 |
| code | string | 异常码 |
| data | object | 数据 |
| ├─user_id | string | 用户ID |
| ├─user_name | string | 用户名称 |
| ├─user_face_s | string | 小尺寸用户头像 |
| ├─user_face_m | string | 中尺寸用户头像 |
| ├─user_face_l | string | 大尺寸用户头像 |
| ├─rt_space_info | object | 用户空间信息 |
| ├─all_total | object | 用户总空间 |
| ├─size | int | 用户总空间大小(字节) |
| ├─size_format | string | 用户总空间大小(格式化) |
| ├─all_remain | object | 用户剩余空间 |
| ├─size | int | 用户剩余空间大小(字节) |
| ├─size_format | string | 用户剩余空间大小(格式化) |
| ├─all_use | object | 用户已使用空间 |
| ├─size | int | 用户已使用空间大小(字节) |
| ├─size_format | string | 用户已使用空间大小(格式化) |
| ├─vip_info | object | 用户vip等级信息 |
| ├─level_name | string | vip等级名称；原石会员、尝鲜VIP、体验VIP、月费VIP、年费VIP、年费VIP高级版、年费VIP特级版、超级VIP、长期VIP； |
| ├─expire | int | 过期时间戳 |

---

### 文件管理-文件上传

#### 上传流程

1，请求【**文件上传**】接口初始化上传文件，接口返回的status=2时，表示秒传成功，上传结束。非秒传status=1的情况将会返回预上传callback**。，**
2，请求【**文件上传**】接口初始化上传文件的时候提示需要二次认证，需要做二次认证，二次认证成功的时候返回相应的结果
3，非秒传时，携带步骤1返回的callback以及其他参数和请求[**获取上传凭证**]接口返回的数据调用【**对象存储**】开始上传文件。
4，断点续传，携带步骤1返回的*pickcode*和要上传文件信息，请求【**断点续传**】接口获取返回的callback以及其他参数和**上传凭证数据**调用【**对象存储**】开始上传文件。
5，【**对象存储**】返回上传成功，这时非秒传和断点续传文件上传成功，详见对象存储API文档
​

---

#### 获取上传凭证

### 基本信息
**Path：** GET 域名 + /open/upload/get_token
**Method：** GET
**接口描述：获取上传凭证**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean | 状态；true：正常；false：错误 |  |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | object [] | ​ | item 类型: object |
| ├─ endpoint | string | 上传域名 |  |
| ├─ AccessKeySecrett | string | 上传凭证-密钥 |  |
| ├─ SecurityToken | string | 上传凭证-token |  |
| ├─ Expiration | string | 上传凭证-过期日期 |  |
| ├─ AccessKeyId | string | 上传凭证-ID |  |

---

#### 文件上传

### 基本信息
**Path：** POST 域名 + /open/upload/init
**Method：** POST
**接口描述：断点续传上传初始化调度接口**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| file_name | 是 | 图片.jpg | 文件名 |
| file_size | 是 | 5335 | 文件大小(字节) |
| target | 是 | U_1_0 | 文件上传目标约定 U_1_0 其中U_1_是固定约定，0代表根目录，其他数字是文件夹ID |
| fileid | 是 |  | 文件sha1值 |
| preid | 否 |  | 文件前128Ksha1 |
| pick_code | 否 |  | 上传任务key[非秒传的调度接口返回的pick_code字段] |
| topupload | 否 | 0 | 上传调度文件类型调度标记【0：单文件上传任务标识一条单独的文件上传记录 1：文件夹任务调度的第一个子文件上传请求标识一次文件夹上传记录。 2：文件夹任务调度的其余后续子文件不作记作单独上传的上传记录 -1：没有该参数】 |
| sign_key | 否 | ​ | 二次认证需要 |
| sign_val | 否 | ​ | 二次认证需要(大写) |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean | 状态；true：正常；false：错误 |  |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | object [] | ​ | item 类型: object |
| ├─ pick_code | string | 上传任务唯一ID,用于续传 |  |
| ├─ status | number | 上传状态；1：非秒传；2：秒传 |  |
| ​├─ sign_key | string | 本次计算的sha1标识（二次认证 ） | 参照下面【二次认证】说明 |
| ├─ sign_check | string | 本次计算本地文件sha1区间范围（二次认证) | 参照下面【二次认证】说明 |
| ├─ file_id | string | 秒传成功返回的新增文件ID | ​ |
| ├─ target | string | 文件上传目标约定 |  |
| ├─ bucket | string | 上传的bucket |  |
| ├─ object | string | OSS objectID |  |
| ├─ callback | object [] | 上传时间 | tem 类型: object |
| ├─callback | string | 上传完回调信息 |  |
| ├─callback_var | string | 上传完回调参数 |  |

### 二次认证

| code | status | 说明 | 后续处理 | 第二次验证 |
| --- | --- | --- | --- | --- |
| 700 | 6 | 签名认证后失败 | sign_check （用下划线隔开,截取上传文内容的sha1）(单位是byte)"sign_check": "2392148-2392298"取2392148-2392298之间的内容(包含2392148、2392298)的sha1例如：获取0-99字节（包含0和99）共100个字节 | sign_keysign_val (根据sign_check计算的值大写的sha1值) |
| 701 | 7 | 需要认证签名 |  |  |
| 702 | 8 | 签名认证失败 |  |  |

​

---

#### 断点续传

### 基本信息
**Path：** POST 域名 + /open/upload/resume
**Method：** POST
**接口描述：断点续传上传续传调度接口**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| file_size | 是 | 5335 | 文件大小(字节) |
| target | 是 | U_1_0 | 文件上传目标约定 |
| fileid | 是 |  | 文件sha1值 |
| pick_code | 是 |  | 上传任务key[非秒传的调度接口返回的pick_code字段] |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean | 状态；true：正常；false：错误 |  |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | object [] | ​ | item 类型: object |
| ├─ pick_code | string | 上传任务唯一ID,用于续传 |  |
| ├─ target | string | 文件上传目标约定 |  |
| ├─ version | string | 接口版本 |  |
| ├─ bucket | string | 上传的bucket |  |
| ├─ object | string | OSS objectID |  |
| ├─ callback | object [] | 上传时间 | tem 类型: object |
| ├─callback | string | 上传完回调信息 |  |
| ├─callback_var | string | 上传完回调参数 |  |

---

### 文件管理-文件操作

#### 新建文件夹

### 基本信息
**Path：** POST 域名 + /open/folder/add
**Method：** POST
**接口描述：新建文件夹**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**​**

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pid | text | 是 | 3073323042143855813 | 新建文件夹所在的父目录ID (根目录的ID为0) |
| file_name | text | 是 | 新建文件夹名称 | 新建文件夹名称，限制255个字符 |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean | 接口状态；true正常；false：异常 |  |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | object |  |  |
| ├─ file_name | string | 新建的文件夹名称 |  |
| ├─ file_id | string | 新建的文件夹ID |  |

---

#### 获取文件列表

### 基本信息
**Path：** GET 域名 + /open/ufile/files
**Method：** GET
**接口描述：**获取文件列表

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| cid | 否 |  | 目录ID，对应parent_id |
| type | 否 |  | 文件类型；1.文档；2.图片；3.音乐；4.视频；5.压缩；6.应用；7.书籍 |
| limit | 否 |  | 查询数量，默认20，最大1150 |
| offset | 否 |  | 查询起始位，默认0 |
| suffix | 否 |  | 文件后缀名 |
| asc | 否 |  | 排序，1：升序 0：降序 |
| o | 否 |  | 排序字段，file_name：文件名 file_size：文件大小 user_utime：更新时间 file_type 文件类型 |
| custom_order | 否 |  | 是否使用记忆排序。1 使用自定义排序，不使用记忆排序,0 使用记忆排序，自定义排序失效,2自定义排序，非文件夹置顶 |
| stdir | 否 |  | 筛选文件时， 是否显示文件夹；1:要展示文件夹 0不展示 |
| star | 否 |  | 筛选星标文件，1:是 0全部 |
| cur | 否 |  | 是否只显示当前文件夹内文件 |
| show_dir | 否 |  | 是否显示目录；0 或 1，默认为0 |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| data | object [] |  | item 类型: object |
| ├─ fid | string | 文件ID |  |
| ├─ aid | string | 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除 |  |
| ├─ pid | string | 父目录ID |  |
| ├─ fc | string | 文件分类。0 文件夹，1 文件 |  |
| ├─ fn | string | 文件(夹)名称 |  |
| ├─ fco | string | 文件夹封面 |  |
| ├─ ism | string | 是否星标，1：星标 |  |
| ├─ isp | number | 是否加密；1：加密 |  |
| ├─ pc | string | 文件提取码 |  |
| ├─ upt | int | 修改时间 |  |
| ├─ uet | int | 修改时间 |  |
| ├─ uppt | int | 上传时间 |  |
| ├─ cm | number |  |  |
| ├─ fdesc | string | 文件备注 |  |
| ├─ ispl | number | 是否统计文件夹下视频时长开关  |  |
| ├─ fl | object [] | 文件标签 | item 类型: object |
| ├─ id | string | 文件标签id |  |
| ├─ name | string | 文件标签名称 |  |
| ├─ sort | string | 文件标签排序 |  |
| ├─ color | string | 文件标签颜色 |  |
| ├─ is_default | int | 文件标签类型；0：最近使用；1：非最近使用；2：为默认标签 |  |
| ├─ update_time | int | 文件标签更新时间 |  |
| ├─ create_time | int | 文件标签创建时间 |  |
| ├─ sha1 | string | sha1值 |  |
| ├─ fs | int | 文件大小 |  |
| ├─ fta | string | 文件状态 0/2 未上传完成，1 已上传完成  |  |
| ├─ ico | string | 文件后缀名 |  |
| ├─ fatr | string | 音频长度 |  |
| ├─ isv | number | 是否为视频 |  |
| ├─ def | number | 视频清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画 |  |
| ├─ def2 | number | 视频清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画 |  |
| ├─ play_long | number | 音视频时长 |  |
| ├─ v_img | string |  |  |
| ├─ thumb | string | 图片缩略图 |  |
| ├─ uo | string | 原图地址 |  |
| count | number | 当前目录文件数量 |  |
| sys_count | number | 系统文件夹数量 |  |
| offset | number | 偏移量 |  |
| limit | number | 分页量 |  |
| aid | string | 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除 |  |
| cid | number | 父目录id |  |
| is_asc | number | 排序，1：升序 0：降序 |  |
| min_size | number |  |  |
| max_size | number |  |  |
| sys_dir | string |  |  |
| hide_data | string | 是否返回文件数据 |  |
| record_open_time | string | 是否记录文件夹的打开时间 |  |
| star | number | 是否星标；1：星标；0：未星标 |  |
| type | number | 一级筛选大分类，1：文档，2：图片，3：音乐，4：视频，5：压缩包，6：应用 |  |
| suffix | string | 一级筛选选其他时填写的后缀名 |  |
| path | object [] | 父目录树 | item 类型: object |
| ├─ name | string | 父目录文件名称 |  |
| ├─ aid | number,string |  |  |
| ├─ cid | number,string |  |  |
| ├─ pid | number,string |  |  |
| ├─ isp | number,string |  |  |
| ├─ p_cid | string |  |  |
| ├─ fv | string |  |  |
| cur | number |  |  |
| stdir | number |  |  |
| fields | string |  |  |
| order | string | 排序 |  |
| state | boolean |  |  |
| code | number |  |  |
| message | number |  |  |

## 更新记录

| 更新时间 | 更新内容 |
| --- | --- |
| 2025年4月3日周四 | 新增补充fl(文件标签)字段和其二级结构详情 |

---

#### 获取文件(夹)详情

### 基本信息
**Path：** GET 域名 + /open/folder/get_info
**Method：** GET/POST
**接口描述：**获取文件(夹)详情

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| file_id | 否 | 1288444975268439877 | 文件(夹)ID。和path需必传一个 |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| path  | string | 否 | /a/b/c.png 或>a>b>c | 文件路径；分隔符支持 / > 两种符号，最前面需分隔符开头，以分隔符分隔目录层级；和file_id必传一个 |

### 返回数据

| 名称 | 类型 | 备注 |
| --- | --- | --- |
| state | boolean |  |
| message | string |  |
| code | number |  |
| data | object |  |
| ├─ count | string | 包含文件总数量 |
| ├─ size | string | 文件(夹)总大小 |
| ├─ size_byte  | int | 文件(夹)总大小(字节单位) |
| ├─folder_count | string | 包含文件夹总数量 |
| ├─ play_long | number | 视频时长；-1：正在统计，其他数值为视频时长的数值(单位秒) |
| ├─show_play_long | number | 是否开启展示视频时长 |
| ├─ ptime | string | 上传时间 |
| ├─ utime | string | 修改时间 |
| ├─ file_name | string | 文件名 |
| ├─pick_code | string | 文件提取码 |
| ├─ sha1 | string | sh1值 |
| ├─ file_id | string | 文件(夹)ID |
| ├─ is_mark | string | 是否星标 |
| ├─open_time | number | 文件(夹)最近打开时间 |
| ├─file_category | string | 文件属性；1；文件；0：文件夹 |
| ├─ paths | object [] | 文件(夹)所在的路径 |
| ├─ file_id | number | 父目录ID |
| ├─file_name | string | 父目录名称 |

## 更新记录

| 更新时间 | 更新内容 |
| --- | --- |
| 2025年6月6日周五 | 新增根据文件路径获取文件(夹)详情、新增size_byte字段返回 |

---

#### 文件搜索

### 基本信息
**Path：** GET 域名 + /open/ufile/search
**Method：** GET
**接口描述：根据文件名搜索文件(夹)**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| search_value | 是 |  | 查找关键字 |
| limit | 是 |  | 单页记录数，默认20，offset+limit最大不超过10000 |
| offset | 是 |  | 数据显示偏移量 |
| file_label | 否 |  | 支持文件标签搜索 |
| cid | 否 |  | 目标目录cid=-1时，表示不返回列表任何内容 |
| gte_day | 否 |  | 搜索结果匹配的开始时间；格式：2020-11-19 |
| lte_day | 否 |  | 搜索结果匹配的结束时间；格式：2020-11-20 |
| fc | 否 |  | 只显示文件或文件夹。1 只显示文件夹，2 只显示文件 |
| type | 否 |  | 一级筛选大分类，1：文档，2：图片，3：音乐，4：视频，5：压缩包，6：应用 |
| suffix | 否 |  | 一级筛选选其他时填写的后缀名 |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| count | number | 搜索符合条件的文件(夹)总数 |  |
| data | object [] |  | item 类型: object |
| ├─ file_id | string | 文件ID |  |
| ├─ user_id | string | 用户ID |  |
| ├─ sha1 | string | 文件sha1值 |  |
| ├─file_name | string | 文件名称 |  |
| ├─file_size | string | 文件大小 |  |
| ├─user_ptime | string | 上传时间 |  |
| ├─user_utime | string | 更新时间 |  |
| ├─pick_code | string | 文件提取码 |  |
| ├─parent_id | string | 父目录ID |  |
| ├─area_id | string | 文件的状态，aid 的别名。1 正常，7 删除(回收站)，120 彻底删除 |  |
| ├─is_private | number | 文件是否隐藏。0 未隐藏，1 已隐藏 |  |
| ├─file_category | string | 1：文件；0；文件夹 |  |
| ├─ ico | string | 文件后缀 |  |
| limit | number | 分页获取值 |  |
| offset | number | 偏移值 |  |
| state | boolean | 状态；true：正常；false：错误 |  |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |

---

#### 文件复制

### 基本信息
**Path：** POST 域名 + /open/ufile/copy
**Method：** POST
**接口描述：批量复制文件**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pid | text | 是 | 1054251402869818368 | 目标目录，即所需移动到的目录 |
| file_id | text | 是 | 2323423573680609857 | 所复制的文件和目录ID，多个文件和目录请以 , 隔开 |
| nodupli | text | 否 | 1 | 复制的文件在目标目录是否允许重复，默认0：0：可以；1：不可以 |

### 返回数据

| state | boolean | 接口状态；true正常；false：异常 |  |
| --- | --- | --- | --- |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | array | 数据 |  |

---

#### 文件移动

### 基本信息
**Path：** POST 域名 + /open/ufile/move
**Method：** POST
**接口描述：批量移动文件**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| file_ids | string | 是 | 3073323042143855813,3073323042143855822 | 需要移动的文件(夹)ID |
| to_cid | string | 是 | 3073311192547189943 | 要移动所在的目录ID，根目录为0 |

### 返回数据

| state | boolean | 接口状态；true正常；false：异常 | 备注 |
| --- | --- | --- | --- |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | array | 数据 |  |

---

#### 获取文件下载地址

### 基本信息
**Path：** POST 域名 + /open/ufile/downurl
**Method：** POST
**接口描述：**根据文件提取码取文件下载地址

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 |  |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pick_code | text | 是 | dtctprlmfkl4exiok | 文件提取码 |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean |  |  |
| message | string |  |  |
| code | number |  |  |
| data | object |  |  |
| ├─2323423573680609857 | object | 文件ID |  |
| ├─ file_name | string | 文件名 |  |
| ├─ file_size | number | 文件大小 |  |
| ├─ pick_code | string | 文件提取码 |  |
| ├─ sha1 | string | 文件sha1值 |  |
| ├─ url | object |  |  |
| ├─ url | string | 文件下载地址 |  |

---

#### 文件(夹)更新

### 基本信息
**Path：** POST 域名 + /open/ufile/update
**Method：** POST
**接口描述：更新文件名或星标文件**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| file_id | text | 是 | 3073323042143855813 | 需要更改名字的文件(夹)ID |
| file_name | text | 否 | 新的名字 | 新的文件(夹)名字(文件夹名称限制255字节) |
| star | text | 否 | 1 | 是否星标；1：星标；0；取消星标 |

### 返回数据

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean | 状态；true：正常；false：异常 |  |
| message | string | 异常信息 |  |
| code | number | 异常码 |  |
| data | object |  |  |
| ├─ file_name | string | 新的文件(夹)名字 |  |
| ├─ star | string | 文件星标状态 |  |

---

#### 删除文件

### 基本信息
**Path：** POST 域名 + /open/ufile/delete
**Method：** POST
**接口描述：批量删除文件(夹)**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 备注示例 |  |
| --- | --- | --- | --- | --- |
| file_ids | text | 是 | 3073323042143855813,3073323042143855822 | 需要删除的文件(夹)ID |
| parent_id | text | 否 | 3073311192547189943 | 删除的文件(夹)ID所在的父目录ID |

### 返回数据

| 名称 | 类型 | 备注 |
| --- | --- | --- |
| state | boolean | 接口状态；true正常；false：异常 |
| message | string | 异常信息 |
| code | string | 异常码 |
| data | string [] | 数据 |

---

### 文件管理-回收站

#### 回收站列表

### 基本信息
**Path：** GET 域名 + /open/rb/list
**Method：** GET
**接口描述：回收站列表**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| limit | 是 | 30 | 单页记录数，int，默认30，最大200 |
| offset | 是 | 0 | 数据显示偏移量 |

### 返回数据

| 名称 | 类型 | 备注 |
| --- | --- | --- |
| state | boolean |  |
| message | string |  |
| code | number |  |
| data | object |  |
| ├─ offset | number | 偏移量 |
| ├─ limit | number | 分页量 |
| ├─ count | string | 回收站文件总数 |
| ├─ rb_pass | number | 是否设置回收站密码 |
| ├─3074054555277845747 | object | 文件(夹)回收站ID |
| ├─ id | string | 文件(夹)回收站ID |
| ├─ file_name | string | 文件(夹)名称 |
| ├─ type | string | 类型（1：文件，2：目录 |
| ├─ file_size | string | 文件大小 |
| ├─ dtime | string | 删除日期 |
| ├─ thumb_url | string | 缩略图地址 |
| ├─ status | string | 还原状态，-1 表示还原中，0 表示正常状态 |
| ├─ cid | number | 原文件的父目录id |
| ├─ parent_name | string | 原文件的父目录名称 |
| ├─ pick_code | string | 文件提取码 |

---

#### 回收站还原

### 基本信息
**Path：** POST 域名 + /open/rb/revert
**Method：** POST
**接口描述：还原回收站文件(夹)**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| tid | text | 是 | 111,222,333,444 | 需要还原的ID，可多个，用半角逗号分开，最多1150个 |

### **返回数据**

| 名称 | 类型 | 备注 |
| --- | --- | --- |
| data | object |  |
| ├─ 2064439620648589453 | object | 还原的回收站ID |
| ├─ state | boolean |  |
| ├─ error | string |  |
| ├─ errno | number |  |
| state | boolean |  |
| message | string |  |
| code | number |  |

**

**

---

#### 删除/清空回收站

### 基本信息
**Path：** POST 域名 + /open/rb/del
**Method：** POST
**接口描述：批量删除回收站文件、清空回收站**

### 请求参数
**Headers**

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| tid | text | 否 | 1,2.3.4 | 需要删除的文件的Id,如若不传就是清空回收站(最多支持1150个) |

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | boolean | 非必须 |  | 接口状态；true正常；false：异常 |  |
| message | string | 非必须 |  | 异常信息 |  |
| code | numnber | 非必须 |  | 异常码 |  |
| data | string [] | 非必须 |  | 数据 | ​ |

---

### 视频播放

#### 记忆视频播放进度

### 基本信息
Path：域名 + /open/video/history 

Method：POST 

接口描述：记录视频播放进度 

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### ​**Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pick_code | string | 是 | b53gu6z3hvqji8wrm | 文件提取码 |
| time | int | 否 | 10 | 视频播放进度时长(单位秒) |
| watch_end | int | 否 | 0 | 视频是否播放播放完毕；默认为：0；1：播放完毕 |

### 注意事项
time和watch_end都不传的时候，默认都是为0；即视频没有播放完毕，视频播放进度为0，
返回结果

### 返回数据说明

| 字段 | 类型及范围 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功时返回0 |
| data | array | 数据 |

---

#### 视频字幕列表

### 基本信息
**Path：** 域名 + /open/video/subtitle
**Method：** GET
**接口描述：**

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### 请求参数
**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| pick_code | 是 | b53gu6z3hvqji8wrm | 视频文件提取码 |

### 返回数据

| 名称 | 类型 | 是否必须 | 默认值 | 备注 | 其他信息 |
| --- | --- | --- | --- | --- | --- |
| state | boolean | 非必须 |  |  |  |
| message | string | 非必须 |  |  |  |
| code | number | 非必须 |  |  |  |
| data | object | 非必须 |  |  |  |
| ├─ autoload | object | 非必须 |  | 自动载入字幕列表 |  |
| ├─ sid | string | 非必须 |  |  |  |
| ├─ language | string | 非必须 |  |  |  |
| ├─ title | string | 非必须 |  |  |  |
| ├─ url | string | 非必须 |  |  |  |
| └─ type | string | 非必须 |  |  |  |
| ├─ list | object [] | 非必须 |  | 字幕列表 | item 类型: object |
| ├─ sid | string | 必须 |  |  |  |
| ├─ language | string | 必须 |  |  |  |
| ├─ title | string | 必须 |  | 字幕标题 |  |
| ├─ url | string | 必须 |  | 字幕文件地址 |  |
| ├─ type | string | 必须 |  | 字幕文件类型 |  |
| ├─ sha1 | string | 必须 |  | 字幕文件哈希值 |  |
| ├─ file_id | string | 必须 |  | 字幕文件id |  |
| ├─ file_name | string | 必须 |  | 外挂字幕文件名 |  |
| ├─ pick_code | string | 必须 |  | 外挂字幕文件提取码 |  |
| ├─ caption_map_id | string | 必须 |  |  |  |
| ├─ is_caption_map | number | 必须 |  |  |  |
| ├─ sync_time | number | 必须 |  | 字幕同步时间 |  |

---

#### 获取视频播放进度

### 基本信息
Path：域名 + /open/video/history 

Method：GET

接口描述：获取视频播放进度 

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pick_code | string | 是 | b53gu6z3hvqji8wrm | 文件提取码 |

### 注意事项
视频是否播放完毕在文件列表的played_end展示,1为播放完毕；

### 返回数据说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功时返回0 |
| data | array | 数据 |
| ├─add_time | int | 记录添加时间 |
| ├─file_id | string | 文件id |
| ├─file_name | string | 文件名称 |
| ├─hash | string | 文件哈希值 |
| ├─pick_code | string | 文件提取码 |
| ├─time | string | 记录的已播放时长 |

---

#### 获取视频在线播放地址

### 基本信息
Path：域名 + /open/video/play 

Method：GET

接口描述：获取视频播放地址和视频文件相关数据 

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pick_code | string | 是 | b53gu6z3hvqji8wrm | 文件提取码 |

### 注意事项
1、需切换音轨时，在请求返回的播放地址中增加参数 audio_track = 1(整型)，值就是multitrack_list的key，value是音轨名称，下标从0开始
例如：http://videoplay.115rc.com:80/m3u8/cw2eh308u3n0xdblc?filesha1=35DD3E560F0C613D042CF05AD58E7C1A798CCD78&userid=110003887&rsa=a3680f3650348a2283df0f3a1a19dae1&definition=4&prefix=ANDROID&time=1742469152&pickcode=cw2eh308u3n0xdblc&**audio_track=0**

### 返回数据说明

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | 异常信息，成功时空值 |  |
| code | int | 异常码，成功时返回0 |
| data | object | 数据 |
| ├─file_id | string | 文件id |
| ├─parent_id | string | 文件父目录id |
| ├─file_name | string | 文件名称 |
| ├─file_size | string | 文件大小 |
| ├─file_sha1 | string | 文件哈希值 |
| ├─file_type | string | 文件类型 |
| ├─is_private | string | 文件是否加密隐藏；0：否；1：是 |
| ├─play_long | string | 视频文件时长 |
| ├─ user_def | int | 视频文件记忆选中的清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画 |
| ├─ user_rotate | int | 记忆视频旋转角度；0, 90, 180, 270 |
| ├─ user_turn | int | 视频翻转方向：0：不翻转；1：水平翻转；2：垂直翻转 |
| ├─ multitrack_list | array | 视频多音轨列表 |
| ├─ title | string | 音轨标题 |
| ├─ is_selected | string | 音轨是否上次选中；1：选中 |
| ├─ definition_list_new | array | 视频所有用可切换的清晰度列表;1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画 |
| ├─ video_url | array | 视频各清晰度的播放地址信息 |
| ├─ url | string | 播放地址 |
| ├─ height | int | 视频高度 |
| ├─ width | int | 视频宽度 |
| ├─ definition | int | 视频清晰度 |
| ├─ title | int | 视频清晰度名称 |
| ├─ definition_n | int | 视频清晰度(新) |

## 更新记录

| 更新时间 | 更新内容 |
| --- | --- |
| 2025年6月16日周周一 | 去掉share_id 请求参数，暂不支持 |
| 2025年4月17日周四 | 年费VIP以下用户不支持播放4K视频，当切换4K清晰度时会播放引导升级的视频开发者可接入增值服务引导升级VIP |
| 2025年4月9日周三 | 去除down_url 字段返回 |

---

#### 提交视频转码

### 基本信息
Path 域名 + /open/video/video_push 

Method POST 

接口描述：提交视频通过 vip等级/枫叶 来加速转

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| pick_code | string | 是 | b53gu6z3hvqji8wrm | 文件提取码 |
| op | string | 是 | vip_push | 提交视频加速转码方式；vip_push：根据vip等级加速；pay_push：枫叶加速 |

### 返回数据说明

| 字段 | 类型及范围 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功时返回0 |
| data | array | 数据 |

---

### 云下载

#### 解析BT种子

### 基本信息
Path 域名 + /open/offline/torrent

Method：POST

接口描述：解析BT种子

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 备注 |
| --- | --- | --- | --- |
| torrent_sha1 | string | 是 | BT种子文件sha1，需先上传到“云下载/种子文件”文件夹下(非硬性要求) |
| pick_code | string | 是 | BT种子文件提取码 |

### 返回数据说明

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean |  |  |
| message | string |  |  |
| code | number |  |  |
| data | array |  |  |
| ├─ file_size | int | 任务大小 |  |
| ├─torrent_name | string | 任务名 |  |
| ├─ file_count | int | 文件数 |  |
| ├─ info_hash | string | 任务sha1 |  |
| ├─torrent_filelist | array | 文件列表 |  |
| ├─size | int | 文件大小 |  |
| ├─path | string | 文件路径 |  |
| ├─wanted | int | 文件是否默认选中 |  |

---

#### 获取用户云下载任务列表

### 基本信息
Path 域名 + /open/offline/get_task_list 

Method：GET 

接口描述：获取用户云下载任务列表

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### 请求参数
**Query**

| 参数名称 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- |
| page | 否 | 1 | 获取列表第几，如不传默认第一页 |

####注意事项 无

### 返回数据说明

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean |  |  |
| message | string |  |  |
| code | number |  |  |
| data | array |  |  |
| ├─ page | int | 当前第几页 |  |
| ├─ page_count | int | 总页数 |  |
| ├─ count | int | 总数量 |  |
| ├─ tasks | array | 云下载任务列表 |  |
| ├─ info_hash | string | 任务sha1 |  |
| ├─ add_time | int | 任务添加时间戳 |  |
| ├─ percentDone | int | 任务下载进度 |  |
| ├─ size | int | 任务总大小（字节） |  |
| ├─ name | string | 任务名 |  |
| ├─ last_update | int | 任务最后更新时间戳 |  |
| ├─ file_id | string | 任务源文件（夹）对应文件（夹）id |  |
| ├─delete_file_id | string | 删除任务需删除源文件（夹）时，对应需传递的文件（夹）id |  |
| ├─ status | int | 任务状态：-1下载失败；0分配中；1下载中；2下载成功 |  |
| ├─ url | string | 链接任务url |  |
| ├─ wp_path_id | string | 任务源文件所在父文件夹id |  |
| ├─ def2 | int | 视频清晰度；1:标清 2:高清 3:超清 4:1080P 5:4k;100:原画 |  |
| ├─ play_long | int | 视频时长 |  |
| ├─ can_appeal | int | 是否可申诉 |  |

---

#### 获取云下载配额信息

### 基本信息
Path 域名 + /open/offline/get_quota_info

Method：GET 

接口描述：获取当前用户各个配额类型明细

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

####注意事项 无

### 返回数据说明

| 名称 | 类型 | 备注 | 其他信息 |
| --- | --- | --- | --- |
| state | boolean |  |  |
| message | string |  |  |
| code | number |  |  |
| data | array |  |  |
| ├─ package | array | 配额类型 |  |
| ├─ surplus | int | 该类型剩余配额 |  |
| ├─ used | int | 该类型已用配额 |  |
| ├─ count | int | 该类型总配额 |  |
| ├─ name | string | 该类型配额名称 |  |
| ├─ expire_info | array | 该类型明细项过期信息 |  |
| ├─ surplus | int | 明细项剩余配额 |  |
| ├─ expire_time | int | 明细项过期时间 |  |
| ├─ count | int | 用户总配额数量 |  |
| ├─ surplus | int | 用户总剩余配额数量 |  |
| ├─ used | int | 用户总已用配额数量 |  |

---

#### 清空云下载任务

### 基本信息
Path 域名 + /open/offline/clear_task 

Method POST 

接口描述：清空云下载任务

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| flag | int | 是 | 1 | 清除任务类型：0清空已完成、1清空全部、2清空失败、3清空进行中、4清空已完成任务并清空对应源文件、5清空全部任务并清空对应源文件 |

### 返回数据说明

| 字段 | 类型及范围 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功返回0 |
| data | array | 数据 |

---

#### 添加云下载链接任务

### 基本信息
Path 域名 + /open/offline/add_task_urls 

Method POST 

接口描述：添加云下载链接任务,支持多个链接url,换行符分隔，支持HTTP(S)、FTP、磁力链和电驴链接

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| urls | string | 是 |  | 多个链接url,换行符分隔，支持HTTP(S)、FTP、磁力链和电驴链接 |
| wp_path_id | string | 否 | 0 | 保存目标文件夹id；如不传或传0，则默认保存到根目录下面 |

### 返回数据说明

| 字段 | 类型及范围 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功时返回0 |
| data | array | 数据 |
| ├─state | bool | 链接任务添加状态，成功true；失败false |
| ├─code | int | 链接任务状态码，成功返回0 |
| ├─message | string | 链接任务状态描述，成功返回空字符串 |
| ├─info_hash | string | 链接任务sha1，只有任务成功的时候才会返回 |
| ├─url | string | 链接任务ur |

---

#### 删除用户云下载任务

### 基本信息
Path 域名 + /open/offline/del_task 

Method POST 

接口描述：删除用户云下载任务务

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 示例 | 备注 |
| --- | --- | --- | --- | --- |
| info_hash | string | 是 |  | 需删除的任务hash |
| del_source_file | int | 否 |  | 是否删除源文件：1删除；0不删除 |

### 返回数据说明

| 字段 | 类型及范围 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功时返回0 |
| data | array | 数据 |

---

#### 添加云下载BT 任务

### 基本信息
Path 域名 + /open/offline/add_task_bt 

Method POST 

接口描述：添加云下载链接任务

### 请求参数

### head

| 参数名称 | 参数值 | 是否必须 | 示例 |
| --- | --- | --- | --- |
| Authorization | Bearer access_token | 是 | Bearer abcdefghijklmnopqrstuvwxyz |

### **Body(form-data)**

| 参数名称 | 参数类型 | 是否必须 | 备注 |
| --- | --- | --- | --- |
| info_hash | string | 是 | BT任务hash |
| wanted | string | 是 | BT任务选中下载文件索引，半角逗号隔开 |
| save_path | string | 是 | BT任务文件保存路径 |
| torrent_sha1 | string | 是 | BT种子sha1 |
| pick_code | string | 是 | BT种子的提取码 |
| wp_path_id | string | 否 | 保存目标文件夹id |

### 注意事项
wp_path_id 不传默认到根目录下面
save_path 传的是wp_path_id所在文件夹下面的路径
如wp_path_id不传或传云下载的文件夹ID，save_path传 A/B 最终下载的文件的路径为根目录/A/B/

### 返回数据说明

| 字段 | 类型及范围 | 说明 |
| --- | --- | --- |
| state | bool | 操作结果状态值，true成功，false失败 |
| message | string | 操作返回消息，成功时空值 |
| code | int | 操作返回号码，成功时返回0 |
| data | array | 数据 |

---

