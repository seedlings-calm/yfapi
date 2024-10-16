package enum

// 验证码类型
const (
	// SmsCodeLogin 登录验证码
	SmsCodeLogin = 1
	// SmsCodeResetPass 重置密码
	SmsCodeResetPass = 2
	// SmsCodeForgetPass 忘记密码
	SmsCodeForgetPass = 3
	//SmsCodeSetPass 设置密码
	SmsCodeSetPass = 5
	// SmsCodeUserDeleteApply 账号注销申请
	SmsCodeUserDeleteApply = 4
	// SmsCodeOperationLogin 运营后台登录
	SmsCodeOperationLogin = 6
	//验证原始手机号
	SmsCodeVerifyMobile = 7
	//更换手机号
	SmsCodeChangeMobile = 8
	// SmsCodeGuildAdminLogin 公会后台登录
	SmsCodeGuildAdminLogin = 9
	//SmsCodeGuildAdminLogin 房主后台登录
	SmsCodeRoomAdminLogin = 10
	//SmsCodeBindBankCard 绑定银行卡
	SmsCodeBindBankCard = 11
)
