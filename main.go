/* Copyright 2020 Peiit.com . All Rights Reserved. */
/* main.go -  */
/*
modification history
--------------------
2020/10/16 12:23:15, by jiangpei@peiit.com, create
*/
/*
DESCRIPTION
main.go to do
*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/urfave/cli/v2"
)

const (
	CodeSuc = iota
	CodeGetLocalIp
	CodeAliDnsClient
	CodeAliDnsGetDomainRecords
	CodeAliDnsUpdateDomainRecord
	CodeAliDnsAddDomainRecord
)

func getDomainRecords(domainName string, client *alidns.Client) (records map[string]alidns.Record, err error) {
	records = make(map[string]alidns.Record)
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = domainName

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return records, err
	}
	for _, record := range response.DomainRecords.Record {
		records[record.RR] = record
	}
	return records, err
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (ip net.IP, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP
	return
}

func setLocalIP2DNS(c *cli.Context) {

	//提取参数
	accessKeyId := c.String("ak")
	accessKeySecret := c.String("sk")
	domainName := c.String("d")
	subRR := c.String("r")

	//fmt.Println(accessKeyId, accessKeySecret,domainName,subRR)
	//os.Exit(1)

	outip, err := GetOutboundIP()
	if err != nil {
		log.Fatalln("获取本机ip失败", err)
		os.Exit(CodeGetLocalIp)
	}
	localIpStr := outip.String()

	//创建client
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalln("初始化alidns client失败", err)
		os.Exit(CodeAliDnsClient)
	}

	//获取当前记录
	records, err := getDomainRecords(domainName, client)
	if err != nil {
		log.Fatalln("获取解析记录列表失败", err)
		os.Exit(CodeAliDnsGetDomainRecords)
	}

	if _, ok := records[subRR]; ok {
		//解析记录已经存在,需要更新
		oriRecord := records[subRR]
		if localIpStr == oriRecord.Value {
			//当前解析记录跟本机ip一致
			log.Println(fmt.Sprintf("记录[%s.%s]已经是[%s]", subRR, domainName, localIpStr))
			os.Exit(CodeSuc)
		} else {
			logStr := fmt.Sprintf("更新记录[%s.%s]为[%s]", subRR, domainName, localIpStr)
			log.Println("需要" + logStr)
			request := alidns.CreateUpdateDomainRecordRequest()
			request.Scheme = "https"
			request.RecordId = oriRecord.RecordId
			request.RR = oriRecord.RR
			request.Type = oriRecord.Type
			request.Value = localIpStr

			_, err := client.UpdateDomainRecord(request)
			if err != nil {
				fmt.Print(logStr+"失败", err.Error())
				os.Exit(CodeAliDnsUpdateDomainRecord)
			}
			log.Println(logStr + "成功")
		}
	} else {
		//解析记录不存在,需要新增
		logStr := fmt.Sprintf("新增记录[%s.%s]为[%s]", subRR, domainName, localIpStr)
		log.Println(logStr)
		request := alidns.CreateAddDomainRecordRequest()
		request.Scheme = "https"

		request.DomainName = domainName
		request.RR = subRR
		request.Type = "A"
		request.Value = localIpStr

		_, err := client.AddDomainRecord(request)
		if err != nil {
			fmt.Print(logStr+"失败", err.Error())
			os.Exit(CodeAliDnsAddDomainRecord)
		}
		log.Println(logStr + "成功")
	}
}

func main() {
	app := &cli.App{
		Name:      "myhost",
		Version:   "0.0.1",
		Usage:     "将本机的ip直接设置到域名",
		UsageText: "myhost [--ak 阿里云AK] [--sk 阿里云SK] [-d example.com] [-r mylocalip]",
		HelpName:  "myhost",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required: true,
				Name:     "accesskeyid",
				Aliases:  []string{"ak"},
				Usage:    "阿里云ak",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "accesskeysecret",
				Aliases:  []string{"sk"},
				Usage:    "阿里云sk",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "domainname",
				Aliases:  []string{"d"},
				Usage:    "域名",
			},
			&cli.StringFlag{
				Required: true,
				Name:     "rr",
				Aliases:  []string{"r"},
				Usage:    "主机记录",
			},
			//&cli.StringFlag{
			//	Required : true,
			//	Name:    "type",
			//	Value: "A",
			//	Aliases: []string{"t"},
			//	Usage:   "记录类型,A",
			//},
		},
		Action: func(c *cli.Context) error {
			setLocalIP2DNS(c)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
