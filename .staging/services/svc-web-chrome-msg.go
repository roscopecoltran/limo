package service

// ref. https://github.com/sauyon/go-chromemessage/blob/master/chromemsg/chromemsg.go

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"unsafe"
	"github.com/sirupsen/logrus"
)

var nativeEndian binary.ByteOrder 	= 	endianness()
var defaultMsgr 					= 	MessengerChromeNativeMsg{bufio.NewReadWriter(
											bufio.NewReader(os.Stdin),
											bufio.NewWriter(os.Stdout))}

type MessengerChromeNativeMsg struct {
	port *bufio.ReadWriter
}

func NewChromeNativeMsg(port *bufio.ReadWriter) *MessengerChromeNativeMsg {
	return &MessengerChromeNativeMsg{port}
}

func ReadChromeNativeMsg(data interface{}) error {
	return defaultMsgr.ReadChromeNativeMsg(data)
}

func WriteChromeNativeMsg(msg interface{}) error {
	return defaultMsgr.WriteChromeNativeMsg(msg)
}

func (msgr *MessengerChromeNativeMsg) ReadChromeNativeMsg(data interface{}) error {
	lengthBits := make([]byte, 4)
	_, err := msgr.port.ReadChromeNativeMsg(lengthBits)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"service/svc-web-chrome.go", 
							"method_name": 	"ReadChromeNativeMsg(data interface{})", 
							"driver": 		"chrome", 
		                    "feature":      "native-messages",
							"action": 		"msgr.port.Read(lengthBits)",
							}).Warn("error while trying to connect to chrome.")
		return err
	}
	length := nativeToInt(lengthBits)
	content := make([]byte, length)
	_, err = msgr.port.Read(content)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"service/svc-web-chrome.go", 
							"method_name": 	"ReadChromeNativeMsg(data interface{})", 
							"driver": 		"chrome", 
		                    "feature":      "native-messages",
							"action": 		"msgr.port.Read(content)",
							}).Warn("error while reading native message sent by chrome browser.")
		return err
	}
	json.Unmarshal(content, data)
	return nil
}

func (msgr *MessengerChromeNativeMsg) WriteChromeNativeMsg(msg interface{}) error {
	json, err := json.Marshal(msg)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"service/svc-web-chrome.go", 
							"method_name": 	"WriteChromeNativeMsg(msg interface{})", 
							"driver":		"chrome", 
		                    "feature":      "native-messages",
							"action": 		"json.Marshal(msg)",
							}).Warn("error while writing native message to the chrome browser.")
		return err
	}
	length := len(json)
	bits := make([]byte, 4)
	buf := bytes.NewBuffer(bits)
	err = binary.Write(buf, nativeEndian, length)
	if err != nil {
		return err
	}
	_, err = msgr.port.WriteChromeNativeMsg(bits)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"service/svc-web-chrome.go", 
							"method_name": 	"WriteChromeNativeMsg(msg interface{})", 
							"driver": 		"chrome", 
		                    "feature":      "native-messages",
							"action": 		"msgr.port.Write(bits)",
							}).Warn("error while writing native message to the chrome browser.")
		return err
	}
	_, err = msgr.port.WriteChromeNativeMsg(json)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": 		"service/svc-web-chrome.go", 
							"method_name": 	"WriteChromeNativeMsg(msg interface{})", 
							"driver": 		"chrome", 
		                    "feature":      "native-messages",
							"action": 		"msgr.port.Write(json)",
							}).Warn("error while writing native message to the chrome browser.")
		return err
	}
	return nil
}

func nativeToInt(bits []byte) int {
	var length uint32
	buf := bytes.NewBuffer(bits)
	binary.Read(buf, nativeEndian, &length)
	return int(length)
}

func endianness() binary.ByteOrder {
	var i int = 1
	bs := (*[unsafe.Sizeof(0)]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}

