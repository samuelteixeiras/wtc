// Example WhatsApp client using the whatsmeow library with command-line arguments
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"github.com/mdp/qrterminal/v3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow/proto/waE2E"
)

// Global variables
var client *whatsmeow.Client
var log waLog.Logger
var messageSent bool
var qrScanned bool

// eventHandler handles incoming events from WhatsApp silently
func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		// When we get a QR code, print it to the terminal (necessary for login)
		fmt.Println("QR code received, scan it with your phone!")
		qrterminal.GenerateHalfBlock(v.Codes[0], qrterminal.L, os.Stdout)
		
	case *events.PairSuccess:
		qrScanned = true
	}
}

// sendMessage sends a text message to the specified JID
func sendMessage(recipient string, message string) error {
	// Parse the JID (phone number)
	jid, err := types.ParseJID(recipient)
	if err != nil {
		return fmt.Errorf("invalid JID: %v", err)
	}
	
	// Create the message content
	msg := &waE2E.Message{
		Conversation: proto.String("[wtc]" + message),
	}
	
	// Send the message
	_, err = client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	
	fmt.Printf("Message sent to %s\n", recipient)
	messageSent = true
	return nil
}

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./wtc <phone_number> <message>")
		fmt.Println("Example: ./wtc 1234567890 'Hello, world!'")
		os.Exit(1)
	}
	
	// Get phone number and message from command line arguments
	recipient := os.Args[1]
	message := strings.Join(os.Args[2:], " ")
	
	// Format the phone number as a JID if needed
	// Automatically append @s.whatsapp.net if the recipient doesn't already contain an @ symbol
	if !strings.Contains(recipient, "@") {
		recipient = recipient + "@s.whatsapp.net"
	}
	
	// Set up logging with minimal output
	log = waLog.Stdout("WTS", "ERROR", false)
	
	// Create a database connection for storing session data
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:whatsmeow.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	
	// Get device store - we use the first device in the database
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		os.Exit(1)
	}
	
	// Create the client
	client = whatsmeow.NewClient(deviceStore, log)
	
	// Set the event handler
	client.AddEventHandler(eventHandler)
	
	// Connect to WhatsApp
	if err := client.Connect(); err != nil {
		log.Errorf("Failed to connect: %v", err)
		os.Exit(1)
	}
	
	// Check if we're logged in
	if client.Store.ID == nil {
		// We're not logged in, so we need to wait for the QR code event
		log.Infof("Not logged in. Waiting for QR code scan...")
		
		// Wait for QR code to be scanned
		timeout := time.After(2 * time.Minute)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		for !qrScanned {
			select {
			case <-ticker.C:
				// Check every second if QR was scanned
			case <-timeout:
				log.Errorf("Timeout waiting for QR code scan")
				os.Exit(1)
			}
		}
	}
	
	// Send the message
	err = sendMessage(recipient, message)
	if err != nil {
		log.Errorf("Error sending message: %v", err)
		os.Exit(1)
	}
	
	// Wait a moment to ensure message is sent
	time.Sleep(3 * time.Second)
	
	// Disconnect and exit
	client.Disconnect()
}