// Package handlers 处理HTTP请求
package handlers

import (
	"encoding/json"
	"github-proxy-registry/models"
	"github-proxy-registry/utils/nodereg"
	"net/http"
	"time"
)

// Register 处理节点注册请求
// 1. 验证域名合法性
// 2. 新节点生成Token+Challenge，已有节点复用Token并重新生成Challenge
// 3. 返回Token、Challenge和其他节点列表
func Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Domain == "" {
		http.Error(w, "domain is required", http.StatusBadRequest)
		return
	}

	if !nodereg.IsValidDomain(req.Domain) {
		http.Error(w, "invalid domain", http.StatusBadRequest)
		return
	}

	var token string
	var challenge string

	existingInst, exists := models.GlobalRegistry.FindByDomain(req.Domain)
	if exists {
		token = existingInst.Token
		challenge = nodereg.GenerateChallenge(existingInst.Token)
	} else {
		token = nodereg.GenerateToken()
		challenge = nodereg.GenerateChallenge(token)
	}

	inst := &nodereg.Instance{
		Token:        token,
		Domain:       req.Domain,
		RegisteredAt: time.Now(),
		LastSeen:     time.Now(),
		Challenge:    challenge,
		TokenExpiry:  time.Now().Add(24 * time.Hour),
	}
	models.GlobalRegistry.Register(inst)

	nodes := models.GlobalRegistry.ListNodes(req.Domain)

	resp := map[string]interface{}{
		"token":     token,
		"challenge": challenge,
		"nodes":     nodes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
