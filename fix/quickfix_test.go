package fix

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"

	fix42nos "github.com/quickfixgo/fix42/newordersingle"

	"github.com/shopspring/decimal"
)

func Test_QuickFix(t *testing.T) {

	cfgFileName := path.Join("config", "tradeclient.cfg")
	cfg, err := os.Open(cfgFileName)
	if err != nil {
		t.Errorf("Error opening %v, %v\n", cfgFileName, err)
	}
	defer cfg.Close()

	stringData, readErr := ioutil.ReadAll(cfg)
	if readErr != nil {
		t.Errorf("Error reading cfg: %s,", readErr)
	}

	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if err != nil {
		t.Errorf("Error reading cfg: %s,", err)
	}

	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		t.Errorf("Error creating file log factory: %s,", err)
	}

	app := TradeClient{
		ApiSecret:  os.Getenv("FTX_API_SECRET"),
		SubAccount:  os.Getenv("FTX_SUB_ACCOUNT"),
	}
	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		t.Errorf("Unable to create Initiator: %s\n", err)
	}

	err = initiator.Start()
	if err != nil {
		t.Errorf("Unable to start Initiator: %s\n", err)
	}

	order := queryNewOrderSingle42()

	err = quickfix.Send(order)
	if err != nil {
		t.Errorf("Unable to create Initiator: %s\n", err)
	}

	time.Sleep(10 * time.Second)
}

func queryNewOrderSingle42() (msg *quickfix.Message) {
	var ordType field.OrdTypeField
	order := fix42nos.New(field.NewClOrdID("ClOrdID"), field.NewHandlInst("1"), querySymbol(), querySide(), field.NewTransactTime(time.Now()), queryOrdType(&ordType))
	order.Set(queryOrderQty())

	switch ordType.Value() {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		order.Set(queryPrice())
	}

	switch ordType.Value() {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		order.Set(queryStopPx())
	}

	order.Set(queryTimeInForce())
	msg = order.ToMessage()
	queryHeader(&msg.Header)
	return
}

func queryOrdType(f *field.OrdTypeField) field.OrdTypeField {
	f.FIXString = quickfix.FIXString(enum.OrdType_MARKET)
	return *f
}

func queryTimeInForce() field.TimeInForceField {
	return field.NewTimeInForce(enum.TimeInForce_GOOD_TILL_CANCEL)
}

func queryOrderQty() field.OrderQtyField {
	return field.NewOrderQty(decimal.NewFromFloat(1), 2)
}

func queryPrice() field.PriceField {
	return field.NewPrice(decimal.NewFromFloat(1), 2)
}

func queryStopPx() field.StopPxField {
	return field.NewStopPx(decimal.NewFromFloat(1), 2)
}

func queryClOrdID() field.ClOrdIDField {
	return field.NewClOrdID("ClOrdID")
}

func queryOrigClOrdID() field.OrigClOrdIDField {
	return field.NewOrigClOrdID("OrigClOrdID")
}

func querySymbol() field.SymbolField {
	return field.NewSymbol("Symbol")
}

func querySide() field.SideField {
	return field.NewSide(enum.Side_BUY)
}

type header interface {
	Set(f quickfix.FieldWriter) *quickfix.FieldMap
}

func queryHeader(h header) {
	h.Set(querySenderCompID())
	h.Set(queryTargetCompID())
}

func querySenderCompID() field.SenderCompIDField {
	return field.NewSenderCompID("MDK")
}

func queryTargetCompID() field.TargetCompIDField {
	return field.NewTargetCompID("FTX")
}
