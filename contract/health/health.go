/*
SPDX-License-Identifier: Apache-2.0
*/

// 1. import
package main

import (
	"encoding/json"
	"fmt"
	//"strconv"
	"time"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 2. chaincode 援ъ“泥�
// SmartContract provides functions for managing health
type SmartContract struct {
	contractapi.Contract
}

// 3. health, query result 援ъ“泥�
// health describes basic details of what makes up health
type Health struct {
	Number   string `json:"number"`
	Name  string `json:"name"`
	DateOfBirth string `json:"dateofbirth"`
	KCD  string `json:"KCD"`
	Doctor  string `json:"doctor"`
	DateOfDiagnosis  string `json:"dateofdiganosis"`
}

// 4.1 health_register (number, name, dateofbirth, kcd, doctor, dateofdiagnosis)
func (s *SmartContract) Health_register(ctx contractapi.TransactionContextInterface, number string, name string, dateofbirth string, kcd string, doctor string, dateofdiagnosis string) error {
	health := Health{
		Number: number,
		Name: name,
		DateOfBirth: dateofbirth,
		KCD: kcd,
		Doctor: doctor,
		DateOfDiagnosis: dateofdiagnosis,
	}

	healthAsBytes, _ := json.Marshal(health)

	return ctx.GetStub().PutState(number, healthAsBytes)
}

// 4.2 health_query (number)
func(s *SmartContract) Health_query(ctx contractapi.TransactionContextInterface, number string) (*Health, error) {
	healthAsBytes, err := ctx.GetStub().GetState(number) // state -> JSON format

	if err != nil { // GetState, GetStub, ctx李몄“ �ㅻ쪟 留뚮궗�� ��
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if healthAsBytes == nil { // key媛� ���λ맂�곸씠 �녾굅�� delete�� 寃쎌슦
		return nil, fmt.Errorf("%s does not exist", number)
	}

	// �뺤긽�곸쑝濡� 議고쉶�� 寃쎌슦
	health := new(Health) // 媛앹껜�� SON -> 援ъ“泥�
	_ = json.Unmarshal(healthAsBytes, health) // call by Reference

	return health, nil
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Health    `json:"record"`
	TxId     string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// health_history (number)
func (t *SmartContract) Health_history(ctx contractapi.TransactionContextInterface, number string) ([]HistoryQueryResult, error) {
	log.Printf("health_history: number %v", number) // 泥댁씤肄붾뱶 而⑦뀒�대꼫 -> docker logs dev-asset1...

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(number)
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

		var health Health
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &health)
			if err != nil {
				return nil, err
			}
		} else {
			health = Health{
				Number: number,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}
		
		record := HistoryQueryResult{
			TxId:     response.TxId,
			Timestamp: timestamp,
			Record:    &health,
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
		fmt.Printf("Error create health chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting health chaincode: %s", err.Error())
	}
}

/* ex
type Svc struct {
	Serial string
	SvcID stirng
	SvcRequestInfo string
	SvcInvoice string
	SvcCompleteInfo string
}
*/

// // QueryAllHealth
// func (s *SmartContract) QueryAllHealth(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
// 	startKey := ""
// 	endKey := ""

// 	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	results := []QueryResult{}

// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()

// 		if err != nil {
// 			return nil, err
// 		}

// 		car := new(Car)
// 		_ = json.Unmarshal(queryResponse.Value, car)

// 		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
// 		results = append(results, queryResult)
// 	}

// 	return results, nil
// }


// // ChangeDoctor
// func (s *SmartContract) ChangeDoctor(ctx contractapi.TransactionContextInterface, number string, newDoctor string) error {

// 	health, err := s.Health_query(ctx, number)

// 	if err != nil {
// 		return err
// 	}
// 	// 寃�利� status == registered, paid
// 	// �덉떆
// 	///*
// 	//if goods.Status == "registered" || goods.Status == "paid" {
// 	//	goods.Status = "requested"
//     //
// 	//	// (TODO) SVC Info State Generation
//     //
// 	//	goodsAsBytes, _ := json.Marshal(goods)
// 	//	return ctx.GetStubd().PutState(serial, goodsAsBytes)
// 	//} else {
// 	//	return fmt.Errorf("Goods is not in registered or paid: %s", serial)
// 	//}
// 	//*/


// 	// health.Doctor = newDoctor
// 	health.Status = "requested"
// 	healthAsBytes, _ := json.Marshal(health) // 吏곷젹�� (援ъ“泥� -> JSON �щĸ�� Byte[])

// 	return ctx.GetStub().PutState(number, healthAsBytes)

// }





/*
// 4.1 InitLedger

// 4.2 CreateCar

// 4.3 QueryCar

// 4.4. ChangeCarOwner

// 4.5 QueryAllCars

// 5. main

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}

*/