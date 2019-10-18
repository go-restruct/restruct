// Package png implements some of the PNG format using Restruct.
package png

// ColorType is used to specify the color format of a PNG.
type ColorType byte

// Enumeration of valid ColorTypes.
const (
	ColorGreyscale      ColorType = 0
	ColorTrueColor      ColorType = 2
	ColorIndexed        ColorType = 3
	ColorGreyscaleAlpha ColorType = 4
	ColorTrueColorAlpha ColorType = 6
)

// File contains the data of an image.
type File struct {
	Magic  [8]byte
	Header Chunk
	Chunks []Chunk `struct-while:"!_eof"`
}

// Chunk contains the data of a single chunk.
type Chunk struct {
	Len  uint32
	Type string     `struct:"[4]byte"`
	IHDR *ChunkIHDR `struct-if:"Type == $'IHDR'" json:",omitempty"`
	IDAT *ChunkIDAT `struct-if:"Type == $'IDAT'" json:",omitempty"`
	IEND *ChunkIEND `struct-if:"Type == $'IEND'" json:",omitempty"`
	CRC  uint32
}

// ChunkIHDR contains the body of a IHDR chunk.
type ChunkIHDR struct {
	Width             uint32
	Height            uint32
	BitDepth          byte
	ColorType         ColorType
	CompressionMethod byte
	FilterMethod      byte
	InterlaceMethod   byte
}

// ChunkIDAT contains the body of a IDAT chunk.
type ChunkIDAT struct {
	Parent *Chunk `struct:"parent" json:"-"`

	Data []byte `struct-size:"Parent.Len"`
}

// ChunkIEND contains the body of a IEND chunk.
type ChunkIEND struct {
}

// ChunkPLTE contains the body of a PLTE chunk.
type ChunkPLTE struct {
}
