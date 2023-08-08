package fans

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Bluebugs/rpi-poe-fan/pkg/ioutil"
)

type ControlMode int

const (
	// ControlModeDisabled completely disables control, resulting in a 100% voltage/PWM signal output
	ControlModeDisabled ControlMode = iota
	// ControlModePWM enables manual, fixed speed control via setting the pwm value
	ControlModePWM ControlMode = 1
	// ControlModeAutomatic enables automatic control by the integrated control of the mainboard
	ControlModeAutomatic ControlMode = 2
)

const hwmonPath = "/sys/class/hwmon"

type Fan interface {
	fmt.Stringer
	SetSpeed(percentage uint8) error
	Speed() (uint8, error)
	Control() (ControlMode, error)
	Name() string
}

type PoEFan struct {
	name          string
	pwmPath       string
	pwmEnablePath string
}

var _ Fan = (*PoEFan)(nil)

func HwMon() (*PoEFan, error) {
	hwmons, err := os.ReadDir(hwmonPath)
	if err != nil {
		return nil, err
	}

	for _, hwmon := range hwmons {
		path := filepath.Join(hwmonPath, hwmon.Name(), "name")
		name, err := ioutil.ReadStringFromFile(path)
		if err != nil {
			log.Println(path)
			continue
		}
		if name != "pwmfan" {
			continue
		}

		pwmPath := filepath.Join(hwmonPath, hwmon.Name(), "pwm1")
		pwmEnablePath := filepath.Join(hwmonPath, hwmon.Name(), "pwm1_enable")

		return &PoEFan{
			name:          name,
			pwmPath:       pwmPath,
			pwmEnablePath: pwmEnablePath,
		}, nil
	}

	return nil, fmt.Errorf("no fan found")
}

func (f *PoEFan) SetSpeed(percentage uint8) error {
	if percentage > 100 {
		return fmt.Errorf("percentage must be between 0 and 100, got %d", percentage)
	}
	return ioutil.WriteIntToFile(f.pwmPath, int(percentage))

}

func (f *PoEFan) Speed() (uint8, error) {
	speed, err := ioutil.ReadIntFromFile(f.pwmPath)
	if err != nil {
		return 0, err
	}
	return uint8(speed), nil
}

func (f *PoEFan) Control() (ControlMode, error) {
	enabled, err := ioutil.ReadIntFromFile(f.pwmEnablePath)
	if err != nil {
		return 0, err
	}
	return ControlMode(enabled), nil
}

func (f *PoEFan) Name() string {
	return f.name
}

func (f *PoEFan) String() string {
	if f == nil {
		return "nil"
	}

	speed, err := f.Speed()
	if err != nil {
		return fmt.Sprintf("%s: %s", f.Name(), err)
	}
	return fmt.Sprintf("%s: %d%%", f.Name(), speed)
}
