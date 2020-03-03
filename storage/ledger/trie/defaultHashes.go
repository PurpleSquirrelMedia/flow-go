package trie

var defaultLeafHash []byte = Hash(EmptySlice)
// we are currently supporting keys of a size up to 32 bytes. I.e. path length from the root of a fully expanded tree to the leaf node is 256. A path of length k is comprised of k+1 vertices. Hence, we need 257 default hashes. 
var defaultHashes [257][]byte 

func init() {
	// Creates the Default hashes from base to level height
	defaultHashes[0] = defaultLeafHash
	for i := 1; i < len(defaultHashes); i++ {
		defaultHashes[i] = Hash(defaultHashes[i-1], defaultHashes[i-1])
	}
}

// GetDefaultHashes returns the default hashes of the SMT.
//
// For each tree level N, there is a default hash equal to the chained
// hashing of the default value N times.
func GetDefaultHashes() [257][]byte {
	return defaultHashes
}

// GetDefaultHashForHeight returns the default hashes of the SMT at a specified height.
//
// For each tree level N, there is a default hash equal to the chained
// hashing of the default value N times.
func GetDefaultHashForHeight(height int) []byte {
	return defaultHashes[height]
}
