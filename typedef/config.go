package typedef

type Config struct {
	APP            string  //项目名称
	ENV            string  //环境值
	Debug          bool    //debug模式
	LogFileName    string  //日志文件名
	LogPath        string  //日志目录
	LogLevel       string  //日志输出等级
	ImagePrefix    string  //图片前缀
	Mysql          []Mysql //数据库
	Redis          []Redis //Redis
	Http           Http    //http
	JwtSecret      string  //token密钥
	Language       string  //消息提示语言
	Kafka          Kafka
	Pgsql          []Pgsql
	Region         string //国家地区
	ShuMei         ShuMei
	RiskSwitch     bool //风控开关
	Oss            Oss
	Av             Av
	Ios            Ios
	InnerSecret    InnerSecret
	WxPay          WxPay
	InnerIp        []string
	AggregationPay AggregationPay
	AnotherPay     AnotherPay
}

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Name     string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	Database string
	Name     string
}

type Http struct {
	Port uint16
}

type Kafka struct {
	Action      KafkaConf //从业者主题
	Gift        KafkaConf //礼物主题
	PrivateChat KafkaConf //私聊主题
	PublicChat  KafkaConf //公屏主题
}

type KafkaConf struct {
	Addr  []string
	Topic string
}

type Pgsql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Name     string
}

type ShuMei struct {
	AccessKey            string
	AppId                string
	MomentsImageCallback string
	AudioSignCallback    string
	MomentsVideoCallback string
	NetworkLine          string
}

type Oss struct {
	AccessKey       string
	SecretKey       string
	DefaultBucket   string
	DefaultEndPoint string
	DefaultRegionId string
	DefaultArn      string
	DefaultRoleName string
}

type Av struct {
	AppId       string
	Certificate string
	Supplier    string
}

type Ios struct {
	BundleId string
}

type InnerSecret struct {
	AESEncryptKey1 string
	AESEncryptKey2 string
}

type WxPay struct {
	Appid    string
	Mchid    string
	SerialNo string
	ApiV3Key string
}

type AggregationPay struct {
	Secret       string
	Url          string
	NotifyUrl    string
	MerchantCode string
}

type AnotherPay struct {
	Secret       string
	Url          string
	NotifyUrl    string
	MerchantCode string
}
