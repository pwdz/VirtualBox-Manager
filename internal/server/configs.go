package server

type ConfigSet struct{
	Port	string	`env:"Cloud_Port" env-default:"8000"`
	Host	string	`env:"Cloud_Host" env-default:"localhost"`
	
}

var Cfg ConfigSet