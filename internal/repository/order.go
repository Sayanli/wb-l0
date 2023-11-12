package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"wb-l0/internal/cache"
	"wb-l0/internal/models"
	"wb-l0/pkg/database"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type OrderRepository struct {
	db    database.PGXQuerier
	cache *cache.Cache
}

func NewOrderRepository(db database.PGXQuerier, cache *cache.Cache) *OrderRepository {
	return &OrderRepository{
		db:    db,
		cache: cache,
	}
}

func (r *OrderRepository) FindByUid(ctx context.Context, uid string) (models.Order, error) {
	if ord, b := r.cache.Get(uid); b {
		return ord, nil
	}

	rowsI, err := r.db.Query(context.Background(),
		`
		SELECT
			i.chrt_id, i.track_number, i.price,
			i.rid, i."name", i.sale, i.size,
			i.total_price, i.nm_id, i.brand, i."status"
		FROM items AS i
		WHERE order_id = $1
		`, uid,
	)
	if err != nil {
		logrus.Error(err)
		return models.Order{}, err
	}
	defer rowsI.Close()

	itms := make([]models.Item, 0)
	for rowsI.Next() {
		var itm models.Item
		err := rowsI.Scan(
			&itm.Chrt_id, &itm.Track_number,
			&itm.Price, &itm.Rid, &itm.Name, &itm.Sale,
			&itm.Size, &itm.Total_price, &itm.Nm_id,
			&itm.Brand, &itm.Status,
		)
		if err != nil {
			logrus.Error(err)
			return models.Order{}, err
		}
		itms = append(itms, itm)
	}

	var ord models.Order
	err = r.db.QueryRow(context.Background(),
		"SELECT * FROM orders WHERE order_uid = $1", uid,
	).Scan(
		&ord.Order_uid, &ord.Track_number, &ord.Entry,
		&ord.Delivery, &ord.Payment, &ord.Locale,
		&ord.Internal_signature, &ord.Customer_id,
		&ord.Delivery_service, &ord.Shard_key,
		&ord.Sm_id, &ord.Date_created, &ord.Oof_shard,
	)
	if err != nil {
		logrus.Info(err)
		return models.Order{}, err
	}
	ord.Items = itms

	return ord, nil
}

func (r *OrderRepository) getOrdersCount() (int, error) {
	var count int
	if err := r.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM orders").Scan(&count); err != nil {
		logrus.Error(err)
		return 0, err
	}
	return count, nil
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]models.Order, error) {
	rowsO, err := r.db.Query(context.Background(), "SELECT * FROM orders")
	if err != nil {
		logrus.Error(err)
		return []models.Order{}, err
	}
	defer rowsO.Close()

	countOrds, err := r.getOrdersCount()
	if err != nil {
		logrus.Info("There are no orders.")
		return []models.Order{}, err
	}
	ords := make([]models.Order, 0, countOrds)
	for rowsO.Next() {
		var ord models.Order
		err = rowsO.Scan(
			&ord.Order_uid, &ord.Track_number, &ord.Entry,
			&ord.Delivery, &ord.Payment, &ord.Locale,
			&ord.Internal_signature, &ord.Customer_id,
			&ord.Delivery_service, &ord.Shard_key,
			&ord.Sm_id, &ord.Date_created, &ord.Oof_shard,
		)
		if err != nil {
			logrus.Error(err)
			return []models.Order{}, err
		}

		rowsI, err := r.db.Query(context.Background(),
			`
			SELECT
					i.chrt_id, i.track_number, i.price,
					i.rid, i."name", i.sale, i.size,
					i.total_price, i.nm_id, i.brand, i."status"
				FROM items AS i
				WHERE order_id = $1
			`, ord.Order_uid,
		)
		if err != nil {
			logrus.Error(err)
			return []models.Order{}, err
		}
		defer rowsI.Close()

		itms := make([]models.Item, 0)
		for rowsI.Next() {
			var itm models.Item
			err := rowsI.Scan(
				&itm.Chrt_id, &itm.Track_number,
				&itm.Price, &itm.Rid, &itm.Name, &itm.Sale,
				&itm.Size, &itm.Total_price, &itm.Nm_id,
				&itm.Brand, &itm.Status,
			)
			if err != nil {
				logrus.Error(err)
				return []models.Order{}, err
			}
			itms = append(itms, itm)
		}
		ord.Items = itms
		ords = append(ords, ord)
	}
	return ords, nil
}

func (r *OrderRepository) CreateOrder(msg *stan.Msg) {
	fmt.Println("CreateOrder")
	ord, err := r.parseJsonToModel(msg)
	if err != nil {
		return
	}

	if ord.Validator() != true {
		fmt.Println("Invalid validator order")
		return
	}

	r.cache.Set(ord.Order_uid, ord)

	jsonDelivery, _ := json.Marshal(ord.Delivery)
	jsonPayment, _ := json.Marshal(ord.Payment)

	_, err = r.db.Exec(context.Background(),
		`
		INSERT INTO orders (
			order_uid, track_number, entry,
			delivery, payment, locale,
			internal_signature, customer_id,
			delivery_service, shardkey, sm_id,
			date_created, off_shard
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13
		)
		`,
		ord.Order_uid, ord.Track_number, ord.Entry,
		jsonDelivery, jsonPayment, ord.Locale,
		ord.Internal_signature, ord.Customer_id,
		ord.Delivery_service, ord.Shard_key, ord.Sm_id,
		ord.Date_created, ord.Oof_shard,
	)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, item := range ord.Items {
		_, err := r.db.Exec(context.Background(),
			` 
			INSERT INTO items (
				order_id, chrt_id, track_number,
				price, rid, "name", sale, size,
				total_price, nm_id, brand, "status"
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7,
				$8, $9, $10, $11, $12
			)
			`,
			ord.Order_uid, item.Chrt_id, item.Track_number,
			item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.Total_price, item.Nm_id, item.Brand, item.Status,
		)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	logrus.Info("Заказ размещен")

	return
}

func (r *OrderRepository) parseJsonToModel(msg *stan.Msg) (models.Order, error) {
	var order models.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		logrus.Error(err)
		return order, err
	}
	return order, nil
}
