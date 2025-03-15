package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/db/util"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
func TestGetUserAPI(t *testing.T) {

	user,password := randomUser(t)

	testCases := []struct{
		name string
		body gin.H
		buildStatus func(store *mockdb.MockStore)
		checkResponse func(t *testing.T,recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body : gin.H{
				"username": user.Username,
				"password": password,
				"full_name":user.FullName,
				"email":user.Email,
			},
			buildStatus: func(store *mockdb.MockStore){
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email: user.Email,
				}
				store.EXPECT().
			CreateUser(gomock.Any(),EqCreateUserParams(arg,password)).
			Times(1).
			Return(user,nil)
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusOK,recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body : gin.H{
				"username": "",
				"password": password,
				"full_name":user.FullName,
				"email":user.Email,
			},
			buildStatus: func(store *mockdb.MockStore){
				store.EXPECT().
				CreateUser(gomock.Any(),gomock.Any()).
			Times(0)
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusBadRequest,recorder.Code)
			} ,
		},
		{
			name: "InternalError",
			body : gin.H{
				"username": user.Username,
				"password": password,
				"full_name":user.FullName,
				"email":user.Email,
			},
			buildStatus: func(store *mockdb.MockStore){
				store.EXPECT().
				CreateUser(gomock.Any(),gomock.Any()).
			Times(1).
			Return(db.User{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			} ,
		},
		{
			name: "DuplicateUsername",
			body : gin.H{
				"username": user.Username,
				"password": password,
				"full_name":user.FullName,
				"email":user.Email,
			},
			buildStatus: func(store *mockdb.MockStore){
				store.EXPECT().
				CreateUser(gomock.Any(),gomock.Any()).
			Times(1).
			Return(db.User{}, &pq.Error{Code:"23505"})
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusForbidden,recorder.Code)
			} ,
		},
		{
			name: "InvalidEmail",
			body : gin.H{
				"username": user.Username,
				"password": password,
				"full_name":user.FullName,
				"email":"Invalid-mail",
			},
			buildStatus: func(store *mockdb.MockStore){
				store.EXPECT().
				CreateUser(gomock.Any(),gomock.Any()).
			Times(0)
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusBadRequest,recorder.Code)
			} ,
		},
		{
			name: "ShortPassword",
			body : gin.H{
				"username": user.Username,
				"password": "123",
				"full_name":user.FullName,
				"email":user.Email,
			},
			buildStatus: func(store *mockdb.MockStore){
				store.EXPECT().
				CreateUser(gomock.Any(),gomock.Any()).
			Times(0)
			},
			checkResponse: func(t *testing.T,recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusBadRequest,recorder.Code)
			} ,
		},
	}
	
	for i := range testCases{
	
		tc := testCases[i];
		
		t.Run(tc.name,func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tc.buildStatus(store)

		//start test server and send request
		server := newTestServer(t,store)
		recorder := httptest.NewRecorder()

		data,err := json.Marshal(tc.body)
		require.NoError(t,err)

		url := "/users"
		req, err:= http.NewRequest(http.MethodPost,url,bytes.NewReader(data))
		require.NoError(t,err)

		server.router.ServeHTTP(recorder,req)
		tc.checkResponse(t,recorder)
		})
		
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.Password(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwnerName(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwnerName(),
		Email:          util.RandomEmail(),
	}
	return
}
