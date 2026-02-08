package request

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type RejectLoanRequest struct {
    Reason string `json:"reason"`
}

func (r *RejectLoanRequest) BindAndValidate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        return errors.New("payload tidak valid")
    }

    return validation.ValidateStruct(r,
        validation.Field(&r.Reason,
            validation.Required.Error("alasan penolakan wajib diisi"),
            validation.Length(5, 255).Error("alasan penolakan 5–255 karakter"),
        ),
    )
}


type RequestRevisionRequest struct {
    Note string `json:"note"`
}

func (r *RequestRevisionRequest) BindAndValidate(c *gin.Context) error {
    if err := c.ShouldBindJSON(r); err != nil {
        return errors.New("payload tidak valid")
    }

    return validation.ValidateStruct(r,
        validation.Field(&r.Note,
            validation.Required.Error("catatan revisi wajib diisi"),
            validation.Length(5, 255).Error("catatan revisi 5–255 karakter"),
        ),
    )
}

type AdminLoanFilter struct {
    LoanStatus     string `form:"loan_status"`
    LoanFlowStatus string `form:"loan_flow_status"`
    UserID         uint   `form:"user_id"`
}

func (f *AdminLoanFilter) BindAndValidate(c *gin.Context) error {
    if err := c.ShouldBindQuery(f); err != nil {
        return errors.New("query parameter tidak valid")
    }

    return nil
}
