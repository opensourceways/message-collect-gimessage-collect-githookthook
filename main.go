package main

import (
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	framework "github.com/opensourceways/community-robot-lib/robot-gitee-framework"
	"github.com/opensourceways/message-collect-githook/kafka"
	"github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
)

type options struct {
	service liboptions.ServiceOptions
	gitee   liboptions.GiteeOptions
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.gitee.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options
	o.service.AddFlags(fs)
	err := fs.Parse(args)
	if err != nil {
		return options{}
	}
	return o
}

func main() {
	logrusutil.ComponentInit(botName)
	log := logrus.NewEntry(logrus.StandardLogger())
	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	cfg := Init()
	if err := kafka.Init(&cfg.Kafka, log, false); err != nil {
		logrus.Errorf("init kafka failed, err:%s", err.Error())
		return
	}
	p := newRobot()

	framework.Run(p, o.service)
}

type Config struct {
	Kafka kafka.Config `json:"kafka"`
}

func Init() *Config {
	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	cfg := new(Config)
	logrus.Info(os.Args[1:])
	if err := utils.LoadFromYaml(o.service.ConfigFile, cfg); err != nil {
		logrus.Error("Config初始化失败, err:", err)
		return nil
	}
	logrus.Infof("the version is %v", cfg.Kafka.Version)
	return cfg
}
