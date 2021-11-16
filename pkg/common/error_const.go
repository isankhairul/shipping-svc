package common

const (
	ERRCODE_SUCCESS       int = 1
	ERRCODE_VALIDATE      int = 2
	ERRCODE_RESTCLIENT    int = 2
	ERRCODE_DB            int = 2
	ERRCODE_NOWITHDATA    int = 2
	ERRCODE_REDIRECTLOGIN int = 4
	ERRCODE_WITHDATA      int = 3
	ERRCODE_ALREADYEXISTS int = 300
	ERRCODE_BADROUTING    int = 3
	RESCODE_WARNING       int = 3
	RESCODE_OTP_EXPIRED   int = 5
)

const (
	ERRMSG_SUCCESS                    string = "success"
	ERRMSG_NOTFOUND                   string = "not found"
	ERRMSG_ALREADYEXISTS              string = "already exists"
	ERRMSG_BADROUTING                 string = "inconsistent mapping between route and handler"
	ERRMSG_INVALIDREQUEST             string = "Request tidak valid"
	ERRMSG_INVALID_REQUEST_CHANGE_MDN string = "Request Anda sudah tidak valid, silahkan request ulang untuk ganti MDN"
	ERRMSG_DB                         string = "ERR from db query data"
	ERRMSG_CLIENT                     string = "ERRR client rest api"
	ERRMSG_PUBLIC                     string = "Ada kesalahan, silahkan coba kembali"
	ERRMSG_PUBLIC_WITH_LOGID		  string = "Ada kesalahan, silahkan coba kembali ERR: %s"
	ERRMSG_REDIRECT                   string = "Sihlakan Login Ulang kembali"
	ERRMSG_OLD_PASSWORD               string = "Password lama tidak sesuai"
	ERRMSG_NEW_PASSWORD               string = "Password baru tidak boleh sama"
	ERRMSG_OTP_EXPIRED                string = "Kode OTP yang anda masukan sudah expired, silahkan request ulang OTP"
)

//Voucher
const (
	ERRMSG_VOUCHER               string = "Voucher Nominal %s Tersedia"
	ERRMSG_VOUCHERUSED           string = "Voucher Sudah pernah di pakai"
	ERRMSG_VOUCHERNOTEXIST       string = "Voucher Belum Tersedia"
	ERRMSG_VOUCHERMDNNOTEXIST    string = "Nomor handphone tidak terdaftar/sudah terminate"
	ERRMSG_VOUCHERMDNNOTACTIVE   string = "Nomor handphone belum aktif"
	ERRMSG_VOUCHERMDNSUSPEND     string = "Nomor handphone pelanggan dalam kondisi suspend"
	ERRMSG_VOUCHERSUCCESRECHARGE string = "Pengisian Voucher data ke Nomor %s sebesar %s berhasil. Aktif sampai dengan %s"
	ERRMSG_POSTPAID              string = "Pengisian tidak dapat dilakukan untuk nomor pascabayar. Silakan hubungi gallery smartfren terdekat."
)

//package
const (
	ERRMSG_PACKAGEMDNNOTEXIST  string = "Nomor handphone tidak terdaftar/sudah terminate"
	ERRMSG_PACKAGEMDNNOTACTIVE string = "Nomor handphone belum aktif"
	ERRMSG_PACKAGEMDNSUSPEND   string = "Nomor handphone pelanggan dalam kondisi suspend"
	ERRMSG_PACKAGENOTEXIST     string = "Package Tidak Tersedia"
	ERRMSG_PACKAGEGOON         string = "Pembelian paket gagal karena Anda masih berlangganan %s"
)

//Pulsa
const (
	ERRMSG_PULSAMDNNOTEXIST     string = "Nomor tidak terdaftar/sudah terminate"
	ERRMSG_PULSARMDNNOTACTIVE   string = "Nomor handphone pelanggan belum aktif"
	ERRMSG_PULSAMDNSUSPEND      string = "Nomor handphone pelanggan dalam kondisi suspend"
	ERRMSG_PULSAPACKAGENOTEXIST string = "Package Tidak Tersedia"
	ERRMSG_ITEM_NOT_FOUND       string = "Item tidak ditemukan"
)

