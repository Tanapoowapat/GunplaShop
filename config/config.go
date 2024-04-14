package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("Error loading .env file Fail %v", err)
	}
	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("Error  Fail to load port ENV %v", err)
				}
				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			readTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("Error  Fail to load readTimeout ENV %v", err)
				}
				return time.Duration(int64(t) * int64(math.Pow10(9)))
			}(),
			writeTimeout: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_WRTIE_TIMEOUT"])
				if err != nil {
					log.Fatalf("Error  Fail to load writeTimeout ENV %v", err)
				}
				return time.Duration(int64(t) * int64(math.Pow10(9)))
			}(),
			bodyLimit: func() int {
				b, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("Error  Fail to load bodyLimit ENV %v", err)
				}
				return b
			}(),
			fileLimit: func() int {
				f, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("Error  Fail to load fileLimit ENV %v", err)
				}
				return f
			}(),
			gcpbucket: envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("Error  Fail to load port ENV %v", err)
				}
				return p
			}(),
			protocal: envMap["DB_PROTOCAL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnection: func() int {
				m, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("Error  Fail to load maxConnection ENV %v", err)
				}
				return m
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			sercetKey: envMap["JWT_SERCET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpiresAt: func() int {
				a, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("Error  Fail to load accessExpiresAt ENV %v", err)
				}
				return a
			}(),
			refreshExpiresAt: func() int {
				r, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("Error  Fail to load refreshExpiresAt ENV %v", err)
				}
				return r
			}(),
		},
	}
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

type IAppConfig interface {
	Url() string //host:port
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	GcpBucket() string
	Host() string
	Port() int
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int //bytes
	fileLimit    int //bytes
	gcpbucket    string
}

func (c *config) App() IAppConfig {
	return c.app
}

func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) FileLimit() int              { return a.fileLimit }
func (a *app) GcpBucket() string           { return a.gcpbucket }
func (a *app) Host() string                { return a.host }
func (a *app) Port() int                   { return a.port }

type IDbConfig interface {
	Url() string //host:port
	MaxConnection() int
}

type db struct {
	host          string
	port          int
	protocal      string
	username      string
	password      string
	database      string
	sslMode       string
	maxConnection int
}

func (c *config) Db() IDbConfig {
	return c.db
}

func (d *db) Url() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}

func (d *db) MaxConnection() int { return d.maxConnection }

type IJwtConfig interface {
	SercetKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetAccessExpiresAt(int)
	SetRefreshExpiresAt(int)
	GetKeyInfo() string
}

type jwt struct {
	adminKey         string
	sercetKey        string
	apiKey           string
	accessExpiresAt  int //sec
	refreshExpiresAt int //sec
}

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}

func (j *jwt) SercetKey() []byte         { return []byte(j.sercetKey) }
func (j *jwt) AdminKey() []byte          { return []byte(j.adminKey) }
func (j *jwt) ApiKey() []byte            { return []byte(j.apiKey) }
func (j *jwt) AccessExpiresAt() int      { return j.accessExpiresAt }
func (j *jwt) RefreshExpiresAt() int     { return j.refreshExpiresAt }
func (j *jwt) SetAccessExpiresAt(a int)  { j.accessExpiresAt = a }
func (j *jwt) SetRefreshExpiresAt(r int) { j.refreshExpiresAt = r }
func (j *jwt) GetKeyInfo() string {
	return fmt.Sprintf("adminKey=%s, apiKey=%s", j.adminKey, j.apiKey)
}
