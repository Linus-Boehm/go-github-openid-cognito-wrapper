package github_openid_wrapper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetAuthorizeUrl(t *testing.T){
	req := &AuthorizeRequest{
		ClientID: "clientId",
		Scope: "openId user",
		ResponseType: "code",
		State: "state",
	}
	fmt.Println(req.GetAuthorizeUrl())
}

func TestUnmarshalMap(t *testing.T){
	queryParams := map[string]string{
		"client_id": "test",
		"scope": "openid user:email",
		"state": "ExampleState",
		"response_type": "authorization_code",
	}
	req := &AuthorizeRequest{}
	err := req.UnmarshalMap(queryParams)
	assert.NoError(t, err)
	fmt.Println(*req)
	iVal := reflect.ValueOf(req).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		mapKey := strings.Split(typ.Field(i).Tag.Get("json"),",")[0]
		mapVal, ok := queryParams[mapKey]
		assert.True(t, ok)
		assert.Equal(t, mapVal, fmt.Sprint(iVal.Field(i)))

	}
}

func TestMakeIDToken(t *testing.T){
	req := TokenResponse{
		AccessToken: "AccessToken",
		ExpiresAt: time.Now().AddDate(0,0,7).Unix(),
	}
	token, err := MakeIDToken(req)
	assert.NoError(t, err)
	fmt.Println(token)
}