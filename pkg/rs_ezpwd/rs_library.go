package rs_ezpwd

func EncodeWrapper(n, f int64, data []byte) [][]byte {
	chunkSize := n - 2*f
	chunksCount := 1 + ((int64(len(data)) - 1) / chunkSize)

	result := make([][]byte, n)
	for i := range result {
		result[i] = make([]byte, chunksCount)
	}

	for i := int64(0); i < chunksCount; i++ {
		toEncode := NewVecUInt8(chunkSize)
		for j, b := range data[i*chunkSize : (i+1)*chunkSize] {
			toEncode.Set(j, b)
		}
		encoded := Encode(uint(f), toEncode)
		for j := 0; j < int(encoded.Size()); j++ {
			result[j][i] = encoded.Get(j)
		}
		DeleteVecUInt8(encoded)
		DeleteVecUInt8(toEncode)
	}

	return result
}

func DecodeWrapper(n, f int64, data [][]byte, missing []int) []byte {
	missingVec := NewVecInt()
	for _, m := range missing {
		missingVec.Add(m)
	}

	result := make([]byte, len(data[0])*int(n-2*f))

	for i := 0; i < len(data[0]); i++ {
		toDecode := NewVecUInt8(n)
		for j := 0; j < len(data); j++ {
			if data[j] != nil {
				toDecode.Set(j, data[j][i])
			}
		}
		decoded := Decode(uint(f), toDecode, missingVec)
		if decoded.Size() == 0 {
			// Fail
			DeleteVecUInt8(decoded)
			DeleteVecInt(missingVec)
			DeleteVecUInt8(toDecode)
			return []byte{}
		}

		for j := 0; j < int(decoded.Size()); j++ {
			result[i*(int(n-2*f))+j] = decoded.Get(j)
		}
		DeleteVecUInt8(decoded)
		DeleteVecUInt8(toDecode)
	}

	DeleteVecInt(missingVec)
	return result
}

func bytesToVec(data []byte) VecUInt8 {
	vec := NewVecUInt8(int64(len(data)))
	for i, b := range data {
		vec.Set(i, b)
	}
	return vec
}

func bytesToVecVec(data [][]byte) VecVecUInt8 {
	vec := NewVecVecUInt8(int64(len(data)))
	for i, v := range data {
		vec.Set(i, bytesToVec(v))
	}
	return vec
}

func vecToBytes(vec VecUInt8) []byte {
	data := make([]byte, vec.Size())
	for i := 0; i < int(vec.Size()); i++ {
		data[i] = vec.Get(i)
	}
	return data
}

//func main() {
//	vec := NewVecUInt8(int64(0))
//	for i := 0; i < 128; i++ {
//		vec.Add(byte(i))
//	}
//	println(vec.Size())
//	encoded := Encode(6, vec)
//	decoded := Decode(6, encoded, []int{})
//
//	println(decoded.Size())
//	println(decoded)
//}
