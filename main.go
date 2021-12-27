package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/scrypt"
)

type ScryptParams struct {
	Salt   string `json:"salt"`
	N      int    `json:"n"`
	R      int    `json:"r"`
	P      int    `json:"p"`
	KeyLen int    `json:"keyLen"`
}

type KeyStore struct {
	Version      int          `json:"version"`
	ScryptParams ScryptParams `json:"params"`
	ETHAddress   string       `json:"ETHAddress"`
}

type Key struct {
	Private big.Int
	Store   KeyStore
}

func newKey(password string, params ScryptParams) Key {
	privateBytes, _ := scrypt.Key([]byte(password), []byte(params.Salt), params.N, params.R, params.P, params.KeyLen)
	private := new(big.Int).SetBytes(privateBytes)

	return Key{
		Private: *private,
		Store: KeyStore{
			Version:      1,
			ScryptParams: params,
			ETHAddress:   computeETHAddress(private),
		},
	}
}

func (key Key) saveKeyStore(filename string) {
	bytes, _ := json.Marshal(key.Store)
	_ = ioutil.WriteFile(filename, bytes, 0644)
}

func loadKeyStore(filename string) (*KeyStore, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var store KeyStore
	err = json.Unmarshal(bytes, &store)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func computeETHAddress(private *big.Int) string {
	keyPair, _ := crypto.ToECDSA(private.Bytes())
	publicKey := keyPair.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}

// returns a random number pair (a, b) such that n = a + b
func splitInt(n *big.Int) (*big.Int, *big.Int) {
	a, _ := rand.Int(rand.Reader, n)
	b := new(big.Int).Sub(n, a)
	return a, b
}

func sumInts(a *big.Int, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func humanReadableKey(key *big.Int, separate bool) string {
	hex := strings.ToUpper(key.Text(16))
	for len(hex) != 64 {
		hex = "0" + hex
	}

	if !separate {
		return hex
	}
	pairs := make([]string, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		pairs[i/2] = hex[i : i+2]
	}
	return strings.Join(pairs, " ")
}

func askForPassword() string {
	fmt.Print("Input a passphrase: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func main() {

	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Bad arguments")
		return
	}

	command := args[0]

	switch command {
	case "generate":
		{
			password := askForPassword()
			name := args[1]
			salt := make([]byte, 32)
			rand.Read(salt)

			params := ScryptParams{
				Salt:   hex.EncodeToString(salt),
				N:      262144,
				R:      8,
				P:      1,
				KeyLen: 32,
			}
			key := newKey(password, params)

			filename := name + "-" + key.Store.ETHAddress + ".json"
			key.saveKeyStore(filename)
			fmt.Println("Your keystore has been generated.")
			fmt.Println("File: ", filename)
		}
	case "reveal":
		{
			filename := args[1]
			loadedKeyStore, e := loadKeyStore(filename)
			if e != nil {
				fmt.Println("File load failed.")
			}
			password := askForPassword()
			key := newKey(password, loadedKeyStore.ScryptParams)
			if loadedKeyStore.ETHAddress != key.Store.ETHAddress {
				fmt.Println("Wrong password.")
				return
			}

			a, b := splitInt(&key.Private)

			fmt.Println("Success")
			fmt.Println("ETH address :", key.Store.ETHAddress)
			fmt.Println("------------")
			fmt.Println("Private key =", humanReadableKey(&key.Private, false))
			fmt.Println("")
			fmt.Println("A           =", humanReadableKey(a, true))
			fmt.Println("B           =", humanReadableKey(b, true))
			fmt.Println("")
			fmt.Println("Private key = A + B")
		}
	}

}
