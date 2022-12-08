package main

// 1. import
import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 2. chaincode 구조체
// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// 3. goods 구조체
type Goods struct {
	Serial  string `json:"serial"`
	Info  	string `json:"info"`
	Owner  	string `json:"owner"`
	Status	string `json:"status"` // registered, svc_requested , svc_responded, svc_completed, svc_paid 
}

type Svc struct {
	Serial string
	SvcID string
	SvcRequestInfo string
	SvcInvoice string
	SvcCompleteInfo string
}

type HistoryQueryResult struct {
	Record    *Goods    `json:"record"`
	TxId     string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// goods_register ( serial, customer, purchase_info )
func (s *SmartContract) Goods_register(ctx contractapi.TransactionContextInterface, serial string, customer string, purchase_info string) error {
	goods := Goods{
		Serial: serial,
		Info: purchase_info,
		Owner: customer,
		Status: "registered",
	}

	goodsAsBytes, _ := json.Marshal(goods)

	return ctx.GetStub().PutState(serial, goodsAsBytes)
}

// goods_svc_request ( serial, svc_info )
func (s *SmartContract) Goods_svc_request (ctx contractapi.TransactionContextInterface, serial string, svc_info string) error {

	goods, err := s.Goods_query(ctx, serial)

	if err != nil {
		return err
	}
	// 검증  status == registered, paid
	if goods.Status == "registered" || goods.Status == "paid" {
		goods.Status = "requested"

		// (TODO) SVC Info State Generation

		goodsAsBytes, _ := json.Marshal(goods)
	
		return ctx.GetStub().PutState(serial, goodsAsBytes)
	} else {
		return fmt.Errorf("Goods is not in registered or paid: %s", serial)
	}
}

// goods_svc_respond ( serial, charge, invoice )

// goods_svc_complete ( serial, price, repairinfo, price)

// goods_svc_pay  ( serial )

// goods_query (serial )
func(s *SmartContract) Goods_query(ctx contractapi.TransactionContextInterface, serial string) (*Goods, error) {
	goodsAsBytes, err := ctx.GetStub().GetState(serial)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if goodsAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist.", serial)
	}

	goods := new(Goods)
	_ = json.Unmarshal(goodsAsBytes, goods)

	return goods, nil
}

// goods_history (serial )
func (t *SmartContract) Goods_history(ctx contractapi.TransactionContextInterface, serial string) ([]HistoryQueryResult, error) {
	log.Printf("Goods_history: ID %v", serial) // 체인코드 컨테이너 -> docker logs dev-asset1...

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(serial)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var goods Goods
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &goods)
			if err != nil {
				return nil, err
			}
		} else {
			goods = Goods{
				Serial: serial,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &goods,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// 5. main
func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create goods chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting goods chaincode: %s", err.Error())
	}
}