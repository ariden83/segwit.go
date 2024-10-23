# SegWit Wallet - Go Library

This Go library provides functionality for creating and managing Hierarchical Deterministic (HD) wallets based on the SegWit protocol. It supports key derivation, mnemonic phrase validation (BIP39), and address generation for both Mainnet and Testnet Bitcoin networks.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Creating a Wallet](#creating-a-wallet)
- [Deriving Keys](#deriving-keys)
- [Getting Public Key and Address](#getting-public-key-and-address)
- [Validating a Bitcoin Address](#validating-a-bitcoin-address)
- [Fetching Private Key](#fetching-private-key)
- [Configuration](#configuration)
- [Errors](#errors)
- [Contributing](#contributing)
- [License](#license)

## Features

- Supports Bitcoin Mainnet and Testnet.
- BIP39 mnemonic phrase validation.
- HD wallet creation and management.
- SegWit address generation (P2WPKH).
- Public and private key retrieval.
- Bitcoin address validation.
- Extended public key (xpub) export.

## Installation

To use this package, you will need to have Go installed on your system. You can install the package using go get:

```bash
go get github.com/ariden83/p2pkh.go
```

Then, import it in your Go code:

```go
import "github.com/ariden83/p2pkh.go"
```

## Usage

### Creating a Wallet

You can create a wallet by providing a mnemonic phrase, a derivation path, and selecting the network (Mainnet or Testnet):

```go
config := &segwit.Config{
    Mnemonic: "your mnemonic phrase here",
    Path:     "m/84'/0'/0'/0",
    Network:  segwit.NetworkMainnet,
}

wallet, err := segwit.New(config)
if err != nil {
    log.Fatalf("Failed to create wallet: %v", err)
}
```

### Deriving Keys

You can derive new keys from the wallet using a specific index:

```go
derivedWallet, err := wallet.Derive(0)
if err != nil {
    log.Fatalf("Failed to derive wallet: %v", err)
}
fmt.Println("Derived Address:", derivedWallet.AddressHex())
```

### Getting Public Key and Address

Once the wallet is created, you can retrieve the public key and address associated with it:

```go
publicKey := wallet.PublicKey()
fmt.Println("Public Key:", publicKey)

address := wallet.AddressHex()
fmt.Println("Address:", address)
```

### Validating a Bitcoin Address

To validate if a Bitcoin address is valid for the current network (Mainnet or Testnet), you can use the `ValidateAddress` function:

```go
isValid, err := wallet.ValidateAddress("your Bitcoin address here")
if err != nil {
    log.Fatalf("Address validation failed: %v", err)
}

if isValid {
    fmt.Println("Address is valid.")
} else {
    fmt.Println("Address is not valid.")
}
```

### Fetching Private Key

To fetch the private key associated with the wallet in Wallet Import Format (WIF):

```go
privateKey, err := wallet.PrivateKey()
if err != nil {
    log.Fatalf("Failed to fetch private key: %v", err)
}
fmt.Println("Private Key (WIF):", privateKey)
```

## Configuration

### Config Struct

The `Config` struct is used to create a new wallet. It requires the following fields:

- **Mnemonic**: A valid BIP39 mnemonic phrase.
- **Path**: The derivation path (e.g., m/84'/0'/0'/0/0 for Bitcoin Mainnet).
- **Network**: Either NetworkMainnet or NetworkTestnet.

### Example:

```go
config := &segwit.Config{
    Mnemonic: "romance trash engine during cliff verify tunnel memory vault chief fluid fox",
    Path:     `m/84'/1'/0'/0/0`,  // Testnet derivation path
    Network:  p2pkh.NetworkTestnet,
}
```

## Errors
The library provides descriptive error messages for common issues:

- `ErrInvalidMnemonic`: Thrown when the mnemonic phrase is invalid or missing.
- `ErrUnsupportedNet`: Thrown when an unsupported network type is selected.
- `ErrInvalidPath`: Thrown when the derivation path cannot be parsed.
- `ErrKeyDerivation`: Thrown during key derivation failure.
- `ErrIndexNegative`: Thrown when a negative index is provided for key derivation.
- `ErrUnsupportedIndex`: Thrown when an unsupported index type is used.

### Example Error Handling

```go
wallet, err := segwit.New(config)
if err != nil {
    if errors.Is(err, segwit.ErrInvalidMnemonic) {
    log.Fatalf("Invalid mnemonic provided")
    }
    log.Fatalf("Failed to create wallet: %v", err)
}
```
## Wallet Methods

The `Wallet` struct provides the following methods:

- `PublicKey()`: Returns the wallet's ECDSA public key.
- `PrivateKey()`: Returns the wallet's private key in Wallet Import Format (WIF).
- `Address()`: Returns the wallet's P2PKH Bitcoin address (native btcutil format).
- `AddressHex()`: Returns the wallet's Bitcoin address in a hexadecimal string format.
- `ValidateAddress(address string)`: Validates if the provided address belongs to the current network.
- `ExtendedPublicKey()`: Returns the extended public key (xpub).
- `Derive(index interface{})`: Derives a new wallet based on the provided index from the current wallet.

### Example: Retrieving the Private Key

```go
privateKey, err := wallet.PrivateKey()
if err != nil {
    fmt.Println("Error retrieving private key:", err)
} else {
    fmt.Println("Private Key (WIF):", privateKey)
}
```

### Example: Validating an Address

```go
isValid, err := wallet.ValidateAddress("1QHTz6wMURLy8DT6aeGAVbF2UvtuWZKozr")
if err != nil {
    fmt.Println("Error validating address:", err)
} else if isValid {
    fmt.Println("Address is valid")
} else {
    fmt.Println("Address is invalid")
}
```

## Testing

The package includes a set of unit tests that can be run using the go test command. The tests cover the core functionality of the wallet, including key and address generation, derivation paths, and validation.

To run the tests, simply run:

```bash
go test ./...
```

Example Test

```go
func Test_InvalidMnemonic(t *testing.T) {
    config := &segwit.Config{
        Mnemonic: "invalid mnemonic phrase",
        Path:     `m/44'/0'/0'/0`,
        Network:  p2pkh.NetworkMainnet,
    }

    wallet, err := segwit.New(config)
    assert.Nil(t, wallet)
    assert.EqualError(t, err, segwit.ErrInvalidMnemonic)
}
```

## Contributing

If you'd like to contribute to this project, feel free to submit a pull request or open an issue on GitHub.

## License

This project is licensed under the MIT License.
