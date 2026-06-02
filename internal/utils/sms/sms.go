package sms

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WSSE_HEADER_FORMAT 无需修改,用于格式化鉴权头域,给"X-WSSE"参数赋值
const WSSE_HEADER_FORMAT = "UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\""

// AUTH_HEADER_VALUE 无需修改,用于格式化鉴权头域,给"Authorization"参数赋值
const AUTH_HEADER_VALUE = "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\""
const (
	//必填,请参考"开发准备"获取如下数据,替换为实际值
	apiAddress    = "https://smsapi.cn-south-1.myhuaweicloud.com:443/sms/batchSendSms/v1" //APP接入地址(在控制台"应用管理"页面获取)+接口访问URI
	appKey        = "3sH8I0SWreLVY1nh7A0RL2TUPdi1"                                        //APP_Key
	appSecret     = "Dy20t33SB3NQGlNxtl0kFdBUkgij"                                        //APP_Secret
	sender        = "8823051811871"                                                       //国内短信签名通道号或国际/港澳台短信通道号
	ContractPhone = "13635243974"
)

func sendSms(templateId, receiver, templateParas string) error {
	receiver = "+86" + receiver

	//条件必填,国内短信关注,当templateId指定的模板类型为通用模板时生效且必填,必须是已审核通过的,与模板类型一致的签名名称
	//国际/港澳台短信不用关注该参数
	signature := "一起趣浪" //签名名称

	statusCallBack := ""

	/*
	 * 选填,使用无变量模板时请赋空值 string templateParas = "";
	 * 单变量模板示例:模板内容为"您的验证码是${1}"时,templateParas可填写为"[\"369751\"]"
	 * 双变量模板示例:模板内容为"您有${1}件快递请到${2}领取"时,templateParas可填写为"[\"3\",\"人民公园正门\"]"
	 * 模板中的每个变量都必须赋值，且取值不能为空
	 * 查看更多模板和变量规范:产品介绍>模板和变量规范
	 */
	//	templateParas := "[\"369751\"]" //模板变量，此处以单变量验证码短信为例，请客户自行生成6位验证码，并定义为字符串类型，以杜绝首位0丢失的问题（例如：002569变成了2569）。

	body := buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Authorization"] = AUTH_HEADER_VALUE
	headers["X-WSSE"] = buildWsseHeader(appKey, appSecret)
	body, err := post(apiAddress, []byte(body), headers)
	if err != nil {
		return err
	}
	logx.Info(body)
	return nil
}

// SendOutGroupSms 退团申请成功
func SendOutGroupSms(SendPhone, GroupTitle string) {
	/*
	 * 选填,使用无变量模板时请赋空值 string templateParas = "";
	 * 单变量模板示例:模板内容为"您的验证码是${1}"时,templateParas可填写为"[\"369751\"]"
	 * 双变量模板示例:模板内容为"您有${1}件快递请到${2}领取"时,templateParas可填写为"[\"3\",\"人民公园正门\"]"
	 * 模板中的每个变量都必须赋值，且取值不能为空
	 * 查看更多模板和变量规范:产品介绍>模板和变量规范
	 */
	var templateParas = "[\"" + GroupTitle + "\",\"" + ContractPhone + "\"]"
	err := sendSms("22e85859e03c43ccb59401498a96b166", SendPhone, templateParas)
	if err != nil {
		logx.Error(err)
		return
	}
}

// SendFinalPaySuccess 私家团支付尾款
func SendFinalPaySuccess(SendPhone string) {
	var templateParas = "[\"" + ContractPhone + "\"]"
	err := sendSms("dcf3c83fee3642eea48862987d598577", SendPhone, templateParas)
	if err != nil {
		logx.Error(err)
		return
	}
}

// SendGroupSuccess 已经成团
func SendGroupSuccess(SendPhone, GroupTitle, UserMail string) {
	var templateParas = "[\"" + GroupTitle + "\",\"" + UserMail + "\",\"" + ContractPhone + "\"]"
	err := sendSms("795d590e7e764b73aac89cc9b34f7b53", SendPhone, templateParas)
	if err != nil {
		logx.Error(err)
		return
	}
}

/**
 * sender,receiver,templateId不能为空
 */
func buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature string) string {
	param := "from=" + url.QueryEscape(sender) + "&to=" + url.QueryEscape(receiver) + "&templateId=" + url.QueryEscape(templateId)
	if templateParas != "" {
		param += "&templateParas=" + url.QueryEscape(templateParas)
	}
	if statusCallBack != "" {
		param += "&statusCallback=" + url.QueryEscape(statusCallBack)
	}
	if signature != "" {
		param += "&signature=" + url.QueryEscape(signature)
	}
	return param
}

func post(url string, param []byte, headers map[string]string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func buildWsseHeader(appKey, appSecret string) string {
	var cTime = time.Now().Format("2006-01-02T15:04:05Z")
	var nonce = uuid.NewV4().String()
	nonce = strings.ReplaceAll(nonce, "-", "")

	h := sha256.New()
	h.Write([]byte(nonce + cTime + appSecret))
	passwordDigestBase64Str := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf(WSSE_HEADER_FORMAT, appKey, passwordDigestBase64Str, nonce, cTime)
}
