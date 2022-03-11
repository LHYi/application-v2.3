# application-v2.3

This repo contains the applications developed for the response credit scenario with Fabric v2.3.

## Running the GO program directly

The usage of the applications should be based on an Hyperledger Fabric blockchain network with response credit contract implemented.

Once the blockchain network has been brought up and the chaincode has been successfully committed, one can run the following command to start the application.

``` shell
cd rc-VPPO
go run rc-VPPO.go
```

## Running the applications as executable files

Compile the VPPO application with the go build command.

```shell
cd rc-VPPO
go build
```

Once the program is successfully built, there will exist an executable file named "application-gateway-VPPO" in the folder (it will be an .exe file in Windows). Run it and you will be able to interact with the application via terminal.

## Interacting with the blockchain network as VPPO

In this section you will learn how to use the VPPO application.

### Environment setup

The blockchain network can either be brought up within a single machine (local host) or within a network consisting of multiple machines.

If the blockchain network is implemented in a single machine, setting ```DISCOVERY_AS_LOCALHOST``` to ```true```. Otherwise, set ```DISCOVERY_AS_LOCALHOST``` to ```false```.

The choice of the above setting will be asked at the beginning of the application.

### Enroll

At current stage, since the application is tested with the test network provided by fabric samples, VPPO is interacting with the blockchain network as Organization 1. Please use ```appUser``` as the username. This can be changed accordingly later when implemented within a customized network.

The application will first ask for the username, enter ```appUser``` and confirm your input with ```y```.

### Connecting to the gateway

The application will connect to the gateway with the pre-specified peer endpoint (an IP address or a localhost address), using the identity files generated by the MSP.

The network name is the name of the channel and the contract name is the name of the chaincode. The default channel name in the test network is ```mychannel``` and the default contract name is ```basic``` according to the setting of the smart contract tutorial.

## Invoking chaincode functions

There are several chaincode functions that could be invoked by the VPPO. Enter ```help``` to print a table of the chaincode functions available.

## Cleaning up

The enrollment process will generate two directories, named ```wallet``` and ```keystore``` respectively, to put the credential files related to the current user. If you want them to be deleted, enter ```y``` when ending the application, then these directories will be recursively removed.
