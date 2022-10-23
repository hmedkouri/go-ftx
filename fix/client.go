package fix

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

// TradeClient implements the quickfix.Application interface
type TradeClient struct {
	ApiSecret  string
	SubAccount string
}

// OnCreate implemented as part of Application interface
func (ts TradeClient) OnCreate(sessionID quickfix.SessionID) {
	fmt.Printf("OnCreate %s\n", sessionID)
}

// OnLogon implemented as part of Application interface
func (ts TradeClient) OnLogon(sessionID quickfix.SessionID) {
	fmt.Printf("OnLogon %s\n", sessionID)
}

// OnLogout implemented as part of Application interface
func (ts TradeClient) OnLogout(sessionID quickfix.SessionID) {
	fmt.Printf("OnLogout %s\n", sessionID)
}

// FromAdmin implemented as part of Application interface
func (ts TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromAdmin %v %v\n", msg, sessionID)
	return nil
}

// ToAdmin implemented as part of Application interface
func (ts TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	fmt.Printf("ToAdmin %v %v\n", msg, sessionID)
	isLogonMessage := msg.IsMsgTypeOf("A")
	if isLogonMessage {
		ts.handleLogonMessage(msg)
		return
	}
}

func (ts TradeClient) handleLogonMessage(msg *quickfix.Message) {
	msg.Body.SetInt(98, 0)
	msg.Body.SetInt(108, 30)
	msg.Body.SetString(96, ts.getLogonRawData(msg))
	msg.Body.SetString(8013, "S")
	msg.Body.SetString(1, ts.SubAccount)
}

func (ts TradeClient) getLogonRawData(msg *quickfix.Message) string {
	sendingtime, _ := msg.Header.GetString(tag.SendingTime)
	msgType := "A"
	msgSeqNum, _ := msg.Header.GetString(tag.MsgSeqNum)
	senderCompID, _ := msg.Header.GetString(tag.SenderCompID)
	targetCompID, _ := msg.Header.GetString(tag.TargetCompID)

	presign := strings.Join(
		[]string{
			sendingtime,
			msgType,
			msgSeqNum,
			senderCompID,
			targetCompID,
		},
		string("\x01"),
	)

	encoded := hmacEncrypt(sha256.New, presign, ts.ApiSecret)

	return encoded
}

func hmacEncrypt(pfn func() hash.Hash, data, key string) string {
	h := hmac.New(pfn, []byte(key))
	if _, err := h.Write([]byte(data)); err == nil {
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
}

// ToApp implemented as part of Application interface
func (ts TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (tsrr error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

// FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (ts TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}
