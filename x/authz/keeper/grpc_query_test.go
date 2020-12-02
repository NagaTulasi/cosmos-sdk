package keeper_test

import (
	gocontext "context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz/types"
)

func (suite *TestSuite) TestGRPCQueryAuthorization() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs
	var (
		req              *types.QueryAuthorizationRequest
		expAuthorization types.Authorization
	)
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		postTest func(res *types.QueryAuthorizationResponse)
	}{
		{
			"fail invalid granter addr",
			func() {
				req = &types.QueryAuthorizationRequest{}
			},
			false,
			func(res *types.QueryAuthorizationResponse) {},
		},
		{
			"fail invalid grantee addr",
			func() {
				req = &types.QueryAuthorizationRequest{
					Granter: addrs[0].String(),
				}
			},
			false,
			func(res *types.QueryAuthorizationResponse) {},
		},
		{
			"fail invalid msg-type",
			func() {
				req = &types.QueryAuthorizationRequest{
					Granter: addrs[0].String(),
					Grantee: addrs[1].String(),
				}
			},
			false,
			func(res *types.QueryAuthorizationResponse) {},
		},
		{
			"Succcess",
			func() {
				now := ctx.BlockHeader().Time
				newCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 100))
				expAuthorization = &types.SendAuthorization{SpendLimit: newCoins}
				app.MsgAuthKeeper.Grant(ctx, addrs[0], addrs[1], expAuthorization, now.Add(time.Hour))
				req = &types.QueryAuthorizationRequest{
					Granter: addrs[1].String(),
					Grantee: addrs[0].String(),
					MsgType: expAuthorization.MethodName(),
				}
			},
			true,
			func(res *types.QueryAuthorizationResponse) {
				var auth types.Authorization
				err := suite.app.InterfaceRegistry().UnpackAny(res.Authorization.Authorization, &auth)
				suite.Require().NoError(err)
				suite.Require().NotNil(auth)
				suite.Require().Equal(auth.String(), expAuthorization.String())
			},
		},
	}
	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()
			result, err := queryClient.Authorization(gocontext.Background(), req)
			if testCase.expPass {
				fmt.Println(result)
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
			testCase.postTest(result)
		})
	}
}

func (suite *TestSuite) TestGRPCQueryAuthorizations() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs
	var (
		req              *types.QueryAuthorizationsRequest
		expAuthorization types.Authorization
	)
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		postTest func(res *types.QueryAuthorizationsResponse)
	}{
		{
			"fail invalid granter addr",
			func() {
				req = &types.QueryAuthorizationsRequest{}
			},
			false,
			func(res *types.QueryAuthorizationsResponse) {},
		},
		{
			"fail invalid grantee addr",
			func() {
				req = &types.QueryAuthorizationsRequest{
					Granter: addrs[0].String(),
				}
			},
			false,
			func(res *types.QueryAuthorizationsResponse) {},
		},
		{
			"Succcess",
			func() {
				now := ctx.BlockHeader().Time
				newCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 100))
				expAuthorization = &types.SendAuthorization{SpendLimit: newCoins}
				app.MsgAuthKeeper.Grant(ctx, addrs[0], addrs[1], expAuthorization, now.Add(time.Hour))
				req = &types.QueryAuthorizationsRequest{
					Granter: addrs[1].String(),
					Grantee: addrs[0].String(),
				}
			},
			true,
			func(res *types.QueryAuthorizationsResponse) {
				var auth types.Authorization
				suite.Require().Equal(1, len(res.Authorizations))
				err := suite.app.InterfaceRegistry().UnpackAny(res.Authorizations[0].Authorization, &auth)
				suite.Require().NoError(err)
				suite.Require().NotNil(auth)
				suite.Require().Equal(auth.String(), expAuthorization.String())
			},
		},
	}
	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()
			result, err := queryClient.Authorizations(gocontext.Background(), req)
			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
			testCase.postTest(result)
		})
	}
}
