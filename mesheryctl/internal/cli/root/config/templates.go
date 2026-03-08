package config

var TemplateToken = Token{
	Name:     "default",
	Location: "",
}

var TemplateContext = Context{
	Endpoint: "",
	Platform: "docker",
	Channel:  "stable",
	Version:  "latest",
}
