package forms

type ShopCartItemForm struct {
	GoodsId int32  `form:"goods_id" json:"goods_id" binding:"required"`
	Nums    int32  `form:"nums" json:"nums" binding:"required"`
}
type ShopCartItemUpdateForm struct {
	Nums    int32  `form:"nums" json:"nums" binding:"required"`
	Checked *bool `json:"checked"`
}