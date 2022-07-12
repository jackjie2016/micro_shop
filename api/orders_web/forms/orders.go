package forms

type CreateOrderForm struct {
	//GoodsId int32  `form:"goods_id" json:"goods_id" binding:"required"`
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Mobile string `json:"mobile" binding:"required,min=3,max=11,mobile"`
	Address    string `json:"address" binding:"required,min=10,max=100"`
	Post    string `json:"post" binding:"required,min=10,max=100"`
}



