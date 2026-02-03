package f5xc

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
)

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
// For example, `cloudflare` may be used as the name of a solver.
func (solver *F5XCDNSProviderSolver) Name() string {
	return "f5-xc"
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (solver *F5XCDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	fmt.Printf("Config %+v\n", cfg)
	fmt.Printf("Challenge %+v\n", ch)

	// TODO: add code that sets a record in the DNS provider's console
	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (solver *F5XCDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	// TODO: add code that deletes a record from the DNS provider's console
	return nil
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (solver *F5XCDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	k8sClient, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	solver.k8sClient = k8sClient

	return nil
}

func (solver *F5XCDNSProviderSolver) getSecret(namespace string, ref corev1.SecretKeySelector) (string, error) {
	if ref.Name == "" {
		return "", fmt.Errorf("Missing secret ref name")
	}
	if ref.Key == "" {
		return "", fmt.Errorf("Missing secret key name")
	}

	secret, err := solver.k8sClient.CoreV1().Secrets(namespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	bytes, ok := secret.Data[ref.Key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret '%s/%s'", ref.Key, namespace, ref.Name)
	}

	return string(bytes), nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (F5XCDNSProviderConfig, error) {
	cfg := F5XCDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		klog.Infof("Empty config loaded")
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %+v", err)
	}

	return cfg, nil
}
