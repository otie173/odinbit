package compress

import "github.com/minio/minlz"

func CompressPkt(pkt []byte) ([]byte, error) {
	compressedPkt, err := minlz.Encode(nil, pkt, minlz.LevelSmallest)
	if err != nil {
		return nil, err
	}
	return compressedPkt, nil
}
