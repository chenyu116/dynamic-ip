package providers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	http_client "github.com/chenyu116/http-client"
	"github.com/chenyu116/node-dynamic-ip/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
	"time"
)

func NewIp138(cfg ProviderConfig) *ip138 {
	return &ip138{cfg: cfg}
}

type ip138 struct {
	cfg ProviderConfig
}

func (i *ip138) Decode(buf io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return "", err
	}
	title := doc.Find("title").Text()
	reg := regexp.MustCompile(pattern)
	find := reg.FindStringSubmatch(title)
	if len(find) == 0 {
		return "", errors.New("not found")
	}
	for _, ip := range find {
		matched, err := regexp.MatchString(pattern, ip)
		if err != nil || !matched {
			continue
		}
		return ip, nil
	}
	return "", errors.New("not found")
}
func (i *ip138) Sync() {
	if i.cfg.URL == "" {
		logger.Zap.Error("[ip138] URL not found")
		return
	}
	timer := time.NewTimer(0)
	buf := new(bytes.Buffer)
	lastIp := ""
	var err error
	step := stepPending
	for {
		select {
		case <-timer.C:
			logger.Zap.Info("[ip138] Start")
			if step < stepGetContent {
				code, err := http_client.New().Get(i.cfg.URL).Send(buf)
				if code == http.StatusNotFound {
					logger.Zap.Error("[ip138] Not Found", zap.String("step", "stepGetContent"), zap.String("next", fmt.Sprintf("%+v", time.Second*30)))
					timer.Reset(time.Second * 30)
					continue
				}
				if err != nil {
					logger.Zap.Error("[ip138]", zap.String("step", "stepGetContent"), zap.String("next", fmt.Sprintf("%+v", time.Second*10)), zap.Error(err))
					timer.Reset(time.Second * 10)
					continue
				}
				step = stepGetContent
			}
			if step < stepDecode {
				lastIp, err = i.Decode(buf)
				if err != nil {
					logger.Zap.Error("[ip138]", zap.String("step", "stepDecode"), zap.String("next", fmt.Sprintf("%+v", time.Second*10)), zap.Error(err))
					timer.Reset(time.Second * 10)
					continue
				}
				buf.Reset()
				logger.Zap.Info("[ip138]", zap.String("ip", lastIp))
				step = stepDecode
			}
			if step < stepPatched {
				err = patch(i.cfg.NodeName, lastIp)
				if err != nil {
					logger.Zap.Error("[ip138]", zap.String("step", "stepPatched"), zap.String("next", fmt.Sprintf("%+v", time.Second*5)), zap.Error(err))
					timer.Reset(time.Second * 5)
					continue
				}
			}
			step = stepPending
			lastIp = ""
			logger.Zap.Info("[ip138] Finished", zap.String("next", fmt.Sprintf("%+v", i.cfg.CheckInterval*time.Second)))
			timer.Reset(i.cfg.CheckInterval * time.Second)
		}
	}

}
