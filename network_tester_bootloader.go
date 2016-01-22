package main

import (
  "flag"
  "log"
  //"io"
  "os"

  "github.com/tarm/serial"
  //"github.com/Omegaice/go-xmodem/xmodem"
)

var errorLog, warningLog, infoLog, debugLog *log.Logger

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
    }

    if *device != "nil" && *firmware != "nil" {
        success = upload_firmware(*device, *firmware)
    } else { success = false }

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

    infoLog.Println("Opening", dev_path)

    c := &serial.Config { Name: dev_path, Baud: 115200 }

    s, err := serial.OpenPort(c)
    if err != nil {
        errorLog.Println(err)
    }

    n, err := s.Write([]byte("test"))
    if n!=2|| err != nil {
        errorLog.Println(err)
    }

    return true
}
