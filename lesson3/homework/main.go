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

// Options хранит настройки чтения и записи
type Options struct {
	From      string          // путь к исходному файлу
	To        string          // путь к копии
	Conv      map[string]bool // параметры форматирования
	Offset    int64           // кол-во игнорируемых байт
	Limit     int64           // максимальное кол-во считываемых байт
	BlockSize int64           // максимальное кол-во байт, хранимых в памяти одновременно
}

// SetDefault устанавливает значения настроек по умолчанию
func (opts *Options) SetDefault() {
	opts.From = ""
	opts.To = ""
	opts.Conv = make(map[string]bool)
	opts.Offset = 0
	opts.Limit = -1
	opts.BlockSize = -1
}

// validateFlags проверяет все флаги на валидность
func validateFlags(options *Options) error {
	if options.From != "" {
		// проверка того, что файл с данными существует
		_, err := os.Stat(options.From)
		if errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("the file %s is not exist."+
				" there is no way to read the file", options.From))
		}

		// проверка того, что размер файла больше кол-ва пропускаемых байт
		f, _ := os.Stat(options.From)
		if f.Size() < options.Offset && options.Offset > 0 {
			return errors.New("offset is bigger than file size." +
				" there is no way to read the file")
		}
	}

	if options.To != "" {
		// проверка того, что не существует файла, в который будут записываться данные
		_, err := os.Stat(options.To)
		if !errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("the file on path %s is already exist."+
				" there is no way to write data", options.To))
		}
	}

	// проверка того, что все опции, переданные во флаг conv валидны
	validArgs := map[string]bool{"trim_spaces": true, "upper_case": true, "lower_case": true}
	for arg := range options.Conv {
		if !validArgs[arg] {
			return errors.New(fmt.Sprintf("conv arg <%s> is not validate", arg))
		}
	}

	// проверки того, что данные во флагах имеют смысл
	if options.Offset < 0 {
		return errors.New("the value of offset must be positive")
	}

	if options.Limit != -1 && options.Limit < 0 {
		return errors.New("the count of bytes to read must be positive")
	}

	if options.BlockSize != -1 && options.BlockSize <= 0 {
		return errors.New("the block size must be positive")
	}

	return nil
}

// ParseFlags парсит параметры и возвращает их вместе с ошибкой
func ParseFlags() (*Options, error) {
	var opts Options
	opts.SetDefault()

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "count of bytes to skip. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", -1, "max count of bytes to read. by default - -1")
	flag.Int64Var(&opts.BlockSize, "block-size", -1, "max count of bytes to read and write."+
		" by default - -1")
	conv := *flag.String("conv", "", "conversions over text: upper_case, lower_case and trim_spaces")

	flag.Parse()

	for _, arg := range strings.Split(conv, ",") {
		opts.Conv[arg] = true
	}

	err := validateFlags(&opts)

	return &opts, err
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

func MustReadData(filePath string, offset int64, limit int64) string {
	stream, err := OpenFile(filePath, true)

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not open a file:", filePath)
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(stream)
	_, err = reader.Discard(int(offset))
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
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	data := MustReadData(opts.From, opts.Offset, opts.Limit)
	data = MustConvertData(data, opts.Conv)
	MustWriteData(opts.To, []byte(data))
}
