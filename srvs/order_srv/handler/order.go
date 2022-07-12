package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"order_srv/global"
	"order_srv/model"
	"order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}
type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	ctx         context.Context
}

func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

//事务提交，如果是失败或者commit 不回调检查，如果是未知错误就进入checkLocalTransaction
// 是 rollback 还是 commit 根据库存是否扣减，如果库存扣减了，那commit之后，发送一个order_reback 消息
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	var shopCarts []model.ShoppingCart

	parentSpan := opentracing.SpanFromContext(o.ctx) //获取父节点

	_ = json.Unmarshal(msg.Body, &orderInfo)
	//获取购物车产品id
	var goodsIDs []int32

	goodsNumsMap := make(map[int32]int32)
	shopcartSpan := opentracing.GlobalTracer().StartSpan("shopcart", opentracing.ChildOf(parentSpan.Context()))
	//1. 从购物车中获取到选中的商品
	if res := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); res.RowsAffected == 0 {
		o.Code = codes.NotFound
		o.Detail = "购物车中没有商品"
		shopcartSpan.Finish()
		return primitive.RollbackMessageState
	}
	shopcartSpan.Finish()

	for _, shopCart := range shopCarts {
		goodsIDs = append(goodsIDs, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	//查询商品微服务获取价格
	GoodsSrvSpan := opentracing.GlobalTracer().StartSpan("goods-srv", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIDs})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	GoodsSrvSpan.Finish()

	//计算库存
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += float32(goodsNumsMap[good.Id]) * good.ShopPrice
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	//跨服务调用库存微服务进行库存扣减
	queryInvSpan := opentracing.GlobalTracer().StartSpan("inv-srv", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsInvInfo, OrderSn: orderInfo.OrderSn}); err != nil {
		//需要深入判断是什么原因导致的失败是网络问题还是程序问题，或者数据问题
		zap.S().Errorf("扣减库存失败:%v", err.Error())
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}
	queryInvSpan.Finish()

	//生成订单表
	//20210308xxxx
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	SaveOrderSpan := opentracing.GlobalTracer().StartSpan("save-order", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		zap.S().Errorf("创建订单失败:%v", err.Error())
		return primitive.CommitMessageState //创建失败立即回滚
	}
	SaveOrderSpan.Finish()

	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	//批量插入orderGoods
	SaveGoodsSpan := opentracing.GlobalTracer().StartSpan("save-goods", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单商品明细失败"
		zap.S().Errorf("创建订单商品明细失败:%v", err)
		return primitive.CommitMessageState //创建失败立即回滚

	}
	SaveGoodsSpan.Finish()

	DeleteCartSpan := opentracing.GlobalTracer().StartSpan("delete-cart", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车失败"
		zap.S().Errorf("删除购物车失败:%v", err.Error())
		return primitive.CommitMessageState //创建失败立即回滚

	}
	DeleteCartSpan.Finish()

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	o.Code = codes.OK

	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMq.Host, global.ServerConfig.RocketMq.Port)})),
		producer.WithRetry(2),
	)

	if err != nil {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "生成延迟消息producer失败"
		zap.S().Errorf("生成 producer 失败: %s", err.Error())
		return primitive.CommitMessageState //创建失败立即回滚
	}

	err = p.Start()
	if err != nil {
		//os.Exit(1)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "启动延迟消息producer失败"
		zap.S().Errorf("启动延迟消息producer失败: %s", err.Error())
		return primitive.CommitMessageState //创建失败立即回滚
	}

	msg2 := &primitive.Message{
		Topic: "order_timeout",
		Body:  msg.Body,
	}
	msg2.WithDelayTimeLevel(4) //跟普通比就多一句这个
	_, err = p.SendSync(context.Background(), msg2)
	if err != nil {
		//os.Exit(1)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延迟消息失败"
		zap.S().Errorf("发送延迟消息失败: %s", err.Error())
		return primitive.CommitMessageState //创建失败立即回滚
	}
	//提交事务
	tx.Commit()
	return primitive.RollbackMessageState
}

