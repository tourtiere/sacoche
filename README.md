# Sacoche

Sacoche is a command line interface Ethereum wallet.

# Usage

1. Find a good passphrase.
2. Generate a keystore file from your passphrase. Save it on your computer and in a password manager.
3. Retrieve your private key from your keystore file.

```
sacoche <cmd> [arg]

Commands:
  sacoche generate [name]                Generate a keystore file from a passphrase.
  sacoche reveal   [keystore filename]   Show the private key.
```

# Example

```
$ sacoche generate hello

Input a passphrase: I love chocolate.
Your keystore has been generated.
File: hello-0xf1f11424AB293CBBa48602A479109c5C2457d5Cf.json

```

```
$ sacoche reveal hello-0xf1f11424AB293CBBa48602A479109c5C2457d5Cf.json

Input a passphrase: I love chocolate.
Success
ETH address : 0xf1f11424AB293CBBa48602A479109c5C2457d5Cf
---
Private key = E8888E2211F964B4D34C03F6288179C059C830D2E356320A9B5BB7988CDA9667

A = 5C 43 35 D5 A6 AF 59 D6 14 8F F5 6E 1F 28 80 9E 95 22 A0 AA 5A 89 2C 2A 59 A8 BD 76 5D 09 7A 54
B = 8C 45 58 4C 6B 4A 0A DE BE BC 0E 88 09 58 F9 21 C4 A5 90 28 88 CD 05 E0 41 B2 FA 22 2F D1 1C 13

Private key = A + B
```

Content of `hello-0xf1f11424AB293CBBa48602A479109c5C2457d5Cf.json`:

```
{
  "version": 1,
  "params": {
    "salt": "e2f2a81d6a88ab10d52ef5958d051471bfc612607f85aff390d30994b30303a9",
    "n": 262144,
    "r": 8,
    "p": 1,
    "keyLen": 32
  },
  "ETHAddress": "0xf1f11424AB293CBBa48602A479109c5C2457d5Cf"
}
```
