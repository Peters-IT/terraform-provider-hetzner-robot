package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#boot-configuration

import (
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type BootProfile struct {
	ActiveProfile   string // linux/rescue/...
	Architecture    string
	AuthorizedKeys  []string
	HostKeys        []string
	Language        string
	OperatingSystem string
	Password        string
	ServerID        int
	ServerIPv4      string
	ServerIPv6      string
}

func (c *HetznerRobotClient) getBoot(serverID int) (*BootProfile, error) {
	bytes, err := c.makeAPICall("GET", fmt.Sprintf("%s/boot/%d", c.url, serverID), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}

	jsonStr := string(bytes)
	bootProfile := BootProfile{}
	activeBoot := ""

	if gjson.Get(jsonStr, "boot.linux.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.linux").String()
		bootProfile.ActiveProfile = "linux"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "dist").String()
	}
	if gjson.Get(jsonStr, "boot.rescue.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.rescue").String()
		bootProfile.ActiveProfile = "rescue"
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "os").String()
	}

	bootProfile.Architecture = gjson.Get(activeBoot, "arch").String()
	// bootProfile.AuthorizedKeys = gjson.Get(activeBoot, "authorised_keys").Array()
	// bootProfile.HostKeys = gjson.Get(activeBoot, "host_keys").Array()
	bootProfile.Password = gjson.Get(activeBoot, "password").String()
	bootProfile.ServerID = int(gjson.Get(activeBoot, "server_num").Int())
	bootProfile.ServerIPv4 = gjson.Get(activeBoot, "server_ip").String()
	bootProfile.ServerIPv6 = gjson.Get(activeBoot, "server_ipv6_net").String()

	return &bootProfile, nil
}

func (c *HetznerRobotClient) setBootProfile(serverID int, activeBootProfile string, arch string, os string, lang string, authorizedKeys []string) (*BootProfile, error) {
	formParams := url.Values{}
	formParams.Set("arch", arch)
	for _, key := range authorizedKeys {
		formParams.Add("authorized_key", key)
	}
	if activeBootProfile == "linux" {
		formParams.Set("dist", os)
		formParams.Set("lang", lang)
	}
	if activeBootProfile == "rescue" {
		formParams.Set("os", os)
	}
	encodedParams := formParams.Encode()
	log.Println(encodedParams)

	bytes, err := c.makeAPICall("POST", fmt.Sprintf("%s/boot/%d/%s", c.url, serverID, activeBootProfile), strings.NewReader(encodedParams), http.StatusOK)
	if err != nil {
		if strings.Contains(err.Error(), "BOOT_ALREADY_ENABLED") {
			return c.getBoot(serverID)
		}
		return nil, err
	}

	jsonStr := string(bytes)
	bootProfile := BootProfile{}
	activeBoot := ""

	if gjson.Get(jsonStr, "boot.linux.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.linux").String()
		bootProfile.ActiveProfile = "linux"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "dist").String()
	}
	if gjson.Get(jsonStr, "boot.rescue.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.rescue").String()
		bootProfile.ActiveProfile = "rescue"
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "os").String()
	}

	bootProfile.Architecture = gjson.Get(activeBoot, "arch").String()
	// bootProfile.AuthorizedKeys = gjson.Get(activeBoot, "authorised_keys").Array()
	// bootProfile.HostKeys = gjson.Get(activeBoot, "host_keys").Array()
	bootProfile.Password = gjson.Get(activeBoot, "password").String()
	bootProfile.ServerID = int(gjson.Get(activeBoot, "server_num").Int())
	bootProfile.ServerIPv4 = gjson.Get(activeBoot, "server_ip").String()
	bootProfile.ServerIPv6 = gjson.Get(activeBoot, "server_ipv6_net").String()

	return &bootProfile, nil
}
