package configs

//todo iotaをintに修正
//todo config.goに持っていく
//constファイルをなくす方向に持っていきたい。(それぞれ適切な場所で宣言するべき)

//log level
const (
	DEBUG = 1
	INFO  = 2
	WARN  = 3
	ERROR = 4
	FATAL = 5
)

//time.timeを利用しやすい形に整形するための設定値です。
const (
	FormatDatetime = "2006/01/02 15:04:05"
	FormatDate     = "2006/01/02"
	FormatTime     = "15:04:05"
)

//user-agentからデバイスを判定するための設定値です。
const (
	IOS     = "iOS"
	IPhone  = "iPhone"
	IPad    = "iPad"
	IPod    = "iPod"
	Android = "android"
)

//todo エラーパッケージに持っていく
//error code list (利用していないので削除)
const (
	CodeBadRequest = "E0001"
)

//CryptoKey は環境変数に設定されている暗号化キーです。
const CryptoKey = "AES256KEY"

//mysql Error を判定するためのエラーコードです。
const (
	MYSQL_ER_DUP_ENTRY = 1062
)
