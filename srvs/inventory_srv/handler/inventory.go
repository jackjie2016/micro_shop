package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gorm.io/gorm"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"inventory_srv/global"
	"inventory_srv/model"
	"inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{GoodsID: req.GoodsId}).First(&inv)
	inv.GoodsID = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

/*
获取库存
*/
func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory

	if res := global.DB.Where(&model.Inventory{GoodsID: req.GoodsId}).First(&inv); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.GoodsID,
		Num:     inv.Stocks,
	}, nil
}

//var m sync.Mutex
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	//redis 分布式锁
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	//redis 分布式锁
	tx := global.DB.Begin()
	var GoodsDetail model.GoodsDetailList
	//m.Lock()
	for _, goods := range req.GoodsInfo {
		mutex := rs.NewMutex("stock_sell")
		if err := mutex.Lock(); err != nil {
			zap.S().Error("redis 加锁失败")
		}
		var inv model.Inventory
		var version int32
		if res := tx.Where(&model.Inventory{GoodsID: goods.GoodsId}).First(&inv); res.RowsAffected == 0 {
			zap.S().Infof("回滚吧")
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}

		version = inv.Version
		if inv.Stocks < goods.Num {
			zap.S().Infof("回滚吧")
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		inv.Stocks -= goods.Num
		tx.Where(&model.Inventory{Version: version}).Save(&inv)
		GoodsDetail = append(GoodsDetail, model.GoodsDetail{
			Goods: goods.GoodsId,
			Nums:  goods.Num,
		})

		if ok, err := mutex.Unlock(); !ok || err != nil {
			zap.S().Error("redis 解锁失败")
		}

		//OrderDetai:=model.StockSellDetail{
		//OrderSn: "12222",
		//Status:  1,
		//Detail:  model.GoodsDetailList{{1,3},{2,3}},
		//}
		//db.Create(OrderDetai)

		//乐观锁逻辑
		//for{
		//	var inv model.Inventory
		//	var version int32
		//	//if res := tx.Where(&model.Inventory{GoodsID: goods.GoodsId}).Clauses(clause.Locking{Strength: "UPDATE"}).First(&inv); res.RowsAffected == 0 {
		//	//	zap.S().Infof("回滚吧")
		//	//	tx.Rollback()
		//	//	return nil, status.Errorf(codes.NotFound, "没有库存信息")
		//	//}
		//	if res := tx.Where(&model.Inventory{GoodsID: goods.GoodsId}).First(&inv); res.RowsAffected == 0 {
		//		zap.S().Infof("回滚吧")
		//		tx.Rollback()
		//		return nil, status.Errorf(codes.NotFound, "没有库存信息")
		//	}
		//	version=inv.Version
		//	if inv.Stocks < goods.Num{
		//		zap.S().Infof("回滚吧")
		//		tx.Rollback()
		//		return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		//	}
		//	inv.Stocks -= goods.Num
		//	inv.Version+=1
		//
		//	//tx.Model(&model.Inventory{}).Where("version = ? and goods=?", inv,Version,goods.GoodsId).Updates(&model.Inventory{Version:inv,Version+1,Stocks:inv.Stocks})
		//	if res:=tx.Where(&model.Inventory{Version:version }).Save(&inv);res.RowsAffected==0{
		//		zap.S().Infof("更新失败")
		//	}else{
		//		break
		//	}
		//}
	}
	fmt.Printf("len=%d cap=%d slice=%v\n", len(GoodsDetail), cap(GoodsDetail), GoodsDetail)
	if l := len(GoodsDetail); l > 0 {
		var stockSellDetail = model.StockSellDetail{
			OrderSn: req.OrderSn,
			Status:  1,
			Detail:  GoodsDetail,
		}
		if result := tx.Create(stockSellDetail); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
		}
	}

	zap.S().Infof("回滚了不会执行这边了的吧")
	tx.Commit()
	//m.Unlock()
	return &emptypb.Empty{}, nil
}
func (i *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	tx := global.DB.Begin()
	for _, goods := range req.GoodsInfo {
		var inv model.Inventory
		if res := global.DB.Where(&model.Inventory{GoodsID: goods.GoodsId}).First(&inv); res.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}

		inv.Stocks += goods.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		//	既然是归还库存，那么我应该具体的知道每件商品应该归还多少，但是有一个问题是说明？重复归还的问题
		//	所以说这个接口应该确保幂等性，你不能因为消息的重复发送导致应该订单的库存多少归还，没有扣减的库存你别归还
		//	如果确保这些都没有问题，新建一张表，这张表记录了详细的订单扣减细节，以及归还细节
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("json 解析失败：%v", err.Error())
			return consumer.ConsumeSuccess, nil
		}

		//库存归还，更改记录表
		tx := global.DB.Begin()
		var sellerDetail model.StockSellDetail
		if result := tx.Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellerDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		for _, orderGoods := range sellerDetail.Detail {
			//更新
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{GoodsID: orderGoods.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGoods.Nums)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil //重试
			}
		}

		sellerDetail.Status = 2
		//更新sellerstockdetail 记录
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil //重试
		}
		tx.Commit()
	}
	return consumer.ConsumeSuccess, nil
}
