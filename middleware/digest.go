package middleware

import (
	"github.com/labstack/echo"
	digest "github.com/abbot/go-http-auth"
	"net/http"
	"net/url"
	"fmt"
	"strings"
)

type Digest struct {
	auth  *digest.DigestAuth
	/**
	 * 白名单+黑名单
	 * 在白名单中，同时不在黑名单中的，不需要认证
	 */
	allow digestRouter
	deny  digestRouter
}

/**
 * 创建
 */
func NewDigest(realm string, secrets digest.SecretProvider) *Digest {
	auth := digest.NewDigestAuthenticator(realm, secrets)
	// secrets返回的是明文，而非摘要
	auth.PlainTextSecrets = true
	return &Digest{
		auth: auth,
		allow: digestRouter{},
	}
}

/**
 * 中间件处理
 */
func (self *Digest) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 跳过白名单
		if self.check(c.Request().Method, c.Request().URL) {
			return next(c)
		}
		// 判断是否放行url
		//
		w := c.Response()
		r := c.Request()
		if username, nextAuth := self.auth.CheckAuth(c.Request()); username == "" {
			self.auth.RequireAuth(w, r)
			return echo.NewHTTPError(http.StatusUnauthorized, "401 Unauthorized\n");
		} else {
			if (nil != nextAuth) {
				w.Header().Set("Authentication-Info", *nextAuth);
			}
			// 在echo.Context中缓存认证信息
			c.Set("digest-user", username);
			// 执行下一个
			return next(c)
		}
	}
}

/**
 * 白名单
 */
func (self *Digest) Deny(method, url string) {
	self.deny.add(method, url)
}

/**
 * 黑名单
 */
func (self *Digest) Allow(method, url string) {
	self.allow.add(method, url)
}

func (self *Digest) check(method string, url *url.URL) bool {
	path := strings.Split(url.Path, "/")
	pLen := len(path)
	/* 判断白名单 */
	whiteList := false
	if (self.allow.in(method, url.Path)) {
		whiteList = true
	}
	for i := pLen - 1;i > 0;i-- {
		if self.allow.in(method, strings.Join(path[:i], "/") + "/*") {
			whiteList = true
		}
	}
	// 不在白名单中
	if !whiteList {
		return false
	}
	/* 判断黑名单 */
	if self.deny.in(method, url.Path) {
		return false
	}
	for i := pLen - 1; i > 0; i-- {
		if self.deny.in(method, strings.Join(path[:i], "/") + "/*") {
			return false
		}
	}
	return true;
}

// 路径
type digestRouter map[string]*digestURL

type digestURL struct {
	path   string
	method []string
}

/**
 * 增加路径
 */
func (self digestRouter) add(method, path string) {
	if v := self[path]; nil == v {
		// 新增
		self[path] = &digestURL{
			path: path,
			method: []string{method},
		}
	} else {
		v.method = sortMethod(append(v.method, method))
	}
}

func (self digestRouter) remove(method, path string) {
	if v := self[path]; nil != v {
		if "*" == v.method[0] && "*" != method {
			panic(fmt.Errorf("存在通配符方法时请先删除通配符"))
		}
		for i, m := range v.method {
			if m == method {
				v.method = make([]string, 0, len(v.method) - 1)
				v.method = append(v.method, v.method[:i]...)
				v.method = append(v.method, v.method[i + 1:]...)
			}
		}
	}
}

func (self digestRouter) in(method, path string) bool {
	if v := self[path]; nil != v {
		// *通配符通过
		if v.method[0] == "*" {
			return true
		}
		// 遍历符合方法
		for _, m := range v.method {
			if m == method {
				return true
			}
		}
	}
	return false
}

/**
 * 对方法排序
 */
func sortMethod(methods []string) []string {
	result := make([]string, 0, len(methods))
	for _, m := range methods {
		// *通配符
		if m == "*" {
			return []string{"*"};
		}
		for _, r := range result {
			if r == m {
				goto next
			}
		}
		result = append(result, m)
		next:
	}
	return result
}