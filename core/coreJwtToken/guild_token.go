package coreJwtToken

import (
	gojwt "github.com/dgrijalva/jwt-go"
)

// GuildClaim是一些实体（通常指的用户）的状态和额外的元数据
type GuildClaims struct {
	UserId  string `json:"userId"`
	GuildId string `json:"guildId"`
	Mobile  string `json:"mobile"`
	gojwt.StandardClaims
}

// 根据用户的用户名和密码产生token
func GuildEncode(c GuildClaims, jwtSecret []byte, expireTime int64) (string, error) {
	c.ExpiresAt = expireTime
	tokenClaims := gojwt.NewWithClaims(gojwt.SigningMethodHS256, c)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func GuildDecode(token string, jwtSecret []byte) (*GuildClaims, error) {
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := gojwt.ParseWithClaims(token, &GuildClaims{}, func(token *gojwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*GuildClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
