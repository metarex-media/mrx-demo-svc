package register

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPrintRegisters(t *testing.T) {
	// create a blank file to stop duplication
	os.Create("registerEntries/register.db")
	db, err := gorm.Open(sqlite.Open("registerEntries/register.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the Metarex schema
	db.AutoMigrate(&MetarexRegister{})

	// Save each register as an individual file
	for ID, reg := range register {
		f, _ := os.Create(fmt.Sprintf("./registerEntries/%v.json", ID))
		regbytes, _ := json.MarshalIndent(reg, "", "    ")
		f.Write(regbytes)

		// Create the register entry
		db.Create(&MetarexRegister{MRXID: ID, Reg: string(regbytes), OwnerID: 0})
	}
}
