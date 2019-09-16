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

package g_metrics

// all metrics in this file

const (
	ReceivedWaitVHashCount  = "wait_v_block_hash_count"
	ReceivedWaitVBlockCount = "wait_v_block_block_count"
	//NewTransportCount       = "new_transport_count"
	//DeleteTransportCount    = "delete_transport_count"

	ReceivedHashCount  = "v_block_hash_count"
	ReceivedBlockCount = "v_block_block_count"
	//FailedInsertBlockCount  = "failed_insert_block_count"

	// peers metrics
	NorPeerSetGauge   = "normal_peer_set"
	CurPeerSetGauge   = "cur_peer_set"
	NextPeerSetGauge  = "next_peer_set"
	VBootPeerSetGauge = "v_boot_peer_set"
	// The number of  handle peers have been processed
	TotalHandledPeer = "pm_total_handle"
	// The number of successful handshakes
	TotalSuccessHandle = "pm_total_handle_success"
	// The number of failed handshaking
	TotalFailedHandle = "pm_total_handle_failed"
	// Current peer processing statistics
	CurHandelPeer = "pm_cur_handle_peer"
	// running peer count
	RunningPeerGauge = "p2p_running_peers"

	// BFT fetch block
	FetchBlockGoCount = "pbft_fetch_block_go_count"
	BftCurStateGauge  = "pbft_state"
	BftCurRoundGauge  = "pbft_round"
	BftTimeoutCount   = "pbft_timeout"

	// v halt check
	CurBlockNumberGauge = "verBootNodeBlockNumber"

	PendingTxCountInPool = "pending_tx_count_in_pool"
	QueuedTxCountInPool  = "queued_tx_count_in_pool"

	CurChainHeight         = "cur_height"
	FailedInsertBlockCount = "failed_insert_block_count"
)

// call this after NewPrometheusMetricsServer
func InitCSMetrics() {
	CreateCounter(ReceivedWaitVBlockCount, "", nil)
	CreateCounter(ReceivedHashCount, "", nil)
	CreateCounter(ReceivedBlockCount, "", nil)

	CreateGauge(NorPeerSetGauge, "trace normal peer set len", nil)
	CreateGauge(CurPeerSetGauge, "trace cur peer set len", nil)
	CreateGauge(NextPeerSetGauge, "trace next peer set len", nil)
	CreateGauge(VBootPeerSetGauge, "v_boot_peer_set", nil)
	CreateCounter(TotalHandledPeer, "trace pm total handle peer", nil)
	CreateCounter(TotalSuccessHandle, "trace total handle success", nil)
	CreateCounter(TotalFailedHandle, "trace total handle failed", nil)
	CreateGauge(CurHandelPeer, "", nil)
	CreateGauge(RunningPeerGauge, "p2p_running_peers", nil)

	CreateGauge(FetchBlockGoCount, "pbft_fetch_block_go_count", nil)
	CreateGauge(BftCurStateGauge, "trace pbft state", nil)
	CreateGauge(BftCurRoundGauge, "trace pbft cur round", nil)
	CreateCounter(BftTimeoutCount, "trace timeout count", []string{"state_name"})

	CreateGauge(CurBlockNumberGauge, "trace pbft state", nil)

	CreateGauge(PendingTxCountInPool, "trace tx count", nil)
	CreateGauge(QueuedTxCountInPool, "trace tx count", nil)
	CreateGauge(CurChainHeight, "chain height", nil)
	CreateCounter(FailedInsertBlockCount, "trace failed insert block", nil)
}
