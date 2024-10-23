package segwit

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
)

// Helper function to generate a random mnemonic for testing
func createTestMnemonic(t *testing.T) string {
	entropy, err := bip39.NewEntropy(128)
	assert.NoError(t, err, "Failed to generate entropy")
	mnemonic, err := bip39.NewMnemonic(entropy)
	assert.NoError(t, err, "Failed to generate mnemonic")
	return mnemonic
}

// Helper function to generate a sample wallet for testing
func createTestWallet(t *testing.T, network Network, path string) *Wallet {
	mnemonic := createTestMnemonic(t)
	config := &Config{
		Mnemonic: mnemonic,
		Path:     path,
		Network:  network,
	}
	wallet, err := New(config)
	assert.NoError(t, err, "Failed to create wallet")
	assert.NotNil(t, wallet, "Wallet should not be nil")
	return wallet
}

func Test_SelectDerivationPath(t *testing.T) {
	tests := []struct {
		name        string
		network     Network
		path        string
		expected    string
		expectError bool
	}{
		{"Mainnet Default Path", NetworkMainnet, "", `m/84'/0'/0'/0`, false},
		{"Testnet Default Path", NetworkTestnet, "", `m/84'/1'/0'/0`, false},
		{"Invalid Network", Network("invalid"), "", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := selectDerivationPath(test.network, test.path)
			if test.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, test.expected, result, "Derivation path mismatch")
			}
		})
	}
}

func Test_SelectNetworkParams(t *testing.T) {
	tests := []struct {
		name        string
		network     Network
		expectError bool
	}{
		{"Mainnet Params", NetworkMainnet, false},
		{"Testnet Params", NetworkTestnet, false},
		{"Invalid Network", Network("invalid"), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := selectNetworkParams(test.network)
			if test.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
			}
		})
	}
}

func Test_GenerateMasterKey(t *testing.T) {
	seed := []byte("test seed")
	params := &chaincfg.MainNetParams
	_, err := generateMasterKey(seed, params)
	assert.Error(t, err, "Expected an error due to invalid seed length")
}

func Test_Address(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0/0`)
	address := wallet.Address()
	assert.NotNil(t, address, "Address should not be nil")
}

func Test_AddressHex(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0/0`)
	addressHex := wallet.AddressHex()
	assert.NotEmpty(t, addressHex, "Address hex should not be empty")
	assert.Equal(t, 42, len(addressHex), "Bitcoin address should be 42 characters long")
}

func Test_Path(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0/0`)
	path := wallet.Path()
	assert.Equal(t, `m/84'/0'/0'/0/0`, path, "Derivation path should match")
}

func Test_Mnemonic(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0`)
	mnemonic := wallet.Mnemonic()
	assert.NotEmpty(t, mnemonic, "Mnemonic should not be empty")

	t.Run("empty mnemomic", func(t *testing.T) {
		config := &Config{
			Mnemonic: "",
			Path:     `m/84'/0'/0'/0`,
			Network:  NetworkMainnet,
		}
		wallet, err := New(config)
		assert.EqualError(t, err, ErrInvalidMnemonic)
		assert.Nil(t, wallet)
	})

	t.Run("invalid mnemomic", func(t *testing.T) {
		config := &Config{
			Mnemonic: "invalid invalid invalid invalid invalid",
			Path:     `m/84'/0'/0'/0`,
			Network:  NetworkMainnet,
		}
		wallet, err := New(config)
		assert.EqualError(t, err, ErrInvalidMnemonic)
		assert.Nil(t, wallet)
	})
}

