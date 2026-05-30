// Package handlers 处理HTTP请求
package handlers

import (
	"encoding/json"
	"net/http"

	"github-proxy-registry/models"
)

// GetNodes 处理节点列表查询请求
// 节点端调用此接口获取其他节点的信息，使用Token认证
func GetNodes(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is required", http.StatusUnauthorized)
		return
	}

	inst, ok := models.GlobalRegistry.FindByToken(token)
	if !ok {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	nodes := models.GlobalRegistry.ListNodes(inst.Domain)
	version := models.GlobalRegistry.GetVersion()

	resp := map[string]interface{}{
		"version": version,
		"nodes":   nodes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetPublicNodes 公开节点列表查询，无需认证，供前端展示
func GetPublicNodes(w http.ResponseWriter, r *http.Request) {
	nodes := models.GlobalRegistry.ListAllNodes()

	resp := map[string]interface{}{
		"nodes": nodes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
