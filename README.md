# WhatsApp Terminal Client (wtc)

A standalone command-line WhatsApp client that allows sending messages from the terminal without requiring Go installation on the target machine.


This is a simple client for the WhatsApp Web API using the [whatsmeow](https://github.com/tulir/whatsmeow) library.

## Overview

This package contains two main components:

1. `whatsapp_client` - The compiled binary containing the WhatsApp client implementation
2. `wtc` - A shell script wrapper that provides a simple interface to the binary

## Installation

1. Copy both the `whatsapp_client` binary and the `wtc` script to the target machine
2. Make sure both files have executable permissions:
   ```
   chmod +x wtc whatsapp_client
   ```
3. Place both files in the same directory

## Usage

Send a WhatsApp message using the following syntax:

```
./wtc <phone_number> <message>
```

Examples:

```
# Using phone number with country code
./wtc 1234567890 'Hello, world!'

# Using full JID format
./wtc 1234567890@s.whatsapp.net 'Hello, world!'
```

## First-time Setup

When running for the first time on a new machine:

1. The application will display a QR code in the terminal
2. Scan this QR code with your WhatsApp mobile app:
   - Open WhatsApp on your phone
   - Go to Settings > Linked Devices
   - Tap on "Link a Device"
   - Scan the QR code displayed in your terminal
3. After successful pairing, the message will be sent

## Notes

- The application creates a `whatsmeow.db` file in the current directory to store session data
- Once paired, you won't need to scan the QR code again on the same machine
- The binary is compiled for macOS arm64 architecture

## Dependencies

The compiled binary includes all necessary dependencies. No additional installation is required on the target machine.

Source code is main.go

To rebuild client:
```go
go build -o whatsapp_client -ldflags="-s -w" main.go
```