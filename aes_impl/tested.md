# AES implementation in Golang.

## Tested
|      | 128bit  | 192bit | 256bit |
| ---- | ------- | ------ | ------ |
| ECB  |   ✅    |   ✅   |   ✅  |
| CBC  |   ✅    |   ✅   |   ✅  |
| CFB  |   ✅    |   ✅   |   ✅  |
| OFB  |   ✅    |   ✅   |   ✅  |
| CTR  |   ✅    |   ✅   |   ✅  |
| GCM  |   ✅    |   ✅   |   ✅  |

All results are the same as [CyberChef](https://github.com/gchq/CyberChef)

## How To Use
Open this project in GoLand and build `main.go`.

encrypt

`encrypt -m GCM -p aes_plain1.txt -k aes_key.txt -v aes_iv.txt -c aes_cipher.txt -a aes_auth.txt --tag aes_tag.txt`

decrypt

`decrypt -m GCM -p aes_plain1.txt -k aes_key.txt -v aes_iv.txt -c aes_cipher.txt -a aes_auth.txt --tag aes_tag.txt` 

help

`encrypt help` or `decrypt help`

## Arguments
Capital or lower hex letters doesn't matter.

### key
Not being 128, 192, or 256 bits leads to an error. 

### initial vector
Shorter than 128 bits, error. 
Add a `0` to the last second position when the
length is odd. Example:`5072656E7469636548616C6C496E632` considered as `5072656E7469636548616C6C496E6302` 

Longer than 128 bit, only first 128 bits count. 

Specially, iv length can be any in GCM mode.

### plaintext length
Since padding function is used, the length does not matter. Also, add a `0` to the last second position when the 
length is odd. 

### authentication message (GCM)
Can be any length.

### tag (GCM)
Be specified as 128bit.