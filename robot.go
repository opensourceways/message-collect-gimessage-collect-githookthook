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
	log.Info("handle pr event,send kafka gitee_pr_raw")
	return kfklib.Publish("gitee_pr_raw", nil, body)
}
func (bot *robot) handleNoteEvent(e *sdk.NoteEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle note event,send kafka gitee_note_raw")
	return kfklib.Publish("gitee_note_raw", nil, body)
}

func (bot *robot) handlePushEvent(e *sdk.PushEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle push event,send kafka gitee_push_raw")
	return kfklib.Publish("gitee_push_raw", nil, body)
}

func (bot *robot) handleIssueEvent(e *sdk.IssueEvent, c config.Config, log *logrus.Entry) error {
	body, _ := json.Marshal(e)
	log.Info("handle issue event,send kafka gitee_issue_raw")
	if e.Issue.TypeName == "CVE和安全问题" {
		return kfklib.Publish("cve_issue_raw", nil, body)
	} else {
		return kfklib.Publish("gitee_issue_raw", nil, body)
	}
}
