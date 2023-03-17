package feishuRobotGo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/saodd/alog"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Robot struct {
	Secret string
	Hook   string
}

func (c *Robot) send(ctx context.Context, reqBody []byte) error {
	ctx2, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx2, "POST", c.Hook, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		alog.CE(ctx, err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		alog.CE(ctx, err)
		return err
	}

	var data = new(RobotResponse)
	if err := json.Unmarshal(body, data); err != nil {
		alog.CE(ctx, err, alog.V{"body": string(body)})
		return err
	} else if data.RobotResponseError.Code != 0 {
		err2 := &data.RobotResponseError
		alog.CE(ctx, err2, alog.V{"body": string(body)})
		return err2
	}
	return nil
}

// SendPost 发送富文本消息 https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN#f62e72d5
func (c *Robot) SendPost(ctx context.Context, content *RobotContent) error {
	ts := time.Now().Unix()
	sign, _ := GenSign(c.Secret, ts)

	req := &RobotRequest{
		Timestamp: strconv.FormatInt(ts, 10),
		Sign:      sign,
		MsgType:   "post",
		Content:   content,
	}
	body, err := json.Marshal(req)
	if err != nil {
		alog.CE(ctx, err, alog.V{"content": content})
		return err
	}
	return c.send(ctx, body)
}

type RobotRequest struct {
	Timestamp string        `json:"timestamp"` // "1599360473"
	Sign      string        `json:"sign"`
	MsgType   string        `json:"msg_type"` // "post"
	Content   *RobotContent `json:"content"`
}
type RobotContent struct {
	Post RobotPostContent `json:"post"`
}
type RobotPostContent struct {
	ZhCn RobotPostContentGroup `json:"zh_cn"`
	//EnUs RobotPostContentGroup `json:"en_us"`
}
type RobotPostContentGroup struct {
	Title   string                           `json:"title"`
	Content [][]RobotPostContentGroupContent `json:"content"`
}
type RobotPostContentGroupContent struct {
	Tag    string `json:"tag"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

type RobotResponse struct {
	RobotResponseSuccess
	RobotResponseError
}
type RobotResponseSuccess struct {
	Extra         interface{} `json:"Extra"`
	StatusCode    int         `json:"StatusCode"`
	StatusMessage string      `json:"StatusMessage"`
}
type RobotResponseError struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *RobotResponseError) Error() string {
	return fmt.Sprintf("[%d]%s", r.Code, r.Msg)
}

// GenSign 签名算法，来自于官方文档：https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN
func GenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
