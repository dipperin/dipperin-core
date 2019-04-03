# Distribute your token

Dipperin supports formalized ERC20 token smart contract. ERC20 token can be deployed through three tools:

- Command line
- Wallet application
- Dipperin JavaScript API(dipperin.js)

## Deploy token through command line

```shell
# Start command line
$ dipperincli
# Deploy ERC20 token
rpc -m AnnounceERC20 -p [owner_address], [token_name], [token_symbol], [token_total_supply], [decimal],[transactionFee]
# Example
rpc -m AnnounceERC20 -p 0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,chain,stack,5,3,0.00001
```

[See more details for Command Line Tool](../design/commands#ERC20)

## Deploy token through wallet application

Download and install wallet.

<!-- TODO: Github 钱包仓库的 Release 页面 -->

After you created your account, jump to the contract page.

Click create contract and turn to a create contract page.

Fill in the informations and click create. Done.

## Deploy token through JavaScript API

Import dipperin.js in your JavaScript file.

```javascript
dipperin
import dipperin， { Contract， Accounts } from '@dipperin/dipperin.js'
dipperin
const dipperin = new Dipperin("$YOUR_RPC_PROVIDER")
// Deploy token contract
const contract = Contract.createContract(
  {
    owner: $YOUR_CONTRACT_OWNER,
    tokenDecimals: $YOUR_TOKEN_DECIMALS,
    tokenName: $YOUR_TOKEN_NAME,
    tokenSymbol: $YOUR_TOKEN_SYMBOL,
    tokenTotalSupply: $YOUR_TOKEN_TOTAL_SUPPLY
  },
  $YOUR_TOKEN_TYPE,
  $YOUR_TOKEN_ADDRESS
)
// Create a transaction
const signedTransaction = Accounts.signTransaction(
  {
    extraData: contract.contractData,
    fee: $TRANSACTION_FEE,
    nonce: $YOUR_ACCOUNT_NONCE,
    to: $YOUR_TOKEN_ADDRESS,
    value: '0',
  },
  $YOUR_PRIVATE_KEY
)
// Send transaction
dipperin.dr.sendSignedTransaction(signedTransaction.raw)
  .then(transactionHash => {
    // Do something
  })
// Done.
```
