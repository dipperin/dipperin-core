#include "convert.hpp"

void convert::init() {}

void convert::toString()
{
    {
        float f = 1.000001f;
        char buf[10];
        ::snprintf(buf, sizeof(buf), "%f", f);
        DipcAssert(strcmp(buf, "1.000001") == 0);
    }

    {
        float f = 2e-6f / 3.0f;
        char buf[10];
        ::snprintf(buf, sizeof(buf), "%f", f);
        DipcAssert(strcmp(buf, "0.000001") == 0);
    }

    {
        double f = 1.000001;
        char buf[10];
        ::snprintf(buf, sizeof(buf), "%lf", f);
        DipcAssert(strcmp(buf, "1.000001") == 0);
    }

    {
        double f = 2e-6 / 3.0;
        char buf[10];
        ::snprintf(buf, sizeof(buf), "%lf", f);
        DipcAssert(strcmp(buf, "0.000001") == 0);
    }

    {
        long double f = 1.000001l;
        char buf[10];
        ::snprintf(buf, sizeof(buf), "%Lf", f);
        DipcAssert(strcmp(buf, "1.000001") == 0, buf);
    }

    {
        long double f = 2e-6l / 3.0l;
        char buf[10];
        ::snprintf(buf, sizeof(buf), "%Lf", f);
        DipcAssert(strcmp(buf, "0.000001") == 0);
    }
}

void convert::getBlockInfo()
{
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


    DIPC_EMIT_EVENT(uint64, &("blockNum"[0]), num);
    DIPC_EMIT_EVENT(uint64, &("GasLimit"[0]), gasLimit);

    DIPC_EMIT_EVENT(int64, &("GasPrice"[0]), gasPrice);
    DIPC_EMIT_EVENT(int64, &("TimeStamp"[0]), timestamp);
    DIPC_EMIT_EVENT(int64, &("Nonce"[0]), nonce);

    DIPC_EMIT_EVENT(string, &("blockHash"[0]), &(blockHash[0]));
    DIPC_EMIT_EVENT(string, &("CoinBase"[0]), &(coinbase[0]));
    DIPC_EMIT_EVENT(string, &("Balance"[0]), &(balance[0]));
    DIPC_EMIT_EVENT(string, &("Origin"[0]), &(origin[0]));
    DIPC_EMIT_EVENT(string, &("Caller"[0]), &(caller[0]));
    DIPC_EMIT_EVENT(string, &("CallerValue"[0]), &(callValue[0]));
    DIPC_EMIT_EVENT(string, &("Address"[0]), &(address[0]));
}

void convert::rlpTest()
{
    std::string data = "9600005586b883ec6dd4f8c26063e18eb4bd228e59c3e9";
    RLPStream stream;
    stream << Address2("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9", true);
    std::string result = toHex(stream.out());
    DipcAssertEQ(data, result);
    DIPC_EMIT_EVENT(string, &(data[0]), &(result[0]));
}

void convert::printTest()
{
    {
        prints("hello");
        prints_l("world", 5);
        printi(-99999999998);
        printui(99999999998);
        print(true);
    }

    // test int128 and uint128
    {
        __int128 i = std::numeric_limits<__int128>::lowest();
        printi128(&i);

        unsigned __int128 u = std::numeric_limits<unsigned __int128>::max();
        printui128(&u);
    }

    // test float
    {
        float f = 1.0f / 2.0f;
        printsf(f);

        f = 5.0f * -0.75f;
        printsf(f);

        f = 2e-6f / 3.0f;
        printsf(f);
    }

    // test double
    {
        double f = 1.0 / 2.0;
        printdf(f);

        f = 5.0 * -0.75;
        printdf(f);

        f = 2e-6 / 3.0;
        printdf(f);
    }

    // test long double
    {
        long double f = 1.0l / 2.0l;
        printqf(&f);

        f = 5.0l * -0.75l;
        printqf(&f);

        f = 2e-6l / 3.0l;
        printqf(&f);
    }

    // test hex
    {
        FixedHash<4> f("d3adb33f");
        printhex(f.data(), f.size());
    }
}
