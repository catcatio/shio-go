package kernel

type ServiceOptions struct {
	*LineChatOptions
	*DialogflowOptions
}

type LineChatOptions struct {
	ChannelSecret      string
	ChannelAccessToken string
}

type DialogflowOptions struct {
	*GCPOptions
}

type GCPOptions struct {
	ProjectName    string
	CredentialJson string
}

type ServiceConfigBase struct {
	HttpPort string
	GrpcPort string
	Enable   bool
}
