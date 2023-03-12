package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From      string
	To        string
	Conv      map[string]bool
	Offset    int
	Limit     int
	BlockSize int
}

func ParseFlags() *Options {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "count of bytes to skip. by default - 0")
	flag.IntVar(&opts.Limit, "limit", -1, "max count of bytes to read. by default - -1")
	flag.IntVar(&opts.BlockSize, "block-size", -1, "max count of bytes to read and write."+
		" by default - -1")
	conv := flag.String("conv", "", "conversions over text: upper_case, lower_case and trim_spaces")

	flag.Parse()

	valideArgs := map[string]bool{"trim_spaces": true, "upper_case": true, "lower_case": true}

	opts.Conv = make(map[string]bool)
	for _, arg := range strings.Split(*conv, ",") {
		opts.Conv[arg] = true
		if !valideArgs[arg] && arg != "" {
			_, _ = fmt.Fprintln(os.Stderr, arg, " is not a validate arg")
			os.Exit(1)
		}
	}

	return &opts
}

func OpenFile(filePath string, read bool) (io.ReadWriteCloser, error) {
	var (
		stream io.ReadWriteCloser
		err    error
	)

	if filePath == "" {
		if read {
			stream = os.Stdin
		} else {
			stream = os.Stdout
		}
	} else {
		if read {
			stream, err = os.Open(filePath)
		} else {
			_, err := os.Stat(filePath)
			if !errors.Is(err, os.ErrNotExist) {
				_, _ = fmt.Fprintln(os.Stderr, "file is exist", err)
				os.Exit(1)
			}
			stream, err = os.Create(filePath)
		}
	}
	return stream, err
}

func MustReadData(filePath string, offset int, limit int) string {
	stream, err := OpenFile(filePath, true)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not open a file:", filePath)
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(stream)
	_, err = reader.Discard(offset)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: offset is larger than file size")
		os.Exit(1)
	}

	var data []byte
	for limit != 0 {
		readedByte, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not read a data:", err)
			os.Exit(1)
		}
		data = append(data, readedByte)
		limit--
	}

	defer func(stream io.ReadWriteCloser) {
		err := stream.Close()
		if err != nil {

		}
	}(stream)

	return string(data)
}

func MustConvertData(data string, conv map[string]bool) string {
	if conv["upper_case"] != false && conv["lower_case"] != false {
		_, _ = fmt.Fprintln(os.Stderr, "it is not possible to convert text"+
			" to upper and lower case at the same time")
		os.Exit(1)
	}

	if conv["trim_spaces"] {
		data = strings.TrimSpace(data)
	}

	if conv["upper_case"] {
		data = strings.ToUpper(data)
	}

	if conv["lower_case"] {
		data = strings.ToLower(data)
	}
	return data
}

func MustWriteData(filePath string, data []byte) {
	stream, err := OpenFile(filePath, false)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not open a file:", filePath)
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(stream)
	writer.Write(data)
	defer writer.Flush()
}

func main() {
	opts := ParseFlags()
	data := MustReadData(opts.From, opts.Offset, opts.Limit)
	data = MustConvertData(data, opts.Conv)
	MustWriteData(opts.To, []byte(data))
}
