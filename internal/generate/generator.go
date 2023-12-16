package generate

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/mkorobovv/L0/internal/models"
)

func Generate() *models.OrderJSON {

	order_uid := hash32()
	var orderCount = 1 + rand.Intn(2)
	items := make([]models.Item, orderCount)
	for i := 0; i < orderCount; i++ {
		items[i] = models.Item{0, "", 0, "", "", 0, "", 0, 0, "", 0}
	}
	order := &models.OrderJSON{
		//Random values
		Order_uid:          order_uid[:len(order_uid)-15],
		Track_number:       "WBILMTESTTRACK",
		Entry:              "WBIL",
		Delivery:           models.Delivery{Name: "test", Phone: "+79999999999", Zip: "0", City: "moscow", Adress: "adress", Region: "msk", Email: "email"},
		Payments:           models.Payment{Transaction: "", Request_id: "", Currency: "", Provider: "", Amount: 0, Payment_dt: 0, Bank: "", Delivery_cost: 0, Goods_total: 0, Custom_fee: 0},
		Items:              items,
		Locale:             "en",
		Internal_siganture: "",
		Customer_id:        "test",
		Delivery_service:   "some service",
		Shardkey:           "9",
		Sm_id:              0,
		Date_created:       time.Now().Format(time.RFC3339),
		OOF_shard:          "0",
	}
	generateOrderPayment(order)
	generateDelivery(order)
	generateItem(order)

	return order
}

func generateOrderPayment(order *models.OrderJSON) {
	currency := []string{"USD", "RUB", "EUR"}
	banks := []string{"sber", "alpha", "tinkoff"}
	var amount float32 = 0
	for i := range order.Items {
		amount += order.Items[i].Total_price
	}
	deliveryCost := float32(rand.Intn(1500))
	order.Payments.Transaction = order.Order_uid + order.Customer_id
	order.Payments.Request_id = ""
	order.Payments.Currency = currency[rand.Intn(len(currency))]
	order.Payments.Provider = "wbpay"
	order.Payments.Amount = amount + deliveryCost
	order.Payments.Payment_dt = uint32(1000000000 + rand.Intn(1000000000))
	order.Payments.Bank = banks[rand.Intn(len(banks))]
	order.Payments.Delivery_cost = uint32(deliveryCost)
	order.Payments.Goods_total = amount
	order.Payments.Custom_fee = 0
}

func generateDelivery(order *models.OrderJSON) {
	names := []string{"Lazy Gopher", "Ivanov Ivan", "Happy Python", "Elon Musk"}
	addresses := []string{"Komarova 31", "Somestreet -101", "OKStreet 202", "Chongarskii 11"}

	order.Delivery.Name = names[rand.Intn(len(names))]
	order.Delivery.Phone = "+" + strconv.Itoa(1000000000+rand.Intn(8000000000))
	order.Delivery.Zip = strconv.Itoa(100000 + rand.Intn(150000))
	order.Delivery.City = "Moscow"
	order.Delivery.Adress = addresses[rand.Intn(len(addresses))]
	order.Delivery.Region = "Moscow"
	order.Delivery.Email = "example@gmail.com"

}

func generateItem(order *models.OrderJSON) {
	for i := 0; i < len(order.Items); i++ {
		amount := float32(1500 + rand.Intn(10000))
		sale := float32(rand.Intn(50))
		total := ((100 - sale) / 100.0) * amount
		totalPrice := math.Round(float64(total)*10) / 10
		order.Items[i] = models.Item{
			uint32(rand.Intn(1000000)),
			order.Track_number,
			uint16(amount),
			hash32()[:len(hash32())-15] + order.Customer_id,
			"Mascaras",
			uint16(sale),
			"0",
			float32(totalPrice),
			uint32(rand.Intn(1000000)),
			"Vivienne Sabo",
			202,
		}
	}
}

func hash32() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(150000))))
	return hex.EncodeToString(sum[:])
}
