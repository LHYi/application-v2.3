package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	// the mspID should be identical to the one used when calling cryptogen to generate credential files
	// mspID         = "Org1MSP"
	// the path of the certificates
	cryptoPath  = "../../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath    = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath     = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	// an IP address to access the peer node, it is a localhost address when the network is running in a single machine
	peerEndpoint = "localhost:7051"
	// name of the peer node
	gatewayPeer = "peer0.org1.example.com"
	// the channel name and the chaincode name should be identical to the ones used in blockchain network implementation, the following are the default values
	// these information have been designed to be entered by the application user
	// channelName   = "mychannel"
	// chaincodeName = "basic"
)

func main() {
	log.Println("============ application-golang starts ============")
	log.Println("============ The application will end when you enter exit ============")
	// DISCOVERY_AS_LOCALHOST should be set to "false" if the network is deployed on other computers
	for {
		log.Println("============ setting DISCOVERY_AS_LOCALHOST ============")
		fmt.Print("-> Do you want to set DISCOVERY_AS_LOCALHOST to true? [y/n]: ")
		// catchOneInput() catches one line of the terminal input, see more details in function definition
		DAL := catchOneInput()
		// determining whether DAL is yes or no and conduct corresponding operations
		if isNo(DAL) {
			log.Println("-> Setting DISCOVERY_AS_LOCALHOST to false")
			err := os.Setenv("DISCOVERY_AS_LOCALHOST", "false")
			if err != nil {
				log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
				// an exit code that is nonzero indicates that there exist an error
				os.Exit(1)
			}
			log.Println("-> Success")
			break
		} else if isYes(DAL) {
			log.Println("-> Setting DISCOVERY_AS_LOCALHOST to true")
			err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
			if err != nil {
				log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
				os.Exit(1)
			}
			log.Println("-> Success")
			break
		} else {
			log.Println("-> Wrong input, please try again or input exit")
		}
	}

	log.Println("============ Creating wallet ============")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	log.Println("============ Wallet created ============")

	// define the variable outside the loop so that it can be used in the following connection configuration
	var userName string
	// label of the code block is useful when the process of running the code is relatec to the user's selection
	// labels are only used with goto, continue and break
userNameLoop:
	for {
		log.Println("-> Please enter your username:")
		userName = catchOneInput()
	userNameConfirmLoop:
		for {
			// formatted output, %s prints out the value of a string
			fmt.Printf("-> Please confirm your username is %s, [y/n]: ", userName)
			userNameConfirm := catchOneInput()
			if isYes(userNameConfirm) {
				break userNameLoop
			} else if isNo(userNameConfirm) {
				break userNameConfirmLoop
			} else {
				fmt.Println("->Wrong input! Please try again.")
			}
		}
	}
	log.Printf("-> Your username is %s.", userName)
	if !wallet.Exists(userName) {
		err = populateWallet(wallet, userName)
		if err != nil {
			log.Fatalf("->Failed to populate wallet contents: %v", err)
		}
		log.Printf("-> Successfully add user %s to wallet \n", userName)
	} else {
		log.Printf("->  User %s already exists", userName)
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"fabric-samples-2.3",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)
	log.Println("============ connecting to gateway ============")
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	log.Println("============ Successfully connected to gateway ============")

	// the logic is similar to the code above
	var networkName string
	log.Println("============ connecting to network ============")
networkNameLoop:
	for {
		log.Println("-> Please enter the name of the network:")
		networkName = catchOneInput()
	networkNameConfirmLoop:
		for {
			fmt.Printf("-> Please confirm your network name is: %s, [y/n]: ", networkName)
			networkNameConfirm := catchOneInput()
			if isYes(networkNameConfirm) {
				break networkNameLoop
			} else if isNo(networkNameConfirm) {
				break networkNameConfirmLoop
			} else {
				fmt.Println("->Wrong input! Please try again.")
			}
		}
	}
	log.Printf("-> Your network name is %s.", networkName)

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	log.Println("============ successfully connected to network", networkName, "============")

	// the logic is similar to the code above
	var contractName string
	log.Println("============ getting contract ============")
contractNameLoop:
	for {
		log.Println("-> Please enter the name of the contract:")
		contractName = catchOneInput()
	contractNameConfirmLoop:
		for {
			fmt.Printf("-> Please confirm your contract name is: %v, [y/n]: ", contractName)
			contractNameConfirm := catchOneInput()
			if isYes(contractNameConfirm) {
				break contractNameLoop
			} else if isNo(contractNameConfirm) {
				break contractNameConfirmLoop
			} else {
				fmt.Println("->Wrong input! Please try again.")
			}
		}
	}
	log.Printf("-> Your contract name is %s.", contractName)
	contract := network.GetContract("basic")
	log.Println("============ successfully got contract", contractName, "============")
	for {
		fmt.Println("-> Please enter the name of the smart contract function you want to invoke, enter help to print the functions available")
		scfunction := catchOneInput()
		invokeChaincode(contract, scfunction, userName)
		// here provides another way to exit the application after every invocation of the smart contract function
	scContinueConfirmLoop:
		for {
			fmt.Print("Do you want to continue? [y/n]: ")
			continueConfirm := catchOneInput()
			if isYes(continueConfirm) {
				fmt.Println("Preparing for invoking next smart contract function")
				break scContinueConfirmLoop
			} else if isNo(continueConfirm) {
				fmt.Print("Do you want to clean up the wallet? [y/n]: ")
				cleanUpConfirm := catchOneInput()
				if isYes(cleanUpConfirm) {
					log.Println("-> Cleaning up wallet...")
					if _, err := os.Stat("wallet"); err == nil {
						e := os.RemoveAll("wallet")
						if e != nil {
							log.Fatal(e)
						}
					}
					if _, err := os.Stat("keystore"); err == nil {
						e := os.RemoveAll("keystore")
						if e != nil {
							log.Fatal(e)
						}
					}
					log.Println("-> Wallet cleaned up successfully")
				}
				exitApp()
			} else {
				fmt.Println("->Wrong input! Please try again.")
				continue scContinueConfirmLoop
			}
		}
	}
}

