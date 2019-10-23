package models

import "errors"

type Cuenta struct {
    ID                 int     `json:"id"`
    NumeroDeCuenta     int     `json:"numero_de_cuenta"`
    Saldo              float32 `json:"saldo"`
    IDCliente          int      `json:"id_cliente"`
    IDTipoDeCuenta     int     `json:"id_tipo_de_cuenta"`
    habilitado         int
    fechaCreacion      string
}

var cuentaSchemeSQLITE string = `CREATE TABLE IF NOT EXISTS cuentas(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    numero_de_cuenta INTEGER NOT NULL UNIQUE,
    saldo REAL DEFAULT 0.0,
    id_cliente TEXT,
    id_tipo_de_cuenta INTEGER
    habilitado INTEGER DEFAULT 0,
    fecha_creacion TEXT);`


func nuevaCuenta(numeroDeCuenta, idCliente, idTipoDeCuenta int) *Cuenta {
    cuenta := &Cuenta {
        NumeroDeCuenta: numeroDeCuenta,
        Saldo:          0.0,
        IDCliente:      idCliente,
        IDTipoDeCuenta: idTipoDeCuenta,
        habilitado:     0,
        fechaCreacion:  ObtenerFechaHoraActualString(),
    }
    return cuenta
}

func AltaCuenta(numeroDeCuenta, idCliente, idTipoDeCuenta int) (*Cuenta, error) {
    cuenta := nuevaCuenta(numeroDeCuenta, idCliente, idTipoDeCuenta )
    err := cuenta.Guardar()
    return cuenta, err
}

func GetCuentaByNumeroCuenta(numeroDeCuenta int) (*Cuenta, error) {
	cuenta := &Cuenta{}
	query := "SELECT id, numero_de_cuenta, saldo, id_cliente, id_tipo_de_cuenta, habilitado, fecha_creacion FROM cuentas WHERE numero_de_cuenta = ?"
	rows, err := Query(query, numeroDeCuenta)
	if err != nil {
		return cuenta, err
	}
	for rows.Next() {
		rows.Scan(&cuenta.ID, &cuenta.NumeroDeCuenta, &cuenta.Saldo, &cuenta.IDCliente, &cuenta.IDTipoDeCuenta, 
            &cuenta.habilitado, &cuenta.fechaCreacion)
	}
	return cuenta, nil
}

func (cuenta *Cuenta) Depositar(monto float32) error {
	cuenta.Saldo = monto
	return cuenta.Guardar()
}

func (cuenta *Cuenta) Retirar(monto float32) error {
	if cuenta.Saldo >= monto {
		cuenta.Saldo -= monto
		return cuenta.Guardar()
	} else {
		return errors.New("Saldo insuficiente")
	}
}

func (cuenta *Cuenta) Transferir(numeroCuentaDestino int, monto float32) error {
	cuentaDestino := &Cuenta{}
	err := errors.New("")

	if cuentaDestino, err = GetCuentaByNumeroCuenta(numeroCuentaDestino); err != nil {
		return err
	}
	if err = cuenta.Retirar(monto); err != nil {
		return err
	}
	err = cuentaDestino.Depositar(monto)
	return err
}

func (cuenta *Cuenta) SolicitarSaldo() (float32, error) {
	var saldo float32
	query := "SELECT saldo FROM cuentas WHERE numero_de_cuenta = ?"
	rows, err := Query(query, cuenta.NumeroDeCuenta)
	if err != nil {
		return saldo, err
	}
	for rows.Next() {
		rows.Scan(&saldo)
	}
	return saldo, nil
}

func (cuenta *Cuenta) Guardar() error {
	if cuenta.ID == 0 {
		return cuenta.registrar()
	} 

	return cuenta.actualizar()
}

func (cuenta *Cuenta) registrar() error {
	query := "INSERT INTO cuentas(numero_de_cuenta, saldo, id_cliente, id_tipo_de_cuenta, habilitado, fecha_creacion) VALUES(?,?,?,?,?,?);"
	cuentaID, err := InsertData(query, cuenta.NumeroDeCuenta, cuenta.Saldo, cuenta.IDCliente, cuenta.IDTipoDeCuenta, 
        cuenta.habilitado, cuenta.fechaCreacion)
	cuenta.ID = int(cuentaID)
	return err
}

func (cuenta *Cuenta) actualizar() error {
	query := "UPDATE cuentas SET saldo=?, id_cliente=?, id_tipo_de_cuenta=?, habilitado=?, fecha_creacion=? WHERE numero_de_cuenta=?;"
	_, err := Exec(query, cuenta.Saldo, cuenta.IDCliente, cuenta.IDTipoDeCuenta, cuenta.habilitado, cuenta.fechaCreacion, 
        cuenta.NumeroDeCuenta)
	return err
}
