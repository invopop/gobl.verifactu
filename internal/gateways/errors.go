package gateways

import (
	"errors"
	"strings"
)

// Error codes and their descriptions from VeriFactu
var ErrorCodes = map[string]string{
	// Errors that cause rejection of the entire submission
	"4102": "El XML no cumple el esquema. Falta informar campo obligatorio.",
	"4103": "Se ha producido un error inesperado al parsear el XML.",
	"4104": "Error en la cabecera: el valor del campo NIF del bloque ObligadoEmision no está identificado.",
	"4105": "Error en la cabecera: el valor del campo NIF del bloque Representante no está identificado.",
	"4106": "El formato de fecha es incorrecto.",
	"4107": "El NIF no está identificado en el censo de la AEAT.",
	"4108": "Error técnico al obtener el certificado.",
	"4109": "El formato del NIF es incorrecto.",
	"4110": "Error técnico al comprobar los apoderamientos.",
	"4111": "Error técnico al crear el trámite.",
	"4112": "El titular del certificado debe ser Obligado Emisión, Colaborador Social, Apoderado o Sucesor.",
	"4113": "El XML no cumple con el esquema: se ha superado el límite permitido de registros para el bloque.",
	"4114": "El XML no cumple con el esquema: se ha superado el límite máximo permitido de facturas a registrar.",
	"4115": "El valor del campo NIF del bloque ObligadoEmision es incorrecto.",
	"4116": "Error en la cabecera: el campo NIF del bloque ObligadoEmision tiene un formato incorrecto.",
	"4117": "Error en la cabecera: el campo NIF del bloque Representante tiene un formato incorrecto.",
	"4118": "Error técnico: la dirección no se corresponde con el fichero de entrada.",
	"4119": "Error al informar caracteres cuya codificación no es UTF-8.",
	"4120": "Error en la cabecera: el valor del campo FechaFinVeriFactu es incorrecto, debe ser 31-12-20XX, donde XX corresponde con el año actual o el anterior.",
	"4121": "Error en la cabecera: el valor del campo Incidencia es incorrecto.",
	"4122": "Error en la cabecera: el valor del campo RefRequerimiento es incorrecto.",
	"4123": "Error en la cabecera: el valor del campo NIF del bloque Representante no está identificado en el censo de la AEAT.",
	"4124": "Error en la cabecera: el valor del campo Nombre del bloque Representante no está identificado en el censo de la AEAT.",
	"4125": "Error en la cabecera: el campo RefRequerimiento es obligatorio.",
	"4126": "Error en la cabecera: el campo RefRequerimiento solo debe informarse en sistemas No VERIFACTU.",
	"4127": "Error en la cabecera: la remisión voluntaria solo debe informarse en sistemas VERIFACTU.",
	"4128": "Error técnico en la recuperación del valor del Gestor de Tablas.",
	"4129": "Error en la cabecera: el campo FinRequerimiento es obligatorio.",
	"4130": "Error en la cabecera: el campo FinRequerimiento solo debe informarse en sistemas No VERIFACTU.",
	"4131": "Error en la cabecera: el valor del campo FinRequerimiento es incorrecto.",
	"4132": "El titular del certificado debe ser el destinatario que realiza la consulta, un Apoderado o Sucesor",
	"3500": "Error técnico de base de datos: error en la integridad de la información.",
	"3501": "Error técnico de base de datos.",
	"3502": "La factura consultada para el suministro de pagos/cobros/inmuebles no existe.",
	"3503": "La factura especificada no pertenece al titular registrado en el sistema.",

	// Errors that cause rejection of the invoice or entire request if in header
	"1100": "Valor o tipo incorrecto del campo.",
	"1101": "El valor del campo CodigoPais es incorrecto.",
	"1102": "El valor del campo IDType es incorrecto.",
	"1103": "El valor del campo ID es incorrecto.",
	"1104": "El valor del campo NumSerieFactura es incorrecto.",
	"1105": "El valor del campo FechaExpedicionFactura es incorrecto.",
	"1106": "El valor del campo TipoFactura no está incluido en la lista de valores permitidos.",
	"1107": "El valor del campo TipoRectificativa es incorrecto.",
	"1108": "El NIF del IDEmisorFactura debe ser el mismo que el NIF del ObligadoEmision.",
	"1109": "El NIF no está identificado en el censo de la AEAT.",
	"1110": "El NIF no está identificado en el censo de la AEAT.",
	"1111": "El campo CodigoPais es obligatorio cuando IDType es distinto de 02.",
}

// Standard gateway error responses
var (
	ErrConnection = newError("connection")
	ErrInvalid    = newError("invalid")
	ErrDuplicate  = newError("duplicate")
)

// Error allows for structured responses from the gateway to be able to
// response codes and messages.
type Error struct {
	key     string
	code    string
	message string
	cause   error
}

// Error produces a human readable error message.
func (e *Error) Error() string {
	out := []string{e.key}
	if e.code != "" {
		out = append(out, e.code)
	}
	if e.message != "" {
		out = append(out, e.message)
	}
	return strings.Join(out, ": ")
}

// Key returns the key for the error.
func (e *Error) Key() string {
	return e.key
}

// Message returns the human message for the error.
func (e *Error) Message() string {
	return e.message
}

// Code returns the code provided by the remote service.
func (e *Error) Code() string {
	return e.code
}

func newError(key string) *Error {
	return &Error{key: key}
}

// withCode duplicates and adds the code to the error.
func (e *Error) withCode(code string) *Error {
	e = e.clone()
	e.code = code
	return e
}

// withMessage duplicates and adds the message to the error.
func (e *Error) withMessage(msg string) *Error {
	e = e.clone()
	e.message = msg
	return e
}

func (e *Error) withCause(err error) *Error {
	e = e.clone()
	e.cause = err
	e.message = err.Error()
	return e
}

func (e *Error) clone() *Error {
	ne := new(Error)
	*ne = *e
	return ne
}

// Is checks to see if the target error is the same as the current one
// or forms part of the chain.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.cause, target)
	}
	return e.key == t.key
}
