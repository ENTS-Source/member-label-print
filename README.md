# member-label-print

Prints a label with the member's name on it when they scan their fob. Backed by [aMember Pro](https://www.amember.com/).

## Hardware

* Raspberry Pi 4
* [5 inch touchscreen](https://www.amazon.ca/dp/B0B455LDKH)
* [USB RFID reader](https://www.amazon.ca/dp/B0CDBJKMVT)
* [Receipt printer](https://www.amazon.ca/dp/B0BS2CJN1H)

## Setup

1. Enable the "Rest API Module" addon in aMember Pro
2. Create a new user field in aMember Pro with the following (minimum) details:
   * Field Name: fob
   * Field Type: SQL
   * SQL field type: Text (string data)
   * Display Type: text
3. Create a new "Remote API key" with the following permissions:
   * Users: get

## Installation

Set your environment variables of interest:

```bash
TODO
```

Then run `./app.exe` (or whatever the executable's name is). Downloads are available from the github releases.
