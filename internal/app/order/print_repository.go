package order

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hennedo/escpos"
	"github.com/qiniu/iconv"
)

type printRepositoryInterface interface {
	printTable(printer *escpos.Escpos, name string) error
	printPrinterName(printer *escpos.Escpos, name string) error
	printTableCreation(printer *escpos.Escpos, date time.Time) error
	printCourse(printer *escpos.Escpos, number int64) error
	printLine(printer *escpos.Escpos) error
	printItem(printer *escpos.Escpos, quantity int64, name string) error
	printItemAndPrice(printer *escpos.Escpos, quantity int64, name string, price int64) error
	printTotalPrice(printer *escpos.Escpos, price int64) error
	printAndCut(printer *escpos.Escpos) error
}

type printRepository struct {
}

func newPrintRepository() printRepositoryInterface {
	return printRepository{}
}

func (r printRepository) printTable(printer *escpos.Escpos, name string) error {
	if _, err := printer.Bold(true).Reverse(false).Size(3, 2).Justify(escpos.JustifyCenter).Write(fmt.Sprintf("TAV. %s\n", name)); err != nil {
		return err
	}
	return nil
}

func (r printRepository) printPrinterName(printer *escpos.Escpos, name string) error {
	if _, err := printer.Bold(true).Reverse(true).Size(2, 2).Justify(escpos.JustifyCenter).Write(fmt.Sprintf(" %s \n\n", name)); err != nil {
		return err
	}
	return nil
}

func (r printRepository) printTableCreation(printer *escpos.Escpos, date time.Time) error {
	// Convert date in Rome Timezone
	location, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		location = time.Local
	}
	date = date.In(location)
	if _, err := printer.Bold(false).Reverse(false).Size(1, 1).Justify(escpos.JustifyCenter).Write(fmt.Sprintf("Creato il %s\n", date.Format("02/01/2006 15:04"))); err != nil {
		return err
	}
	return nil
}

func (r printRepository) printCourse(printer *escpos.Escpos, number int64) error {
	if err := r.printLine(printer); err != nil {
		return err
	}
	if _, err := printer.Bold(true).Reverse(false).Size(2, 2).Justify(escpos.JustifyLeft).Write(fmt.Sprintf("PORTATA %d\n", number)); err != nil {
		return err
	}
	if err := r.printLine(printer); err != nil {
		return err
	}
	return nil
}

func (r printRepository) printLine(printer *escpos.Escpos) error {
	_, err := printer.Bold(false).Reverse(false).Size(2, 2).Justify(escpos.JustifyLeft).Write("------------------------\n")
	return err
}

func (r printRepository) printItem(printer *escpos.Escpos, quantity int64, name string) error {
	_, err := printer.Bold(false).Reverse(false).Size(1, 2).Justify(escpos.JustifyLeft).Write(fmt.Sprintf("%2d x %s\n\n", quantity, name))
	return err
}

func (r printRepository) printItemAndPrice(printer *escpos.Escpos, quantity int64, name string, price int64) error {
	partial := quantity * price
	c, err := iconv.Open("cp858", "utf-8")
	if err != nil {
		return err
	}
	defer c.Close()
	leftString := fmt.Sprintf("%2d x %s", quantity, name)
	rightString := fmt.Sprintf("%2d x %.2f€\n", quantity, float64(price)/100)
	totalString := fmt.Sprintf("= %.2f€\n", float64(partial)/100)
	toRepeat := 49 - utf8.RuneCountInString(leftString) - utf8.RuneCountInString(rightString)
	spaceString := ""
	if toRepeat > 0 {
		spaceString = strings.Repeat(" ", toRepeat)
	}
	str := c.ConvString(fmt.Sprintf("%s%s%s", leftString, spaceString, rightString))
	_, err = printer.Bold(false).Reverse(false).Size(1, 1).Justify(escpos.JustifyLeft).Write(str)
	totalStr := c.ConvString(fmt.Sprintf("%s", totalString))
	_, err = printer.Bold(true).Reverse(false).Size(1, 1).Justify(escpos.JustifyRight).Write(totalStr)
	return err
}

func (r printRepository) printTotalPrice(printer *escpos.Escpos, price int64) error {
	c, err := iconv.Open("cp858", "utf-8")
	if err != nil {
		return err
	}
	defer c.Close()
	text := fmt.Sprintf("%.2f€\n\n", float64(price)/100)
	convertedText := c.ConvString(text)
	_, err = printer.Bold(true).Reverse(false).Size(1, 1).Justify(escpos.JustifyRight).Write(convertedText)
	return err
}

func (r printRepository) printAndCut(printer *escpos.Escpos) error {
	return printer.PrintAndCut()
}
