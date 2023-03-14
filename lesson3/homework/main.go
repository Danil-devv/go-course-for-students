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

// ConvWriteReader это структура, позволяющая форматированно
// считывать и записывать данные в поток ввода/вывода stream
type ConvWriteReader struct {
	stream io.ReadWriter
}

// NewConvWriteReader создает из потока ввода/вывода структуру ConvWriteReader
func NewConvWriteReader(stream io.ReadWriter) *ConvWriteReader {
	return &ConvWriteReader{stream}
}

// Write записывает все байты из слайса data,
// при этом предварительно форматирует data согласно флагам в conv
func (u *ConvWriteReader) Write(data []byte, conv *map[string]bool) (int, error) {
	buf := string(data)
	if (*conv)["trim_spaces"] {
		buf = strings.TrimSpace(buf)
	}

	if (*conv)["upper_case"] {
		buf = strings.ToUpper(buf)
	}

	if (*conv)["lower_case"] {
		buf = strings.ToLower(buf)
	}
	return u.stream.Write([]byte(buf))
}

// Read считывает limit байт, пропуская при этом offset байт
func (u *ConvWriteReader) Read(offset int64, limit int64) ([]byte, error) {
	var (
		err  error
		data []byte
	)

	reader := bufio.NewReader(u.stream)
	// пропуск offset байт из начала ввода
	_, err = reader.Discard(int(offset))
	if err != nil {
		return data, errors.New("offset is bigger than file size." +
			" unable to read the file")
	}

	// по умолчанию limit равен -1, поэтому стоит такое условие
	// точно известно, что если limit > 0 изначально, то в какой-то момент либо закончится файл,
	// либо limit станет равен нулю и ввод завершится
	for limit != 0 {
		readedByte, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		if err != nil {
			return data, errors.New("unable to read byte")
		}
		data = append(data, readedByte)
		limit--
	}

	return data, err
}

// validateFlags проверяет все флаги на валидность
func validateFlags(options *Options) error {
	if options.From != "" {
		// проверка того, что файл с данными существует
		_, err := os.Stat(options.From)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("the file %s is not exist."+
				" unable to read the file", options.From)
		}

		// проверка того, что размер файла больше кол-ва игнорируемых байт
		f, _ := os.Stat(options.From)
		if f.Size() < options.Offset && options.Offset > 0 {
			return errors.New("offset is bigger than file size." +
				" unable to read the file")
		}
	}

	if options.To != "" {
		// проверка того, что не существует файла, в который будут записываться данные
		_, err := os.Stat(options.To)
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("the file on path %s is already exist."+
				" unable to write data", options.To)
		}
	}

	// проверка того, что все опции, переданные во флаг conv корректны
	validArgs := map[string]bool{"trim_spaces": true, "upper_case": true, "lower_case": true}
	for arg := range options.Conv {
		if !validArgs[arg] {
			return fmt.Errorf("conv arg <%s> is not correct", arg)
		}
	}

	// невозможно одновременно привести текст и к верхнему, и к нижнему регистру
	if options.Conv["upper_case"] && options.Conv["lower_case"] {
		return errors.New("unable to convert data in upper and lower case at the same time")
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

// ParseFlags парсит параметры и возвращает ошибку, если они не валидны
func ParseFlags() (*Options, error) {
	var opts Options
	opts.SetDefault()

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.Int64Var(&opts.Offset, "offset", 0, "count of bytes to skip. by default - 0")
	flag.Int64Var(&opts.Limit, "limit", -1, "max count of bytes to read. by default - -1")
	flag.Int64Var(&opts.BlockSize, "block-size", -1, "max count of bytes to read and write."+
		" by default - -1")
	conv := flag.String("conv", "", "conversions over text: upper_case, lower_case and trim_spaces")

	flag.Parse()

	for _, arg := range strings.Split(*conv, ",") {
		if arg != "" {
			opts.Conv[arg] = true
		}
	}

	err := validateFlags(&opts)

	return &opts, err
}

// OpenFile возвращает поток ввода/вывода типа io.ReadWriteCloser
// filepath - путь к файлу
// если readMode == true, возвращается поток ввода, иначе - поток вывода
func OpenFile(filePath string, readMode bool) (io.ReadWriteCloser, error) {
	var (
		stream io.ReadWriteCloser
		err    error
	)

	if readMode {
		if filePath == "" {
			stream = os.Stdin // по умолчанию считываем из stdin
		} else {
			stream, err = os.Open(filePath)
		}
	} else {
		if filePath == "" {
			stream = os.Stdout // по умолчанию выводим в stdout
		} else {
			stream, err = os.Create(filePath)
		}
	}

	return stream, err
}

func main() {
	// парсим флаги
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	// создаем поток ввода inputStream
	inputStream, err := OpenFile(opts.From, true)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "unable to read from input: ", err)
		os.Exit(1)
	}

	// data - массив данных из потока ввода
	var data []byte

	// считываем данные из потока ввода
	reader := NewConvWriteReader(inputStream)
	data, err = reader.Read(opts.Offset, opts.Limit)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error with reading", err)
		os.Exit(1)
	}

	// создаем поток вывода outputStream
	outputStream, err := OpenFile(opts.To, false)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "unable to write in output", err)
		os.Exit(1)
	}

	// записываем данные в поток вывода
	writer := NewConvWriteReader(outputStream)
	_, err = writer.Write(data, &opts.Conv)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error with writing in output:", err)
		os.Exit(1)
	}
}
