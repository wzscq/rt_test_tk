package robot

type ServerHeader struct {
	Token string  `json:"token"`
}

type UploadMapReq struct {
	RobotId string `form:"robotId"`
	PictureId string	`form:"pictureId"`
	PictureName string	`form:"pictureName"`
	BuildingCode string	`form:"buildingCode"`
	Floor string	`form:"floor"`
}