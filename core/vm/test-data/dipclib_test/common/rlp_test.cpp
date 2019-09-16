#include "rlp_test.hpp"

DIPC_EVENT(rlp_string, const char *, const char *);
DIPC_EVENT(rlp_double, const char *, double);

void rlptest::rlpTest()
{
    print("\r\n dipc lib test rlpTest start\r\n");
    {
        std::string data = "c3010203";
        RLPStream stream(3);
        stream << 1 << 2 << 3;
        std::string result = toHex(stream.out());
        DipcAssertEQ(data, result);
    }

    {
        std::string data = "9443355c787c50b647c425f594b441d4bd751951c1";
        RLPStream stream;
        stream << Address("0x43355c787c50b647c425f594b441d4bd751951c1", true);
        std::string result = toHex(stream.out());
        DipcAssertEQ(data, result);
    }

    {
        std::string data =
            "aa307834333335356337383763353062363437633432356635393462343431643462643735313935316331";
        RLPStream stream;
        stream << "0x43355c787c50b647c425f594b441d4bd751951c1";
        std::string result = toHex(stream.out());
        DipcAssertEQ(data, result);
    }
    print("\r\n dipc lib test rlpTest success\r\n");
}