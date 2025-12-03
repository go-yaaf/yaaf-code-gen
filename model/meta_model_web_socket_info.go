package model

// region Web Socket Info structure ------------------------------------------------------------------------------------

// MessageInfo web Socket Message information
type MessageInfo struct {
	Name      string     // Name of message
	Docs      []string   // Class documentation
	IsRequest bool       // Is request message (client-server) or response(server-client)
	Message   *ClassInfo // List of class fields
}

func NewMessageInfo(name string) *MessageInfo {
	return &MessageInfo{
		Name: name,
		Docs: make([]string, 0),
	}
}

// WebSocketInfo web Socket information
type WebSocketInfo struct {
	Name    string        // Name of the socket
	TsName  string        // TypeScript socket name (small caps)
	Group   string        // Name of the socket group
	Path    string        // Web socket URI path
	Usage   string        // Web socket Usage sample
	Docs    []string      // Web socket Documentation
	Methods []*MethodInfo // List of class fields
}

func NewWebSocketInfo(name string) *WebSocketInfo {
	return &WebSocketInfo{
		Name:    name,
		TsName:  SmallCaps(name),
		Methods: make([]*MethodInfo, 0),
		Docs:    make([]string, 0),
	}
}

// endregion
