package userRouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wlbwlbwlb/lightgate/config"
	"github.com/wlbwlbwlb/lightgate/kvstore"
	"github.com/wlbwlbwlb/lightgate/sessions"

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

		req := struct {
			Token string `json:"token"`
		}{}
		if e := c.Bind(&req); e != nil {
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

		resp2 := struct {
			Code   int   `json:"code"`
			UserId int64 `json:"user_id"`
		}{}
		if e = json.Unmarshal(body, &resp2); e != nil {
			return
		}

		conn := kvstore.RedisPool.Get()
		defer conn.Close()

		keyStr := fmt.Sprintf("user:%d:loc", resp2.UserId)

		loc, e2 := redis.String(conn.Do("GET", keyStr))

		//已经登录本服
		b := nil == e2 && loc == config.TOML.Addr

		//已经登录其他服
		b2 := nil == e2 && loc != config.TOML.Addr

		if b2 {
			sessions.OnLogoutNotify(resp2.UserId, loc)
		}

		if b {
			sessions.OnKickout(resp2.UserId)
		}

		sessions.OnLogin(resp2.UserId, c.Session())
		conn.Do("SET", keyStr, config.TOML.Addr)

		resp.Code = 0
	})

}