func Test_PrivateKey(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0/0`)

	privateKeyWIF, err := wallet.PrivateKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, privateKeyWIF)

	// Vérification basique que la clé WIF commence par un préfixe valide (comme 5, L, ou K pour les clés WIF de Bitcoin).
	assert.True(t, privateKeyWIF[0] == '5' || privateKeyWIF[0] == 'L' || privateKeyWIF[0] == 'K')
}

func Test_ValidateAddress(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0/0`)

	// Utilisation d'une adresse valide pour le réseau principal (Mainnet).
	validAddress := "1QHTz6wMURLy8DT6aeGAVbF2UvtuWZKozr"
	isValid, err := wallet.ValidateAddress(validAddress)
	assert.NoError(t, err)
	assert.True(t, isValid)

	// Utilisation d'une adresse invalide (mauvais format ou non supportée pour ce réseau).
	invalidAddress := "InvalidBitcoinAddress"
	isValid, err = wallet.ValidateAddress(invalidAddress)
	assert.Error(t, err)
	assert.False(t, isValid)
}

func Test_ExtendedPublicKey(t *testing.T) {
	wallet := createTestWallet(t, NetworkMainnet, `m/84'/0'/0'/0/0`)

	extendedPublicKey, err := wallet.ExtendedPublicKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, extendedPublicKey)

	// Vérification basique que la clé publique étendue commence par "xpub" (Mainnet) ou "tpub" (Testnet).
	assert.True(t, extendedPublicKey[:4] == "xpub" || extendedPublicKey[:4] == "tpub")
}

func Test_New_with_mainnet(t *testing.T) {
	mnemonic := "romance trash engine during cliff verify tunnel memory vault chief fluid fox"
	root, err := New(&Config{
		Mnemonic: mnemonic,
		Network:  NetworkMainnet,
	})
	assert.NoError(t, err, "Failed to create wallet on mainnet")
	assert.Equal(t, root.Path(), `m/84'/0'/0'/0`, "Root path mismatch")
	assert.Equal(t, root.AddressHex(), "bc1qam3nmdeg6n9nms98hjyvpq822ktgpr0nera5wq", "Address hex mismatch")

	wallet, err := root.Derive(0)
	assert.NoError(t, err, "Failed to derive wallet")
	assert.Equal(t, wallet.Path(), `m/84'/0'/0'/0/0`, "Derived path mismatch")
	assert.Equal(t, wallet.AddressHex(), "bc1qggu5huqlzd7sjterrewumyak5h58g96yqvvylk", "Derived address hex mismatch")

	wallet, err = root.Derive(1)
	assert.NoError(t, err, "Failed to derive wallet")
	assert.Equal(t, wallet.Path(), `m/84'/0'/0'/0/1`, "Derived path mismatch")
	assert.Equal(t, wallet.AddressHex(), "bc1q7nj4055s85wzau55zju8fnuecxywpgfry2wae2", "Derived address hex mismatch")

	wallet, err = root.Derive(2)
	assert.NoError(t, err, "Failed to derive wallet")
	assert.Equal(t, wallet.Path(), `m/84'/0'/0'/0/2`, "Derived path mismatch")
	assert.Equal(t, wallet.AddressHex(), "bc1qvmx430ek35u29h8e3p8mk77fd5lhqz3dn52gd5", "Derived address hex mismatch")
}

func Test_New_with_testnet(t *testing.T) {
	mnemonic := "romance trash engine during cliff verify tunnel memory vault chief fluid fox"
	root, err := New(&Config{
		Mnemonic: mnemonic,
		Network:  NetworkTestnet,
	})
	assert.NoError(t, err, "Failed to create wallet on testnet")
	assert.Equal(t, root.Path(), `m/84'/1'/0'/0`, "Root path mismatch")
	assert.Equal(t, root.AddressHex(), "tb1qy68ml50k70lj4cuhamkhzhzxue9rmh84t94kcf", "Address hex mismatch")

	wallet, err := root.Derive(0)
	assert.NoError(t, err, "Failed to derive wallet")
	assert.Equal(t, wallet.Path(), `m/84'/1'/0'/0/0`, "Derived path mismatch")
	assert.Equal(t, wallet.AddressHex(), "tb1qjjt24zm4su62zydw58f78nf70g63x3q6j58vz2", "Derived address hex mismatch")
}
