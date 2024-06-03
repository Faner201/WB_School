package repository

import (
	"L0/internal/entity"
	"encoding/json"
	"math/rand"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/icrowley/fake"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
)

func (r *Repository) NatsSubcribe(ns stan.Conn) error {
	_, err := ns.Subscribe("orders", func(msg *stan.Msg) {
		order := entity.Order{}
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Err(err).Msg("it was not possible to convert the data from the message broker to the json standard")
			return
		}

		if err := checkValidDataNats(order); err != nil {
			log.Err(err).Msg("the data has not been validated")
			return
		}
		if _, ok := r.cache.GetOrderByUID(order.OrderUID); !ok {
			if ok := r.сheckOrderDB(&order); !ok {
				r.cache.SetOrder(order)
				if err := r.insertOrder(order); err != nil {
					log.Err(err).Msg("data could not be saved to the database")
					return
				}
			}
		}

	}, stan.DurableName("orders"))

	if err != nil {
		return err
	}

	return err
}

func (r *Repository) NatsGenerateDate(ns stan.Conn) error {
	file, err := checkExistenceFile("OrderUID.txt")
	if err != nil {
		return err
	}
	for i := 0; i < 50; i++ {
		order := createFakeOrder()
		if err := saverOrderUIDforFile(order.OrderUID, file); err != nil {
			return err
		}
		orderJson, err := json.Marshal(order)

		if err != nil {
			return err
		}

		if err := ns.Publish("orders", orderJson); err != nil {
			return err
		}
	}
	file.Close()
	return nil
}

func checkValidDataNats(data entity.Order) error {
	validate := validator.New()

	err := validate.Struct(data)
	if err != nil {
		return err
	}

	return nil
}

func createFakeOrder() entity.Order {
	item := []entity.Item{}
	for i := 0; i < 5; i++ {
		item = append(item, entity.Item{
			ChrtID:      rand.Intn(4),
			TrackNumber: fake.CharactersN(10),
			Price:       rand.Intn(6),
			Rid:         fake.DigitsN(10),
			Name:        fake.ProductName(),
			Sale:        rand.Intn(5),
			Size:        "0",
			TotalPrice:  rand.Intn(500),
			NmID:        rand.Intn(200),
			Brand:       fake.Brand(),
			Status:      202,
		})
	}

	return entity.Order{
		OrderUID:    fake.DigitsN(19),
		TrackNumber: fake.CharactersN(14),
		Entry:       fake.CharactersN(4),
		Delivery: entity.Delivery{
			Name:    fake.FemaleFullName(),
			Phone:   fake.Phone(),
			Zip:     fake.DigitsN(7),
			City:    fake.City(),
			Address: fake.StreetAddress(),
			Region:  fake.CharactersN(10),
			Email:   fake.EmailAddress(),
		},
		Payment: entity.Payment{
			Transaction:  fake.DigitsN(19),
			RequestID:    fake.CharactersN(5),
			Currency:     fake.CharactersN(5),
			Provider:     fake.Brand(),
			Amount:       rand.Intn(5),
			PaymentDT:    int(time.Now().Unix()),
			Bank:         fake.Brand(),
			DeliveryCost: rand.Intn(5000),
			GoodsTotal:   rand.Intn(400),
			CustomFee:    0,
		},
		Items:             item,
		Locale:            "en",
		InternalSignature: fake.CharactersN(5),
		CustomerID:        fake.CharactersN(6),
		DeliveryService:   "meest",
		Shardkey:          fake.CharactersN(3),
		SmID:              rand.Intn(6),
		DateCreated:       time.Now(),
		OofShard:          fake.CharactersN(6),
	}

}

func checkExistenceFile(name string) (*os.File, error) {
	f := new(os.File)

	fs, err := os.Stat(name)
	if err != nil {
		return nil, err
	}

	if !fs.IsDir() {
		f, err = os.Create(name)
		if err != nil {
			return nil, err
		}
	} else {
		if err := os.Truncate(name, 0); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func saverOrderUIDforFile(orderUID string, file *os.File) error {
	orderUID = orderUID + "\n"

	if _, err := file.Write([]byte(orderUID)); err != nil {
		return err
	}
	return nil
}
