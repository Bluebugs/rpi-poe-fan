package cpu

import "github.com/Bluebugs/rpi-poe-fan/pkg/ioutil"

type Temp interface {
	Read() (float32, error)
}

type RPiTemp struct {
}

var _ Temp = (*RPiTemp)(nil)

const cpuTempSysPath = "/sys/class/thermal/thermal_zone0/temp"

func NewRPiTemp() *RPiTemp {
	return &RPiTemp{}
}

func (t *RPiTemp) Read() (float32, error) {
	temp, err := ioutil.ReadIntFromFile(cpuTempSysPath)
	if err != nil {
		return 0, err
	}
	return float32(temp) / 1000, nil
}
