// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package chain_communication

//import (
//	"github.com/dipperin/dipperin-core/common"
//	"github.com/dipperin/dipperin-core/common/g-error"
//	"github.com/hashicorp/golang-lru"
//)
//
//const (
//	isVCacheSize = 10
//)
//
//// use this for get self is verifier, when get verifiers not in cache
//func NewIsVerifierCacheSet() *IsVerifierCacheSet {
//	return &IsVerifierCacheSet{
//		cache: make(map[common.Address]*isVerifierCache),
//	}
//}
//
//type IsVerifierCacheSet struct {
//	cache map[common.Address]*isVerifierCache
//}
//
//func (s *IsVerifierCacheSet) IsVerifier(address common.Address, slot uint64) (bool, error) {
//	c := s.cache[address]
//	if c == nil {
//		return false, g_error.ErrNotInIsVerifierCache
//	}
//	return c.IsVerifier(slot)
//}
//
//func (s *IsVerifierCacheSet) SetIsVerifier(address common.Address, slot uint64, isV bool) {
//	if s.cache[address] == nil {
//		s.cache[address] = NewIsVerifierCache()
//	}
//	s.cache[address].SetIsVerifier(slot, isV)
//}
//
//func NewIsVerifierCache() *isVerifierCache {
//	v, _ := lru.New(isVCacheSize)
//	return &isVerifierCache{
//		isV: v,
//	}
//}
//
//type isVerifierCache struct {
//	isV *lru.Cache
//}
//
//func (c *isVerifierCache) IsVerifier(slot uint64) (bool, error) {
//	if v, ok := c.isV.Get(slot); ok {
//		return v.(bool), nil
//	}
//	return false, g_error.ErrNotInIsVerifierCache
//}
//
//func (c *isVerifierCache) SetIsVerifier(slot uint64, isV bool) {
//	c.isV.Add(slot, isV)
//}
