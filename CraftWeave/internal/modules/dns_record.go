package modules

import (
	"context"
	"fmt"
	"os"

	"craftweave/core/parser"
	"craftweave/internal/ssh"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	alidns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	cloudflare "github.com/cloudflare/cloudflare-go"
)

func dnsRecordHandler(ctx Context, task parser.Task) ssh.CommandResult {
	if task.DNSRecord == nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: "missing dns_record parameters"}
	}

	zone := task.DNSRecord.Zone
	provider := task.DNSRecord.Provider
	recordType := task.DNSRecord.Type
	name := task.DNSRecord.Name
	value := task.DNSRecord.Value
	ttl := task.DNSRecord.TTL
	state := task.DNSRecord.State
	if state == "" {
		state = "present"
	}
	if ttl == 0 {
		ttl = 300
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
		zoneID, err2 := api.ZoneIDByName(zone)
		if err2 != nil {
			return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err2.Error()}
		}
		z := cloudflare.ZoneIdentifier(zoneID)
		if state == "present" {
			_, err = api.CreateDNSRecord(ctxbg, z, cloudflare.CreateDNSRecordParams{
				Type:    recordType,
				Name:    name,
				Content: value,
				TTL:     ttl,
			})
		} else {
			recs, _, err2 := api.ListDNSRecords(ctxbg, z, cloudflare.ListDNSRecordsParams{Name: name + "." + zone, Type: recordType})
			if err2 != nil {
				return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err2.Error()}
			}
			for _, r := range recs {
				if err = api.DeleteDNSRecord(ctxbg, z, r.ID); err != nil {
					break
				}
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
			req := alidns.CreateAddDomainRecordRequest()
			req.DomainName = zone
			req.RR = name
			req.Type = recordType
			req.Value = value
			req.TTL = requests.NewInteger(ttl)
			_, err = client.AddDomainRecord(req)
		} else {
			req := alidns.CreateDeleteSubDomainRecordsRequest()
			req.DomainName = zone
			req.RR = name
			req.Type = recordType
			_, err = client.DeleteSubDomainRecords(req)
		}
	default:
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: fmt.Sprintf("unsupported provider %s", provider)}
	}

	if err != nil {
		return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "FAILED", ReturnCode: 1, Output: err.Error()}
	}

	return ssh.CommandResult{Host: ctx.Host.Name, ReturnMsg: "OK", ReturnCode: 0, Output: fmt.Sprintf("record %s %s.%s", state, name, zone)}
}

func init() { Register("dns_record", dnsRecordHandler) }
