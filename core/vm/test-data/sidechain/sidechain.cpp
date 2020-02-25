#include "sidechain.hpp"

PAYABLE void sidechain::init(char *addr, uint64_t id)
{
    isZeroAddress(addr);
    std::string originStr = dipc::origin().toString();
    uint64_t value = dipc::callValueUDIP();
    uint64_t blockNumber = ::number();

    *from = CharToAddress2Str(addr);
    *owner = originStr;
    *amount = value;
    *chainid = id;
    *height = blockNumber + 1000;
}

EXPORT void sidechain::transfer(char *proof, char *oracleAddr)
{
    // check caller
    DipcAssert(from.get() == caller().toString());

    // get header
    uint64_t num = dipc::getProofHeight(proof);
    DipcAssert(num != 0);
    DeployedContract contract(oracleAddr);
    std::string header = contract.callString("getHeader", num);

    // check SPV proof
    Address fromAddr = Address(&from.get()[0], 22);
    Address toAddr = Address(&owner.get()[0], 22);
    DipcAssert(dipc::validateSPVProof(proof, &header[0], fromAddr, toAddr, amount.get(), chainid.get(), height.get()) == 0);

    // check contract balance and withdraw
    Address contractAddr = address();
    u256 originBalance = dipc::balance(contractAddr);
    DipcAssert(originBalance > 0);
    DipcAssert(callTransfer(fromAddr, originBalance) == 0);
}

EXPORT void sidechain::withdraw()
{
    // check caller
    Address ownerAddr = Address(&owner.get()[0], 22);
    DipcAssert(owner.get() == caller().toString());

    // check height
    uint64_t blockNumber = ::number();
    DipcAssert(blockNumber > height.get());

    // check contract balance and withdraw
    Address contractAddr = address();
    u256 originBalance = dipc::balance(contractAddr);
    DipcAssert(originBalance > 0);
    DipcAssert(callTransfer(ownerAddr, originBalance) == 0);
}
