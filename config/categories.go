package config

// Definição das categorias como constantes
const (
	CategoryTransport     = "Transporte"
	CategoryConvenience   = "Bebidas e Conveniência"
	CategoryFood          = "Restaurantes e Alimentação"
	CategoryMarket        = "Mercado e Compras"
	CategorySubscriptions = "Assinaturas e Serviços Digitais"
	CategoryEntertainment = "Entretenimento e Eventos"
	CategoryPharmacy      = "Farmácias e Saúde"
	CategoryClothing      = "Vestuário e Cosméticos"
	CategoryBarbershop    = "Barbearia e Beleza"
	CategoryElectronics   = "Eletrônicos e Tecnologia"
	CategoryOptical       = "Ótica e Acessórios"
)

var CategoryMap = map[string]string{
	// Transporte
	"Uber - NuPay": CategoryTransport,
	"99app *99app": CategoryTransport,

	// Bebidas e Conveniência
	"Zé Delivery - NuPay": CategoryConvenience,
	"Pontinho do Beco":    CategoryConvenience,
	"Mesconvenienciae":    CategoryConvenience,

	// Restaurantes e Alimentação
	"Ifd*Ftr Restaurante":    CategoryFood,
	"Superheroisburger":      CategoryFood,
	"Emporio Art Cafe":       CategoryFood,
	"Yoidon Restaurante":     CategoryFood,
	"Restaurante Veg":        CategoryFood,
	"Ifd*Le Gourmet Comerci": CategoryFood,
	"Sabor e Ar":             CategoryFood,
	"Pane Di Giovanni":       CategoryFood,
	"Bolo la Dcasa":          CategoryFood,

	// Mercado e Compras
	"Muffato Jk":            CategoryMarket,
	"Muffato Madre":         CategoryMarket,
	"Caheteltg Comercio de": CategoryMarket,

	// Assinaturas e Serviços Digitais
	"Ebanx*Crunchyroll": CategorySubscriptions,
	"Google One":        CategorySubscriptions,
	"Netflix.Com":       CategorySubscriptions,
	"Dm*Spotify":        CategorySubscriptions,
	"Ifd*Ifood Club":    CategorySubscriptions,

	// Entretenimento e Eventos
	"Muvuka Eventos":         CategoryEntertainment,
	"Moviesystem Cinematogr": CategoryEntertainment,
	"EVENTIMCOMBR":           CategoryEntertainment,

	// Farmácias e Saúde
	"Farmacia e Drogaria Ni": CategoryPharmacy,
	"Panvel Farmacias":       CategoryPharmacy,

	// Vestuário e Cosméticos
	"Brs*Sheincom": CategoryClothing,
	"Ec *Sallve":   CategoryClothing,
	"Bawclothing":  CategoryClothing,

	// Barbearia e Beleza
	"Mi Casa Barbearia Ba": CategoryBarbershop,

	// Eletrônicos e Tecnologia
	"Hub*Kabum": CategoryElectronics,

	// Ótica e Acessórios
	"Duty Otica I": CategoryOptical,
}
