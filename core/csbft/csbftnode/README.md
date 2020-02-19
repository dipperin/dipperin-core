# csbftnode

Complete other functions in PBFT other than status maintenance

1. The block is quickly propagated to other nodes after receiving a new block to be verified. First pass hash to confirm that there is no block in the other party's block pool and then pass the block to the other party.