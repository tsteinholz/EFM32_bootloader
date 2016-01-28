package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Omegaice/go-xmodem/xmodem"
	"github.com/tarm/serial"
)

var errorLog, warningLog, infoLog, debugLog *log.Logger
var debug bool

func main() {
	firmware := flag.String("firmware", "nil", "Path to the firmware.")
	device := flag.String("device", "nil", "Path to the device on which we should commmunicate")
	verbose := flag.Bool("verbose", false, "Whether to show verbose/debug log or not.")
	flag.Parse()

	errorLog = log.New(os.Stdout, "ERROR: ", log.Ltime)
	warningLog = log.New(os.Stdout, "WARNING: ", log.Ltime)
	infoLog = log.New(os.Stdout, "INFO: ", log.Ltime)

	if *verbose {
		debugLog = log.New(os.Stdout, "DEBUG: ", log.Ltime)

		debugLog.Println("Arguments:")
		debugLog.Println("   Firmware Path:", *firmware)
		debugLog.Println("   Device Path:", *device)
		debugLog.Println("   Verbose:", *verbose)

		debug = true
	}

	if *device != "nil" && *firmware != "nil" {
		upload_firmware(*device, *firmware)
	} else {
		errorLog.Println("Must have device and firmware arguments to run..")
		infoLog.Println("Program Usage:")
		flag.PrintDefaults()
	}
}

//------------------------------------------------------------------------------
// Purpose: Uploads the firmware to the devices via xmodem.
//
// Param dev_path: The device path to commmunicate on.
// Param firmware_path: The location on disk of the firmware that is to be
//                      installed.
//------------------------------------------------------------------------------
func upload_firmware(dev_path, firmware_path string) bool {
	data, err := ioutil.ReadFile(firmware_path)
	check(err)

	infoLog.Println("Opening", dev_path)

	// TODO : Upload firmware to multiple devices in goroutines simultaneously

	config := &serial.Config{Name: dev_path, Baud: 115200, ReadTimeout: time.Second * 5}

	port, err := serial.OpenPort(config)
	check(err)

	debug_log("Sending xmodem request to serial")
	_, err = port.Write([]byte("U"))
	check(err)
	//debug_log("Verrifing write")
	verify_write(port)
	//read_buff := make([]byte, 10)
	//n, err := port.Read(read_buff)
	//check(err)
	//infoLog.Println(read_buff)
	//infoLog.Println(n)
	//check(port.Flush())

	_, err = port.Write([]byte("u"))
	//check(port.Flush())
	check(err)

	debug_log("Done sending xmodem request to serial")

	infoLog.Println("Starting XMODEM transfer for", dev_path)
	check(xmodem.ModemSend(port, data))
	// TODO : add timeout
	infoLog.Println("Testing for feedback")
	read_buff := make([]byte, 10)
	_, err = port.Read(read_buff)
	check(err)
	infoLog.Println(read_buff)
	return true
}

//------------------------------------------------------------------------------
// Purpose: To handle the debug log and only publish the log when the user
// specifies a need for it.
//
// Param text: The text that should output
//------------------------------------------------------------------------------
func debug_log(text string) {
	if debug {
		debugLog.Println(text)
	}
}

//------------------------------------------------------------------------------
// Purpose: Quick error checking to fix the excess amount of error checking we
// need to do
//
// Param err: The error we are checking
//------------------------------------------------------------------------------
func check(err error) {
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------
// Purpose: To make sure that the xmodem calls were read by the device without
//          clearing the data
//
// Param port: The serial port that we should utilize
//------------------------------------------------------------------------------
func verify_write(port *serial.Port) {
	var err error
	n, read_buff := 1, make([]byte, 5)
	for n > 0 && err != io.EOF {
		n, err = port.Read(read_buff)
	}
}
