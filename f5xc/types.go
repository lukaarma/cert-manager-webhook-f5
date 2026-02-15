package f5xc

import (
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type F5XCDNSProviderSolver struct {
	k8sClient *kubernetes.Clientset
}

type f5xcDNSProviderConfig struct {
	// Change the two fields below according to the format of the configuration
	// to be decoded.
	// These fields will be set by users in the
	// `issuer.spec.acme.dns01.providers.webhook.config` field.

	ApiKeySecretRef corev1.SecretKeySelector `json:"apiKeySecretRef"`
	TenantName      string                   `json:"tenantName"`
	ZoneName        string                   `json:"zoneName"`
	RRGroupName     string                   `json:"rrGroupName"`
	RRName          string                   `json:"rrName"`
}

type updateMode int

type f5xcClient struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

type f5xcTXTResouceRecord struct {
	DnsZoneName string       `json:"dns_zone_name"`
	GroupName   string       `json:"group_name"`
	Namespace   string       `json:"namespace"`
	RecordName  string       `json:"record_name"`
	Type        string       `json:"type"`
	RRSet       f5xcTXTrrset `json:"rrset"`
}

type f5xcTXTrrset struct {
	Description string        `json:"description"`
	TTL         int           `json:"ttl"`
	TXTRecord   f5xcTXTRecord `json:"txt_record"`
}

type f5xcTXTRecord struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type f5xcTXTResouceRecordCreation struct {
	DnsZoneName string       `json:"dns_zone_name"`
	GroupName   string       `json:"group_name"`
	RRSet       f5xcTXTrrset `json:"rrset"`
}
