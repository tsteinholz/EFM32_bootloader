package main

import (
  "flag"
  "log"
  "io/ioutil"
  "os"

  "github.com/tarm/serial"
  "github.com/Omegaice/go-xmodem/xmodem"
)

var errorLog, warningLog, infoLog, debugLog *log.Logger
var debug bool

func main() {
    firmware := flag.String("firmware", "nil", "Path to the firmware.")
    device   := flag.String("device", "nil", "Path to the device on which we should commmunicate")
    verbose  := flag.Bool("verbose", false, "Whether to show verbose/debug log or not.")
    flag.Parse()

    errorLog   = log.New(os.Stdout, "ERROR: ",   log.Ltime)
    warningLog = log.New(os.Stdout, "WARNING: ", log.Ltime)
    infoLog    = log.New(os.Stdout, "INFO: ",    log.Ltime)

    success := true

    if *verbose {
        debugLog   = log.New(os.Stdout, "DEBUG: ",   log.Ltime)

        debugLog.Println("Arguments:")
        debugLog.Println("   Firmware Path:", *firmware)
        debugLog.Println("   Device Path:", *device)
        debugLog.Println("   Verbose:", *verbose)

        debug = true
    }

    if success && *device != "nil" && *firmware != "nil" {
        success = upload_firmware(*device, *firmware)
    } else {
        success = false
        errorLog.Println("Must have device and firmware arguments to run..")
    }

    if !success {
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

    config := &serial.Config { Name: dev_path, Baud: 115200 }

    port, err := serial.OpenPort(config)
    check(err)

    debug_log("Sending xmodem request to serial")
    _, err = port.Write([]byte("U"))
    check(err)
    check(port.Flush())

    _, err = port.Write([]byte("u"))
    check(err)
    check(port.Flush())

    debug_log("Done sending xmodem request to serial")

    infoLog.Println("Starting XMODEM transfer for", dev_path)
    err = xmodem.ModemSend(port, data)
    check(err)
    // TODO : add timeout
    return true
}

//------------------------------------------------------------------------------
// Purpose: To handle the debug log and only publish the log when the user
// specifies a need for it.
//
// Param text: The text that should output
//------------------------------------------------------------------------------
func debug_log(text string) {
    if (debug) { debugLog.Println(text) }
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
