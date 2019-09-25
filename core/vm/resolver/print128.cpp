// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
#include "print128.h"

#include <stdint.h>
#include <iterator>

const char* printi128(uint64_t lo, uint64_t hi) {
    __int128 u = hi;
    u <<= 64;
    u |= lo;

    unsigned __int128 tmp = u < 0 ? -u : u;
    static char buffer[128+1];
    buffer[sizeof(buffer)-1] = '\0';
    //char* d = std::end(buffer)-1;
    char* d = &buffer[sizeof(buffer)-2];
    do
    {
        --d;
        *d = "0123456789"[ tmp % 10 ];
        tmp /= 10;
    } while ( tmp != 0 );
    if ( u < 0 ) {
        --d;
        *d = '-';
    }

    return d;
}

const char* printui128(uint64_t lo, uint64_t hi) {
    unsigned __int128 u = hi;
    u <<= 64;
    u |= lo;

    unsigned __int128 tmp = u;
    static char buffer[128+1];
    buffer[sizeof(buffer)-1] = '\0';
    //char* d = std::end(buffer)-1;
    char* d = &buffer[sizeof(buffer)-2];
    do
    {
        --d;
        *d = "0123456789"[ tmp % 10 ];
        tmp /= 10;
    } while ( tmp != 0 );
    return d;
}

