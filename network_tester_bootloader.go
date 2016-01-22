package main

import (
  "flag"
  "log"
  //"io"
  "os"

  //"github.com/Omegaice/go-xmodem/xmodem"
)

var errorLog, warningLog, infoLog, debugLog *log.Logger

func main() {
    firmware := flag.String("firmware", "null", "Path to the firmware.")
    device := flag.String("device", "null", "Path to the device on which we should commmunicate")
    verbose := flag.Bool("verbose", false, "Whether to show verbose/debug log or not.")

    flag.Parse()

    errorLog   = log.New(os.Stdout, "ERROR: ",   log.Ltime)
    warningLog = log.New(os.Stdout, "WARNING: ", log.Ltime)
    infoLog    = log.New(os.Stdout, "INFO: ",    log.Ltime)
    debugLog   = log.New(os.Stdout, "DEBUG: ",   log.Ltime)

    debugLog.Println("Arguments:")
    debugLog.Println("   Firmware Path:", *firmware)
    debugLog.Println("   Device Path:", *device)
    debugLog.Println("   Verbose:", *verbose)
    // TODO : Loop through all the found devices and upload firmware
}

//------------------------------------------------------------------------------
// Purpose: Uploads the firmware to the devices via xmodem.
//
// Param dev_path: The device path to commmunicate on.
// Param firmware_path: The location on disk of the firmware that is to be
//                      installed.
//------------------------------------------------------------------------------
func upload_firmware(dev_path, firmware_path string) {
    // TODO : Implement fucntion
}
