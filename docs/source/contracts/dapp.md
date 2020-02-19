# Dipperin Dapp Development

## How to develop a Dapp with Dipperin Wallet Extension

Dipperin Wallet Extension supply a set of interfaces for Dapp developer, which makes developing a Dapp more easily.

### How Dipperin Wallet Extension works?

If you have already installed Dipperin Wallet Extension in your Chrome，it will inject all Dipperin supplied interfaces into all web pages in your browser. By this way, Dapp can interact with Dipperin network. Developers can get these interface by following ways.

```ts
window.dipperinEx
```

### Interfaces

DipperinEx supplied 5 interfaces，they have functions as follow：

```js
window.dipperinEx.isApproved; // Get the authorization state of Dapp
window.dipperinEx.approve; // Authorize the Dapp
window.dipperinEx.send; // Send transactions
window.dipperinEx.getActiveAccount; // Get accounts of users.
window.dipperinEx.on; // Listen for the message from the wallet extension.
```

### dipperinEx.isApproved

dipperinEx.isApproved supply the Dapp authorization state.

```js
const dapp_name = "Your Dapp's name";
/**
 * @param {string} dappName
 * @returns {Promise<{isApproved: boolean}>}
 */
window.dipperinEx
  .isApproved(dappName)
  .then(res => console.log(res)) // { isApproved: true } 
  .catch(e => console.log(e));
```

If the value isApproved is true, it represent that Dapp is authorized。

### dipperinEx.approve

``dipperinEx.approve`` is used for Dipperin Wallet Extension to authorize Dapp. Function can be called as follow.

```ts
/**
 * @interface ApproveRes
 */
interface ApproveRes {
  popupExist: boolean;
  isHaveWallet: boolean;
  isUnlock: boolean;
}
/**
 * @param {string} dappName
 * @returns {Promise<ApproveRes>}
 * @throws {ApproveRes}
 */
window.dipperinEx
  .approve(dappName)
  .then(res => console.log(res)) // {popupExist: false, isHaveWallet: true, isUnlock: true}
  .catch(e => console.log(e)); // {popupExist: false, isHaveWallet: true, isUnlock: false}
```

After call the function, there will shown dialog to request for user's authorization.

The return value of isHaveWallet represents whether user have accounts in this extension. IsUnlocked means the wallet is unlocked. The result of user authorization can get by call ``dipperinEx.isApproved ``.

### dipperinEx.send

Dipperin Wallet Extension supply ``dipperinEx.send`` for user to send transactions.

```ts
type Send = (
  name: string,
  to: string,
  value: string,
  extraData: string
) => Promise<SendResFailed | string>;
```

There are 4 input parameters, ``name`` for the name of the Dapp, ``to`` for the receiving address. ``value`` for the money to send, and ``extraData`` for extra data.

```ts
const address = "0x00003A9A328170b650E89F2C28F2E61364d2aEdC292e";
const amount = "1";
const extraData = "The message your dapp need";
/**
 * @interface SendResFailed
 */
interface SendResFailed {
  isApproved: boolean;
  isHaveWallet: boolean;
  isUnlock: boolean;
  info: string;
}

/**
 * @param {string} address
 * @param {string} amount
 * @param {string} extraData
 * @returns {Promise<string>}
 * @throws {ApproveRes}
 */
windox.dipperinEx
  .send(APP_NAME, DEAULT_ADDRESS, amount, extraData)
  .then(res => console.log(res)) // 0x8d303cb0b24fd332614a02c477605255e6a29afc3d477086603583f8aea5ddff
  .catch(e => console.log(e));   // {popupExist: false, isApproved: true, isHaveWallet: true, isUnlock: true, info: "send tx failed"}
```
If success it will return a transaction hash, otherwise return an error message.

### dipperinEx.getActiveAccount

``dipperinEx.getActiveAccount`` can get authorized account address.

```js
/**
 * @param {string} dappName
 * @returns {Promise<string>}
 */
window.dipperinEx.getActiveAddress(dappName);   // 0x00008522edBC22d9db52fa3AF05C2093dfFbFFF9DdBD
```

If success it will return an address, otherwise return an empty string.

### dipperinEx.on

"dipperinEx.on" is used by Dapp to get Dipperin wallet extension messages. 

```js
window.dipperinEx.on("changeActiveAccount", () => {
  console.log("Have changed active account!");
});
```


