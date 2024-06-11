package main

import (
	"encoding/json"
	"fmt"
	kfklib "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/community-robot-lib/config"
	framework "github.com/opensourceways/community-robot-lib/robot-gitee-framework"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

const botName = "message-collect-githook"

func newRobot() *robot {
	return &robot{}
}

type robot struct{}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) getConfig(cfg config.Config, org, repo string) (*botConfig, error) {
	c, ok := cfg.(*configuration)
	if !ok {
		return nil, fmt.Errorf("can't convert to configuration")
	}
	if bc := c.configFor(org, repo); bc != nil {
		return bc, nil
	}

	return nil, fmt.Errorf("no config for this repo:%s/%s", org, repo)
}

func (bot *robot) RegisterEventHandler(p framework.HandlerRegitster) {
	p.RegisterIssueHandler(bot.handleIssueEvent)
	p.RegisterPullRequestHandler(bot.handlePREvent)
	p.RegisterNoteEventHandler(bot.handleNoteEvent)
	p.RegisterPushEventHandler(bot.handlePushEvent)
}
func (bot *robot) handlePREvent(e *sdk.PullRequestEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle pr event,send kafka message_collect_pr")
	return kfklib.Publish("message_collect_pr", nil, body)
}
func (bot *robot) handleNoteEvent(e *sdk.NoteEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle note event,send kafka message_collect_note")
	return kfklib.Publish("message_collect_note", nil, body)
}

func (bot *robot) handlePushEvent(e *sdk.PushEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle push event,send kafka message_collect_push")
	return kfklib.Publish("message_collect_push", nil, body)
}

func (bot *robot) handleIssueEvent(e *sdk.IssueEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle issue event,send kafka message_collect_issue")
	return kfklib.Publish("message_collect_issue", nil, body)
}
