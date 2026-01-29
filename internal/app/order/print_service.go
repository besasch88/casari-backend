package order

import (
	"net"

	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hennedo/escpos"
	"gorm.io/gorm"
)

type printServiceInterface interface {
	print(ctx *gin.Context, input printOrderInputDto) error
}

type printService struct {
	storage         *gorm.DB
	pubSubAgent     *ceng_pubsub.PubSubAgent
	repository      orderRepositoryInterface
	printRepository printRepositoryInterface
}

func newPrintService(storage *gorm.DB, pubSubAgent *ceng_pubsub.PubSubAgent, repository orderRepositoryInterface, printRepository printRepositoryInterface) printServiceInterface {
	return printService{
		storage:         storage,
		pubSubAgent:     pubSubAgent,
		repository:      repository,
		printRepository: printRepository,
	}
}

func (s printService) print(ctx *gin.Context, input printOrderInputDto) error {
	switch input.Target {
	case "order":
		return s.printOrder(uuid.MustParse(input.TableID))
	case "course":
		return s.printCourse(uuid.MustParse(input.TableID), uuid.MustParse(*input.CourseID))
	case "bill":
		return s.printBill(uuid.MustParse(input.TableID))
	case "payment":
		return s.printPayment(uuid.MustParse(input.TableID))

	default:
		return errInvalidPrintRequest
	}
}

func (s printService) printOrder(tableId uuid.UUID) error {
	items, err := s.repository.getOrderDetailByTableID(s.storage, tableId)
	if err != nil {
		return err
	}
	return s.printItems(items)
}

func (s printService) printCourse(tableId uuid.UUID, courseId uuid.UUID) error {
	items, err := s.repository.getOrderDetailByTableIDAndCourseID(s.storage, tableId, courseId)
	if err != nil {
		return err
	}
	return s.printItems(items)
}

func (s printService) printBill(tableId uuid.UUID) error {
	items, err := s.repository.getPricedOrderByTableID(s.storage, tableId)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}
	conn, err := net.Dial("tcp", items[0].PrinterURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	printer := escpos.New(conn)
	s.printRepository.printTable(printer, items[0].TableName)
	s.printRepository.printPrinterName(printer, items[0].PrinterTitle)
	s.printRepository.printTableCreation(printer, items[0].TableCreatedAt)
	s.printRepository.printLine(printer)
	total := int64(0)
	for _, item := range items {
		if item.MenuOptionTitle != nil {
			err = s.printRepository.printItemAndPrice(printer, item.Quantity, *item.MenuOptionTitle, *item.MenuOptionPrice)
			total = total + (item.Quantity * *item.MenuOptionPrice)
		} else {
			err = s.printRepository.printItemAndPrice(printer, item.Quantity, item.MenuItemTitle, item.MenuItemPrice)
			total = total + (item.Quantity * item.MenuItemPrice)
		}
		if err != nil {
			return err
		}
	}
	s.printRepository.printLine(printer)
	s.printRepository.printTotalPrice(printer, total)
	s.printRepository.printLine(printer)
	s.printRepository.printRecipeCollection(printer)
	s.printRepository.printAndCut(printer)
	return nil
}

func (s printService) printPayment(tableId uuid.UUID) error {
	item, err := s.repository.getTotalPriceAndPaymentByTableID(s.storage, tableId)
	if err != nil {
		return err
	}
	if ceng_utils.IsEmpty(item) {
		return nil
	}
	conn, err := net.Dial("tcp", item.PrinterURL)
	if err != nil {
		return err
	}
	defer conn.Close()
	printer := escpos.New(conn)
	s.printRepository.printTable(printer, item.TableName)
	s.printRepository.printPrinterName(printer, item.PrinterTitle)
	s.printRepository.printTableCreation(printer, item.TableCreatedAt)
	s.printRepository.printLine(printer)
	s.printRepository.printPaymentMethod(printer, item.TablePayment, item.PriceTotal)
	s.printRepository.printAndCut(printer)
	return nil
}

func (s printService) printItems(items []OrderDetailEntity) error {
	var conn net.Conn
	var err error
	var printer *escpos.Escpos
	lastPrinterTitle := ""
	lastCourseID := ""
	for _, item := range items {
		if item.PrinterTitle != lastPrinterTitle {
			// if the printer is changing, send the print and cut and close the previous connection
			if printer != nil {
				s.printRepository.printAndCut(printer)
			}
			if conn != nil {
				conn.Close()
			}
			// So create a connection to the new printer
			conn, err = net.Dial("tcp", item.PrinterURL)
			if err != nil {
				return err
			}
			printer = escpos.New(conn)
			lastPrinterTitle = item.PrinterTitle
			lastCourseID = ""
			s.printRepository.printTable(printer, item.TableName)
			s.printRepository.printPrinterName(printer, item.PrinterTitle)
			s.printRepository.printTableCreation(printer, item.TableCreatedAt)
		}
		if item.CourseID != lastCourseID {
			lastCourseID = item.CourseID
			// The course changed, so print it on paper
			s.printRepository.printCourse(printer, item.CourseNumber)
		}
		// Now for each element, I can print them
		if item.MenuOptionTitle != nil {
			s.printRepository.printItem(printer, item.Quantity, *item.MenuOptionTitle)
		} else {
			s.printRepository.printItem(printer, item.Quantity, item.MenuItemTitle)
		}
	}
	// At the end, if needed, print and cut and close the connection
	if conn != nil && printer != nil {
		s.printRepository.printAndCut(printer)
		conn.Close()
	}
	return nil
}
