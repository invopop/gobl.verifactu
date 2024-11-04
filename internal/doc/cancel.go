package doc

type RegistroAnulacion struct {
	IDVersion                string         `xml:"IDVersion"`
	IDFactura                IDFactura      `xml:"IDFactura"`
	RefExterna               string         `xml:"RefExterna,omitempty"`
	SinRegistroPrevio        string         `xml:"SinRegistroPrevio"`
	RechazoPrevio            string         `xml:"RechazoPrevio,omitempty"`
	GeneradoPor              string         `xml:"GeneradoPor"`
	Generador                *Tercero       `xml:"Generador"`
	Encadenamiento           Encadenamiento `xml:"Encadenamiento"`
	SistemaInformatico       Software       `xml:"SistemaInformatico"`
	FechaHoraHusoGenRegistro string         `xml:"FechaHoraHusoGenRegistro"`
	TipoHuella               string         `xml:"TipoHuella"`
	Huella                   string         `xml:"Huella"`
	Signature                string         `xml:"Signature"`
}
