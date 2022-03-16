package testcommon

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/config"
	contextmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// InstanceTestSmartContract represents the config data for the smart contract instance to be tested
type InstanceTestSmartContract struct {
	testSmartContract
	code []byte
}

// CreateInstanceContract build a contract to be used in a test creted with BuildInstanceCallTest
func CreateInstanceContract(address []byte) *InstanceTestSmartContract {
	return &InstanceTestSmartContract{
		testSmartContract: testSmartContract{
			address: address,
		},
	}
}

// WithBalance provides the balance for the InstanceTestSmartContract
func (mockSC *InstanceTestSmartContract) WithBalance(balance int64) *InstanceTestSmartContract {
	mockSC.balance = balance
	return mockSC
}

// WithConfig provides the config object for the InstanceTestSmartContract
func (mockSC *InstanceTestSmartContract) WithConfig(config interface{}) *InstanceTestSmartContract {
	mockSC.config = config
	return mockSC
}

// WithCode provides the code for the InstanceTestSmartContract
func (mockSC *InstanceTestSmartContract) WithCode(code []byte) *InstanceTestSmartContract {
	mockSC.code = code
	return mockSC
}

// InstancesTestTemplate holds the data to build a contract call test
type InstancesTestTemplate struct {
	testTemplateConfig
	contracts     []*InstanceTestSmartContract
	gasSchedule   config.GasScheduleMap
	setup         func(arwen.VMHost, *contextmock.BlockchainHookStub)
	assertResults func(arwen.VMHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)
}

// BuildInstanceCallTest starts the building process for a contract call test
func BuildInstanceCallTest(tb testing.TB) *InstancesTestTemplate {
	return &InstancesTestTemplate{
		testTemplateConfig: testTemplateConfig{
			tb:       tb,
			useMocks: false,
		},
		setup: func(arwen.VMHost, *contextmock.BlockchainHookStub) {},
	}
}

// WithContracts provides the contracts to be used by the contract call test
func (callerTest *InstancesTestTemplate) WithContracts(usedContracts ...*InstanceTestSmartContract) *InstancesTestTemplate {
	callerTest.contracts = usedContracts
	return callerTest
}

// WithInput provides the ContractCallInput to be used by the contract call test
func (callerTest *InstancesTestTemplate) WithInput(input *vmcommon.ContractCallInput) *InstancesTestTemplate {
	callerTest.input = input
	return callerTest
}

// WithSetup provides the setup function to be used by the contract call test
func (callerTest *InstancesTestTemplate) WithSetup(setup func(arwen.VMHost, *contextmock.BlockchainHookStub)) *InstancesTestTemplate {
	callerTest.setup = setup
	return callerTest
}

// WithGasSchedule provides gas schedule for the test
func (callerTest *InstancesTestTemplate) WithGasSchedule(gasSchedule config.GasScheduleMap) *InstancesTestTemplate {
	callerTest.gasSchedule = gasSchedule
	return callerTest
}

// AndAssertResults starts the test and asserts the results
func (callerTest *InstancesTestTemplate) AndAssertResults(assertResults func(arwen.VMHost, *contextmock.BlockchainHookStub, *VMOutputVerifier)) {
	callerTest.assertResults = assertResults
	runTestWithInstances(callerTest)
}

func runTestWithInstances(callerTest *InstancesTestTemplate) {
	host, blockchainHookStub := defaultTestArwenForContracts(callerTest.tb, callerTest.contracts, callerTest.gasSchedule)
	defer func() {
		host.Reset()
	}()

	callerTest.setup(host, blockchainHookStub)

	vmOutput, err := host.RunSmartContractCall(callerTest.input)

	allErrors := host.Runtime().GetAllErrors()

	verify := NewVMOutputVerifierWithAllErrors(callerTest.tb, vmOutput, err, allErrors)
	callerTest.assertResults(host, blockchainHookStub, verify)
}
