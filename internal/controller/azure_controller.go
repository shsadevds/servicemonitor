package controller

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	//"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	//"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	//"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"time"
	"context"
)

var (
	//resourceGroupName   = "base"
	interval            = 5 * time.Second
	ctx                 = context.TODO()

	//AZURE_SUB = "9ba71a52-6a1d-485c-b7e5-932e78366990"//订阅id
)

var DNSARecord = make(map[string]string,0)

func GETDNS(){
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	// 创建客户端
	client, err := armprivatedns.NewRecordSetsClient("9ba71a52-6a1d-485c-b7e5-932e78366990", cred, nil)
	if err != nil {
		panic(err)
	}
	getAllRecords(client)
}

func getAllRecords(client *armprivatedns.RecordSetsClient) (int, error) {
	ctx := context.Background()
	pager := client.NewListByTypePager(
		"base",
		"pix.com",
		"A",
		//armprivatedns.RecordTypeA,
		nil,
	)

	countSum := 0
	for pager.More() {
		//var recordList []ARecordSet
		fmt.Println(countSum)
		page, err := pager.NextPage(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get next page: %w", err)
		}

		if len(page.Value) == 0 {
			// 如果当前页没有记录，可能是数据结束，退出循环
			fmt.Println("No more records found.")
			break
		}

		for _, record := range page.Value {
			//只处理A记录
			if record.Properties != nil && record.Properties.ARecords != nil {

				for _, aRecord := range record.Properties.ARecords {
					DNSARecord[*record.Name] =  *aRecord.IPv4Address
				}
			}
			//fmt.Printf("Record: %v\n", record.(*armprivatedns.ARecordSet))
			fmt.Println(DNSARecord)

		}

		countSum++
	}

	return countSum, nil
}

func addDNSRecord(client *armprivatedns.RecordSetsClient,name,record string) {
	ctx := context.Background()
	//判断已有的map是否存在
	if _,ok := DNSARecord[name];!ok{
		// 创建一个新的A记录
		_, err := client.CreateOrUpdate(
			ctx,
			"base",
			"pix.com",
			armprivatedns.RecordTypeA,
			name,
			armprivatedns.RecordSet{
				Properties: &armprivatedns.RecordSetProperties{
					TTL: 10,
					ARecords: []*armprivatedns.ARecord{
						{
							IPv4Address: record,
						},
					},
				},
			},
			nil,
		)


		if err != nil {
			log.Fatalf("failed to add DNS record: %v", err)
		}else{
			//保留
			DNSARecord[name] = record
		}
	}


	fmt.Println("DNS record added successfully")
}


