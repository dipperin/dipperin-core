# BlockValidation

1. Verify header

CsBFT does not modify data when verifying state roots and registering roots
SaveBlock confirms data and save data when verifying status root and registration root

|completeness rate|block type|consensus|method|
|-|-|-|-|
|complete|both|block configuration correct:version ID,chain ID, etc|ValidateBlockVersion|
|complete|both|verify continuum of block height|ValidateBlockNumber|
|complete|both|verify correctness of block hash|ValidateBlockHash|
|complete|both|verify block size smaller than default maximum|ValidateBlockSize|
|complete|normal block|verify block difficulty meets the need|ValidateBlockDifficulty|
|complete|normal block|verify block hash satisfies the difficulty|ValidHashForDifficulty|
|complete|both|verify correctness of block seed|ValidateSeed|
|complete|special block|verify the address of the miner of the block is boot node|ValidateBlockCoinBase|
|complete|both|verify vote root hash|validVerificationRoot|
|complete|both|verify transaction root hash|ValidateBlockTxs|
|complete|both|verify state root hash|validStateRoot|
|complete|both|verify registration root hash|validBlockVerifier|
|incomplete|both|verify interlink root hash|？？？|

2. Verify votes

When CsBFT verifies the vote, it only needs to verify the verification in the body. The object is the last block.
When SaveBlock verifies the vote, you need to verify the two sets of votes, except for verification, the other group is see commit, and the object is the current block.

|completeness rate|block type|consensus|method|
|-|-|-|-|
|complete|normal block|check at least 2/3 of verifiers vote|validVotesForBlock|
|complete|special block|check there is angel node who has voted|validVotesForBlock|
|complete|special block|check vote type is correct|HaltedVoteValid|
|complete|both|verify there is no repeated vote|sameVote|
|complete|both|check the voter is the current verifier|verificationSignerInVerifiers|
|complete|both|verify the correctness of signature of verifiers|ver.Valid()|
|complete|both|verify the correctness of verifiers|validVotesForBlock|

3. Verify Txs

|completeness rate|block type|consensus|method|
|-|-|-|-|
|complete|special block|verify transaction list of the special block is empty|ValidateBlockTxs|
|complete|normal block|verify the correctness of signature of transactions|ValidateBlockTxs|
|complete|normal block|verify transactions of all types meet the need of safety|txValidators|
|complete|normal block|verify the sender's balance is sufficient|ValidTxSender|
|complete|normal block|verify the size of each transaction is smaller than the default maximal value|ValidTxSize|