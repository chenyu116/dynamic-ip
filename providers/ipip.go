package providers

import (
	"bytes"
	"errors"
	"fmt"
	http_client "github.com/chenyu116/http-client"
	"github.com/chenyu116/node-dynamic-ip/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
	"time"
)

func NewIPIP(cfg ProviderConfig) *ipip {
	return &ipip{cfg: cfg}
}

type ipip struct {
	cfg ProviderConfig
}

func (i *ipip) Decode(r io.Reader) (string, error) {
	buf := r.(*bytes.Buffer)
	reg := regexp.MustCompile(pattern)
	find := reg.FindStringSubmatch(buf.String())
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
func (i *ipip) Sync() {
	if i.cfg.URL == "" {
		logger.Zap.Error("[ipip] URL not found")
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
			logger.Zap.Info("[ipip] Start")
			if step < stepGetContent {
				code, err := http_client.New().Get(i.cfg.URL).Send(buf)
				if code == http.StatusNotFound {
					logger.Zap.Error("[ipip] Not Found", zap.String("step", "stepGetContent"), zap.String("next", fmt.Sprintf("%+v", time.Second*30)))
					timer.Reset(time.Second * 30)
					continue
				}
				if err != nil {
					logger.Zap.Error("[ipip]", zap.String("step", "stepGetContent"), zap.String("next", fmt.Sprintf("%+v", time.Second*10)), zap.Error(err))
					timer.Reset(time.Second * 10)
					continue
				}
				step = stepGetContent
			}
			if step < stepDecode {
				lastIp, err = i.Decode(buf)
				if err != nil {
					logger.Zap.Error("[ipip]", zap.String("step", "stepDecode"), zap.String("next", fmt.Sprintf("%+v", time.Second*10)), zap.Error(err))
					timer.Reset(time.Second * 10)
					continue
				}
				buf.Reset()
				logger.Zap.Info("[ipip]", zap.String("ip", lastIp))
				step = stepDecode
			}
			if step < stepPatched {
				err = patch(i.cfg.NodeName, lastIp)
				if err != nil {
					logger.Zap.Error("[ipip]", zap.String("step", "stepPatched"), zap.String("next", fmt.Sprintf("%+v", time.Second*5)), zap.Error(err))
					timer.Reset(time.Second * 5)
					continue
				}
			}
			step = stepPending
			lastIp = ""
			logger.Zap.Info("[ipip] Finished", zap.String("next", fmt.Sprintf("%+v", i.cfg.CheckInterval*time.Second)))
			timer.Reset(i.cfg.CheckInterval * time.Second)
		}
	}

}
