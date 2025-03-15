package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
)

type CreateAccountStruct struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"` //currency is nothing but custom-defined gin validatior
}

func (server *Server) createAccount(ctx *gin.Context) {
    var req CreateAccountStruct
    if err := ctx.ShouldBindJSON(&req); err!=nil {
        ctx.JSON(http.StatusBadRequest,errorResponse(err))
        return 
    }


    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
    arg := db.CreateAccountParams{
        Owner: authPayload.Username,
        Currency: req.Currency,
        Balance: 0,
    }

    account, err := server.store.CreateAccount(ctx,arg)
    if err!=nil{
        if pqErr, ok := err.(*pq.Error);ok{
            switch pqErr.Code.Name(){
            case "foreign_key_violation","unique_violation":
                ctx.JSON(http.StatusForbidden,errorResponse(err))
                return
            }
        }
        ctx.JSON(http.StatusInternalServerError,errorResponse(err))
        return
    }
    ctx.JSON(http.StatusOK,account)
}   

type getAccountRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"` //min=1 (value must be greater than 1)
}

func (server *Server) getAccount(ctx *gin.Context) {
    var req getAccountRequest

    if err := ctx.ShouldBindUri(&req); err!=nil {
        ctx.JSON(http.StatusBadRequest,errorResponse(err))
        return 
    }

   account,err := server.store.GetAccount(ctx,req.ID)
   if err != nil {
     if err == sql.ErrNoRows{  //If the entered ID does not found in DB, it will throw "sql: no rows in result set"
        ctx.JSON(http.StatusNotFound,errorResponse(err)) 
        return
     }
     ctx.JSON(http.StatusInternalServerError,errorResponse(err))
     return
   }
   authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload) //Authorization
   if account.Owner != authPayload.Username{
    err := errors.New("account doesnt belong to the authenticated user")
    ctx.JSON(http.StatusUnauthorized,errorResponse(err))
    return
   }
   ctx.JSON(http.StatusOK,account)
}

type listAccountRequest struct {
    PageID int32 `form:"page_id" binding:"required,min=1"` //page number 
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"` //page size
}

func (server *Server) listAllAccounts(ctx *gin.Context) {
    var req listAccountRequest

    if err := ctx.ShouldBindQuery(&req); err!=nil {
        ctx.JSON(http.StatusBadRequest,errorResponse(err))
        return 
    }

    authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload) //Checks for authorization
    args := db.ListAccountsParams{
        Owner: authPayload.Username,
        Limit: req.PageID,
        Offset: (req.PageID-1)*req.PageSize,
    }
   accounts,err := server.store.ListAccounts(ctx,args)
   if err != nil {
     ctx.JSON(http.StatusInternalServerError,errorResponse(err))
     return
   }

   ctx.JSON(http.StatusOK,accounts)
}

//Deleting acc from the DB using endpoint
type deleteReqParams struct{
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
    var req deleteReqParams

    if err := ctx.ShouldBindUri(&req); err!=nil {
        ctx.JSON(http.StatusBadRequest,errorResponse(err))
        return 
    }

   err := server.store.DeleteAccount(ctx,req.ID)
   if err != nil {
     if err == sql.ErrNoRows{  //If the entered ID does not found in DB, it will throw "sql: no rows in result set"
        ctx.JSON(http.StatusNotFound,errorResponse(err)) 
        return
     }
     ctx.JSON(http.StatusInternalServerError,errorResponse(err))
     return
   }

   ctx.JSON(http.StatusOK,"Account deleted successfully")
}