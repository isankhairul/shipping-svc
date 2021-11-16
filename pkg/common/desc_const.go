package common

// Stock In Item Detil
const (
	SCAN_TYPE_SINGLE            string = "single"
	SCAN_TYPE_RANGE             string = "range"
	SCAN_TYPE_BOX               string = "box"
	SCAN_TYPE_SINGLE_ICCID      string = "SINGLE_ICCID"
	SCAN_TYPE_SINGLE_BOX_NUMBER string = "SINGLE_BOX_NUMBER"
)

const (
	INVALID_SCAN_TYPE string = "Param Scan Type tidak valid"
	//INVALID_ITEM              string = "Param Item tidak valid"
	INVALID_ICCID             string = "Param Iccid tidak valid"
	INVALID_BOXNUMBER         string = "Param Box Number tidak valid"
	INVALID_LIMIT_RANGE_ICCID string = "Invalid range iccid from limit"
)
