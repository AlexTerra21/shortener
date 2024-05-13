package ipchecker

import (
	"net/http"
	"net/netip"

	"github.com/AlexTerra21/shortener/internal/app/config"
)

// Проверка разрешенных IP
func CheckIP(config *config.Config, r *http.Request) (bool, error) {
	trustedSubnet := config.TrustedSubnet
	net, err := netip.ParsePrefix(trustedSubnet)
	if err != nil {
		return false, err
	}
	ipString := r.Header.Get("X-Real-IP")
	ip, err := netip.ParseAddr(ipString)
	if err != nil {
		return false, err
	}
	return net.Contains(ip), nil
}
