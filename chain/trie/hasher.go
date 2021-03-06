package trie

import (
	"bytes"
	"sync"

	"github.com/nknorg/nkn/v2/chain/db"
	"github.com/nknorg/nkn/v2/common"
	"github.com/nknorg/nkn/v2/util/log"
)

type hasher struct {
	tmp *bytes.Buffer
	sha common.Uint256
}

var hasherPool = sync.Pool{
	New: func() interface{} {
		return &hasher{tmp: new(bytes.Buffer), sha: common.Uint256{}}
	},
}

func (h *hasher) hash(n node, db Database, force bool) (node, node, error) {
	if hash, dirty := n.cache(); hash != nil {
		if !dirty {
			return hash, n, nil
		}
	}
	collapsed, cached, err := h.hasChildren(n, db)
	if err != nil {
		return nil, n, err
	}
	hashed, err := h.store(collapsed, db, force)
	if err != nil {
		return nil, n, err
	}
	cachedHash, _ := hashed.(hashNode)
	switch cn := cached.(type) {
	case *shortNode:
		cn.flags.hash = cachedHash
		cn.flags.dirty = false
	case *fullNode:
		cn.flags.hash = cachedHash
		cn.flags.dirty = false
	}
	return hashed, cached, nil
}

func (h *hasher) hasChildren(original node, db Database) (node, node, error) {
	var err error
	switch n := original.(type) {
	case *shortNode:
		collapsed, cached := n.copy(), n.copy()
		collapsed.Key = hexToCompact(n.Key)
		cached.Key = copyBytes(n.Key)
		if _, ok := n.Val.(valueNode); !ok {
			collapsed.Val, cached.Val, err = h.hash(n.Val, db, false)
			if err != nil {
				return original, original, err
			}
		}
		if collapsed.Val == nil {
			collapsed.Val = valueNode(nil)
		}
		return collapsed, cached, nil
	case *fullNode:
		collapsed, cached := n.copy(), n.copy()
		for i := 0; i < 16; i++ {
			if n.Children[i] != nil {
				collapsed.Children[i], cached.Children[i], err = h.hash(n.Children[i], db, false)
				if err != nil {
					return original, original, err
				}
			} else {
				collapsed.Children[i] = valueNode(nil)
			}
		}
		cached.Children[16] = n.Children[16]
		if collapsed.Children[16] == nil {
			collapsed.Children[16] = valueNode(nil)
		}
		return collapsed, cached, nil
	default:
		return n, original, nil
	}
}

func (h *hasher) store(n node, st Database, force bool) (node, error) {
	if _, isHash := n.(hashNode); n == nil || isHash {
		return n, nil
	}
	h.tmp.Reset()
	if err := n.Serialize(h.tmp); err != nil {
		log.Fatalf("Trie node enocde error: %v", err)
	}

	if h.tmp.Len() < 32 && !force {
		return n, nil
	}
	hs, _ := n.cache()
	if hs == nil {
		hs = Sha256Key(h.tmp.Bytes())
	}
	if st != nil {
		return hs, st.BatchPut(db.TrieNodeKey(hs), h.tmp.Bytes())
	}
	return hs, nil
}

func newHasher() *hasher {
	h := hasherPool.Get().(*hasher)
	return h
}

func returnHasherToPool(h *hasher) {
	hasherPool.Put(h)
}
