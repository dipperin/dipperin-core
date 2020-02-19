#include "buildins_test.hpp"
#include "rt/compiler_builtins.hpp"

unsigned __int128 operator"" _ULLL(const char *lit)
{
        __int128 ret = 0;
        size_t i = 0;
        bool sign = false;

        if (lit[i] == '-')
        {
                ++i;
                sign = true;
        }

        if (lit[i] == '+')
                ++i;

        for (; lit[i] != '\0'; ++i)
        {
                const char c = lit[i];
                ret *= 10;
                ret += c - '0';
        }

        if (sign)
                ret *= -1;

        return (unsigned __int128)ret;
}

__int128 operator"" _LLL(const char *lit)
{
        __int128 ret = 0;
        size_t i = 0;
        bool sign = false;

        if (lit[i] == '-')
        {
                ++i;
                sign = true;
        }

        if (lit[i] == '+')
                ++i;

        for (; lit[i] != '\0'; ++i)
        {
                const char c = lit[i];
                ret *= 10;
                ret += c - '0';
        }

        if (sign)
                ret *= -1;

        return ret;
}

void buildins::arithmeticTest() {
        #if 0
        print("\r\n dipc lib test arithmeticTest start\r\n");
        {
                __int128 res = 0;
                __int128 val = 1;
                __int128 test_res = 0;

                test_res = 0x8000000000000000;
                test_res <<= 1;

                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 0);
                DipcAssert(res == 1, "__ashlti3 result should be 1");

                print("\r\n dipc lib test callstack overflow -------1\r\n");
                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 1);
                DipcAssert(res == (1 << 1), "__ashlti3 result should be 2");

                print("\r\n dipc lib test callstack overflow -------2\r\n");
                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 31);
                DipcAssert(res == (__int128)2147483648_ULLL, "__ashlti3 result should be 2^31");

                print("\r\n dipc lib test callstack overflow -------3\r\n");
                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 63);
                DipcAssert(res == (__int128)9223372036854775808_ULLL,
                                "__ashlti3 result should be 2^63");

                print("\r\n dipc lib test callstack overflow -------4\r\n");
                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 64);
                DipcAssert(res == test_res, "__ashlti3 result should be 2^64");

                print("\r\n dipc lib test callstack overflow -------5\r\n");
                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 127);
                test_res <<= 63;
                DipcAssert(res == test_res, "__ashlti3 result should be 2^127");

                print("\r\n dipc lib test callstack overflow -------6\r\n");
                __ashlti3(res, uint64_t(val), uint64_t(val >> 64), 128);
                test_res <<= 1;
                // should rollover
                DipcAssert(res == test_res, "__ashlti3 result should be 2^128");
        }

        {
                __int128 res = 0;
                __int128 test = 1;
                __int128 val = -170141183460469231731687303715884105728_LLL;

                test <<= 127;

                __ashrti3(res, uint64_t(val), uint64_t(val >> 64), 0);
                DipcAssert(res == -170141183460469231731687303715884105728_LLL,
                                "__ashrti3 result should be -2^127");

                __ashrti3(res, uint64_t(val), uint64_t(val >> 64), 1);
                DipcAssert(res == -85070591730234615865843651857942052864_LLL,
                                "__ashrti3 result should be -2^126");

                __ashrti3(res, uint64_t(val), uint64_t(val >> 64), 2);
                DipcAssert(res == test >> 2, "__ashrti3 result should be -2^125");

                __ashrti3(res, uint64_t(val), uint64_t(val >> 64), 64);
                DipcAssert(res == test >> 64, "__ashrti3 result should be -2^63");

                __ashrti3(res, uint64_t(val), uint64_t(val >> 64), 95);
                DipcAssert(res == test >> 95, "__ashrti3 result should be -2^31");

                __ashrti3(res, uint64_t(val), uint64_t(val >> 64), 127);
                DipcAssert(res == test >> 127, "__ashrti3 result should be -2^0");
        }

        {
                __int128 res = 0;
                __int128 val = 1;
                __int128 test_res = 0;

                test_res = 0x8000000000000000;
                test_res <<= 1;
     
                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 0);
                DipcAssert(res == 1, "__lshlti3 result should be 1");

                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 1);
                DipcAssert(res == (1 << 1), "__lshlti3 result should be 2");

                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 31);
                DipcAssert((unsigned __int128)res == 2147483648_ULLL,
                                "__lshlti3 result should be 2^31");

                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 63);
                DipcAssert((unsigned __int128)res == 9223372036854775808_ULLL,
                                "__lshlti3 result should be 2^63");

                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 64);
                DipcAssert(res == test_res, "__lshlti3 result should be 2^64");

                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 127);
                test_res <<= 63;
                DipcAssert(res == test_res, "__lshlti3 result should be 2^127");

                __lshlti3(res, uint64_t(val), uint64_t(val >> 64), 128);
                test_res <<= 1;
                // should rollover
                DipcAssert(res == test_res, "__lshlti3 result should be 2^128");
        }

        {
                __int128 res = 0;
                __int128 val = 0x8000000000000000;
                __int128 test_res = 0x8000000000000000;

                val <<= 64;
                test_res <<= 64;
             
                __lshrti3(res, uint64_t(val), uint64_t(val >> 64), 0);
                DipcAssert(res == test_res, "__lshrti3 result should be 2^127");

                __lshrti3(res, uint64_t(val), uint64_t(val >> 64), 1);
                DipcAssert(res == (__int128)85070591730234615865843651857942052864_ULLL,
                                "__lshrti3 result should be 2^126");

                __lshrti3(res, uint64_t(val), uint64_t(val >> 64), 63);
                DipcAssert(res == (__int128)18446744073709551616_ULLL,
                                "__lshrti3 result should be 2^64");

                __lshrti3(res, uint64_t(val), uint64_t(val >> 64), 64);
                DipcAssert(res == (__int128)9223372036854775808_ULLL,
                                "__lshrti3 result should be 2^63");

                __lshrti3(res, uint64_t(val), uint64_t(val >> 64), 96);
                DipcAssert(res == (__int128)2147483648_ULLL, "__lshrti3 result should be 2^31");

                __lshrti3(res, uint64_t(val), uint64_t(val >> 64), 127);
                DipcAssert(res == 0x1, "__lshrti3 result should be 2^0");
        }

        {
                __int128 a = 1;
                __int128 ret = 0;
                __divti3(ret, uint64_t(a), uint64_t(a >> 64), uint64_t(a), uint64_t(a >> 64));

                /*
        * test for negative values
        */
                __int128 res = 0;
                __int128 lhs_a = -30;
                __int128 rhs_a = 100;
                __int128 lhs_b = 100;
                __int128 rhs_b = -30;
                __divti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 0, "__divti3 result should be 0");

                __divti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == -3, "__divti3 result should be -3");

                __divti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == 1, "__divti3 result should be 1");

                /*
        * test for positive values
        */
                __int128 lhs_c = 3333;
                __int128 rhs_c = 3333;

                __divti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 1, "__divti3 result should be 1");

                __divti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_c),
                                uint64_t(rhs_c >> 64));
                DipcAssert(res == 0, "__divti3 result should be 0");

                __divti3(res, uint64_t(lhs_c), uint64_t(lhs_c >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 33, "__divti3 result should be 33");

                /*
        * test identity
        */
                __divti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), 1, 0);
                DipcAssert(res == 100, "__divti3 result should be 100");

                __divti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), 1, 0);
                DipcAssert(res == -30, "__divti3 result should be -30");
        }

        {
                /*
        * test for negative values
        */
                unsigned __int128 res = 0;
                unsigned __int128 lhs_a = (unsigned __int128)-30;
                unsigned __int128 rhs_a = 100;
                unsigned __int128 lhs_b = 100;
                unsigned __int128 rhs_b = (unsigned __int128)-30;

                __udivti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 3402823669209384634633746074317682114_ULLL,
                                "__udivti3 result should be 0");

                __udivti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == 0, "__udivti3 result should be 0");

                __udivti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == 1, "__udivti3 result should be 1");

                /*
        * test for positive values
        */
                __int128 lhs_c = 3333;
                __int128 rhs_c = 3333;

                __udivti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 1, "__divti3 result should be 1");

                __udivti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_c),
                                uint64_t(rhs_c >> 64));
                DipcAssert(res == 0, "__divti3 result should be 0");

                __udivti3(res, uint64_t(lhs_c), uint64_t(lhs_c >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 33, "__divti3 result should be 33");

                /*
        * test identity
        */
                __udivti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), 1, 0);
                DipcAssert(res == 100, "__divti3 result should be 100");

                __udivti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), 1, 0);
                DipcAssert(res == (unsigned __int128)-30, "__divti3 result should be -30");
        }

        {
                __int128 res = 0;
                __int128 lhs_a = -30;
                __int128 rhs_a = 100;
                __int128 lhs_b = 30;
                __int128 rhs_b = -100;

                __modti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == -30, "__modti3 result should be -30");

                __modti3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == 30, "__modti3 result should be 30");

                __modti3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == -30, "__modti3 result should be -30");

                __modti3(res, uint64_t(rhs_a), uint64_t(rhs_a >> 64), uint64_t(lhs_b),
                                uint64_t(lhs_b >> 64));
                DipcAssert(res == 10, "__modti3 result should be 10");

                __modti3(res, uint64_t(rhs_a), uint64_t(rhs_a >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == 0, "__modti3 result should be 0");

                __modti3(res, uint64_t(rhs_a), uint64_t(rhs_a >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 0, "__modti3 result should be 0");

                __modti3(res, 0, 0, uint64_t(rhs_a), uint64_t(rhs_a >> 64));
                DipcAssert(res == 0, "__modti3 result should be 0");
        }

        {
                /*
        * tests for negative values
        */
                __int128 res = 0;
                __int128 lhs_a = -30;
                __int128 rhs_a = 100;
                __int128 lhs_b = 100;
                __int128 rhs_b = -30;

                __multi3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == -3000, "__multi3 result should be -3000");

                __multi3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == -3000, "__multi3 result should be -3000");

                __multi3(res, uint64_t(lhs_a), uint64_t(lhs_a >> 64), uint64_t(rhs_b),
                                uint64_t(rhs_b >> 64));
                DipcAssert(res == 900, "__multi3 result should be 900");

                /*
        * test for positive values
        */
                __multi3(res, uint64_t(lhs_b), uint64_t(lhs_b >> 64), uint64_t(rhs_a),
                                uint64_t(rhs_a >> 64));
                DipcAssert(res == 10000, "__multi3 result should be 10000");

                /*
        * test identity
        */
                __multi3(res, 1, 0, uint64_t(rhs_a), uint64_t(rhs_a >> 64));
                DipcAssert(res == 100, "__multi3 result should be 100");

                __multi3(res, 1, 0, uint64_t(rhs_b), uint64_t(rhs_b >> 64));
                DipcAssert(res == -30, "__multi3 result should be -30");
        } 
        print("\r\n dipc lib test arithmeticTest success \r\n");
        #endif
}
