package compress

import "github.com/minio/minlz"

func CompressPacket(binaryPacket []byte) ([]byte, error) {
	compressedData, err := minlz.Encode(nil, binaryPacket, minlz.LevelSmallest)
	if err != nil {
		return nil, err
	}
	return compressedData, nil
}

func DecompressedPkt(compressedPkt []byte) ([]byte, error) {
	decompressedPkt, err := minlz.Decode(nil, compressedPkt)
	if err != nil {
		return nil, err
	}
	return decompressedPkt, nil
}
