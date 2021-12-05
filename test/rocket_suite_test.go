// +acceptance

package test

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/wfen/go-grpc-services-course/proto/rocket/v1"
)

type RocketTestSuite struct {
	suite.Suite
}

func (s *RocketTestSuite) TestAddRocket() {
	s.T().Run("adds a new rocket succesfully", func(t *testing.T) {
		client := GetClient()
		resp, err := client.AddRocket(
			context.Background(),
			&rocket.AddRocketRequest{
				Rocket: &rocket.Rocket{
					Id:   "c2443bba-99de-485b-858e-92403d029ab8",
					Name: "V1",
					Type: "Falcon Heavy",
				},
			})
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "c2443bba-99de-485b-858e-92403d029ab8", resp.Rocket.Id)
	})

	s.T().Run("validates the id in the new rocket is a uuid", func(t *testing.T) {
		client := GetClient()
		_, err := client.AddRocket(
			context.Background(),
			&rocket.AddRocketRequest{
				Rocket: &rocket.Rocket{
					Id:   "not-a-valid-uuid",
					Name: "V1",
					Type: "Falcon Heavy",
				},
			})
		assert.Error(s.T(), err)
		st := status.Convert(err)
		assert.Equal(s.T(), codes.InvalidArgument, st.Code())
	})
}

func TestRocketService(t *testing.T) {
	suite.Run(t, new(RocketTestSuite))
}
