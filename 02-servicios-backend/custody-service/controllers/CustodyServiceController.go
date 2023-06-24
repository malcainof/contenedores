package controllers

import (
	"context"
	"errors"
	"regexp"

	pb "github.com/malarcon-79/microservices-lab/grpc-protos-go/system/custody"
	"github.com/malarcon-79/microservices-lab/orm-go/dao"
	"github.com/malarcon-79/microservices-lab/orm-go/model"
	"go.uber.org/zap"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Controlador de servicio gRPC
type CustodyServiceController struct {
	logger *zap.SugaredLogger // Logger
	re     *regexp.Regexp     // Expresión regular para validar formato de períodos YYYY-MM
}

// Método a nivel de package, permite construir una instancia correcta del controlador de servicio gRPC
func NewCustodyServiceController() (CustodyServiceController, error) {
	_logger, _ := zap.NewProduction() // Generamos instancia de logger
	logger := _logger.Sugar()

	re, err := regexp.Compile(`^\d{4}\-(0?[1-9]|1[012])$`) // Expresión regular para validar períodos YYYY-MM
	if err != nil {
		return CustodyServiceController{}, err
	}

	instance := CustodyServiceController{
		logger: logger, // Asignamos el logger
		re:     re,     // Asignamos el RegExp precompilado
	}
	return instance, nil // Devolvemos la nueva instancia de este Struct y un puntero nulo para el error
}

func (c *CustodyServiceController) AddCustodyStock(ctx context.Context, msg *pb.CustodyAdd) (*pb.Empty, error) {
	// Implementar este método
	orm := dao.DB.Model(&model.Custody{})
	out := new(pb.Empty)
	if len(msg.Period) ==  0{
		return nil, errors.New("campo periodo es nulo")
	}
	if !c.re.MatchString(msg.Period) {
		return nil, errors.New("campo periodo invalido")
	}
	if len(msg.Stock) ==  0{
		return nil, errors.New("campo Stock es nulo")
	}
	if len(msg.ClientId) ==  0{
		return nil, errors.New("el id de cliente es nulo")
	}
	if msg.Quantity < 0 {
		return nil, errors.New("la cantidad no debe ser negativa")
	}
	custody := &model.Custody{
		Period:        	msg.Period,
		ClientId:      	msg.ClientId,
		Stock: 			msg.Stock,
		Quantity:   	int32(msg.Quantity),
		Market:			"",
		Price:			decimal.NewFromInt(0),
	}

	if err := orm.Create(custody).Error; err != nil {
		c.logger.Error("error al ingresar custodia", err)
		return nil, errors.New("error al guardar")
	}

	
	return out, nil
}

func (c *CustodyServiceController) ClosePeriod(ctx context.Context, msg *pb.CloseFilters) (*pb.Empty, error) {
	return nil, errors.New("no implementado")
}

func (c *CustodyServiceController) GetCustody(ctx context.Context, msg *pb.CustodyFilter) (*pb.Custodies, error) {
	orm := dao.DB.Model(&model.Custody{})
	custodies := []*model.Custody{}

	filter := &model.Custody{
		Period:        	msg.Period,
		ClientId:      	msg.ClientId,
		Stock: 			msg.Stock,
	}

	if err := orm.Find(&custodies, filter).Error; err != nil {
		c.logger.Errorf("no se pudo buscar custodia con filtros %v", filter, err)
		return nil, status.Errorf(codes.Internal, "no se pudo realizar query")
	}

	result := &pb.Custodies{}
	for _, item := range custodies {
		
		result.Items = append(result.Items, &pb.Custodies_Custody{
			Period:        	item.Period,
			Stock:			item.Stock,
			ClientId:      	item.ClientId,
			Market:   		item.Market,
			Price:       	item.Price.InexactFloat64(),
			Quantity:   	int32(item.Quantity),
		})
	}

	return result,nil
}
