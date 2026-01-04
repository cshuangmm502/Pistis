package main

import (
	appClient "fabric_client/client"
	"os"
)

const ccPath = "./config/connection-org1.yaml"
const credPath = "./wallet/msp"

//const credPath = "./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp"

func init() {

}

func main() {
	//// 选择组织和用户
	user := "user1"
	cli := appClient.NewClient("wallet", user, credPath, ccPath, "monichain", "basic")
	//cli.Submit("InitLedger")
	cli.Evaluate("GetAllAssets")
	//// 加载连接配置文件
	//walletPath := filepath.Join("wallet")
	////ccpPath := filepath.Join(".", "config", "connection.yaml")
	//fmt.Println(walletPath)
	//// 打开钱包
	//wallet, err := gateway.NewFileSystemWallet(wtpath)
	//if err != nil {
	//	log.Fatalf("Failed to create wallet: %v", err)
	//}
	//

	defer func() {
		if _, err := os.Stat("keystore"); !os.IsNotExist(err) {
			os.RemoveAll("keystore")
		}
	}()
}
