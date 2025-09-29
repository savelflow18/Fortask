package main

//для удобства использования я работал не с путем к директории файлов а непосредственно с именами файлов
//их имена для себя можно будет поменять в строках 67 и 94 а также 119 и 124 где аргументы функции
//важное по изменениям
//исправил недочеты по типу лишний закоментенный код
//сгенерировал тэстеры по двум основным функциям
//теперь в длине файлов передаются не количество байт а реальное количество символов
//(последнее изменее было внесено в функцию info_size(34 строка)
//увеличил колличество информации для файла log.txt
//lol
import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

func Open(cxt context.Context, s string) (*os.File, error) {
	select {
	case <-cxt.Done():
		fmt.Println("the end")
		return nil, nil
	default:
		file, err := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		} else {
			return file, nil
		}
	}
}
func info_size(cxt context.Context, file *os.File) int {
	select {
	case <-cxt.Done():
		fmt.Println("the end")
		return 0
	default:
		data, err := os.ReadFile(file.Name())
		if err != nil {
			fmt.Println(" файл не найден", err)
			return -1
		} else {
			return utf8.RuneCountInString(string(data))

		}
	}
}
func main() {
	var file_info2 os.FileInfo
	var file_info1 os.FileInfo
	start := time.Now()
	var w sync.WaitGroup
	res_file2 := make(chan *os.File)
	var file2 *os.File
	res_file1 := make(chan *os.File)
	var file1 *os.File
	res_len1 := make(chan int)
	var len1 int
	res_len2 := make(chan int)
	var len2 int
	a := 2
	parent_cxt, cancel := context.WithCancel(context.Background())
	if a == 1 {
		cancel()
	} else {
		cancel = nil
	}
	if cancel != nil {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("cancel ,program is the end")
	}
	w.Add(1)
	go func() {
		defer w.Done()
		file_1, err := Open(parent_cxt, "file1.txt")
		if err != nil {
			fmt.Println(err)
		}
		res_file1 <- file_1
		defer file1.Close()
		close(res_file1)
	}()
	for v := range res_file1 {
		file1 = v
		file_inf1, err := file1.Stat()
		if err != nil {
			fmt.Println("ошибка при извлечении информации из первого файла", err)
		} else {
			file_info1 = file_inf1
		}
	}
	w.Add(1)
	go func() {
		defer w.Done()
		len := info_size(parent_cxt, file1)
		res_len1 <- len
		close(res_len1)
	}()
	for v := range res_len1 {
		len1 += v
	}
	if cancel == nil {
		fmt.Println("размер первого", len1)
	}
	w.Add(1)
	go func() {
		defer w.Done()
		file_2, err := Open(parent_cxt, "test.txt")
		if err != nil {
			fmt.Println(err)
		}
		res_file2 <- file_2
		defer file1.Close()
		close(res_file2)
	}()
	for v := range res_file2 {
		file2 = v
		file_inf2, err := file2.Stat()
		if err != nil {
			fmt.Println("ошибка при извлечении информации из второго файла", err)
		} else {
			file_info2 = file_inf2
		}
	}
	w.Add(1)
	go func() {
		defer w.Done()
		len2 := info_size(parent_cxt, file2)
		res_len2 <- len2
		close(res_len2)
	}()
	for v := range res_len2 {
		len2 += v
	}
	if cancel == nil {
		fmt.Println("размер второго", len2)
	}
	if len1 != len2 && cancel == nil {
		source_file, err := os.Open("file1.txt")
		if err != nil {
			fmt.Println(err)
		}
		defer source_file.Close()
		dest_file, err := os.Create("test.txt")
		if err != nil {
			fmt.Println(err)
		}
		defer dest_file.Close()
		fmt.Println("файлы синхронизированы")
		_, err = io.Copy(dest_file, source_file)
	} else {
		if len1 == len2 && cancel == nil {
			fmt.Println("файлы уже синхронизированы")
		} else {
			fmt.Println("cancel")
		}
	}
	end_time := (time.Since(start).Seconds())
	fmt.Println(end_time)
	w.Wait()
	log_f, err := os.Create("log.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		defer log_f.Close()
	}
	if cancel == nil {
		log_f.WriteString(fmt.Sprint("время выполнения:", end_time, "\n"))
		log_f.WriteString(fmt.Sprint("размер первого файла:", len1, "\n"))
		log_f.WriteString(fmt.Sprint("время последнего изменения первого файла:", file_info1.ModTime(), "\n"))
		log_f.WriteString(fmt.Sprint("Является ли директорией первый файл:: ", file_info1.IsDir(), "\n"))
		log_f.WriteString(fmt.Sprint("Права доступа первого файла:", file_info1.Mode(), "\n"))
		log_f.WriteString(fmt.Sprint("размер второго файла:", len2, "\n"))
		log_f.WriteString(fmt.Sprint("время последнего изменения второго файла:", file_info2.ModTime(), "\n"))
		log_f.WriteString(fmt.Sprint("Является ли директорией второй файл:: ", file_info2.IsDir(), "\n"))
		log_f.WriteString(fmt.Sprint("Права доступа второго файла:", file_info2.Mode(), "\n"))
	} else {
		fmt.Println("cancel", log_f.Name(), "cant work")
	}
}
