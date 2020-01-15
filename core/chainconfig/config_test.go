package chainconfig

import (
	"github.com/dipperin/dipperin-core/third_party/p2p/enode"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func Test_defaultChainConfig(t *testing.T) {
	// inity
	type expectStc struct {
		VerifyNumber int
		NetworkID    uint64
		ChainId      *big.Int
	}
	// test case
	situations := []struct {
		name   string
		given  func() *ChainConfig
		expect expectStc
	}{
		{
			"default environment without set env",
			func() *ChainConfig {
				os.Unsetenv(BootEnvTagName)
				return defaultChainConfig()
			},
			expectStc{
				22,
				1601,
				big.NewInt(1601),
			},
		},
		{
			"mercury environment",
			func() *ChainConfig {
				os.Setenv(BootEnvTagName, BootEnvMercury)
				return defaultChainConfig()
			},
			expectStc{
				22,
				99,
				big.NewInt(1),
			},
		},
		{
			"venus environment",
			func() *ChainConfig {
				os.Setenv(BootEnvTagName, BootEnvVenus)
				return defaultChainConfig()
			},
			expectStc{
				22,
				100,
				big.NewInt(2),
			},
		},
		{
			"mercury environment",
			func() *ChainConfig {
				os.Setenv(BootEnvTagName, BootEnvTest)
				return defaultChainConfig()
			},
			expectStc{
				22,
				1600,
				big.NewInt(1600),
			},
		},
		{
			"mercury environment",
			func() *ChainConfig {
				os.Setenv(BootEnvTagName, BootEnvLocal)
				return defaultChainConfig()
			},
			expectStc{
				4,
				1601,
				big.NewInt(1601),
			},
		},
	}
	// test
	for _, situation := range situations {
		chainConfig := situation.given()
		assert.Equal(t, situation.expect.VerifyNumber, chainConfig.VerifierNumber, situation.name)
		assert.Equal(t, situation.expect.NetworkID, chainConfig.NetworkID, situation.name)
		assert.Equal(t, 0, chainConfig.ChainId.Cmp(situation.expect.ChainId), situation.name)
	}
}

func TestGetChainConfig(t *testing.T) {
	testConfig := GetChainConfig()
	assert.Equal(t, 22, testConfig.VerifierNumber)
	assert.Equal(t, uint64(1601), testConfig.NetworkID)
	assert.Equal(t, 0, testConfig.ChainId.Cmp(big.NewInt(1601)))
}

func TestGetCurBootsEnv(t *testing.T) {
	os.Setenv(BootEnvTagName, "test GetCurBootsEnv")
	env := GetCurBootsEnv()
	assert.Equal(t, "test GetCurBootsEnv", env)
}

func TestInitBootNodes(t *testing.T) {
	// test case
	situations := []struct {
		name                   string
		given                  func() string
		expectVerifierNodesLen int
		expectKBucketNodesLen  int
	}{
		{
			"default env which boot locally",
			func() string {
				os.Unsetenv(BootEnvTagName)
				VerifierBootNodes = []*enode.Node{}
				KBucketNodes = []*enode.Node{}
				return ""
			},
			1,
			1,
		},
		{
			"test env",
			func() string {
				os.Setenv(BootEnvTagName, BootEnvTest)
				VerifierBootNodes = []*enode.Node{}
				KBucketNodes = []*enode.Node{}
				return ""
			},
			4,
			1,
		},
		{
			"mercury env",
			func() string {
				os.Setenv(BootEnvTagName, BootEnvMercury)
				VerifierBootNodes = []*enode.Node{}
				KBucketNodes = []*enode.Node{}
				return ""
			},
			4,
			1,
		},
		{
			"venus env",
			func() string {
				os.Setenv(BootEnvTagName, BootEnvVenus)
				VerifierBootNodes = []*enode.Node{}
				KBucketNodes = []*enode.Node{}
				return ""
			},
			4,
			1,
		},
	}
	// test
	for _, situation := range situations {
		dataDir := situation.given()
		InitBootNodes(dataDir)
		// check result
		assert.Equal(t, situation.expectVerifierNodesLen, len(VerifierBootNodes), situation.name)
		assert.Equal(t, situation.expectKBucketNodesLen, len(KBucketNodes), situation.name)
	}
}