// When no response to prepare(half) message. broker will send check message to check the transaction status, and this
// method will be invoked to get local transaction status.
// 回调检查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("消息会查")
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	if res := global.DB.Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).Find(&orderInfo); res.RowsAffected == 0 {
		return primitive.CommitMessageState //回滚，并不能说明已经扣减库存，需要做好幂等性
	}
	return primitive.RollbackMessageState
}

//购物车列表
func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse

	if res := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); res.Error != nil {
		return nil, res.Error
	} else {
		rsp.Total = int32(res.RowsAffected)
	}

	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return &rsp, nil
}

//购物车添加商品
func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车 1. 购物车中原本没有这件商品 - 新建一个记录 2. 这个商品之前添加到了购物车- 合并
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		//如果记录已经存在，则合并购物车记录, 更新操作
		shopCart.Nums += req.Nums
	} else {
		//插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

//更新购物车
func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	//更新或者插入
	var ShopCart *model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&ShopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	ShopCart.Checked = req.Checked
	if req.Nums > 0 {
		ShopCart.Nums = req.Nums
	}

	global.DB.Save(&ShopCart)

	return &emptypb.Empty{}, nil
}

//删除商品
func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var ShopCart model.ShoppingCart
	if result := global.DB.Where("user=? and goods=?", req.UserId, req.GoodsId).Delete(&ShopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

//创建订单
func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单
			1. 从购物车中获取到选中的商品
			2. 商品的价格自己查询 - 访问商品服务 (跨微服务)
			3. 库存的扣减 - 访问库存服务 (跨微服务)
			4. 订单的基本信息表 - 订单的商品信息表
			5. 从购物车中删除已购买的记录
	*/

	orderListener := OrderListener{ctx: ctx}
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"43.242.33.9:9876"})),
		producer.WithRetry(1),
	)

	if err != nil {
		zap.S().Errorf("生成 producer 失败: %s", err.Error())
		return nil, err
	}

	err = p.Start()

	if err != nil {
		zap.S().Errorf("启动 producer 失败: %s", err.Error())
		//os.Exit(1)
		return nil, err
	}

	topic := "order_reback"

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
		Status:       "PAYING",
	}
	jsonString, _ := json.Marshal(order)
	msg := &primitive.Message{
		Topic: topic,
		Body:  jsonString,
	}
	_, err = p.SendMessageInTransaction(context.Background(), msg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "发送消息失败")
	}

	if orderListener.Code != codes.OK {
		return nil, status.Errorf(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{Id: orderListener.ID, OrderSn: order.OrderSn, Total: orderListener.OrderAmount}, nil

}

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	//分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	//这个订单的id是否是当前用户的订单， 如果在web层用户传递过来一个id的订单， web层应该先查询一下订单id是否是当前用户的
	//在个人中心可以这样做，但是如果是后台管理系统，web层如果是后台管理系统 那么只传递order的id，如果是电商系统还需要一个用户的id
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Nums,
		})
	}

	return &rsp, nil
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	//先查询，再更新 实际上有两条sql执行， select 和 update语句
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}

func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msgs {
		var orderInfo model.OrderInfo

		_ = json.Unmarshal(msgs[i].Body, &orderInfo)
		if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
			zap.S().Info("当前订单已处理完成")
		}
		//判断当前订单状态是否为未支付，如果是执行超时机制
		if orderInfo.Status == "PAYING" {
			tx := global.DB.Begin()

			orderInfo.Status = "TRADE_CLOSED"
			tx.Save(&orderInfo)

			p, _ := rocketmq.NewProducer(
				producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"43.242.33.9:9876"})),
				producer.WithRetry(2),
			)
			err := p.Start()
			if err != nil {
				zap.S().Errorf("start producer error: %s", err.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			msg := &primitive.Message{
				Topic: "order_reback",
				Body:  msgs[i].Body,
			}
			_, err = p.SendSync(context.Background(), msg)

			if err != nil {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}

			tx.Commit()

		}
	}
	return consumer.ConsumeSuccess, nil
}
