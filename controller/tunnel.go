package controller

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/openziti-test-kitchen/zrok/build"
	"github.com/openziti-test-kitchen/zrok/controller/store"
	"github.com/openziti-test-kitchen/zrok/model"
	"github.com/openziti-test-kitchen/zrok/rest_model_zrok"
	"github.com/openziti-test-kitchen/zrok/rest_server_zrok/operations/tunnel"
	"github.com/openziti/edge/rest_management_api_client"
	"github.com/openziti/edge/rest_management_api_client/config"
	"github.com/openziti/edge/rest_management_api_client/service"
	"github.com/openziti/edge/rest_management_api_client/service_edge_router_policy"
	"github.com/openziti/edge/rest_management_api_client/service_policy"
	"github.com/openziti/edge/rest_model"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tunnelHandler struct {
}

func newTunnelHandler() *tunnelHandler {
	return &tunnelHandler{}
}

func (h *tunnelHandler) Handle(params tunnel.TunnelParams, principal *rest_model_zrok.Principal) middleware.Responder {
	tx, err := str.Begin()
	if err != nil {
		logrus.Errorf("error starting transaction: %v", err)
		return tunnel.NewTunnelInternalServerError()
	}
	defer func() { _ = tx.Rollback() }()

	envZId := params.Body.ZID
	envId := 0
	if envs, err := str.FindEnvironmentsForAccount(int(principal.ID), tx); err == nil {
		found := false
		for _, env := range envs {
			if env.ZId == envZId {
				logrus.Debugf("found identity '%v' for user '%v'", envZId, principal.Email)
				envId = env.Id
				found = true
				break
			}
		}
		if !found {
			logrus.Errorf("environment '%v' not found for user '%v'", envZId, principal.Email)
			return tunnel.NewTunnelUnauthorized().WithPayload("bad environment identity")
		}
	} else {
		logrus.Errorf("error finding environments for account '%v'", principal.Email)
		return tunnel.NewTunnelInternalServerError()
	}

	edge, err := edgeClient()
	if err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}
	svcName, err := createServiceName()
	if err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}
	cfgId, err := h.createConfig(envZId, svcName, params, edge)
	if err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}
	svcZId, err := h.createService(envZId, svcName, cfgId, edge)
	if err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}
	if err := h.createServicePolicyBind(envZId, svcName, svcZId, envZId, edge); err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}
	if err := h.createServicePolicyDial(envZId, svcName, svcZId, edge); err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}
	if err := h.createServiceEdgeRouterPolicy(envZId, svcName, svcZId, edge); err != nil {
		logrus.Error(err)
		return tunnel.NewTunnelInternalServerError()
	}

	logrus.Debugf("allocated service '%v'", svcName)

	frontendUrl := h.proxyUrl(svcName)
	sid, err := str.CreateService(envId, &store.Service{
		ZId:      svcZId,
		Name:     svcName,
		Frontend: frontendUrl,
		Backend:  params.Body.Endpoint,
	}, tx)
	if err != nil {
		logrus.Errorf("error creating service record: %v", err)
		_ = tx.Rollback()
		return tunnel.NewUntunnelInternalServerError()
	}
	if err := tx.Commit(); err != nil {
		logrus.Errorf("error committing service record: %v", err)
		return tunnel.NewTunnelInternalServerError()
	}
	logrus.Infof("recorded service '%v' with id '%v' for '%v'", svcName, sid, principal.Email)

	return tunnel.NewTunnelCreated().WithPayload(&rest_model_zrok.TunnelResponse{
		ProxyEndpoint: frontendUrl,
		SvcName:       svcName,
	})
}

func (h *tunnelHandler) createConfig(envZId, svcName string, params tunnel.TunnelParams, edge *rest_management_api_client.ZitiEdgeManagement) (cfgID string, err error) {
	authScheme, err := model.ParseAuthScheme(params.Body.AuthScheme)
	if err != nil {
		return "", err
	}
	cfg := &model.ProxyConfig{
		AuthScheme: authScheme,
	}
	if cfg.AuthScheme == model.Basic {
		cfg.BasicAuth = &model.BasicAuth{}
		for _, authUser := range params.Body.AuthUsers {
			cfg.BasicAuth.Users = append(cfg.BasicAuth.Users, &model.AuthUser{Username: authUser.Username, Password: authUser.Password})
		}
	}
	cfgCrt := &rest_model.ConfigCreate{
		ConfigTypeID: &zrokProxyConfigId,
		Data:         cfg,
		Name:         &svcName,
		Tags:         h.zrokTags(svcName),
	}
	cfgReq := &config.CreateConfigParams{
		Config:  cfgCrt,
		Context: context.Background(),
	}
	cfgReq.SetTimeout(30 * time.Second)
	cfgResp, err := edge.Config.CreateConfig(cfgReq, nil)
	if err != nil {
		return "", err
	}
	logrus.Infof("created config '%v' for environment '%v'", cfgResp.Payload.Data.ID, envZId)
	return cfgResp.Payload.Data.ID, nil
}

