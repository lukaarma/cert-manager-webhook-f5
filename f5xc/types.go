package f5xc

import (
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type F5XCDNSProviderSolver struct {
	k8sClient *kubernetes.Clientset
}

type F5XCDNSProviderConfig struct {
	// Change the two fields below according to the format of the configuration
	// to be decoded.
	// These fields will be set by users in the
	// `issuer.spec.acme.dns01.providers.webhook.config` field.

	ApiKeySecretRef corev1.SecretKeySelector `json:"apiKeySecretRef"`
	TenantName      string                   `json:"tenantName"`
	ZoneName        string                   `json:"zoneName"`
	RRGroupName     string                   `json:"rrGroupName"`
}

type F5XCClient struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}
