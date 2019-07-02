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

package chain

//import (
//	"github.com/dipperin/dipperin-core/core/model"
//	"github.com/ethereum/go-ethereum/rlp"
//	"io"
//)
//
//type LightProofRLP struct {
//	Header *model.Header
//	Link   model.InterLink
//}
//
//func (l LightProof) rlpProof() *LightProofRLP {
//	link := make(model.InterLink, len(l.Link))
//	copy(link, l.Link)
//	return &LightProofRLP{
//		Header: l.Header.(*model.Header),
//		Link:   link,
//	}
//}
//
//func (l LightProofRLP) proof(res *LightProof) {
//	res.Header = l.Header
//	res.Link = make(model.InterLink, len(l.Link))
//	copy(res.Link, l.Link)
//}
//
//func (l *LightProof) DecodeRLP(s *rlp.Stream) error {
//	var proof LightProofRLP
//	err := s.Decode(&proof)
//	proof.proof(l)
//	return err
//}
//
//func (l LightProof) EncodeRLP(w io.Writer) error {
//	return rlp.Encode(w, l.rlpProof())
//}
