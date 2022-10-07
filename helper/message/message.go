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
var ErrUnAuth = Message{Code: 34004, Message: "Unauthorized"}
var ErrInvalidHeader = Message{Code: 34005, Message: "Invalid header"}
var ErrDB = Message{Code: 34005, Message: "Error has been occured while processing database request"}
var ErrNoData = Message{Code: 34005, Message: "Data is not found"}
var ErrSaveData = Message{Code: 34005, Message: "Data cannot be saved, please check your request"}
var ErrImportData = Message{Code: 34005, Message: "Data cannot be saved or updated, please check your import file"}
var ErrReq = Message{Code: 34005, Message: "Required field"}
var ErrCourierServiceID = Message{Code: 34007, Message: "courier_service_uid required"}
var ErrCourierID = Message{Code: 34006, Message: "courier_id required"}
var ErrChannelID = Message{Code: 34007, Message: "channel_id required"}
var ErrChannelCourierID = Message{Code: 34007, Message: "channel_courier_id required"}
var ErrPrioritySort = Message{Code: 34008, Message: "prioriy_sort required"}

// channel-courier 344xx
var ErrChannelCourierFound = Message{Code: 34401, Message: "Channel courier existed."}
var ErrChannelCourierNotFound = Message{Code: 34402, Message: "Channel courier not found."}
var ErrCourierServiceNotMatch = Message{Code: 34402, Message: "Unable to use Courier Service"}
var ErrDuplicatedCourier = Message{Code: 34010, Message: "Duplicated courier"}
var ErrUnableToDeleteChannelCourier = Message{Code: 34406, Message: "Unable to delete channel courier"}
var ErrChannelCourierServiceCreateFailed = Message{Code: 34606, Message: "Unable to create channel courier service"}

// courier: 341xx
var ErrCourierNotFound = Message{Code: 34101, Message: "Courier not found"}
var ErrCourierServiceNotFound = Message{Code: 34101, Message: "Courier service not found"}
var ErrCourierServiceIsRequired = Message{Code: 34101, Message: "Courier service uid is required"}
var ErrDataCourierUIdNotExist = Message{Code: 34001, Message: "CourierUID not exist"}

// courier coverage code : 346xx
var ErrCourierCoverageCodeUidNotExist = Message{Code: 34601, Message: "Courier Coverage Code not exist"}
var ErrCourierCoverageCodeUidExist = Message{Code: 34602, Message: "Courier Coverage Code exists"}
var ErrCourierCoverageCodeExist = Message{Code: 34005, Message: "The combination of courier_uid, country_code postal_code, and subdistrict is exist in database"}
var ErrShippingTypeRequired = Message{Code: 34602, Message: "shipping type is required"}

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

var (
	ErrCourierHasChildCourierService  = Message{Code: 209002, Message: "Can not delete Courier that has one or more Courier Service(s)"}
	ErrCourierHasChildCourierCoverage = Message{Code: 209003, Message: "Can not delete Courier that has one or more Courier Coverage(s)"}
	ErrCourierHasChildChannelCourier  = Message{Code: 209004, Message: "Can not delete Courier that has already assigned to Channel"}
	ErrCourierHasChildShippingStatus  = Message{Code: 209005, Message: "Can not delete Courier that has one or more Shipping Status"}
	ErrCourierServiceHasAssigned      = Message{Code: 209006, Message: "Can not delete Courier Service that has already assigned to Channel"}
	ErrChannelHasCourierAssigned      = Message{Code: 209007, Message: "Can not delete Channel that has already assigned to Courier"}
	ErrChannelHasChildShippingStatus  = Message{Code: 209008, Message: "Can not delete Channel that has one or more Shipping Status"}
	ErrChannelCourierHasChild         = Message{Code: 209009, Message: "Can not delete Channel Courier that has one or more Channel Courier Service(s)"}
)

var ErrShippingRateNotFound = Message{Code: 34602, Message: "shipping rate not found"}
var ErrInvalidCourierType = Message{Code: 34602, Message: "courier type is not valid"}
var ErrInvalidCourierCode = Message{Code: 34602, Message: "courier code is not valid"}
var ErrCreateOrder = Message{Code: 34602, Message: "failed when trying to create order"}
var ErrGetPickUpTimeslot = Message{Code: 34602, Message: "failed when trying to get pickup timeslot"}
var ErrCreatePickUpOrder = Message{Code: 34602, Message: "failed when trying to create pickup order"}
var ErrOrderShippingNotFound = Message{Code: 34602, Message: "order shipping not found"}
var ErrGetOrderDetail = Message{Code: 34602, Message: "failed when trying to get order detail"}
var ErrOrderBelongToAnotherChannel = Message{Code: 34602, Message: "the order is belong to another channel"}
var ErrChannelUIDRequired = Message{Code: 34602, Message: "channel_uid is required"}
var ErrSaveOrderShipping = Message{Code: 34602, Message: "failed when trying to save order shipping"}
var ErrFormatDateYYYYMMDD = Message{Code: 34602, Message: "date format must be YYYY-MM-DD"}
var ErrCancelPickup = Message{Code: 34602, Message: "failed when trying to cancel pickup"}
var ErrCantCancelOrderShipping = Message{Code: 34602, Message: "can't cancel this order"}
var ErrCantCancelOrderCourierService = Message{Code: 34602, Message: "courier service is not cancelable"}
var ErrUpdateOrderShipping = Message{Code: 34602, Message: "error update order shipping"}

var (
	ShippingProviderMsg               = Message{Code: 209002, Message: ""}
	CourierNotActiveMsg               = Message{Code: 209002, Message: "courier is not active"}
	CourierServiceNotActiveMsg        = Message{Code: 209002, Message: "courier service is not active"}
	ChannelCourierNotActiveMsg        = Message{Code: 209002, Message: "channel courier is not active"}
	ChannelCourierServiceNotActiveMsg = Message{Code: 209002, Message: "channel courier service is not active"}
	CourierHiddenInPurposeMsg         = Message{Code: 209002, Message: "courier is hidden in purpose"}
	PrescriptionNotAllowedMsg         = Message{Code: 209002, Message: "prescription is not allowed"}
	WeightExceedsmsg                  = Message{Code: 209002, Message: "final weight exceeds the maximum weight allowed"}
	InvalidCourierTypeMsg             = Message{Code: 209002, Message: "courier type is not valid"}
	InvalidCourierCodeMsg             = Message{Code: 209002, Message: "courier code is not valid"}
	CourierServiceNotFoundMsg         = Message{Code: 209002, Message: "courier service is not found"}
	ChannelNotFoundMsg                = Message{Code: 209002, Message: "channel not found"}
	WeightExceedsMsg                  = Message{Code: 209002, Message: "final weight exceeds the maximum weight allowed"}
	ShippingStatusNotFoundMsg         = Message{Code: 209002, Message: "shipping status not found"}
	OrderNoAlreadyExistsMsg           = Message{Code: 209002, Message: "order no already exists"}
	OriginNotFoundMsg                 = Message{Code: 209002, Message: "origin is not in courier coverage"}
	DestinationNotFoundMsg            = Message{Code: 209002, Message: "destination not in courier coverage"}
	RequestPickupHasBeenMadeMsg       = Message{Code: 209002, Message: "request pickup has been made"}
	OrderHasBeenCancelledMsg          = Message{Code: 209002, Message: "order has been cancelled"}
)
