/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package shim

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	pb "github.com/openblockchain/obc-peer/protos"
)

type TestChaincode struct {
}

// Used by the test chaincode
var A, B string
var Aval, Bval int

func (t *TestChaincode) Run(stub ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error

	// Handle different functions
	if function == "init" {
		if len(args) != 4 {
			return nil, errors.New("Incorrect number of arguments. Expecting 4")
		} 
	
		// Initialize the chaincode
		A = args[0]
		Aval,err = strconv.Atoi(args[1]) 
		if err != nil {
			return nil, errors.New("Expecting integer value for asset holding")
		}
		B = args[2]
		Bval,err = strconv.Atoi(args[3])
		if err != nil {
			return nil, errors.New("Expecting integer value for asset holding")
		}

		/*
		// Write the state to the ledger
		err = stub.PutState(A, []byte(strconv.Itoa(Aval))
		if err != nil {
			return nil, err
		}

		stub.PutState(B, []byte(strconv.Itoa(Bval))	
		err = stub.PutState(B, []byte(strconv.Itoa(Bval))
		if err != nil {
			return nil, err
		}
		*/
	} else if function == "invoke" {
		// Transaction makes payment of X units from A to B
		X, err = strconv.Atoi(args[0])
		Aval = Aval - X
		Bval = Bval + X
	}
	
	return nil, nil
}

func (t *TestChaincode) Query(stub ChaincodeStub, args []byte) ([]byte, error) {
	return nil, nil
}

func TestChaincode(t *testing.T) {
	// Start the chaincode
	go func() {
		err := Start(new(TestChaincode))
		if err != nil {
			t.Logf("Error Start(ing) chaincode: %s", err)
			t.Fail()
		}
	}

	// Invoke deploy
	payload := `{"funcName" : "init", "args" : ["A", "100", "B", "100"]}`
	msg : = &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INIT, Payload: payload} 
	err := handler.FSM.Event(pb.ChaincodeMessage_INIT.String(), msg)

	// Ensure deploy completes
	time.Sleep(2 * time.Second)

	// Invoke transaction
	payload := `{"funcName" : "init", "args" : ["10"]}`
	msg : = &pb.ChaincodeMessage{Type: pb.ChaincodeMessage_INIT, Payload: payload} 
	err := handler.FSM.Event(pb.ChaincodeMessage_TRANSACTION.String(), msg)
	
	// Ensure deploy completes
	time.Sleep(2 * time.Second)

	// Check result
	if Aval != 90 || Bval != 110 {
		t.Error("Transaction did not execute or incorrect execution")	
	}

	fmt.Printf("Transaction executed successfully")
}
