package appClient

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Client struct {
	User   string
	Wallet *gateway.Wallet
	//userIdentity *gateway.X509Identity
	Gw       *gateway.Gateway
	Network  *gateway.Network
	Contract *gateway.Contract
}

func init() {

}

func NewClient(wp, user, credPath, ccPath, channel, ccName string) *Client {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Println("============ application-golang starts ============")
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "false")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}
	wp = filepath.Join(wp, user)
	wallet, err := gateway.NewFileSystemWallet(wp)
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(user) {
		err = populateWallet(wallet, credPath, user)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	log.Println(ccPath)
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccPath))),
		//gateway.WithUser("user1")
		gateway.WithIdentity(wallet, user),
		//gateway.WithDiscovery(gateway.DiscoveryOptions{Enabled: true, AsLocalhost: false}),

	)

	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	log.Println("test2")
	network, err := gw.GetNetwork(channel)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract := network.GetContract(ccName)

	cli := Client{
		User:     user,
		Wallet:   wallet,
		Gw:       gw,
		Network:  network,
		Contract: contract,
	}
	return &cli
}

func populateWallet(wallet *gateway.Wallet, credPath string, user string) error {
	log.Println("============ Populating wallet ============")

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}
	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put(user, identity)
}

func (c *Client) Submit(fc string, params ...string) {
	log.Println("--> Submit Transaction, called function is " + fc)
	result, err := c.Contract.SubmitTransaction(fc, params...)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println("success to invoke transaction,the txid is " + string(result))
	//SaveTime(params)
}

func (c *Client) Evaluate(fc string, params ...string) {
	log.Println("--> Evaluate Transaction, called function is " + fc)
	result, err := c.Contract.EvaluateTransaction(fc, params...)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))
}
