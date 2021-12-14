package message

// Message wrapper.
type Message struct {
	Code    int `json:"code"`
	Message string `json:"message"`
}

//Code 34000 - 3499 User & Auth bn;
var TelErrUserNotFound = Message{Code: 34000, Message: "Pengguna tidak ditemukan pada sistem"}
var TelErrUserRoleExists = Message{Code: 34001, Message: "User sudah di assign untuk role ini. Mohon periksa kembali data Anda"}
var TelErrNonActiveUser = Message{Code: 34002, Message: "Pengguna tidak terdaftar Pada Sistem"}
var TelErrPassword = Message{Code: 34003, Message: "Kata sandi yang di masukkan tidak sesuai"}
var TelErrUserSave = Message{Code: 34004, Message: "Pengguna tidak dapat di simpan. Mohon periksa kembali data Anda"}
var TelErrCodeNotValid = Message{Code: 34005, Message: "Activation Code Tidak Valid"}
var TelErrUserIsActive = Message{Code: 34006, Message: "User Sudah Aktif!"}
var TelErrEmailAlreadyUsed = Message{Code: 34007, Message: "Email sudah ada yang menggunakan"}
var TelErrUserNonActive = Message{Code: 34008, Message: "User belum aktif"}
var TelErrUserUnAuthorize = Message{Code: 34009, Message: "Unauthorized Access"}

// Code 39000 - 39999 Server error
var TelErrRevocerRoute = Message{Code: 39000, Message: "Terjadi kesalahan routing"}
var TelErrPageNotFound = Message{Code: 39404, Message: "Halaman Tidak ditemukan"}
var SuccessMsg = Message{Code: 1111, Message: "Success"}
var FailedMsg = Message{Code: 0000, Message: "Failed"}
var ErrReqParam = Message{Code: 4000, Message: "Invalid Request Parameter"}

const (
	CODE_SUCCESS           = 1000
	CODE_ERR_VALIDATE      = 2000
	CODE_ERR_DB            = 2100
	CODE_MSG_REDIRECTLOGIN = 3000
	CODE_ERR_BADREQUEST    = 4000
	CODE_ERR_NOTFOUND      = 5000
	CODE_ERR_BADROUTING    = 5100
	CODE_ERR_NOAUTH        = 6000
)

//General error/desc message
const (
	MSG_SUCCESS         = "Success"
	MSG_NOTFOUND        = "Not found"
	MSG_ALREADYEXISTS   = "Data already exists"
	MSG_BADROUTING      = "Inconsistent mapping between route and handler"
	MSG_INVALID_REQUEST = "Invalid request parameter(s)"
	MSG_INTERNAL_ERROR  = "Error has been occured while processing request"
	MSG_NO_AUTH         = "No Authorization"
	MSG_INVALID_HEADER  = "Invalid header"
)

//General error/desc database message
const (
	MSG_ERR_DB        = "Error has been occured while processing database request"
	MSG_NO_DATA       = "Data is not found"
	MSG_ERR_SAVE_DATA = "Data cannot be saved, please check your request"
	MSG_ERR_DEL_DATA  = "Data cannot be deleted, please check your request"
)

//General error/desc field validation Message
const (
	MSG_ERR_CHAR_LENGTH = "Character length from %v until %v"
	MSG_ERR_REQUIRED    = "Required field"
)
