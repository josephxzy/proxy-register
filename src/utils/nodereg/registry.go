package nodereg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

// InvalidHosts 不允许作为节点域名的地址列表
var InvalidHosts = map[string]bool{
	"127.0.0.1": true,
	"localhost": true,
	"0.0.0.0":   true,
	"":          true,
}

// Instance 注册表中的单个实例记录
type Instance struct {
	Token        string    // 节点身份令牌
	Domain       string    // 节点域名
	RegisteredAt time.Time // 注册时间
	LastSeen     time.Time // 最后活跃时间
	Challenge    string    // 挑战值，用于域名所有权验证
	TokenExpiry  time.Time // Token过期时间
}

// NodeEntry 返回给客户端的节点信息
type NodeEntry struct {
	Name string `json:"name"` // 节点域名
	URL  string `json:"url"`  // 节点访问地址（不含协议头的域名）
}

// Registry 节点注册表，管理所有已注册的节点实例
type Registry struct {
	mu       sync.RWMutex
	byToken  map[string]*Instance // 按Token索引
	byDomain map[string]*Instance // 按域名索引
	version  int                  // 数据版本号
}

// NewRegistry 创建新的节点注册表
func NewRegistry() *Registry {
	return &Registry{
		byToken:  make(map[string]*Instance),
		byDomain: make(map[string]*Instance),
	}
}

// FindByToken 根据Token查找实例
func (r *Registry) FindByToken(token string) (*Instance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inst, ok := r.byToken[token]
	return inst, ok
}

// FindByDomain 根据域名查找实例
func (r *Registry) FindByDomain(domain string) (*Instance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inst, ok := r.byDomain[domain]
	return inst, ok
}

// Register 注册新实例或更新已有实例
// 如果域名已存在，先移除旧Token，再注册新实例
func (r *Registry) Register(inst *Instance) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if old, ok := r.byDomain[inst.Domain]; ok {
		if old.Token != inst.Token {
			delete(r.byToken, old.Token)
		}
	}

	r.byToken[inst.Token] = inst
	r.byDomain[inst.Domain] = inst
	r.version++
}

// ListNodes 返回所有节点列表，排除指定域名
// 用于向各节点返回其他节点的信息
func (r *Registry) ListNodes(excludeDomain string) []NodeEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make([]NodeEntry, 0, len(r.byDomain))
	for domain := range r.byDomain {
		if domain == excludeDomain {
			continue
		}
		entries = append(entries, NodeEntry{
			Name: domain,
			URL:  domain,
		})
	}

	return entries
}

func (r *Registry) ListAllNodes() []NodeEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries := make([]NodeEntry, 0, len(r.byDomain))
	for domain := range r.byDomain {
		entries = append(entries, NodeEntry{
			Name: domain,
			URL:  domain,
		})
	}

	return entries
}

// GetVersion 获取当前数据版本号
func (r *Registry) GetVersion() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.version
}

// Cleanup 清理超时未活跃的实例
func (r *Registry) Cleanup(timeout time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for token, inst := range r.byToken {
		if now.Sub(inst.LastSeen) > timeout {
			delete(r.byDomain, inst.Domain)
			delete(r.byToken, token)
			r.version++
		}
	}
}

// GenerateToken 生成随机Token，用于节点身份认证
func GenerateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateChallenge 根据Token生成挑战值，用于域名所有权验证
func GenerateChallenge(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// IsValidDomain 检查域名是否合法（仅排除本地地址和空值）
func IsValidDomain(domain string) bool {
	if domain == "" {
		return false
	}
	domain = strings.TrimSpace(domain)
	return !InvalidHosts[strings.ToLower(domain)]
}
