package struct_folder

var photos = []string{
	"http://188.32.24.142:27012/api/photo/12.png",
	"http://188.32.24.142:27012/api/photo/13.png",
}

type Variant struct {
	Id    []int
	Value []int
	Price []int
	Unit  string
}

type Product struct {
	Name           string
	Url            string
	Photo          []string
	Description    string
	Characteristic []map[string]string
	Variants       Variant
}

type MiniProducts struct {
	ID         int
	Name       string
	Url        string
	MainPhoto  string
	Price      string
	Categories []string
}

type BonusInfoProduct struct {
	Decscription   string
	Characteristic []map[string]string
	Photo          []string
	Variants       Variant
}

type OrderItem struct {
	ID        int     `json:"id"`         // Ожидает "id"
	VariantID int     `json:"variant_id"` // Ожидает "variant_id"
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Qty       int     `json:"quantity"` // Ожидает "quantity"
}

type OrderData struct {
	Customer  CustomerInfo `json:"customer"`
	Items     []OrderItem  `json:"items"`
	Total     string       `json:"total"`
	CreatedAt string       `json:"createdAt"`
}

type CustomerInfo struct {
	Name    string `json:"name"`    // соответствует name="name" в HTML-форме
	Phone   string `json:"phone"`   // соответствует name="phone"
	Comment string `json:"comment"` // соответствует name="comment"
}

type OrdersInfo struct {
	Id           int
	Created_at   string
	CustomerName string
	Phone        string
	TotalPrice   int
	Status       string
}

type ProductInfo struct {
	Name  string
	Value string
	Unit  string
	Count int
	Price int
}

type AdminInfo struct {
	OrdersInfo  OrdersInfo
	ProductInfo []ProductInfo
}
