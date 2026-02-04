package f5xc

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

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
	conf, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	apiKey, err := solver.getSecret(ch.ResourceNamespace, conf.ApiKeySecretRef)
	if err != nil {
		return err
	}

	f5xcClient := NewClient(conf.TenantName, apiKey)

	klog.Infof("Got challenge for %s", ch.DNSName)

	resouceRecord, err := f5xcClient.getTXTResourceRecord(conf.ZoneName, conf.RRGroupName, conf.RRName)
	if err != nil {
		return err
	}

	if resouceRecord == nil {
		klog.Infof("No existing TXT record %q found in group %q under zone %q, a new TXT record will be created", resouceRecord.RecordName, conf.RRGroupName, conf.ZoneName)

		resouceRecord, err = f5xcClient.createTXTResourceRecord(conf.ZoneName, conf.RRGroupName, conf.RRName, ch.Key)
	} else {
		klog.Infof("Found existing record %q in group %q under zone %q, updating it", resouceRecord.RecordName, conf.RRGroupName, conf.ZoneName)

		if slices.Contains(resouceRecord.RRSet.TXTRecord.Values, ch.Key) {
			klog.Info("Challenge key already exists in record %q in group %q under zone %q, not applying any change", ch.Key, resouceRecord.RecordName, conf.RRGroupName, conf.ZoneName)

			return nil
		}

		resouceRecord, err = f5xcClient.updateTXTResourceRecord(conf.ZoneName, conf.RRGroupName, conf.RRName, resouceRecord, ch.Key)
	}

	klog.Infof("Created DNS challenge record for %s", ch.DNSName)

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
func loadConfig(confJSON *extapi.JSON) (f5xcDNSProviderConfig, error) {
	conf := f5xcDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if confJSON == nil {
		klog.Infof("Empty config loaded")
		return conf, nil
	}
	if err := json.Unmarshal(confJSON.Raw, &conf); err != nil {
		return conf, fmt.Errorf("error decoding solver config: %+v", err)
	}

	klog.Infof("Config loaded")

	return conf, nil
}
