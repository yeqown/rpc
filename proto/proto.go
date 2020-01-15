// Package proto .
package proto

import (
	"bufio"
	"encoding/binary"
	"errors"
)

const (
	// OpRequest .
	OpRequest uint16 = iota + 1
	// OpResponse .
	OpResponse
)

const (
	// Ver1 .
	Ver1 uint16 = 1
	// Ver2 .
	Ver2 = 2
)

var (
	// ErrProtoHeaderLen .
	ErrProtoHeaderLen = errors.New("not matched proto header len")
	// ErrEmptyReader .
	ErrEmptyReader = errors.New("empty reader")
)

const (
	// size
	_packSize      uint16 = 4
	_headerSize    uint16 = 2 // uint16
	_verSize       uint16 = 2 // uint16
	_opSize        uint16 = 2 // uint16
	_seqSize       uint16 = 2 // uint16
	_rawHeaderSize uint16 = _packSize + _headerSize + _verSize + _opSize + _seqSize

	// offset
	_packOffset   uint16 = 0
	_headerOffset        = _packOffset + _packSize
	_verOffset           = _headerOffset + _headerSize
	_opOffset            = _verOffset + _verSize
	_seqOffset           = _opOffset + _opSize
)

// Proto .
type Proto struct {
	Ver  uint16
	Op   uint16 // Type of Proto
	Seq  uint16 // Seq of message, 0 means done, else means not finished
	Body []byte // Body of Proto
}

// New .
func New() *Proto {
	return &Proto{
		Ver:  Ver1,
		Op:   OpRequest,
		Seq:  0,
		Body: nil,
	}
}

// WriteTCP .
// packLen(32bit):headerLen(16bit):ver(16bit):op(16bit):body
func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {
	var (
		buf     = make([]byte, _rawHeaderSize)
		packLen int
	)

	packLen = int(_rawHeaderSize) + len(p.Body)
	binary.BigEndian.PutUint32(buf[_packOffset:], uint32(packLen))
	binary.BigEndian.PutUint16(buf[_headerOffset:], _rawHeaderSize)
	binary.BigEndian.PutUint16(buf[_verOffset:], p.Ver)
	binary.BigEndian.PutUint16(buf[_opOffset:], p.Op)
	binary.BigEndian.PutUint16(buf[_seqOffset:], p.Seq)

	if _, err = wr.Write(buf); err != nil {
		return
	}

	if p.Body != nil {
		_, err = wr.Write(p.Body)
	}

	// println(wr.Buffered(), len(p.Body))
	return
}

// ReadTCP .
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen uint16
		packLen   int
		buf       []byte
	)

	if buf, err = ReadNBytes(rr, int(_rawHeaderSize)); err != nil {
		return
	}

	packLen = int(binary.BigEndian.Uint32(buf[_packOffset:_headerOffset]))
	headerLen = binary.BigEndian.Uint16(buf[_headerOffset:_verOffset])
	p.Ver = binary.BigEndian.Uint16(buf[_verOffset:_opOffset])
	p.Op = binary.BigEndian.Uint16(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Uint16(buf[_seqOffset:])

	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}

	if bodyLen = packLen - int(headerLen); bodyLen > 0 {
		p.Body, err = ReadNBytes(rr, bodyLen)
	} else {
		p.Body = nil
	}

	return
}

// ReadNBytes . read limitted `N` bytes from bufio.Reader.
func ReadNBytes(rr *bufio.Reader, N int) ([]byte, error) {
	if rr == nil {
		return nil, ErrEmptyReader
	}

	var (
		buf = make([]byte, N)
		err error
	)
	for i := 0; i < N; i++ {
		if buf[i], err = rr.ReadByte(); err != nil {
			return nil, err
		}
	}

	return buf, err
}
