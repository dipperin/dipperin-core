#include "state_test.hpp"

DIPC_EVENT(state_int64, const char *, int64_t);
DIPC_EVENT(state_uint64, const char *, uint64_t);
DIPC_EVENT(state_string, const char *, const char *);
DIPC_EVENT(state_double, const char *, double);

void statetest::stateTests()
{
    print("\r\n dipc lib test stateTests start\r\n");
    uint64_t num = ::number();
    uint64_t gasLimit = ::gasLimit();
    int64_t gasPrice = ::gasPrice();
    int64_t timestamp = ::timestamp();
    int64_t nonce = ::getCallerNonce();

    std::string blockHash = dipc::blockHash(num - 1).toString();
    std::string coinbase = dipc::coinbase2().toString();
    Address2 addr = dipc::address2();
    std::string address = addr.toString();
    std::string balance = dipc::balance(addr).convert_to<std::string>().c_str();
    std::string origin = dipc::origin2().toString();
    std::string caller = dipc::caller2().toString();
    std::string callValue = dipc::callValue().convert_to<std::string>().c_str();

    //test sha3
    sha3Test();

    DIPC_EMIT_EVENT(state_uint64, &("blockNum"[0]), num);
    DIPC_EMIT_EVENT(state_uint64, &("GasLimit"[0]), gasLimit);

    DIPC_EMIT_EVENT(state_int64, &("GasPrice"[0]), gasPrice);
    DIPC_EMIT_EVENT(state_int64, &("TimeStamp"[0]), timestamp);
    DIPC_EMIT_EVENT(state_int64, &("Nonce"[0]), nonce);

    DIPC_EMIT_EVENT(state_string, &("blockHash"[0]), &(blockHash[0]));
    DIPC_EMIT_EVENT(state_string, &("CoinBase"[0]), &(coinbase[0]));
    DIPC_EMIT_EVENT(state_string, &("Balance"[0]), &(balance[0]));
    DIPC_EMIT_EVENT(state_string, &("Origin"[0]), &(origin[0]));
    DIPC_EMIT_EVENT(state_string, &("Caller"[0]), &(caller[0]));
    DIPC_EMIT_EVENT(state_string, &("CallerValue"[0]), &(callValue[0]));
    DIPC_EMIT_EVENT(state_string, &("Address"[0]), &(address[0]));
    print("\r\n dipc lib test stateTests success\r\n");
}

bytes generateHashSrcData(uint64_t datalen)
{
    bytes data(datalen);
    int i = 0;
    for (i = 0; i < datalen; i++)
    {
        data[i] = 0x66;
    }

    return data;
}

void sha3Test()
{
    //test sha3
    bytes dataSrc;
    uint64_t datalen = 1024 * 64;
    h256 correctResult("2c3b4230dea83ea0aa58f43f74f5cf3757f355290ee8b1113b5a4a3b3ff997e3");
    h256 hashData;

    dataSrc = generateHashSrcData(datalen);
    hashData = sha3(dataSrc);
    DipcAssertEQ(correctResult, hashData);
}