package vulnerabilities

import (
	"os"
	"strings"
	"testing"

	"github.com/facebookincubator/nvdtools/cpedict"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mock"
	"github.com/stretchr/testify/require"
)

const xmlCPEDict = `
<?xml version='1.0' encoding='UTF-8'?>
<cpe-list xmlns:config="http://scap.nist.gov/schema/configuration/0.1" xmlns="http://cpe.mitre.org/dictionary/2.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:scap-core="http://scap.nist.gov/schema/scap-core/0.3" xmlns:cpe-23="http://scap.nist.gov/schema/cpe-extension/2.3" xmlns:ns6="http://scap.nist.gov/schema/scap-core/0.1" xmlns:meta="http://scap.nist.gov/schema/cpe-dictionary-metadata/0.2" xsi:schemaLocation="http://scap.nist.gov/schema/cpe-extension/2.3 https://scap.nist.gov/schema/cpe/2.3/cpe-dictionary-extension_2.3.xsd http://cpe.mitre.org/dictionary/2.0 https://scap.nist.gov/schema/cpe/2.3/cpe-dictionary_2.3.xsd http://scap.nist.gov/schema/cpe-dictionary-metadata/0.2 https://scap.nist.gov/schema/cpe/2.1/cpe-dictionary-metadata_0.2.xsd http://scap.nist.gov/schema/scap-core/0.3 https://scap.nist.gov/schema/nvd/scap-core_0.3.xsd http://scap.nist.gov/schema/configuration/0.1 https://scap.nist.gov/schema/nvd/configuration_0.1.xsd http://scap.nist.gov/schema/scap-core/0.1 https://scap.nist.gov/schema/nvd/scap-core_0.1.xsd">
  <generator>
    <product_name>National Vulnerability Database (NVD)</product_name>
    <product_version>4.6</product_version>
    <schema_version>2.3</schema_version>
    <timestamp>2021-07-20T03:50:36.509Z</timestamp>
  </generator>
  <cpe-item name="cpe:/a:vendor:product:1.2.3:~~~macos~~">
    <title xml:lang="en-US">Vendor Product 1.2.3 for MacOS</title>
    <references>
      <reference href="https://someurl.com">Change Log</reference>
    </references>
    <cpe-23:cpe23-item name="cpe:2.3:a:vendor:product:1.2.3:*:*:*:*:macos:*:*"/>
  </cpe-item>
  <cpe-item name="cpe:/a:vendor2:product2:0.3:~~~macos~~" deprecated="true" deprecation_date="2021-06-10T15:28:05.490Z">
    <title xml:lang="en-US">Vendor2 Product2 0.3 for MacOS</title>
    <cpe-23:cpe23-item name="cpe:2.3:a:vendor2:product2:0.3:*:*:*:*:macos:*:*">
      <cpe-23:deprecation date="2021-06-10T11:28:05.490-04:00">
        <cpe-23:deprecated-by name="cpe:2.3:a:vendor2:product3:1.2:*:*:*:*:macos:*:*" type="NAME_CORRECTION"/>
      </cpe-23:deprecation>
    </cpe-23:cpe23-item>
  </cpe-item>
  <cpe-item name="cpe:/a:vendor2:product3:1.2:~~~macos~~" deprecated="true" deprecation_date="2021-06-10T15:28:05.490Z">
    <title xml:lang="en-US">Vendor2 Product3 1.2 for MacOS</title>
    <cpe-23:cpe23-item name="cpe:2.3:a:vendor2:product3:1.2:*:*:*:*:macos:*:*">
      <cpe-23:deprecation date="2021-06-10T11:28:05.490-04:00">
        <cpe-23:deprecated-by name="cpe:2.3:a:vendor2:product4:999:*:*:*:*:macos:*:*" type="NAME_CORRECTION"/>
      </cpe-23:deprecation>
    </cpe-23:cpe23-item>
  </cpe-item>
  <cpe-item name="cpe:/a:vendor2:product4:999">
    <title xml:lang="en-US">Vendor2 Product4 999 for MacOS</title>
    <cpe-23:cpe23-item name="cpe:2.3:a:vendor2:product4:999:*:*:*:*:macos:*:*"/>
  </cpe-item>
</cpe-list>
`

func TestCpeFromSoftware(t *testing.T) {
	tempDir := os.TempDir()

	ds := new(mock.Store)
	ds.AppConfigFunc = func() (*fleet.AppConfig, error) {
		return &fleet.AppConfig{VulnerabilityDatabasesPath: &tempDir}, nil
	}

	require.NoError(t, GenerateCPEDatabaseSkeleton(tempDir))

	items, err := cpedict.Decode(strings.NewReader(xmlCPEDict))
	require.NoError(t, err)

	err = GenerateCPEDB(tempDir, items)
	require.NoError(t, err)

	// checking an non existent version returns empty
	cpe, err := CPEFromSoftware(ds, &fleet.Software{Name: "Vendor Product.app", Version: "2.3.4", Source: "apps"})
	require.NoError(t, err)
	require.Equal(t, "", cpe)

	// checking a version that exists works
	cpe, err = CPEFromSoftware(ds, &fleet.Software{Name: "Vendor Product.app", Version: "1.2.3", Source: "apps"})
	require.NoError(t, err)
	require.Equal(t, "cpe:2.3:a:vendor:product:1.2.3:*:*:*:*:macos:*:*", cpe)

	// follows many deprecations
	cpe, err = CPEFromSoftware(ds, &fleet.Software{Name: "Vendor2 Product2.app", Version: "0.3", Source: "apps"})
	require.NoError(t, err)
	require.Equal(t, "cpe:2.3:a:vendor2:product4:999:*:*:*:*:macos:*:*", cpe)
}

func TestSyncCPEDatabase(t *testing.T) {
}
