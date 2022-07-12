package forms


type SmsSendForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,min=3,max=11,mobile"`
	Type int `form:"type" json:"type" binding:"required,oneof=1 2"`
}