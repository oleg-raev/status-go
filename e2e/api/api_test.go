package api_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"


	"github.com/status-im/status-go/e2e"
	"github.com/status-im/status-go/geth/api"
	"github.com/status-im/status-go/geth/log"
	"github.com/status-im/status-go/geth/params"
	. "github.com/status-im/status-go/testing"
	"github.com/stretchr/testify/suite"
)

const (
	testChatID = "testChat"
)

func TestAPI(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

type APITestSuite struct {
	suite.Suite
	api *api.StatusAPI
}

func (s *APITestSuite) SetupTest() {
	s.api = api.NewStatusAPI()
	s.NotNil(s.api)
}

func (s *APITestSuite) TestCHTUpdate() {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "cht-updates")
	s.NoError(err)
	defer os.RemoveAll(tmpDir)

	configJSON := `{
		"NetworkId": ` + strconv.Itoa(params.RopstenNetworkID) + `,
		"DataDir": "` + tmpDir + `",
		"LogLevel": "INFO",
		"RPCEnabled": true
	}`
	//nodeConfig, err := params.LoadNodeConfig(configJSON)
	_, err = params.LoadNodeConfig(configJSON)
	s.NoError(err)

	// start node
	//nodeConfig.DevMode = true
	//s.api.StartNode(nodeConfig)
	//s.api.StopNode()
	// TODO(tiabc): Test that CHT is really updated.
}

func (s *APITestSuite) TestRaceConditions() {
	cnt := 25
	progress := make(chan struct{}, cnt)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	nodeConfig1, err := e2e.MakeTestNodeConfig(params.RopstenNetworkID)
	s.NoError(err)

	nodeConfig2, err := e2e.MakeTestNodeConfig(params.RinkebyNetworkID)
	s.NoError(err)

	nodeConfigs := []*params.NodeConfig{nodeConfig1, nodeConfig2}

	var funcsToTest = []func(*params.NodeConfig){
		func(config *params.NodeConfig) {
			log.Info("StartNodeAsync()")
			_, err := s.api.StartNodeAsync(config)
			s.T().Logf("StartNodeAsync() for network: %d, error: %v", config.NetworkID, err)
			progress <- struct{}{}
		},
		func(config *params.NodeConfig) {
			log.Info("StopNodeAsync()")
			_, err := s.api.StopNodeAsync()
			s.T().Logf("StopNodeAsync(), error: %v", err)
			progress <- struct{}{}
		},
		func(config *params.NodeConfig) {
			log.Info("RestartNodeAsync()")
			_, err := s.api.RestartNodeAsync()
			s.T().Logf("RestartNodeAsync(), error: %v", err)
			progress <- struct{}{}
		},
		// TODO(adam): quarantined until it uses a different datadir
		// as otherwise it wipes out cached blockchain data.
		// func(config *params.NodeConfig) {
		// 	log.Info("ResetChainDataAsync()")
		// 	_, err := s.api.ResetChainDataAsync()
		// 	s.T().Logf("ResetChainDataAsync(), error: %v", err)
		// 	progress <- struct{}{}
		// },
	}

	// increase StartNode()/StopNode() population
	for i := 0; i < 5; i++ {
		funcsToTest = append(funcsToTest, funcsToTest[0], funcsToTest[1])
	}

	for i := 0; i < cnt; i++ {
		randConfig := nodeConfigs[rnd.Intn(len(nodeConfigs))]
		randFunc := funcsToTest[rnd.Intn(len(funcsToTest))]

		if rnd.Intn(100) > 75 { // introduce random delays
			time.Sleep(500 * time.Millisecond)
		}
		go randFunc(randConfig)
	}

	for range progress {
		cnt -= 1
		if cnt <= 0 {
			break
		}
	}

	time.Sleep(2 * time.Second) // so that we see some logs
	s.api.StopNode()            // just in case we have a node running
}

func (s *APITestSuite) TestJailStop() {
	const itersCount = 5
	var (
		wg        sync.WaitGroup
		require   = s.Require()
		getChatId = func(id int) string {
			return testChatID + strconv.Itoa(id)
		}
	)

	config, err := e2e.MakeTestNodeConfig(params.RopstenNetworkID)
	require.NoError(err)
	err = s.api.StartNode(config)
	require.NoError(err)
	defer s.api.StopNode()

	address1, _, _, err := s.api.AccountManager().CreateAccount(TestConfig.Account1.Password)
	require.NoError(err)

	address2, _, _, err := s.api.AccountManager().CreateAccount(TestConfig.Account2.Password)
	require.NoError(err)

	err = s.api.SelectAccount(address1, TestConfig.Account1.Password)
	require.NoError(err)

	for i := 0; i < itersCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := s.api.JailManager().NewCell(getChatId(id))
			require.NoError(err)
		}(i)
	}
	wg.Wait()

	err = s.api.SelectAccount(address2, TestConfig.Account2.Password)
	require.NoError(err)

	wg = sync.WaitGroup{}
	for i := 0; i < itersCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := s.api.JailManager().Cell(getChatId(id))
			require.Error(err)
		}(i)
	}
	wg.Wait()
}

// TestLogoutRemovesCells we want be sure that
// cells will be removed after the API call "Logout"
func (s *APITestSuite) TestLogoutRemovesCells() {
	var (
		err     error
		require = s.Require()
	)

	config, err := e2e.MakeTestNodeConfig(params.RopstenNetworkID)
	require.NoError(err)
	err = s.api.StartNode(config)
	require.NoError(err)
	defer s.api.StopNode()

	address1, _, _, err := s.api.AccountManager().CreateAccount(TestConfig.Account1.Password)
	require.NoError(err)

	err = s.api.SelectAccount(address1, TestConfig.Account1.Password)
	require.NoError(err)

	s.api.JailManager().Parse(testChatID, ``)

	_, err = s.api.JailManager().NewCell(testChatID)
	require.NoError(err)

	err = s.api.Logout()
	require.NoError(err)

	_, err = s.api.JailManager().Cell(testChatID)
	require.Error(err, "Expected that cells was removed")
}
