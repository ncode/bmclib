package discover

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"gitlab.booking.com/go/bmc/c7000"
	"gitlab.booking.com/go/bmc/errors"
	"gitlab.booking.com/go/bmc/idrac8"
	"gitlab.booking.com/go/bmc/idrac9"
	"gitlab.booking.com/go/bmc/m1000e"

	log "github.com/sirupsen/logrus"
	"gitlab.booking.com/go/bmc/devices"
	"gitlab.booking.com/go/bmc/hp"
	"gitlab.booking.com/go/bmc/httpclient"
	"gitlab.booking.com/go/bmc/ilo"
	"gitlab.booking.com/go/bmc/supermicrox10"
)

// ScanAndConnect will scan the bmc trying to learn the device type and return a working connection
func ScanAndConnect(host string, username string, password string) (bmcConnection interface{}, err error) {
	log.WithFields(log.Fields{"step": "ScanAndConnect", "host": host}).Debug("detecting vendor")

	client, err := httpclient.Build()
	if err != nil {
		return bmcConnection, err
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/res/ok.png", host))
	if err != nil {
		return bmcConnection, err
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.WithFields(log.Fields{"step": "ScanAndConnect", "host": host, "vendor": devices.Cloudline}).Debug("it's a discrete")
		return bmcConnection, errors.ErrVendorNotSupported
	}

	resp, err = client.Get(fmt.Sprintf("https://%s/xmldata?item=all", host))
	if err != nil {
		return bmcConnection, err
	}
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bmcConnection, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		iloXMLC := &hp.Rimp{}
		err = xml.Unmarshal(payload, iloXMLC)
		if err != nil {
			return bmcConnection, err
		}

		if iloXMLC.Infra2 != nil {
			log.WithFields(log.Fields{"step": "ScanAndConnect", "host": host, "vendor": devices.HP}).Debug("it's a chassis")
			return c7000.New(host, username, password)
		}

		iloXML := &hp.RimpBlade{}
		err = xml.Unmarshal(payload, iloXML)
		if err != nil {
			return bmcConnection, err
		}

		if iloXML.HSI != nil {
			if strings.HasPrefix(iloXML.MP.Pn, "Integrated Lights-Out") {
				return ilo.New(host, username, password)
			}

			return bmcConnection, fmt.Errorf("it's an HP, but I cound't not identify the hardware type. Please verify")
		}
	}

	resp, err = client.Get(fmt.Sprintf("https://%s/session?aimGetProp=hostname,gui_str_title_bar,OEMHostName,fwVersion,sysDesc", host))
	if err != nil {
		return bmcConnection, err
	}

	payload, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bmcConnection, err
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return idrac8.New(host, username, password)
	}

	resp, err = client.Get(fmt.Sprintf("https://%s/sysmgmt/2015/bmc/info", host))
	if err != nil {
		return bmcConnection, err
	}

	payload, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bmcConnection, err
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return idrac9.New(host, username, password)
	}

	resp, err = client.Get(fmt.Sprintf("https://%s/cgi-bin/webcgi/login", host))
	if err != nil {
		return bmcConnection, err
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.WithFields(log.Fields{"step": "connection", "host": host, "vendor": devices.Dell}).Debug("it's a chassis")
		return m1000e.New(host, username, password)
	}

	resp, err = client.Get(fmt.Sprintf("https://%s/cgi/login.cgi", host))
	if err != nil {
		return bmcConnection, err
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return supermicrox10.New(host, username, password)
	}

	return bmcConnection, errors.ErrVendorUnknown
}
