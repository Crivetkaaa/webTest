package struct_folder

var photos = []string{
	"http://188.32.24.142:27012/api/photo/12.png",
	"http://188.32.24.142:27012/api/photo/13.png",
}

type Variant struct {
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
