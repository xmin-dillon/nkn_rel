package config

import "github.com/nknorg/nkn/v2/common"

type HeightDependentInt32 struct {
	heights []uint32
	values  []int32
}

func (hdi *HeightDependentInt32) GetValueAtHeight(height uint32) int32 {
	for i, h := range hdi.heights {
		if height >= h {
			return hdi.values[i]
		}
	}
	return 0
}

type HeightDependentInt64 struct {
	heights []uint32
	values  []int64
}

func (hdi *HeightDependentInt64) GetValueAtHeight(height uint32) int64 {
	for i, h := range hdi.heights {
		if height >= h {
			return hdi.values[i]
		}
	}
	return 0
}

type HeightDependentUint256 struct {
	heights []uint32
	values  []common.Uint256
}

func (hdi *HeightDependentUint256) GetValueAtHeight(height uint32) common.Uint256 {
	for i, h := range hdi.heights {
		if height >= h {
			return hdi.values[i]
		}
	}
	return common.EmptyUint256
}

type HeightDependentBool struct {
	heights []uint32
	values  []bool
}

func (hdi *HeightDependentBool) GetValueAtHeight(height uint32) bool {
	for i, h := range hdi.heights {
		if height >= h {
			return hdi.values[i]
		}
	}
	return false
}

type HeightDependentString struct {
	heights []uint32
	values  []string
}

func (hdi *HeightDependentString) GetValueAtHeight(height uint32) string {
	for i, h := range hdi.heights {
		if height >= h {
			return hdi.values[i]
		}
	}
	return ""
}
