package sender

import (
	"bytes"
	"html/template"

	"github.com/ccfos/nightingale/v6/alert/aconf"
	"github.com/ccfos/nightingale/v6/alert/astats"
	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

type (
	// Sender 发送消息通知的接口
	Sender interface {
		Send(ctx MessageContext)
	}

	// MessageContext 一个event所生成的告警通知的上下文
	MessageContext struct {
		Users  []*models.User
		Rule   *models.AlertRule
		Events []*models.AlertCurEvent
		Stats  *astats.Stats
		Ctx    *ctx.Context
	}
)

func NewSender(key string, tpls map[string]*template.Template, smtp ...aconf.SMTPConfig) Sender {
	switch key {
	case models.Dingtalk:
		return &DingtalkSender{tpl: tpls}
	case models.Wecom:
		return &WecomSender{tpl: tpls[models.Wecom]}
	case models.Feishu:
		return &FeishuSender{tpl: tpls[models.Feishu]}
	case models.FeishuCard:
		return &FeishuCardSender{tpl: tpls[models.FeishuCard]}
	case models.Email:
		return &EmailSender{subjectTpl: tpls[models.EmailSubject], contentTpl: tpls[models.Email], smtp: smtp[0]}
	case models.Mm:
		return &MmSender{tpl: tpls[models.Mm]}
	case models.Telegram:
		return &TelegramSender{tpl: tpls[models.Telegram]}
	case models.Lark:
		return &LarkSender{tpl: tpls[models.Lark]}
	case models.LarkCard:
		return &LarkCardSender{tpl: tpls[models.LarkCard]}
	}
	return nil
}

func BuildMessageContext(ctx *ctx.Context, rule *models.AlertRule, events []*models.AlertCurEvent,
	uids []int64, userCache *memsto.UserCacheType, stats *astats.Stats) MessageContext {
	users := userCache.GetByUserIds(uids)
	return MessageContext{
		Rule:   rule,
		Events: events,
		Users:  users,
		Stats:  stats,
		Ctx:    ctx,
	}
}

type BuildTplMessageFunc func(channel string, tpl *template.Template, events []*models.AlertCurEvent) string
type BuildTplsMessageFunc func(channel string, tpl map[string]*template.Template, events []*models.AlertCurEvent) string

var BuildTplMessage BuildTplMessageFunc = buildTplMessage
var BuildTplsMessage BuildTplsMessageFunc = buildTplsMessage

func buildTplMessage(channel string, tpl *template.Template, events []*models.AlertCurEvent) string {
	if tpl == nil {
		return "tpl for current sender not found, please check configuration"
	}

	var content string
	for i, event := range events {
		if i > 0 {
			event.RuleName = ""
		}

		var body bytes.Buffer
		if err := tpl.Execute(&body, event); err != nil {
			return err.Error()
		}
		content += body.String() + "\n\n"
	}

	return content
}



func buildTplsMessage(channel string, tpl map[string]*template.Template, events []*models.AlertCurEvent) string {
	if tpl == nil {
		return "tpl for current sender not found, please check configuration"
	}

	var content, recoverContent string
	var header, recoverHeader = false, false
	for _, event := range events {
		body := bytes.Buffer{}
		isRecovered := event.IsRecovered
		//生成header
		if isRecovered && !recoverHeader {
				if err := tpl[models.Header].Execute(&body, event); err != nil {
					return err.Error()
				}
				recoverContent += body.String() + "\n\n"
				recoverHeader = true
		} else if !header {
			if err := tpl[models.Header].Execute(&body, event); err != nil {
				return err.Error()
			}
			content += body.String() + "\n\n"
			header = true
		}
		
		body.Reset()
		if err := tpl[models.Dingtalk].Execute(&body, event); err != nil {
			return err.Error()
		}

		if isRecovered {
			recoverContent += body.String() + "\n\n"
		} else {
			content += body.String() + "\n\n"
		}
	}

	return content + recoverContent
}