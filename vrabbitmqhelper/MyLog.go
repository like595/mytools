package vrabbitmqhelper


import (
	"encoding/json"
	"fmt"
	"github.com/like595/mytools/vtools"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type MyLog struct {
	sugarLogger *zap.SugaredLogger
	myRabbitMQ MyRabbitMQ
}

func (this *MyLog) Init(mqurl string)  {
	this.initLogger()

	this.myRabbitMQ = MyRabbitMQ{}
	this.myRabbitMQ.Init(mqurl)
}


func (this *MyLog) InfoW(key string,args ...interface{})  {
	this.sugarLogger.Info(args)
	if key != ""{
		msg := ""


		//if len(args) == 1 {
		//	if str, ok := args[0].(string); ok {
		//		msg += "["+str + "]"
		//	}
		//}
		//
		//
		//for _,val := range args{
		//	msg += "["+(val).(string) + "]"
		//}
		msg = this.getMessage("",args)
		createMsg, err := this.createMsg("INFO", msg)
		if err == nil{
			this.myRabbitMQ.Publish("pmLog.topic",key,createMsg)
		}

	}
}

func (this *MyLog) getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}


func (this *MyLog) ErrorW(key string,args ...interface{})  {
	this.sugarLogger.Error(args)
	if key != ""{
		msg := ""
		//for _,val := range args{
		//	msg +=  "["+(val).(string) + "]"
		//}
		msg = this.getMessage("",args)
		createMsg, err := this.createMsg("ERROR", msg)
		if err == nil{
			this.myRabbitMQ.Publish("pmLog.topic",key,createMsg)
		}

	}
}

/*
创建消息
 */
func (this *MyLog) createMsg(leiXing,msg string )  (*[]byte,error){
	logMsg := LogMsg{}
	logMsg.LeiXing = leiXing
	logMsg.ShiJian = time.Now().Format(ShiJianMuBanHM)
	logMsg.NeiRong = msg

	mainMsg := MainMsg{}
	mainMsg.Data = logMsg
	mainMsg.DataType = "log"
	mainMsg.DeviceClass = ""
	mainMsg.DeviceCode = ""
	mainMsg.DeviceType = ""
	bytes, err := json.Marshal(mainMsg)
	return &bytes,err
}


func (this *MyLog) initLogger() {
	iniUtil := vtools.IniUtil{}
	iniUtil.Init("./conf/Config.ini")
	logLevel := iniUtil.GetInt("logs","logLevel")
	logMaxSize := iniUtil.GetInt("logs","logMaxSize")

	hook := getLogWriter(logMaxSize)
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))	// 打印到控制台和文件
	encoder := getEncoder()
	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	//    // debug->info->warn->error
	var le zapcore.LevelEnabler
	switch logLevel {
	case 1:
		le = zapcore.DebugLevel
	case 2:
		le = zapcore.InfoLevel
	case 3:
		le = zapcore.WarnLevel
	case 4:
		le = zapcore.ErrorLevel
	case 5:
		le = zapcore.DPanicLevel
	case 6:
		le = zapcore.PanicLevel
	case 7:
		le = zapcore.FatalLevel
	default:
		le = zapcore.DebugLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, le)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	// 设置初始化字段
	//filed := zap.Fields(zap.String("serviceName", "webgo"))


	logger := zap.New(core, caller, development)
	this.sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(size int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/log.log",  //日志文件路径
		MaxSize:    size,					// 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 5,					// 日志文件最多保存多少个备份
		MaxAge:     30,					// 文件最多保存多少天
		Compress:   false,				// 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

