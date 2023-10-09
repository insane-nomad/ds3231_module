package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/beevik/ntp"
	"github.com/jacobsa/go-serial/serial"
)

var command string
var mod uint

type safeWriter struct {
	w   io.Writer
	err error // Место для хранения первой ошибки
}

func (sw *safeWriter) writeln(s string) {
	if sw.err != nil {
		return // Пропускает запись, если раньше была ошибка
	}
	_, sw.err = fmt.Fprint(sw.w, s) // Записывает строку и затем хранить любую ошибку
}

func saveFile(name, data string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	sw := safeWriter{w: f}
	sw.writeln(data)
	return sw.err // Возвращает ошибку в случае ее возникновения
}

var options serial.OpenOptions = serial.OpenOptions{
	PortName:        "COM6",
	BaudRate:        115200,
	DataBits:        8,
	StopBits:        1,
	MinimumReadSize: 4,
}

func init() {
	flag.StringVar(&command, "command", "", "executable command.\n\t\"set\" - to set time to module\n\t\"compare\" - to compare time and calculate ppm")
	flag.UintVar(&mod, "mod", 0, "module number. Required value\n\nfor example\n\tmain.exe -command compare -mod 1\n\tmain.exe -command set -mod 1")
}

func connectToNTP() (timeToString, timeToFile string) {
	NTPTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		fmt.Println(err)
		return
	}
	timeToString = fmt.Sprint(NTPTime.Format("$2006 01 02 15 04 05;\n"))
	timeToFile = fmt.Sprint(NTPTime.Format("2006-01-02 15:04:05"))
	return
}

func readSerialData() string {

	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	for {
		buf := make([]byte, 50)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]

			if len(buf) > 10 {
				buf = buf[:n-2]
				return string(buf)
			}
		}
	}
}

// func writeSerialData(data []byte) {
// 	options := serial.OpenOptions{
// 		PortName:        "COM6",
// 		BaudRate:        115200,
// 		DataBits:        8,
// 		StopBits:        1,
// 		MinimumReadSize: 4,
// 	}
// 	port, err := serial.Open(options)
// 	if err != nil {
// 		log.Fatalf("serial.Open: %v", err)
// 	}
// 	defer port.Close()
// 	b := []byte{0x00, 0x01, 0x02, 0x03}
// 	n, err := port.Write(data)
// 	if err != nil {
// 		log.Fatalf("port.Write: %v: %v", err, n)
// 	}
// }

func writeTimeToModule(filename string) {

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	for {
		buf := make([]byte, 100)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			//fmt.Print(string(buf))
			timeToString, timeToFile := connectToNTP()

			n, err := port.Write([]byte(timeToString))
			if err != nil {
				log.Fatalf("port.Write: %v: %v", err, n)
			}
			if len(buf) != 0 && buf[len(buf)-1] == '@' {
				saveFile := saveFile(filename, timeToFile)
				if saveFile != nil {
					//fmt.Println(saveFile)
					os.Exit(1)
				}
				break
			}
			//time.Sleep(250 * time.Millisecond)
		}
	}
}

func main() {
	flag.Parse()

	_, err := os.Stat("modules")
	if os.IsNotExist(err) {
		if err := os.Mkdir("modules", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	filename := "modules/mod" + strconv.Itoa(int(mod)) + ".txt"
	switch command {

	case "set":
		writeTimeToModule(filename)
		fmt.Println("\n\n+-----------------------------------------------------+")
		fmt.Println("|\t\t     Time is setted\t\t      |")
		fmt.Println("+-----------------------------------------------------+")
		fmt.Println("")
	case "compare":
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}

		moduleTime, _ := time.Parse("2006-01-02 15:04:05", readSerialData())
		_, timeFromNTP := connectToNTP()
		NTPTime, _ := time.Parse("2006-01-02 15:04:05", timeFromNTP)

		fileTime, _ := time.Parse("2006-01-02 15:04:05", string(fileContents))
		diffNTPFromFile := NTPTime.Sub(fileTime).Seconds()
		diffNTPFromModule := NTPTime.Sub(moduleTime).Seconds()
		accuracy := (diffNTPFromModule / diffNTPFromFile) * 1_000_000

		fmt.Println("\n\n+-----------------------------------------------------+")
		fmt.Println("|\t\t   Time is compared\t\t      |")
		fmt.Println("+-----------------------------------------------------+")
		fmt.Printf("| Time from NTP\t\t%v |\n| Time from file\t%v |\n| Time from module\t%v |\n| Sec from file time\t%-25v sec |\n| Accuracy\t\t%-25v ppm |\n",
			NTPTime, fileTime, moduleTime, diffNTPFromFile, accuracy)
		fmt.Println("+-----------------------------------------------------+")
		fmt.Println("")
	default:
		fmt.Println("\n\n+-----------------------------------------------------+")
		fmt.Println("|\t\t   incorrect command\t\t      |")
		fmt.Println("+-----------------------------------------------------+")
		fmt.Println("")
	}
}
