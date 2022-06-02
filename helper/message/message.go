package message

// Message wrapper.
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var TelErrUserNotFound = Message{Code: 34000, Message: "Not found"}
var ErrDataExists = Message{Code: 34001, Message: "Data already exists"}
var ErrBadRouting = Message{Code: 34002, Message: "Inconsistent mapping between route and handler"}
var ErrInternalError = Message{Code: 34003, Message: "Error has been occured while processing request"}
var ErrNoAuth = Message{Code: 34004, Message: "No Authorization"}
var ErrInvalidHeader = Message{Code: 34005, Message: "Invalid header"}
var ErrDB = Message{Code: 34005, Message: "Error has been occured while processing database request"}
var ErrNoData = Message{Code: 34005, Message: "Data is not found"}
var ErrSaveData = Message{Code: 34005, Message: "Data cannot be saved, please check your request"}
var ErrImportData = Message{Code: 34005, Message: "Data cannot be saved or updated, please check your import file"}
var ErrReq = Message{Code: 34005, Message: "Required field"}
var ErrCourierID = Message{Code: 34006, Message: "courier_id required"}
var ErrChannelID = Message{Code: 34007, Message: "channel_id required"}
var ErrPrioritySort = Message{Code: 34008, Message: "prioriy_sort required"}

// channel-courier 344xx
var ErrChannelCourierFound = Message{Code: 34401, Message: "Channel courier existed."}
var ErrChannelCourierNotFound = Message{Code: 34402, Message: "Channel courier not found."}
var ErrDuplicatedCourier = Message{Code: 34010, Message: "Duplicated courier"}
var ErrUnableToDeleteChannelCourier = Message{Code: 34406, Message: "Unable to delete channel courier"}
var ErrChannelCourierServiceCreateFailed = Message{Code: 34606, Message: "Unable to create channel courier service"}

// courier: 341xx
var ErrCourierNotFound = Message{Code: 34101, Message: "Courier not found"}
var ErrDataCourierUIdNotExist = Message{Code: 34001, Message: "CourierUID not exist"}

// courier coverage code : 346xx
var ErrCourierCoverageCodeUidNotExist = Message{Code: 34001, Message: "Courier Coverage Code not exist"}

//channel : 342xx
var ErrChannelNotFound = Message{Code: 34201, Message: "Channel not found"}
var ErrDataCourierServiceUidNotExist = Message{Code: 34001, Message: "CourierSerivceUID not exist"}

// Shipment predefine section
var ErrShipmentPredefinedNotFound = Message{Code: 34501, Message: "Shipment predefined not found"}
var ErrDataCourierServiceExists = Message{Code: 34001, Message: "Data courier_id/shipping_code already exists"}
var ErrNoDataCourierService = Message{Code: 34005, Message: "Courier Service data not found"}
var ErrDataChannelExists = Message{Code: 34001, Message: "Data channel_code already exists"}
var ErrCourierServiceHasInvalidStatus = Message{Code: 34006, Message: "Courier Service has status = 0 "}

// Code 39000 - 39999 Server error
var ErrRevocerRoute = Message{Code: 39000, Message: "Terjadi kesalahan routing"}
var ErrPageNotFound = Message{Code: 39404, Message: "Halaman Tidak ditemukan"}
var SuccessMsg = Message{Code: 201000, Message: "Success"}
var FailedMsg = Message{Code: 0000, Message: "Failed"}
var ErrReqParam = Message{Code: 4000, Message: "Invalid Request Parameter(s)"}
