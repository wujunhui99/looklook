package main

import (
	"flag"
	"fmt"

	"github.com/wujunhui99/looklook/app/payment/cmd/api/internal/config"
	"github.com/wujunhui99/looklook/app/payment/cmd/api/internal/handler"
	"github.com/wujunhui99/looklook/app/payment/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/payment.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
