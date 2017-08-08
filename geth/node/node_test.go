package node_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/status-im/status-go/geth/node"
	"github.com/status-im/status-go/geth/params"
	. "github.com/status-im/status-go/geth/testing"
)

func TestCHTTestSuite(t *testing.T) {
	suite.Run(t, new(CHTTestSuite))
}

type CHTTestSuite struct {
	BaseTestSuite
}

func (s *CHTTestSuite) SetupTest() {
	s.NodeManager = node.NewNodeManager()
	s.Require().NotNil(s.NodeManager)
	s.Require().IsType(&node.NodeManager{}, s.NodeManager)
}

func (s *CHTTestSuite) TestCHTBoot() {
	require := s.Require()
	require.NotNil(s.NodeManager)

	config, err := MakeTestNodeConfig(params.RinkebyNetworkID)
	require.NoError(err)

	config.IPCEnabled = false
	config.WSEnabled = false
	config.HTTPHost = "" // to make sure that no HTTP interface is started

	nodeStarted, err := s.NodeManager.StartNode(config)
	require.NoError(err)
	require.NotNil(config)

	defer s.NodeManager.StopNode()

	<-nodeStarted
}

// func (s *CHTTestSuite) TestCHTBoot() {
// 	t := s.T()
//
// 	lightConfig := strconv.Quote(`[{
//       "networkID": 3,
//       "genesisHash": "0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d",
//       "prod": {
//           "number": 259,
//           "hash": "91825fffecb5678167273955deaddbf03c26ae04287cfda61403c0bad5ceab8d",
//           "bootnodes": [
//               "enode://7ab298cedc4185a894d21d8a4615262ec6bdce66c9b6783878258e0d5b31013d30c9038932432f70e5b2b6a5cd323bf820554fcb22fbc7b45367889522e9c449@51.15.63.93:30303",
//               "enode://f59e8701f18c79c5cbc7618dc7bb928d44dc2f5405c7d693dad97da2d8585975942ec6fd36d3fe608bfdc7270a34a4dd00f38cfe96b2baa24f7cd0ac28d382a1@51.15.79.88:30303",
//               "enode://e2a3587b7b41acfc49eddea9229281905d252efba0baf565cf6276df17faf04801b7879eead757da8b5be13b05f25e775ab6d857ff264bc53a89c027a657dd10@51.15.45.114:30303",
//               "enode://fe991752c4ceab8b90608fbf16d89a5f7d6d1825647d4981569ebcece1b243b2000420a5db721e214231c7a6da3543fa821185c706cbd9b9be651494ec97f56a@51.15.67.119:30303",
//               "enode://482484b9198530ee2e00db89791823244ca41dcd372242e2e1297dd06f6d8dd357603960c5ad9cc8dc15fcdf0e4edd06b7ad7db590e67a0b54f798c26581ebd7@51.15.75.138:30303",
//               "enode://9e99e183b5c71d51deb16e6b42ac9c26c75cfc95fff9dfae828b871b348354cbecf196dff4dd43567b26c8241b2b979cb4ea9f8dae2d9aacf86649dafe19a39a@51.15.79.176:30303",
//               "enode://0f7c65277f916ff4379fe520b875082a56e587eb3ce1c1567d9ff94206bdb05ba167c52272f20f634cd1ebdec5d9dfeb393018bfde1595d8e64a717c8b46692f@51.15.54.150:30303",
//               "enode://e006f0b2dc98e757468b67173295519e9b6d5ff4842772acb18fd055c620727ab23766c95b8ee1008dea9e8ef61e83b1515ddb3fb56dbfb9dbf1f463552a7c9f@212.47.237.127:30303",
//               "enode://d40871fc3e11b2649700978e06acd68a24af54e603d4333faecb70926ca7df93baa0b7bf4e927fcad9a7c1c07f9b325b22f6d1730e728314d0e4e6523e5cebc2@51.15.132.235:30303",
//               "enode://ea37c9724762be7f668e15d3dc955562529ab4f01bd7951f0b3c1960b75ecba45e8c3bb3c8ebe6a7504d9a40dd99a562b13629cc8e5e12153451765f9a12a61d@163.172.189.205:30303",
//               "enode://88c2b24429a6f7683fbfd06874ae3f1e7c8b4a5ffb846e77c705ba02e2543789d66fc032b6606a8d8888eb6239a2abe5897ce83f78dcdcfcb027d6ea69aa6fe9@163.172.157.61:30303",
//               "enode://ce6854c2c77a8800fcc12600206c344b8053bb90ee3ba280e6c4f18f3141cdc5ee80bcc3bdb24cbc0e96dffd4b38d7b57546ed528c00af6cd604ab65c4d528f6@163.172.153.124:30303",
//               "enode://00ae60771d9815daba35766d463a82a7b360b3a80e35ab2e0daa25bdc6ca6213ff4c8348025e7e1a908a8f58411a364fe02a0fb3c2aa32008304f063d8aaf1a2@163.172.132.85:30303",
//               "enode://86ebc843aa51669e08e27400e435f957918e39dc540b021a2f3291ab776c88bbda3d97631639219b6e77e375ab7944222c47713bdeb3251b25779ce743a39d70@212.47.254.155:30303"
//           ]
//       },
//       "dev": {
//          "number": 259,
//          "hash": "91825fffecb5678167273955deaddbf03c26ae04287cfda61403c0bad5ceab8d",
//          "bootnodes": [
//              "enode://7ab298cedc4185a894d21d8a4615262ec6bdce66c9b6783878258e0d5b31013d30c9038932432f70e5b2b6a5cd323bf820554fcb22fbc7b45367889522e9c449@51.15.63.93:30303",
//              "enode://f59e8701f18c79c5cbc7618dc7bb928d44dc2f5405c7d693dad97da2d8585975942ec6fd36d3fe608bfdc7270a34a4dd00f38cfe96b2baa24f7cd0ac28d382a1@51.15.79.88:30303",
//              "enode://e2a3587b7b41acfc49eddea9229281905d252efba0baf565cf6276df17faf04801b7879eead757da8b5be13b05f25e775ab6d857ff264bc53a89c027a657dd10@51.15.45.114:30303",
//              "enode://fe991752c4ceab8b90608fbf16d89a5f7d6d1825647d4981569ebcece1b243b2000420a5db721e214231c7a6da3543fa821185c706cbd9b9be651494ec97f56a@51.15.67.119:30303",
//              "enode://482484b9198530ee2e00db89791823244ca41dcd372242e2e1297dd06f6d8dd357603960c5ad9cc8dc15fcdf0e4edd06b7ad7db590e67a0b54f798c26581ebd7@51.15.75.138:30303",
//              "enode://9e99e183b5c71d51deb16e6b42ac9c26c75cfc95fff9dfae828b871b348354cbecf196dff4dd43567b26c8241b2b979cb4ea9f8dae2d9aacf86649dafe19a39a@51.15.79.176:30303",
//              "enode://0f7c65277f916ff4379fe520b875082a56e587eb3ce1c1567d9ff94206bdb05ba167c52272f20f634cd1ebdec5d9dfeb393018bfde1595d8e64a717c8b46692f@51.15.54.150:30303",
//              "enode://e006f0b2dc98e757468b67173295519e9b6d5ff4842772acb18fd055c620727ab23766c95b8ee1008dea9e8ef61e83b1515ddb3fb56dbfb9dbf1f463552a7c9f@212.47.237.127:30303",
//              "enode://d40871fc3e11b2649700978e06acd68a24af54e603d4333faecb70926ca7df93baa0b7bf4e927fcad9a7c1c07f9b325b22f6d1730e728314d0e4e6523e5cebc2@51.15.132.235:30303",
//              "enode://ea37c9724762be7f668e15d3dc955562529ab4f01bd7951f0b3c1960b75ecba45e8c3bb3c8ebe6a7504d9a40dd99a562b13629cc8e5e12153451765f9a12a61d@163.172.189.205:30303",
//              "enode://88c2b24429a6f7683fbfd06874ae3f1e7c8b4a5ffb846e77c705ba02e2543789d66fc032b6606a8d8888eb6239a2abe5897ce83f78dcdcfcb027d6ea69aa6fe9@163.172.157.61:30303",
//              "enode://ce6854c2c77a8800fcc12600206c344b8053bb90ee3ba280e6c4f18f3141cdc5ee80bcc3bdb24cbc0e96dffd4b38d7b57546ed528c00af6cd604ab65c4d528f6@163.172.153.124:30303",
//              "enode://00ae60771d9815daba35766d463a82a7b360b3a80e35ab2e0daa25bdc6ca6213ff4c8348025e7e1a908a8f58411a364fe02a0fb3c2aa32008304f063d8aaf1a2@163.172.132.85:30303",
//              "enode://86ebc843aa51669e08e27400e435f957918e39dc540b021a2f3291ab776c88bbda3d97631639219b6e77e375ab7944222c47713bdeb3251b25779ce743a39d70@212.47.254.155:30303"
//          ]
//       }
//   }]`)
//
// 	baseConfigJSON := (`{
//     "NetworkId": 3,
//     "DataDir": "./data_tmp",
//     "Name": "TestStatusNode",
//     "IPCEnabled": false,
//     "WSEnabled": false,
//     "LightEthConfig": {
//       "Enabled": true,
//       "Genesis": ` + lightConfig + `,
//       "DatabaseCache": 16
//     }
//   }`)
//
// 	config, err := params.LoadNodeConfig(baseConfigJSON)
// 	if err != nil {
// 		t.Fatalf("Config failed to be set: %+q", err)
// 	}
//
// 	config.WhisperConfig.Enabled = false
//
// 	cnode, err := MakeNode(config)
// 	if err != nil {
// 		t.Fatalf("Failed to create node with NodeConfig: %+q", err)
// 	}
//
// 	genesis := new(core.Genesis)
// 	if err := json.Unmarshal([]byte(config.LightEthConfig.Genesis), genesis); err != nil {
// 		t.Fatalf("Config failed to load json LES set: %+q", err)
// 	}
//
// 	ethConf := eth.DefaultConfig
// 	ethConf.Genesis = genesis
// 	ethConf.SyncMode = downloader.LightSync
// 	ethConf.NetworkId = config.NetworkID
// 	ethConf.DatabaseCache = config.LightEthConfig.DatabaseCache
// 	ethConf.MaxPeers = config.MaxPeers
//
// 	if err3 := cnode.Register(func(ctx *node.ServiceContext) (node.Service, error) {
// 		fmt.Printf("Receiving write request: \n")
//
// 		ethClient, err := les.New(ctx, &ethConf)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		ethClient.WriteTrustedCht(light.TrustedCht{
// 			Number: 3,
// 			Root:   gethcommon.HexToHash("0x435344333tt34t34t3t3"),
// 		})
//
// 		return ethClient, nil
// 	}); err3 != nil {
// 		t.Fatalf("Config failed to register with LES node: %+q", err3)
// 	}
//
// 	fmt.Printf("Waiting on cht: \n")
//
// 	if err4 := cnode.Start(); err4 != nil {
// 		t.Fatalf("Config failed to start with LES node: %+q", err4)
// 	}
// }
