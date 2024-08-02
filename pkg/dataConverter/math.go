package dataConverter

func FromUint24TRoUint32(v []byte) uint32 {
	return uint32(v[2])<<16 | uint32(v[1])<<8 | uint32(v[0])
}
