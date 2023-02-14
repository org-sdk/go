package gpt

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "github.com/go-resty/resty/v2"
    "time"
)

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTRequestBody 请求体
type ChatGPTRequestBody struct {
    Model       string  `json:"model"`
    Prompt      string  `json:"prompt"`
    MaxTokens   uint    `json:"max_tokens"`
    Temperature float64 `json:"temperature"`
}

// ChatGPTResponseBody 响应体
type ChatGPTResponseBody struct {
    ID      string                 `json:"id"`
    Object  string                 `json:"object"`
    Created int                    `json:"created"`
    Model   string                 `json:"model"`
    Choices []ChoiceItem           `json:"choices"`
    Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
    Text         string `json:"text"`
    Index        int    `json:"index"`
    Logprobs     int    `json:"logprobs"`
    FinishReason string `json:"finish_reason"`
}

/*
{
    "api_key": "xxxxxxxxx",  // openai api_key
    "session_timeout": 60,   // 会话超时时间,默认60秒,在会话时间内所有发送给机器人的信息会作为上下文
    "max_tokens": 1024,      // GPT响应字符数，最大2048，默认值512。值大小会影响接口响应速度，越大响应越慢。
    "model": "text-davinci-003", // GPT选用模型，默认text-davinci-003，具体选项参考官网训练场
    "temperature": 0.9, // GPT热度，0到1，默认0.9。数字越大创造力越强，但更偏离训练事实，越低越接近训练事实
    "session_clear_token": "清空会话" // 会话清空口令，默认`清空会话`
}
*/
type GPT struct {
    ApiKey            string  `json:"api_key"`         // key
    SessionTimeout    uint    `json:"session_timeout"` //
    MaxTokens         uint    `json:"max_tokens"`
    Model             string  `json:"model"`
    Temperature       float64 `json:"temperature"`
    SessionClearToken string  `json:"session_clear_token"`
}

// Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string, token string) (string, error) {

    gpt := GPT{
        ApiKey:            token,
        SessionTimeout:    180,
        MaxTokens:         2000,
        Model:             `text-davinci-003`,
        Temperature:       0.9,
        SessionClearToken: `clear`,
    }
    requestBody := ChatGPTRequestBody{
        Model:       gpt.Model,
        Prompt:      msg,
        MaxTokens:   gpt.MaxTokens,
        Temperature: gpt.Temperature,
    }

    client := resty.New().
        SetRetryCount(2).
        SetRetryWaitTime(1*time.Second).
        SetTimeout(time.Duration(gpt.SessionTimeout)).
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "Bearer "+gpt.ApiKey).
        SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

    rsp, err := client.R().SetBody(requestBody).Post(BASEURL + "completions")
    if err != nil {
        return "", fmt.Errorf("request openai failed, err : %v", err)
    }
    if rsp.StatusCode() != 200 {
        return "", fmt.Errorf("gtp api status code not equals 200, code is %d ,details:  %v ", rsp.StatusCode(), string(rsp.Body()))
    } else {
        fmt.Println(fmt.Sprintf("response gtp json string : %v", string(rsp.Body())))
    }

    gptResponseBody := &ChatGPTResponseBody{}
    err = json.Unmarshal(rsp.Body(), gptResponseBody)
    if err != nil {
        return "", err
    }
    var reply string
    if len(gptResponseBody.Choices) > 0 {
        reply = gptResponseBody.Choices[0].Text
    }
    return reply, nil

}
