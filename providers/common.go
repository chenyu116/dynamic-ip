package providers

import (
	"encoding/json"
	"github.com/chenyu116/node-dynamic-ip/logger"
	"go.uber.org/zap"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

const (
	pattern     = "((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}"
	stepPending = iota
	stepGetContent
	stepDecode
	stepPatched
)

func patch(name, ip string) error {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		log.Fatal(err)
	}
	patchTemplate := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"flannel.alpha.coreos.com/public-ip": ip,
			},
		},
	}
	_ = clientSet
	patchData, _ := json.Marshal(patchTemplate)
	logger.Zap.Info("[patch]", zap.ByteString("patchData", patchData))
	//_, err = clientSet.CoreV1().Nodes().Patch(context.Background(), name, types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	return nil
}

type Provider interface {
	Sync()
	Decode(reader io.Reader) (string, error)
}

type ProviderConfig struct {
	NodeName      string
	CheckInterval time.Duration
	URL           string
}

func NewSyncer() *syncer {
	return new(syncer)
}

type syncer struct {
	providers []Provider
}

func (s *syncer) Register(provider ...Provider) {
	s.providers = append(s.providers, provider...)
}

func (s *syncer) Start() {
	for _, p := range s.providers {
		go p.Sync()
	}
	select {}
}
