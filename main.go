package main

import (
	"fmt"
	"github.com/chenyu116/node-dynamic-ip/config"
	"github.com/chenyu116/node-dynamic-ip/providers"
)

var (
	_version_   = ""
	_branch_    = ""
	_commit_    = ""
	_buildTime_ = ""
)

func main() {
	fmt.Printf("Version: %s, Branch: %s, Commit: %s, BuildTime: %s\n",
		_version_, _branch_, _commit_, _buildTime_)

	cf := config.GetConfig()

	syncer := providers.NewSyncer()
	syncer.Register(providers.NewIPIP(providers.ProviderConfig{NodeName: cf.Common.NodeName, URL: cf.Providers.IPIP.URL, CheckInterval: cf.Providers.IPIP.CheckInterval}))
	syncer.Start()
}
