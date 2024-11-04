package doc

type Cabecera struct {
	Obligado              Obligado
	Representante         *Obligado
	RemisionVoluntaria    *RemisionVoluntaria
	RemisionRequerimiento *RemisionRequerimiento
}

type Obligado struct {
	NombreRazon string
	NIF         string
}

type RemisionVoluntaria struct {
	FechaFinVerifactu string
	Incidencia        string
}

type RemisionRequerimiento struct {
	RefRequerimiento string
	FinRequerimiento string
}
