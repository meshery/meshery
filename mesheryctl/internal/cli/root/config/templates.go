package config

// TemplateToken is the default token entry created during config initialization
var TemplateToken = Token{
	Name:     "default",
	Location: "",
}

// TemplateContext is the default context VALUE (name is provided separately)
var TemplateContext = Context{
	Endpoint: "",
	Platform: "docker",
	Channel:  "stable",
	Version:  "latest",
}
