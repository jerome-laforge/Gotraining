package log

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

/**
 * @author Jérôme LAFORGE - Orange / IMT / OLPS / SOFT
 *         <b>Copyright :</b> Orange 2016<br>
 */

const loggerCtx = "LoggerCtx"

var (
	logger     log15.Logger
	onceLogger sync.Once
)

func GetLogger() log15.Logger {
	onceLogger.Do(func() {
		logger = log15.Root()
		logger.SetHandler(log15.CallerStackHandler("%+v", log15.StdoutHandler))
		logger = logger.New(log15.Ctx{"pid": os.Getpid()})
	})
	return logger
}

func GetLoggerFromContext(ctx context.Context) (logger log15.Logger) {
	var ok bool
	if logger, ok = ctx.Value(loggerCtx).(log15.Logger); !ok {
		logger = log15.New(log15.Ctx{"uuid": "Not found"})
	}

	return logger
}

func SetLogger(ctx *gin.Context) {
	logger := GetLogger()
	uuid := uuid.NewV4()

	ctx.Set(loggerCtx, logger.New(log15.Ctx{"uuid": uuid}))
}
