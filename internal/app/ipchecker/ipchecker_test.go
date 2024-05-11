package ipchecker

import (
	"net/http"
	"strings"
	"testing"

	"github.com/AlexTerra21/shortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тест для CheckIP
func TestCheckIP(t *testing.T) {

	tests := []struct {
		name          string
		trustedSubnet string
		ip            string
		want          bool
	}{
		{
			name:          "trusted",
			trustedSubnet: "192.168.0.0/24",
			ip:            "192.168.0.1",
			want:          true,
		},
		{
			name:          "untrusted",
			trustedSubnet: "192.168.0.0/24",
			ip:            "192.160.0.1",
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
			require.NoError(t, err)

			req.Header.Set("X-Real-IP", tt.ip)

			config := &config.Config{}
			config.TrustedSubnet = tt.trustedSubnet

			isFromTrustedSubnet, _ := CheckIP(config, req)
			assert.Equal(t, tt.want, isFromTrustedSubnet)
		})
	}
}
