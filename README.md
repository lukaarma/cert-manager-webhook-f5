<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
</p>

# cert-manager-webhook-f5

This is a Cert-manager ACME DNS01 webhook provider for F5 tecnologies, currently only F5 Distributed Cloud (F5XC).

## Development

This is based on the example repository and inspired by two other webhook providers that really helped me to make sense of the structure of the plugin: [cert-manager-webhook-hetzner](https://github.com/vadimkim/cert-manager-webhook-hetzner) and [alidns-webhook](https://github.com/pragkent/alidns-webhook/tree/master)

### Running the test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

**It is essential that you configure and run the test suite when creating a
DNS01 webhook.**

An example configuration has been provided in testdata, the API key must be replaced with a valid one.

You can run the test suite with:

```bash
TEST_ZONE_NAME=example.com. make test
```
