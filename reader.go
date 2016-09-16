package main

import (
	"io"
  "bytes"
)

import "gopkg.in/cheggaaa/pb.v1"

// It's proxy reader, implement io.Reader
type Reader struct {
	Reader *io.SectionReader
	bar *pb.ProgressBar
  offset int64
}

func NewReader(buffer []byte, bar *pb.ProgressBar) *Reader {
    reader := bytes.NewReader(buffer)

    return &Reader { io.NewSectionReader(reader, 0, bar.Total), bar, 0 }
}

func (r *Reader) Read(p []byte) (n int, err error) {
  if r.offset >= r.bar.Total {
   return 0, io.EOF
  }

  n, err =  r.Reader.Read(p)

  r.bar.Add(n)
  r.offset += int64(n)

	return n, err
}
