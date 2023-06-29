package routers

import (
	"server-template/internal/models/err_const"
	"strconv"
)

type Params map[string]string

func (params Params) GetId() (int, error) {
	id, ok := params["id"]
	if !ok {
		return 0, err_const.ErrIdMissing
	}

	idUint, err := strconv.Atoi(id)
	if err != nil || idUint < 0 {
		return 0, err_const.ErrIdValidate
	}

	return idUint, nil
}

func (params Params) GetUuid() (string, error) {
	uuidParam, ok := params["uuid"]
	if !ok {
		return "", err_const.ErrUuidMissing
	}

	if len(uuidParam) < 10 {
		return "", err_const.ErrUuidValidate
	}

	return uuidParam, nil
}
