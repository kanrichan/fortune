package main

import (
	ft "fortune/fortune"
	"github.com/wdvxdr1123/ZeroBot"
)

func main() {
	zero.Run(zero.Option{
		Host:          ft.Conf.Host,
		Port:          ft.Conf.Port,
		AccessToken:   ft.Conf.AccessToken,
		NickName:      []string{"ft"},
		CommandPrefix: "",
		SuperUsers:    []string{ft.Conf.Master},
	})
	select {}
}

func init() {
	ft.Init()
}
