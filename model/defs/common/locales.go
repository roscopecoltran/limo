package common 

/*
	Locales definitions
*/
type Locales_Options struct {
	// gorm.Model
	Disable 			bool 						`default:"false" yaml:"disable,omitempty" json:"disable,omitempty"`
	Custom 				map[string]string 			`json:"custom,omitempty" yaml:"custom,omitempty,omitempty"`
	Default 			Locales_DefaultOptions		`json:"default,omitempty" yaml:"default,omitempty,omitempty"`
}

type Locales_DefaultOptions struct {
	Bg   				string 						`json:"bg,omitempty" yaml:"bg,omitempty,omitempty"`
	Cs   				string 						`json:"cs,omitempty" yaml:"cs,omitempty,omitempty"`
	De   				string 						`json:"de,omitempty" yaml:"de,omitempty,omitempty"`
	DeDE 				string 						`json:"de_DE,omitempty" yaml:"de_DE,omitempty"`
	ElGR 				string 						`json:"el_GR,omitempty" yaml:"el_GR,omitempty"`
	En   				string 						`json:"en,omitempty" yaml:"en,omitempty"`
	Eo   				string 						`json:"eo,omitempty" yaml:"eo,omitempty"`
	Es   				string 						`json:"es,omitempty" yaml:"es,omitempty"`
	Fi   				string 						`json:"fi,omitempty" yaml:"fi,omitempty"`
	Fr   				string 						`json:"fr,omitempty" yaml:"fr,omitempty"`
	He   				string 						`json:"he,omitempty" yaml:"he,omitempty"`
	Hu   				string 						`json:"hu,omitempty" yaml:"hu,omitempty"`
	It   				string 						`json:"it,omitempty" yaml:"it,omitempty"`
	Ja   				string 						`json:"ja,omitempty" yaml:"ja,omitempty"`
	Nl   				string 						`json:"nl,omitempty" yaml:"nl,omitempty"`
	Pt   				string 						`json:"pt,omitempty" yaml:"pt,omitempty"`
	PtBR 				string 						`json:"pt_BR,omitempty" yaml:"pt_BR,omitempty"`
	Ro   				string 						`json:"ro,omitempty" yaml:"ro,omitempty"`
	Ru   				string 						`json:"ru,omitempty" yaml:"ru,omitempty"`
	Sk   				string 						`json:"sk,omitempty" yaml:"sk,omitempty"`
	Sv   				string 						`json:"sv,omitempty" yaml:"sv,omitempty"`
	Tr   				string 						`json:"tr,omitempty" yaml:"tr,omitempty"`
	Uk   				string 						`json:"uk,omitempty" yaml:"uk,omitempty"`
	Zh   				string 						`json:"zh,omitempty" yaml:"zh,omitempty"`	
}
