package bundle

import (
	"errors"
	"fmt"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/bundler"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/cli"
	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cloudflare/cfssl/ubiquity"
)

var bundlerUsageText = // Usage text of 'cfssl bundle'
`cfssl bundle -- create a certificate bundle that contains the client cert

Usage of bundle:
	- Bundle local certificate files
        cfssl bundle -cert file [-ca-bundle file] [-int-bundle file] [-metadata file] [-key keyfile] [-flavor optimal|ubiquitous|force]
	- Bundle certificate from remote server.
        cfssl bundle -domain domain_name [-ip ip_address] [-ca-bundle file] [-int-bundle file] [-metadata file]

Flags:
`

// flags used by 'cfssl bundle'
var bundlerFlags = []string{"cert", "key", "ca-bundle", "int-bundle", "flavor", "metadata", "domain", "ip"}

// bundlerMain is the main CLI of bundler functionality.
func bundlerMain(args []string, c cli.Config) (err error) {
	ubiquity.LoadPlatforms(c.Metadata)
	flavor := bundler.BundleFlavor(c.Flavor)
	// Initialize a bundler with CA bundle and intermediate bundle.
	b, err := bundler.NewBundler(c.CABundleFile, c.IntBundleFile)
	if err != nil {
		return
	}

	var bundle *bundler.Bundle
	if c.CertFile != "" {
		if c.CertFile == "-" {
			var certPEM, keyPEM []byte
			certPEM, err = cli.ReadStdin(c.CertFile)
			if err != nil {
				return
			}
			if c.KeyFile != "" {
				keyPEM, err = cli.ReadStdin(c.KeyFile)
				if err != nil {
					return
				}
			}
			bundle, err = b.BundleFromPEM(certPEM, keyPEM, flavor)
			if err != nil {
				return
			}
		} else {
			// Bundle the client cert
			bundle, err = b.BundleFromFile(c.CertFile, c.KeyFile, flavor)
			if err != nil {
				return
			}
		}
	} else if c.Domain != "" {
		bundle, err = b.BundleFromRemote(c.Domain, c.IP, flavor)
		if err != nil {
			return
		}
	} else {
		return errors.New("Must specify bundle target through -cert or -domain")
	}

	marshaled, err := bundle.MarshalJSON()
	if err != nil {
		return
	}
	fmt.Printf("%s", marshaled)
	return
}

// Command assembles the definition of Command 'bundle'
var Command = &cli.Command{UsageText: bundlerUsageText, Flags: bundlerFlags, Main: bundlerMain}
