package internal

import "io"

type ArgsIter struct {
	Args []Message
	Pos  int
	Len  int
}

func NewArgsIterator(args []Message) ArgsIter {
	return ArgsIter{
		Args: args,
		Pos:  0,
		Len:  len(args),
	}
}

func (iter ArgsIter) Peek() (*Message, error) {
	if iter.Pos >= iter.Len {
		return nil, io.EOF
	}
	return &iter.Args[iter.Pos], nil
}

func (iter *ArgsIter) Next() (*Message, error) {
	if iter.Pos >= iter.Len {
		return nil, io.EOF
	}

	val := iter.Args[iter.Pos]
	iter.Pos += 1

	return &val, nil
}
