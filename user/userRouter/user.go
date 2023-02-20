package userRouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wl955/lightgate/config"
	"github.com/wl955/lightgate/kvstore"
	"github.com/wl955/lightgate/sess"

	"github.com/DarthPestilane/easytcp"
	"github.com/gomodule/redigo/redis"
)

var client = &http.Client{}

func Router(serve *easytcp.Server) {

	serve.AddRoute(1001, func(c easytcp.Context) {
		resp := struct {
			Code int `json:"code"`
		}{
			Code: 1,
		}
		defer func() {
			c.SetResponse(1002, &resp)
		}()

		param := struct {
			Token string `json:"token"`
		}{}
		if e := c.Bind(&param); e != nil {
			return
		}

		req, e := http.NewRequest("GET",
			config.TOML.Authapi,
			nil,
		)
		if e != nil {
			return
		}
		req.Header.Add("x-token", param.Token)

		res, e := client.Do(req)
		if e != nil {
			return
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			e = errors.New(res.Status)
			return
		}

		body, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return
		}

		got := struct {
			Code   int   `json:"code"`
			UserId int64 `json:"user_id"`
		}{}
		if e = json.Unmarshal(body, &got); e != nil {
			return
		}

		conn := kvstore.RedisPool.Get()
		defer conn.Close()

		keyStr := fmt.Sprintf("user:%d:loc", got.UserId)

		val, e2 := redis.Int(conn.Do("GET", keyStr))

		//没登录过
		b := redis.ErrNil == e2

		//已经登录本服
		b2 := nil == e2 && val == config.TOML.Port

		//已经登录其他服
		b3 := nil == e2 && val != config.TOML.Port

		if b || b2 {
			sess.OnLoginSuccess(got.UserId, c.Session())
			conn.Do("SET", keyStr, config.TOML.Port)
		}

		if b3 {
			//todo
		}

		fmt.Printf("user %d login\n", got.UserId)

		resp.Code = 0
	})

}
