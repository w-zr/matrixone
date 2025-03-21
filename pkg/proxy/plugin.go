// Copyright 2021 - 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"context"
	"time"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/common/morpc"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/pb/plugin"
	v2 "github.com/matrixorigin/matrixone/pkg/util/metric/v2"
)

// pluginRouter is a router implementation that uses external plugin to select CN server.
type pluginRouter struct {
	sid string
	// Router is a delegated impl that is used when plugin flags a Bypass action
	Router
	// plugin is the plugin that is used to select CN server
	plugin Plugin
}

func newPluginRouter(sid string, r Router, p Plugin) *pluginRouter {
	return &pluginRouter{
		sid:    sid,
		Router: r,
		plugin: p,
	}
}

// Route implements Router.Route.
func (r *pluginRouter) Route(
	ctx context.Context, sid string, ci clientInfo, filter func(uuid string) bool,
) (*CNServer, error) {
	re, err := r.plugin.RecommendCN(ctx, ci)
	if err != nil {
		return nil, err
	}
	if re.Updated {
		// plugin signals that a state updated has happened, request a refresh if the delegated router is refreshable
		if rr, ok := r.Router.(RefreshableRouter); ok {
			rr.Refresh(false)
		}
	}
	switch re.Action {
	case plugin.Select:
		if re.CN == nil {
			return nil, moerr.NewInternalErrorNoCtx("no CN server selected")
		}
		// selected CN should be filtered out, fall back to the delegated router
		if filter != nil && filter(re.CN.SQLAddress) {
			return r.Router.Route(ctx, sid, ci, filter)
		}
		v2.ProxyConnectSelectCounter.Inc()
		return &CNServer{
			reqLabel: ci.labelInfo,
			cnLabel:  re.CN.Labels,
			uuid:     re.CN.ServiceID,
			addr:     re.CN.SQLAddress,
			hash:     ci.hash,
		}, nil
	case plugin.Reject:
		v2.ProxyConnectRejectCounter.Inc()
		return nil, withCode(moerr.NewInfoNoCtx(re.Message),
			codeAuthFailed)
	case plugin.Bypass:
		return r.Router.Route(ctx, r.sid, ci, filter)
	default:
		return nil, moerr.NewInternalErrorNoCtxf("unknown recommended action %d", re.Action)
	}
}

// Plugin is the interface of proxy plugin.
type Plugin interface {
	// RecommendCN returns the recommended CN server.
	RecommendCN(ctx context.Context, client clientInfo) (*plugin.Recommendation, error)
}

type rpcPlugin struct {
	client  morpc.RPCClient
	backend string
	timeout time.Duration
}

func newRPCPlugin(sid string, backend string, timeout time.Duration) (*rpcPlugin, error) {
	codec := morpc.NewMessageCodec(
		sid,
		func() morpc.Message {
			return &plugin.Response{}
		},
	)
	backendOpts := []morpc.BackendOption{
		morpc.WithBackendConnectTimeout(timeout),
		morpc.WithBackendHasPayloadResponse(),
		morpc.WithBackendLogger(logutil.GetGlobalLogger().Named("plugin-backend")),
	}
	bf := morpc.NewGoettyBasedBackendFactory(codec, backendOpts...)

	clientOpts := []morpc.ClientOption{
		morpc.WithClientInitBackends([]string{backend}, []int{1}),
		morpc.WithClientMaxBackendPerHost(10),
		morpc.WithClientLogger(logutil.GetGlobalLogger()),
	}
	cli, err := morpc.NewClient("plugin-client", bf, clientOpts...)
	if err != nil {
		return nil, err
	}
	return &rpcPlugin{client: cli, backend: backend, timeout: timeout}, nil
}

func (p *rpcPlugin) RecommendCN(ctx context.Context, ci clientInfo) (*plugin.Recommendation, error) {

	resp, err := p.request(ctx, &plugin.Request{ClientInfo: &plugin.ClientInfo{
		Tenant:        string(ci.Tenant),
		Username:      ci.username,
		OriginIP:      ci.originIP.String(),
		LabelSelector: ci.labelInfo.allLabels(),
	}})
	if err != nil {
		return nil, err
	}
	return resp.Recommendation, nil
}

func (p *rpcPlugin) Close() error {
	return p.client.Close()
}

func (p *rpcPlugin) request(ctx context.Context, req *plugin.Request) (*plugin.Response, error) {
	cc, cancel := context.WithTimeoutCause(ctx, p.timeout, moerr.CauseRequest)
	defer cancel()
	f, err := p.client.Send(cc, p.backend, req)
	if err != nil {
		return nil, moerr.AttachCause(ctx, err)
	}
	defer f.Close()
	resp, err := f.Get()
	if err != nil {
		return nil, err
	}
	return resp.(*plugin.Response), nil
}
