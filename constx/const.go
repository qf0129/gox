package constx

const (
	KeyOfRequestId     = "reqId"
	KeyOfRequestUser   = "reqUser"
	KeyOfRequestUserId = "reqUserId"
	KeyOfCookieToken   = "t"
	KeyOfCookieUserId  = "u"
)

const (
	DefaultListenAddr           = ":8080"
	DefaultGinMode              = "debug"
	DefaultLogLevel             = "debug"
	DefaultDBLogLevel           = "warn"
	DefaultReadTimeout          = 60
	DefaultWriteTimeout         = 60
	DefaultEncryptSecret        = "DEFAULT_SECRET_X"
	DefaultCookieExpiredSeconds = 3600

	DefaultSqliteFile      = "db.sqlite"
	DefaultModelPrimaryKey = "id"
	DefaultQueryPageSize   = 10
)
