package modules

import (
	"context"
	"fmt"
	"os"

	"craftweave/core/parser"
	"craftweave/internal/ssh"

	alidns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

func dnsZoneHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.DNSZone == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing dns_zone parameters"}
	}
	zone := task.DNSZone.Zone
	provider := task.DNSZone.Provider
	state := task.DNSZone.State
	if state == "" {
		state = "present"
	}

	var err error

	switch provider {
	case "cloudflare":
		apiToken := os.Getenv("CF_API_TOKEN")
		if apiToken == "" {
			return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing CF_API_TOKEN env"}
		}
		api, err2 := cloudflare.NewWithAPIToken(apiToken)
		if err2 != nil {
			return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err2.Error()}
		}
		ctxbg := context.Background()
		if state == "present" {
			_, err = api.CreateZone(ctxbg, zone, false, cloudflare.Account{}, "")
		} else {
			zoneID, err2 := api.ZoneIDByName(zone)
			if err2 != nil {
				err = err2
			} else {
				_, err = api.DeleteZone(ctxbg, zoneID)
			}
		}
	case "aliyun":
		accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
		secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
		region := os.Getenv("ALICLOUD_REGION_ID")
		if region == "" {
			region = "cn-hangzhou"
		}
		if accessKey == "" || secretKey == "" {
			return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing ALICLOUD_ACCESS_KEY or SECRET_KEY"}
		}
		client, err2 := alidns.NewClientWithAccessKey(region, accessKey, secretKey)
		if err2 != nil {
			return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err2.Error()}
		}
		if state == "present" {
			req := alidns.CreateAddDomainRequest()
			req.DomainName = zone
			_, err = client.AddDomain(req)
		} else {
			req := alidns.CreateDeleteDomainRequest()
			req.DomainName = zone
			_, err = client.DeleteDomain(req)
		}
	default:
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("unsupported provider %s", provider)}
	}

	if err != nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err.Error()}
	}

	return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "OK", ReturnCode: 0, Output: fmt.Sprintf("zone %s %s", state, zone)}
}

func init() { Register("dns_zone", dnsZoneHandler) }
