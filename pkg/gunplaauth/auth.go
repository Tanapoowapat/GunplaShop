package gunplaauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Tanapoowapat/GunplaShop/config"
	"github.com/Tanapoowapat/GunplaShop/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokensType string

const (
	AccessToken TokensType = "access"
	RefeshToken TokensType = "refresh"
	AdminToken  TokensType = "admin"
	ApiKeyToken TokensType = "api"
)

type IAuth interface {
	SignToken() string
}

func JwtTimeDuration(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func JwtTimeRepeatAdpter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

// Auth
type Auth struct {
	mapclaims *mapClaims
	cfg       config.IJwtConfig
}

// Admin Auth
type adminAuth struct {
	*Auth
}

type IAdmin interface {
	SignToken() string
}

// Api key
type ApiKey struct {
	*Auth
}

type IApiKey interface {
	SignToken() string
}

type mapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

func NewAuthTokens(tokensType TokensType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuth, error) {

	switch tokensType {
	case AccessToken:
		return NewAcessTokens(cfg, claims), nil
	case RefeshToken:
		return NewRefreshTokens(cfg, claims), nil
	case AdminToken:
		return NewAdminTokens(cfg), nil
	case ApiKeyToken:
		return NewApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("unknows token type")
	}

}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SercetKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseApiKey(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.ApiKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*mapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &mapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*mapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}

// Sign Token
func (a *Auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapclaims)
	ss, _ := token.SignedString(a.cfg.SercetKey())
	return ss
}

func (a *adminAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapclaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func (a *ApiKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapclaims)
	ss, _ := token.SignedString(a.cfg.ApiKey())
	return ss
}

func RepeatyToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	mapclaims := &Auth{
		mapclaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gunpla-shop",
				Subject:   "access-tokens",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: JwtTimeRepeatAdpter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		cfg: cfg,
	}
	return mapclaims.SignToken()
}

func NewAcessTokens(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &Auth{
		cfg: cfg,
		mapclaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gunpla-shop",
				Subject:   "access-tokens",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: JwtTimeDuration(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func NewRefreshTokens(cfg config.IJwtConfig, claims *users.UserClaims) IAuth {
	return &Auth{
		cfg: cfg,
		mapclaims: &mapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "gunpla-shop",
				Subject:   "refresh-tokens",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: JwtTimeDuration(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func NewAdminTokens(cfg config.IJwtConfig) IAuth {
	return &adminAuth{
		&Auth{
			cfg: cfg,
			mapclaims: &mapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "gunpla-shop",
					Subject:   "admin-tokens",
					Audience:  []string{"admin"},
					ExpiresAt: JwtTimeDuration(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func NewApiKey(cfg config.IJwtConfig) IAuth {
	return &ApiKey{
		Auth: &Auth{
			cfg: cfg,
			mapclaims: &mapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "gunpla-shop",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
