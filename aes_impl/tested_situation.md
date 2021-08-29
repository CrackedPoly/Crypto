|      | 128bit  | 192bit | 256bit |
| ---- | ------- | ------ | ------ |
| ECB  |       |        |        |
| CBC  |   ✅   |   ✅   |   ✅   |
| CFB  |       |       |       |
| OFB  |       |       |       |

Capital or lower hex letters doesn't matter.

## key length
Not being 128, 192, or 256 bits leads to an error. ✅

## initial vector length
Shorter than 128 bits, error. Add a `0` to the last second position when the
length is odd. Example:`5072656E7469636548616C6C496E632` considered as `5072656E7469636548616C6C496E6302` ✅

Longer than 128 bit, only first 128 bits count. ✅

## plaintext length
Since padding function is used, the length does not matter. Also, add a `0` to the last second position when the 
length is odd. ✅