// NOO
const (
	ELOAD_NULL              string = "Nomor Eload Kosong"
	USER_ID_NULL            string = "User ID Kosong"
	ELOAD_NOT_FOUND         string = "Nomor Eload tidak terdaftar"
	MDN_NOT_FOUND           string = "Nomor MDN tidak terdaftar"
	OWNER_IMAGE_NULL        string = "Foto Pemilik Kosong"
	OWNER_ID_IMAGE_NULL     string = "Foto KTP Kosong"
	OUTLET_IMAGE_NULL       string = "Foto Outlet Kosong"
	OUTLET_NAME_NULL        string = "Outlet Name Kosong"
	OWNER_NAME_NULL         string = "Owner Name Kosong"
	USER_NOT_FOUND          string = "User Tidak Ditemukan"
	CLUSTER_ID_ERR          string = "Cluster ID Error"
	CLUSTER_NAME_ERR        string = "Cluster Name Error"
	CLUSTER_REGION_ERR      string = "Cluster Region Error"
	REGION_NAME_ERR         string = "Region Name Error"
	USER_GROUP_ERR          string = "User Group Error"
	DELETE_NOTOFICATION_ERR string = "Error Delete Notification"
	COUNT_PENDING_ERR       string = "Error Menghitung Pending Notif Notification"
	PARRENT_OUTLET_ERR      string = "Error Mengambil Parrent Outlet"
	ELOAD_EXISTS            string = "No Eload sudah terdaftar"
	BLACKLIST_ELOAD         string = "Nomor eload dalam daftar blacklist"
)

// API PESAN ERR
const (
	LOGIN_FAILED              string = "Login Gagal"
	ERRMSG_IMSI_NULL          string = "IMSI tidak terbaca"
	ERRMSG_IMEI_NULL          string = "IMEI tidak terbaca"
	ERRMSG_USERDISABLED       string = "Status akun anda tidak aktif, silahkan hubungi Administrator."
	ERRMSG_LOGINNOTALLOWED    string = "Anda tidak diizinkan login pada waktu ini."
	ERRMSG_WRONGPASSWORD      string = "password salah"
	ERRMSG_PASSWORDVALIDATION string = "Password harus memiliki 8 karakter berisi huruf kecil, huruf besar, angka dan simbol"
)

// API PESAN ERR
const (
	INVALID_USER_ID         string = "User ID Tidak Ditemukan"
	INVALID_MESSAGE_TO      string = "Penerima Pesan Kosong"
	INVALID_MESSAGE_SUBJECT string = "Subject Pesan Kosong"
	INVALID_MESSAGE_CONTENT string = "Konten Pesan Kosong"
	INVALID_MESSAGE_SENDER  string = "Pengirim Pesan Kosong"
	INVALID_PARTNER_ID      string = "Partner ID Tidak Ditemukan"
	INVALID_DELETE_MESSAGE  string = "Gagal Menghapus Pesan"
	SUCCESS_DELETE_MESSAGE  string = "Berhasil Menghapus Pesan"
	SUCCESS_SENT_MESSAGE    string = "Berhasil Kirim Pesan"
	PARAM_NOT_COMPLETE      string = "Parameter Tidak Lengkap"
)

// API TRANSAKSI PREPAID
const (
	ERR_COUNT_MDN             string = "Ada kesalahan, silahkan coba kembali"
	INVALID_MANDATORY         string = "Parameter tidak lengkap"
	UNAUTHORIZE_IP_ADDRESS    string = "Unauthorize IP Address"
	ERR_RETURN_CODE           string = "Err Return Code"
	BLACKLIST_CARD_NUMBER     string = "Nomor Identitas tidak diizinkan untuk melakukan registrasi"
	GETSUBINFO_ERR            string = "Ada kesalahan, silahkan coba kembali"
	FAIL_GET_LIMIT_MDN        string = "Ada kesalahan, silahkan coba kembali"
	FAIL_GET_MAX_LIMIT_MDN    string = "Ada kesalahan, silahkan coba kembali"
	ERR_OVER_MASTER_DATA      string = "NIK sudah terdaftar melebihi batas maksimal"
	ERR_CITIZENSHIP           string = "Error Citizenship"
	ERR_TIBCO                 string = "Error Tibco"
	ERR_COUNT_NIK             string = "NIK atau KK Harus 16 Karakter"
	ERROR_PREPAID             string = "Ada kesalahan, silahkan coba kembali"
	FAIL_REGISTRATION         string = "Gagal Mendaftarkan Mdn"
	BLACKLIST_IDENTITY_NUMBER string = " Nomor Identitas tidak diizinkan untuk melakukan registrasi"
	ERR_INVALID_CHARACTERS    string = "Nomor Identitas Harus Angka Dan 16 Digit"
	NULL_DATA                 string = "Data Kosong"
	MITRA_IS_NOT_ALLOWED      string = "Outlet ini tidak dapat melakukan registrasi WNA"
	ITEM_DATA_NULL            string = "item untuk no mdn tersebut tidak ditemukan"
)

const (
	ERRMSG_SELLOUTPARAMSICCID string = "Range Of iccid [1-2000]"
)

const (
	ORDER_STOCK_ERR_PARAM string = "Parameter tidak lengkap"
)

const (
	ACTION_NOO             = "noo"
	ACTION_LOGIN           = "login"
	ACTION_CREATE_PASSWORD = "create_pwd"
	ACTION_FORGOT_PASSWORD = "forgot_pwd"
)
