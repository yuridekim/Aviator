package controller

const (
	apiUrlCreate = "https://ncloud.apigw.ntruss.com/vserver/v2/createServerInstances"
	apiUrlDelete = "https://ncloud.apigw.ntruss.com/vserver/v2/terminateServerInstances"
	apiUrlGet    = "https://ncloud.apigw.ntruss.com/vserver/v2/getMemberServerImageInstanceList"
	apiUrlStop   = "https://ncloud.apigw.ntruss.com/vserver/v2/stopServerInstances"
	apiUrlUpdate = "https://ncloud.apigw.ntruss.com/vserver/v2/changeServerInstanceSpec"
	// error level
	ErrorLevelIsInfo    = 0
	ErrorLevelIsFatal   = 1
	ErrorLevelIsAnError = 2
	ErrorLevelIsWarn    = 3
	ErrorLevelIsDebug   = 5
	ErrorLevelIsTrace   = 6
)
