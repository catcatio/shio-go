package kernel

type ServiceConfiguration struct {
	Chat   *ChatConfiguration
	Intent *IntentConfiguration
}

type ChatConfiguration struct {
	ServiceConfigBase
}

type IntentConfiguration struct {
	ServiceConfigBase
}

type ServiceConfigBase struct {
	HttpPort string
	GrpcPort string
	Enable   bool
}
