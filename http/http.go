package http

import (
    "crypto/tls"
    "fmt"
    "github.com/go-resty/resty/v2"
)

type Resty struct {
    C *resty.Client
}

func (this *Resty) New() *Resty {

    this.C = resty.New()
    return this

}

//
// Default
//  @Description:
//  @param url
//  @param schema
//  @param body
//  @return []byte
//
func Default(url string, body interface{}) []byte {

    client := resty.New().SetRetryCount(2).SetRetryWaitTime(10).SetTimeout(12).
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "Bearer ").
        SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

    rsp, err := client.R().SetBody(body).Post(url)
    if err != nil {
        fmt.Println(err)
    }
    if rsp.StatusCode() != 200 {
        fmt.Println(rsp)
    }
    return rsp.Body()

}


