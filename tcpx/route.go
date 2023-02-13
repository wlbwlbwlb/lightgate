package tcpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wl955/lightgate/config"
	"github.com/wl955/lightgate/kvstore"

	"github.com/DarthPestilane/easytcp"
	"github.com/gomodule/redigo/redis"
)

var client = &http.Client{}

func addRoute(serve *easytcp.Server) {

	serve.AddRoute(1001, func(c easytcp.Context) {
		var e error

		resp := struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 1,
		}
		defer func() {
			if e != nil {
				resp.Msg = e.Error()
			}
			c.SetResponse(1002, &resp)
		}()

		req := struct {
			Token string `json:"token"`
		}{}
		e = c.Bind(&req)
		if e != nil {
			return
		}

		req2, e := http.NewRequest("GET",
			config.TOML.Authapi,
			nil,
		)
		if e != nil {
			return
		}
		req2.Header.Add("x-token", req.Token)

		res, e := client.Do(req2)
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
		//log.Info(string(body))

		got := struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data struct {
				UserId int64 `json:"user_id"`
			}
		}{}
		if e = json.Unmarshal(body, &got); e != nil {
			return
		}

		conn := kvstore.RedisPool.Get()
		defer conn.Close()

		keyStr := fmt.Sprintf("user:%d:loc", got.Data.UserId)

		val, e2 := redis.Int(conn.Do("GET", keyStr))

		//没登录过
		b := redis.ErrNil == e2

		//已经登录本服
		b2 := nil == e2 && val == config.TOML.Port

		//已经登录其他服
		b3 := nil == e2 && val != config.TOML.Port

		if b || b2 {
			sessions.onLoginSuccess(got.Data.UserId, c.Session())
			conn.Do("SET", keyStr, config.TOML.Port)
		}

		if b3 {
			//todo
		}

		log.Printf("user %d signin\n", got.Data.UserId)

		resp.Code = 0
	})

}