// TODO: this function can be further seperated into several functions in the future
func invokeChaincode(contract *gateway.Contract, scfunction string, userName string) {
	// the defer function is important in handling error
	// once the chaincode invocation is unsuccessful, the panic function will be called and the recover function which is deferred before the error exists
	// will allow the program to print out the error and recover to the next line after the line that caused error
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Occured an error while invoking chiancode function: %v...Recovered, please try again.\n", r)
		}
	}()
	switch scfunction {
	case "instantiate", "Instantiate", "INSTANTIATE":
		instantiate(contract)
	case "issue", "Issue", "ISSUE":
		log.Println("============ Issuing a new credit ============")
		// in the issuing process, functions are added to allow the credit details to be automatically generated
	issueLoop:
		for {
			var creditNumber string
		enterCreditNumberLoop:
			for {
				fmt.Println("-> Do you want to assign a specific credit number? [y/n]: ")
				enterConfirm := catchOneInput()
				if isYes(enterConfirm) {
					fmt.Println("-> Please enter the credit number:")
					creditNumber = catchOneInput()
					fmt.Println("-> The credit number you entered is: " + creditNumber)
					break enterCreditNumberLoop
				} else if isNo(enterConfirm) {
					fmt.Println("-> Generating credit number.")
					creditNumber = generateCreditNumber()
					fmt.Println("-> The credit number automatically generated is: " + creditNumber)
					break enterCreditNumberLoop
				} else {
					fmt.Println("->Wrong input! Please try again.")
				}
			}
			var issuer string
		enterIssuerLoop:
			for {
				fmt.Println("-> Do you want to use your username as issuer? [y/n]: ")
				enterConfirm := catchOneInput()
				if isYes(enterConfirm) {
					fmt.Println("-> Using your username as the issuer")
					issuer = userName
					fmt.Println("-> The issuer is: " + userName)
					break enterIssuerLoop
				} else if isNo(enterConfirm) {
					fmt.Println("-> Please enter the issuer: ")
					issuer = catchOneInput()
					fmt.Println("-> The issuer you entered is: " + issuer)
					break enterIssuerLoop
				} else {
					fmt.Println("->Wrong input! Please try again.")
				}
			}
			var issueDateTime string
		enterCreditDateTimeLoop:
			for {
				fmt.Println("-> Do you want to generate the issue date and time of the credit automatically? [y/n]: ")
				enterConfirm := catchOneInput()
				if isYes(enterConfirm) {
					fmt.Println("-> Getting date and time.")
					issueDateTime = generateCreditDateTime()
					fmt.Println("-> The date and time is: " + issueDateTime)
					break enterCreditDateTimeLoop
				} else if isNo(enterConfirm) {
					fmt.Println("-> Please enter the issue date and time:")
					issueDateTime = catchOneInput()
					fmt.Println("-> The issue date and time you entered is: " + issueDateTime)
					break enterCreditDateTimeLoop
				} else {
					fmt.Println("->Wrong input! Please try again.")
				}
			}
		issueConfirmLoop:
			for {
				fmt.Printf("-> Are these inputs correct? [y/n]: ")
				issueConfirm := catchOneInput()
				if isYes(issueConfirm) {
					issue(contract, creditNumber, issuer, issueDateTime)
					break issueLoop
				} else if isNo(issueConfirm) {
					fmt.Println("-> Please enter the details of the credit to issue again.")
					break issueConfirmLoop
				} else {
					fmt.Println("->Wrong input! Please try again.")
				}
			}
		}
	case "query", "Query", "QUERY":
		log.Println("============ Querying a credit ============")
	queryLoop:
		for {
			fmt.Println("-> Please enter the credit number:")
			creditNumber := catchOneInput()
			fmt.Println("-> The credit number you entered is: " + creditNumber)
			fmt.Println("-> Please enter the issuer:")
			issuer := catchOneInput()
			fmt.Println("-> The issuer you entered is: " + issuer)
		queryConfirmLoop:
			for {
				fmt.Printf("-> Are these inputs correct? [y/n]: ")
				queryConfirm := catchOneInput()
				if isYes(queryConfirm) {
					query(contract, creditNumber, issuer)
					break queryLoop
				} else if isNo(queryConfirm) {
					fmt.Println("-> Please enter the details of the credit to query again.")
					break queryConfirmLoop
				} else {
					fmt.Println("->Wrong input! Please try again.")
				}
			}
		}
	case "help", "HELP", "Help", "":
		listFuncs()
	default:
		fmt.Println("->Wrong input! Please try again!")
	}
}

