/*
Package client implements the TICS 02 - Client Semantics specification. This
concrete implementations defines types and method to store and update light
clients which tracks on other chain's state.

The main type is `Client`, which provides `commitment.Root` to verify state proofs and `ConsensusState` to
verify header proofs.
*/
package client
