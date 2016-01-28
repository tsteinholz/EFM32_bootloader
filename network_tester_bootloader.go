package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chrizzzzz/go-xmodem/xmodem"
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
		uploadFirmware(*device, *firmware)
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
//------------------------------------------------------------------------------
//                      installed.
func uploadFirmware(devPath, firmwarePath string) bool {
	logDebug("Reading binary file")
	data, err := ioutil.ReadFile(firmwarePath)

	check(err)
	infoLog.Println("Opening", devPath)

	// TODO : Upload firmware to multiple devices in goroutines simultaneously

	// TODO : Come to a final decision on baud rate
	config := &serial.Config{Name: devPath, Baud: 115200, ReadTimeout: time.Second * 5}

	logDebug("Opening serial port")
	port, err := serial.OpenPort(config)
	check(err)

	logDebug("Sending xmodem request to serial")
	_, err = port.Write([]byte("U"))
	check(err)
	verifyWrite(port)

	_, err = port.Write([]byte("u"))
	check(err)
	verifyWrite(port)

	logDebug("Done sending xmodem request to serial")

	startTime := time.Now()
	infoLog.Println("Starting XMODEM transfer for", devPath)
	check(xmodem.ModemSend(port, data))
	infoLog.Println("Finished XMODEM transfer for", devPath, "in", time.Since(startTime), "seconds")
	return true
}

//------------------------------------------------------------------------------
// Purpose: To handle the debug log and only publish the log when the user
// specifies a need for it.
//
// Param text: The text that should output
//------------------------------------------------------------------------------
func logDebug(text string) {
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
func verifyWrite(port *serial.Port) {
	var err error
	n, readBuff := 1, make([]byte, 5)
	for n > 0 && err != io.EOF {
		n, err = port.Read(readBuff)
	}
}
