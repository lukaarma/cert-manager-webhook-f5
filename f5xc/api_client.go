package f5xc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"k8s.io/klog/v2"
)

const (
	baseUrl = "https://%s.console.ves.volterra.io/api"
)

const (
	add_key    updateMode = iota
	delete_key updateMode = iota
)

func NewClient(tenantName string, apiKey string) *f5xcClient {
	return &f5xcClient{
		BaseURL: fmt.Sprintf(baseUrl, tenantName),
		ApiKey:  apiKey,
		Client:  &http.Client{},
	}

}

func (c *f5xcClient) send(req *http.Request, resData interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("APIToken %s", c.ApiKey))

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
			klog.Error("Error decoding JSON response!")

			return err
		}

		return nil

	case http.StatusNotFound:
		resData = nil

		return nil

	case http.StatusUnauthorized:
		return fmt.Errorf("Credentials not valid!")

	case http.StatusForbidden:
		return fmt.Errorf("Permission missing for requested resource!")
	}

	klog.Errorf("Unexpected response!")

	var unexpectedRes interface{}
	if err = json.NewDecoder(res.Body).Decode(&unexpectedRes); err != nil {
		klog.Error("Error decoding JSON response!")

		return err
	}
	klog.Errorf("%+v", unexpectedRes)

	return fmt.Errorf("Invalid status code: %d", res.StatusCode)
}

func (c *f5xcClient) getTXTResourceRecord(zone string, rrgroup string, rrname string) (*f5xcTXTResouceRecord, error) {
	const endpoint = "config/dns/namespaces/system/dns_zones/%s/rrsets/%s/%s/TXT"

	klog.Infof("Getting record %q from group %q for zone %q", rrname, rrgroup, zone)

	url, err := url.JoinPath(c.BaseURL, fmt.Sprintf(endpoint, zone, rrgroup, rrname))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := &f5xcTXTResouceRecord{}

	err = c.send(req, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *f5xcClient) createTXTResourceRecord(zone string, rrgroup string, rrname string, key string) (*f5xcTXTResouceRecord, error) {
	const endpoint = "config/dns/namespaces/system/dns_zones/%s/rrsets/%s"

	url, err := url.JoinPath(c.BaseURL, fmt.Sprintf(endpoint, zone, rrgroup))
	if err != nil {
		return nil, err
	}

	reqBody := &f5xcTXTResouceRecordCreation{
		DnsZoneName: zone,
		GroupName:   rrgroup,
		RRSet: f5xcTXTrrset{
			TTL: 60,
			TXTRecord: f5xcTXTRecord{
				Name:   rrname,
				Values: []string{key},
			},
		},
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, err
	}

	res := &f5xcTXTResouceRecord{}

	err = c.send(req, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *f5xcClient) updateTXTResourceRecord(zone string, rrgroup string, rrname string, reqBody *f5xcTXTResouceRecord, key string, mode updateMode) (*f5xcTXTResouceRecord, error) {
	const endpoint = "config/dns/namespaces/system/dns_zones/%s/rrsets/%s/%s/TXT"

	url, err := url.JoinPath(c.BaseURL, fmt.Sprintf(endpoint, zone, rrgroup, rrname))
	if err != nil {
		return nil, err
	}

	switch mode {
	case add_key:
		reqBody.RRSet.TXTRecord.Values = append(reqBody.RRSet.TXTRecord.Values, key)
	case delete_key:
		for i, v := range reqBody.RRSet.TXTRecord.Values {
			if v == key {
				reqBody.RRSet.TXTRecord.Values = append(reqBody.RRSet.TXTRecord.Values[:i], reqBody.RRSet.TXTRecord.Values[i+1:]...)
			}
		}
	default:
		return nil, fmt.Errorf("Invalid mode %q", mode)
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, err
	}

	res := &f5xcTXTResouceRecord{}

	err = c.send(req, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *f5xcClient) deleteTXTResourceRecord(zone string, rrgroup string, rrname string) error {
	const endpoint = "config/dns/namespaces/system/dns_zones/%s/rrsets/%s/%s/TXT"

	url, err := url.JoinPath(c.BaseURL, fmt.Sprintf(endpoint, zone, rrgroup, rrname))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	res := &f5xcTXTResouceRecord{}

	err = c.send(req, &res)
	if err != nil {
		return err
	}

	return nil
}