func (h *tunnelHandler) createService(envZId, svcName, cfgId string, edge *rest_management_api_client.ZitiEdgeManagement) (serviceId string, err error) {
	configs := []string{cfgId}
	encryptionRequired := true
	svc := &rest_model.ServiceCreate{
		Configs:            configs,
		EncryptionRequired: &encryptionRequired,
		Name:               &svcName,
		Tags:               h.zrokTags(svcName),
	}
	req := &service.CreateServiceParams{
		Service: svc,
		Context: context.Background(),
	}
	req.SetTimeout(30 * time.Second)
	resp, err := edge.Service.CreateService(req, nil)
	if err != nil {
		return "", err
	}
	logrus.Infof("created zrok service named '%v' (with ziti id '%v') for environment '%v'", svcName, resp.Payload.Data.ID, envZId)
	return resp.Payload.Data.ID, nil
}

func (h *tunnelHandler) createServicePolicyBind(envZId, svcName, svcZId, envId string, edge *rest_management_api_client.ZitiEdgeManagement) error {
	semantic := rest_model.SemanticAllOf
	identityRoles := []string{fmt.Sprintf("@%v", envId)}
	name := fmt.Sprintf("%v-backend", svcName)
	var postureCheckRoles []string
	serviceRoles := []string{fmt.Sprintf("@%v", svcZId)}
	dialBind := rest_model.DialBindBind
	svcp := &rest_model.ServicePolicyCreate{
		IdentityRoles:     identityRoles,
		Name:              &name,
		PostureCheckRoles: postureCheckRoles,
		Semantic:          &semantic,
		ServiceRoles:      serviceRoles,
		Type:              &dialBind,
		Tags:              h.zrokTags(svcName),
	}
	req := &service_policy.CreateServicePolicyParams{
		Policy:  svcp,
		Context: context.Background(),
	}
	req.SetTimeout(30 * time.Second)
	resp, err := edge.ServicePolicy.CreateServicePolicy(req, nil)
	if err != nil {
		return err
	}
	logrus.Infof("created bind service policy '%v' for service '%v' for environment '%v'", resp.Payload.Data.ID, svcZId, envZId)
	return nil
}

func (h *tunnelHandler) createServicePolicyDial(envZId, svcName, svcZId string, edge *rest_management_api_client.ZitiEdgeManagement) error {
	var identityRoles []string
	for _, proxyIdentity := range cfg.Proxy.Identities {
		identityRoles = append(identityRoles, "@"+proxyIdentity)
		logrus.Infof("added proxy identity role '%v'", proxyIdentity)
	}
	name := fmt.Sprintf("%v-dial", svcName)
	var postureCheckRoles []string
	semantic := rest_model.SemanticAllOf
	serviceRoles := []string{fmt.Sprintf("@%v", svcZId)}
	dialBind := rest_model.DialBindDial
	svcp := &rest_model.ServicePolicyCreate{
		IdentityRoles:     identityRoles,
		Name:              &name,
		PostureCheckRoles: postureCheckRoles,
		Semantic:          &semantic,
		ServiceRoles:      serviceRoles,
		Type:              &dialBind,
		Tags:              h.zrokTags(svcName),
	}
	req := &service_policy.CreateServicePolicyParams{
		Policy:  svcp,
		Context: context.Background(),
	}
	req.SetTimeout(30 * time.Second)
	resp, err := edge.ServicePolicy.CreateServicePolicy(req, nil)
	if err != nil {
		return err
	}
	logrus.Infof("created dial service policy '%v' for service '%v' for environment '%v'", resp.Payload.Data.ID, svcZId, envZId)
	return nil
}

func (h *tunnelHandler) createServiceEdgeRouterPolicy(envZId, svcName, svcZId string, edge *rest_management_api_client.ZitiEdgeManagement) error {
	edgeRouterRoles := []string{"#all"}
	semantic := rest_model.SemanticAllOf
	serviceRoles := []string{fmt.Sprintf("@%v", svcZId)}
	serp := &rest_model.ServiceEdgeRouterPolicyCreate{
		EdgeRouterRoles: edgeRouterRoles,
		Name:            &svcName,
		Semantic:        &semantic,
		ServiceRoles:    serviceRoles,
		Tags:            h.zrokTags(svcName),
	}
	serpParams := &service_edge_router_policy.CreateServiceEdgeRouterPolicyParams{
		Policy:  serp,
		Context: context.Background(),
	}
	serpParams.SetTimeout(30 * time.Second)
	resp, err := edge.ServiceEdgeRouterPolicy.CreateServiceEdgeRouterPolicy(serpParams, nil)
	if err != nil {
		return err
	}
	logrus.Infof("created service edge router policy '%v' for service '%v' for environment '%v'", resp.Payload.Data.ID, svcZId, envZId)
	return nil
}

func (h *tunnelHandler) proxyUrl(svcName string) string {
	return strings.Replace(cfg.Proxy.UrlTemplate, "{svcName}", svcName, -1)
}

func (h *tunnelHandler) zrokTags(svcName string) *rest_model.Tags {
	return &rest_model.Tags{
		SubTags: map[string]interface{}{
			"zrok":              build.String(),
			"zrok-service-name": svcName,
		},
	}
}
