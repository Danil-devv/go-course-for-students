package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type ReadWriteSeekCloser interface {
	io.ReadWriter
	io.Seeker
	io.Closer
}

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
	opts.Limit = math.MaxInt64
	opts.BlockSize = 4096
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

	if options.Limit < 0 {
		return errors.New("the count of bytes to read must be positive")
	}

	if options.BlockSize <= 0 {
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
	flag.Int64Var(&opts.Limit, "limit", math.MaxInt64, "max count of bytes to read. by default - -1")
	flag.Int64Var(&opts.BlockSize, "block-size", 4, "max count of bytes to read and write."+
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

// OpenFile возвращает поток ввода/вывода типа ReadWriteSeekCloser
// filepath - путь к файлу
// если readMode == true, возвращается поток ввода, иначе - поток вывода
func OpenFile(filePath string, readMode bool) (ReadWriteSeekCloser, error) {
	var (
		stream ReadWriteSeekCloser
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

func ConvertData(data *[]byte, conv *map[string]bool) {
	buffer := string(*data)

	if (*conv)["upper_case"] {
		buffer = strings.ToUpper(buffer)
	}

	if (*conv)["lower_case"] {
		buffer = strings.ToLower(buffer)
	}

	*data = []byte(buffer)
}

func PipeData(inputStream ReadWriteSeekCloser, outputStream ReadWriteSeekCloser,
	offset int64, limit int64, conv map[string]bool, blockSize int64) error {
	var (
		err  error
		data []byte
	)

	// выполняем смещение потока ввода на offset байт
	data = make([]byte, 1)
	for offset > 0 {
		_, err := inputStream.Read(data)
		offset--
		if err != nil {
			return errors.New("offset is bigger than file size." +
				" unable to read the file")
		}
	}

	if !conv["trim_space"] {
		for limit > 0 {
			data = make([]byte, blockSize)

			size, err := inputStream.Read(data)

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			if int64(size) < blockSize {
				data = data[:size]
			}

			ConvertData(&data, &conv)

			_, err = outputStream.Write(data)

			if err != nil {
				return err
			}

			limit -= int64(size)
		}
	}

	return err
}

func main() {
	// парсим флаги
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "flags are not valid: ", err)
		os.Exit(1)
	}

	// создаем поток ввода inputStream
	inputStream, err := OpenFile(opts.From, true)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "unable to read from input: ", err)
		os.Exit(1)
	}

	// создаем поток вывода outputStream
	outputStream, err := OpenFile(opts.To, false)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "unable to write in output", err)
		os.Exit(1)
	}

	// выполняем чтение и запись информации
	err = PipeData(inputStream, outputStream, opts.Offset,
		opts.Limit, opts.Conv, opts.BlockSize)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "an error occurred while reading/writing data: ", err)
		os.Exit(1)
	}

	// закрытие и сохранение потока вывода и вывода
	defer func(in *ReadWriteSeekCloser, out *ReadWriteSeekCloser) {
		err := (*in).Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not close an input stream:", err)
			os.Exit(1)
		}

		err = (*out).Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not close an output stream:", err)
			os.Exit(1)
		}
	}(&inputStream, &outputStream)
}
