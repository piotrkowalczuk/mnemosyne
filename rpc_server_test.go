package mnemosyne

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestRPCServer_mockedStore(t *testing.T) {
	var (
		err         error
		suite       *integrationSuite
		expectedErr error
		subjectID   string
		bag         map[string]string
		session     *Session
		token       *AccessToken
	)

	suite = &integrationSuite{}
	suite.setup(t)

	Convey("RPCServer", t, func() {
		Convey("Start", func() {
			var (
				req *StartRequest
				res *StartResponse
			)

			expectedErr = nil
			subjectID = "subject_id"
			bag = map[string]string{"key": "value"}
			tk := NewAccessToken([]byte("key"), []byte("hash"))
			token = &tk

			itSuccess := func() {
				Convey("Not return any error", func() {
					So(err, ShouldBeNil)
				})
				Convey("Return session with same bag", func() {
					So(res.Session.Bag, ShouldResemble, req.Bag)
				})
				Convey("Return session with same subject id", func() {
					So(res.Session, ShouldNotBeNil)
					So(res.Session.SubjectId, ShouldEqual, req.SubjectId)
				})
				Convey("Return session with expire at timestamp", func() {
					So(res.Session.ExpireAt, ShouldNotBeNil)
					So(res.Session.ExpireAt.Nanos, ShouldNotEqual, 0)
					So(res.Session.ExpireAt.Seconds, ShouldNotEqual, 0)
				})
				Convey("Return session with token", func() {
					So(res.Session.AccessToken, ShouldBeValidToken)
				})
			}

			Convey("With subject id and bag", func() {
				expireAt, err := ptypes.TimestampProto(time.Now())
				So(err, ShouldBeNil)

				req = &StartRequest{SubjectId: subjectID, Bag: bag}
				session = &Session{AccessToken: token, SubjectId: subjectID, Bag: bag, ExpireAt: expireAt}

				Convey("Without storage error", func() {
					suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(session, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					Convey("Should", itSuccess)
				})
				Convey("With storage postgres error", func() {
					expectedErr = pq.Error{Message: "fake postgres error"}
					suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(nil, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					Convey("Should return grpc error with code 13", func() {
						So(err, ShouldBeGRPCError, codes.Internal, expectedErr.Error())
					})
					Convey("Should return an nil response", func() {
						So(res, ShouldBeNil)
					})
				})
			})
			Convey("With subject and without bag", func() {
				expireAt, err := ptypes.TimestampProto(time.Now())
				So(err, ShouldBeNil)

				req = &StartRequest{SubjectId: subjectID}
				session = &Session{AccessToken: token, SubjectId: subjectID, ExpireAt: expireAt}
				suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Once().
					Return(session, expectedErr)

				res, err = suite.service.Start(context.Background(), req)

				Convey("Should", itSuccess)
			})
			Convey("Without subject and with bag", func() {
				req = &StartRequest{Bag: bag}
				expectedErr = errors.New("mnemosyne: session cannot be started, subject id is missing")
				suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Once().
					Return(session, expectedErr)

				res, err = suite.service.Start(context.Background(), req)

				Convey("Should return an error", func() {
					So(err, ShouldNotEqual, expectedErr)
				})
				Convey("Should return an nil response", func() {
					So(res, ShouldBeNil)
				})
			})
		})
	})

	suite.teardown(t)
}

func TestRPCServer_Start_postgresStore(t *testing.T) {
	Convey("Start", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With subject id", func() {
			sid := "entity:1"
			Convey("Should work", func() {
				resp, err := s.client.Start(context.Background(), &StartRequest{
					SubjectId: sid,
				})

				So(err, ShouldBeNil)
				So(resp, ShouldBeValidStartResponse, sid)
			})
		})
		Convey("Without subject id", func() {
			Convey("Should fail", func() {
				resp, err := s.client.Start(context.Background(), &StartRequest{})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(ErrMissingSubjectID))
			})
		})
	}))
}

func TestRPCServer_Get_postgresStore(t *testing.T) {
	var (
		sid string
		at  *AccessToken
	)
	Convey("Get", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken

			Convey("With proper access token", func() {
				Convey("Should work", func() {
					resp, err := s.client.Get(context.Background(), &GetRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidGetResponse, sid)
				})
			})
			Convey("Without access token", func() {
				Convey("Should work", func() {
					resp, err := s.client.Get(context.Background(), &GetRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(ErrMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should fail", func() {
				access := DecodeAccessTokenString("0000000000test")
				resp, err := s.client.Get(context.Background(), &GetRequest{
					AccessToken: &access,
				})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.NotFound, "mnemosyne: session (get) with access token key:\"0000000000\" hash:\"test\"  does not exists")
			})
		})
	}))
}
