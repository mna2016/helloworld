package main
import (
			"fmt"
			"encoding/json"
			"github.com/hyperledger/fabric/core/chaincode/shim"
			"github.com/hyperledger/fabric/protos/peer"
			"strings"
)
// Main function starts up the chaincode in the container during instantiate
//
type HelloWorld struct {}

type message struct {

						ID string `json:"ID"`
						Value string `json:"value"`
}



func main() {
			if err := shim.Start(new(HelloWorld)); err != nil {
			fmt.Printf("Main: Error starting HelloWorld chaincode: %s", err)
			}
}


// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data. Note that chaincode upgrade also calls this
// function to reset or to migrate data, so be careful to avoid a scenario
// where you inadvertently clobber your ledger's data!
//
func (t *HelloWorld) Init(stub shim.ChaincodeStubInterface) (peer.Response) {
			// Validate supplied init parameters, in this case zero arguments!
			if _, args := stub.GetFunctionAndParameters(); len(args) > 0 {
			return shim.Error("Init: Incorrect number of arguments; none expected.")
			}
			return shim.Success(nil)
}



func (cc *HelloWorld) Invoke(stub shim.ChaincodeStubInterface) (peer.Response) {
// Which function is been called?
			function, args := stub.GetFunctionAndParameters()
			function = strings.ToLower(function)
			// Route call to the correct function
			switch function {
			case "write": return cc.write(stub, args)
			case "read": return cc.read(stub, args);
			default: return shim.Error("Valid methods are 'write' or 'read'!")
			}
}



func (cc *HelloWorld) write(stub shim.ChaincodeStubInterface, args []string)(peer.Response){
			if len(args) != 2 {
			return shim.Error("Write: incorrect arguments; expecting ID & value.")
			}
			id := strings.ToLower(args[0])
			msg := &message{ID: id, Value:args[1]}
			msgJSON, _ := json.Marshal(msg)
			// Validate that this ID does not yet exist
			if messageAsBytes, err := stub.GetState(id); err != nil || messageAsBytes != nil {
			return shim.Error("Write: this ID already has a message assigned.")
			}
			// Write the message
			if err := stub.PutState(id, msgJSON); err != nil {
			return shim.Error(err.Error())
			} else {
			return shim.Success(nil)
			}
}



func (cc *HelloWorld) read(stub shim.ChaincodeStubInterface, args []string)(peer.Response){

			if len(args) != 1 {
			return shim.Error("Read: incorrect number of arguments; expecting only ID.")
			}
			id := strings.ToLower(args[0])
			if value, err := stub.GetState(id); err != nil || value == nil {
			return shim.Error("Read: invalid ID supplied.")
			} else {
			return shim.Success(value)
			}
}