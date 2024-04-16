package gopicontrol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func (ctrl *RevPiControl) WriteDO(name string, v bool) (err error) {
	var (
		sPIValue SPIValue
		variable *SPIVariable
	)

	variable, err = ctrl.GetVariableInfo(name)
	if err != nil {
		return
	}

	sPIValue.I16uAddress = variable.I16uAddress
	sPIValue.I8uBit = variable.I8uBit

	if v {
		sPIValue.I8uValue = 1
	} else {
		sPIValue.I8uValue = 0
	}

	return ctrl.SetBitValue(&sPIValue)
}

func (ctrl *RevPiControl) WriteAO(name string, v uint32) (err error) {
	var variable *SPIVariable

	variable, err = ctrl.GetVariableInfo(name)
	if err != nil {
		return
	}

	var data interface{}
	switch variable.I16uLength {
	case 8:
		data = uint8(v)
	case 16:
		data = uint16(v)
	case 32:
		data = uint32(v)
	}

	b, e := NumToBytes(data)
	if e != nil {
		return e
	}

	_, err = ctrl.Write(uint32(variable.I16uAddress), b)

	return err
}

func (ctrl *RevPiControl) ReadDI(name string) (v bool, err error) {
	var (
		value    SPIValue
		variable *SPIVariable
	)

	variable, err = ctrl.GetVariableInfo(name)
	if err != nil {
		return
	}

	value.I16uAddress = variable.I16uAddress
	value.I8uBit = variable.I8uBit

	err = ctrl.GetBitValue(&value)
	if err != nil {
		return
	}

	return value.I8uValue > 0, err
}

func (ctrl *RevPiControl) ReadAI(name string) (v uint32, err error) {
	var variable *SPIVariable

	variable, err = ctrl.GetVariableInfo(name)
	if err != nil {
		return
	}

	sizeRemainder := variable.I16uLength % 8
	if sizeRemainder != 0 {
		return v, fmt.Errorf("could not read variable %s. Internal Error", name)
	}

	data := make([]byte, variable.I16uLength/8)
	if _, err = ctrl.Read(uint32(variable.I16uAddress), data); err != nil {
		return
	}

	buf := bytes.NewReader(data)
	switch variable.I16uLength {
	case 8:
		var ui8 uint8
		err = binary.Read(buf, binary.LittleEndian, &ui8)
		v = uint32(ui8)
	case 16:
		var ui16 uint16
		err = binary.Read(buf, binary.LittleEndian, &ui16)
		v = uint32(ui16)
	case 32:
		err = binary.Read(buf, binary.LittleEndian, &v)
	}

	return v, err
}