// instantiate function do nothing, but it can be used to verify whether the connection is successful before interacting with the ledger
func instantiate(contract *gateway.Contract) {
	log.Println("Submit Transaction: Instantiate, function calls the instantiate function, with no effect.")

	_, err := contract.SubmitTransaction("Instantiate")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully!\n")
}

// Issuing a new response credit
// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func issue(contract *gateway.Contract, creditNumber string, issuer string, issueDateTime string) {
	log.Println("Submit Transaction: IssueCredit, creates new response credit with credit issuer, credit number and credit issueDateTime.")
	// submit transaction is usually used in the case where an update of the ledger is required
	_, err := contract.SubmitTransaction("Issue", creditNumber, issuer, issueDateTime)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// querying an existing credit
func query(contract *gateway.Contract, creditNumber string, issuer string) {
	fmt.Printf("Evaluate Transaction: QueryCredit, function returns credit attributes\n")
	// evaluate transaction is usually used in the case where only querying the world state is required
	evaluateResult, err := contract.EvaluateTransaction("Query", creditNumber, issuer)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

func populateWallet(wallet *gateway.Wallet, userName string) error {
	credPath := filepath.Join(
		"..",
		"..",
		"fabric-samples-2.3",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

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

	return wallet.Put(userName, identity)
}

// Format JSON data for pretty printing credit details in JSON format
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}

// returns the confirmation
func isYes(s string) bool {
	return strings.Compare(s, "Y") == 0 || strings.Compare(s, "y") == 0 || strings.Compare(s, "Yes") == 0 || strings.Compare(s, "yes") == 0
}

func isNo(s string) bool {
	return strings.Compare(s, "N") == 0 || strings.Compare(s, "n") == 0 || strings.Compare(s, "No") == 0 || strings.Compare(s, "no") == 0
}

func isExit(s string) bool {
	return strings.Compare(s, "Exit") == 0 || strings.Compare(s, "exit") == 0 || strings.Compare(s, "EXIT") == 0
}

// catchOneInput() catches one line of the terminal input, ended with \n, it returns a string where \n is stripped
func catchOneInput() string {
	// instantiate a new reader
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	// get rid of the \n at the end of the string
	s = strings.Replace(s, "\n", "", -1)
	// if the string is exit, exit the application directly
	// this allows the user to exit the application whereever they want and saves the effort of detecting the exit command elsewhere
	if isExit(s) {
		exitApp()
	}
	return s
}

// safely exit application
func exitApp() {
	log.Println("============ application-golang ends ============")
	// exit code zero indicates that no error occurred
	os.Exit(0)
}

// list functions that can be invoked and their arguments
// a individual package is used to pretty print a table, see https://github.com/jedib0t/go-pretty/tree/main/table for details of using this package
func listFuncs() {
	tof := table.NewWriter()
	// directing the output to the system standard output
	tof.SetOutputMirror(os.Stdout)
	// add one row as the table header
	tof.AppendHeader(table.Row{"Commend", "Function discription", "Arguments", "Argument discription"})
	// the beginning of the table content
	tof.AppendRows([]table.Row{
		{"list", "List out all the functions that can be called and the arguments required.", "", ""},
	})
	// add one line of seperators between two rows
	tof.AppendSeparator()
	// multiple lines with no seperators in the middle
	tof.AppendRows([]table.Row{
		{"issue", "The issue function collects the information of a new credit and submit a transaction proposal to the blockchain network to issue a new credit.", "credit number", "Credit number is the unique ID number of a credit."},
		{"", "", "issuer", "Issuer is the unique identity of the entity which issues this credit."},
		{"", "", "issue date and time", "The date and time when the credit is issued."},
	})
	tof.AppendSeparator()
	tof.AppendRows([]table.Row{
		{"query", "The query function collects the information of an existing credit and submit a evaluation proposal to the world state to query the details of that credit.", "credit number", "Credit number is the unique ID number of a credit."},
		{"", "", "issuer", "Issuer is the unique identity of the entity which issues this credit."},
	})
	// print out the formatted table
	tof.Render()
}

// generating the credit number and date time with current time
// TODO: needs to be modified according to the naming rule of the credit
func generateCreditNumber() string {
	var now = time.Now()
	creditNumber := []string{"Credit", fmt.Sprint(now.Unix()*1e3 + int64(now.Nanosecond())/1e6)}
	return strings.Join(creditNumber, "-")
}

func generateCreditDateTime() string {
	var now = time.Now()
	creditNumber := []string{"", fmt.Sprint(now.Unix()*1e3 + int64(now.Nanosecond())/1e6)}
	return strings.Join(creditNumber, "")
}
