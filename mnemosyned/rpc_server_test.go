package mnemosyned

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne"
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
		session     *mnemosyne.Session
		token       *mnemosyne.AccessToken
	)

	suite = &integrationSuite{}
	suite.setup(t)

	Convey("RPCServer", t, func() {
		Convey("Start", func() {
			var (
				req *mnemosyne.StartRequest
				res *mnemosyne.StartResponse
			)

			expectedErr = nil
			subjectID = "subject_id"
			bag = map[string]string{"key": "value"}
			tk := mnemosyne.NewAccessToken([]byte("key"), []byte("hash"))
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

				req = &mnemosyne.StartRequest{SubjectId: subjectID, Bag: bag}
				session = &mnemosyne.Session{AccessToken: token, SubjectId: subjectID, Bag: bag, ExpireAt: expireAt}

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
						So(err, ShouldBeGRPCError, codes.Internal, "mnemosyned: "+expectedErr.Error())
					})
					Convey("Should return an nil response", func() {
						So(res, ShouldBeNil)
					})
				})
			})
			Convey("With subject and without bag", func() {
				expireAt, err := ptypes.TimestampProto(time.Now())
				So(err, ShouldBeNil)

				req = &mnemosyne.StartRequest{SubjectId: subjectID}
				session = &mnemosyne.Session{AccessToken: token, SubjectId: subjectID, ExpireAt: expireAt}
				suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Once().
					Return(session, expectedErr)

				res, err = suite.service.Start(context.Background(), req)

				Convey("Should", itSuccess)
			})
			Convey("Without subject and with bag", func() {
				req = &mnemosyne.StartRequest{Bag: bag}
				expectedErr = errors.New("session cannot be started, subject id is missing")
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
			Convey("Should return newly created session", func() {
				resp, err := s.client.Start(context.Background(), &mnemosyne.StartRequest{
					SubjectId: sid,
				})

				So(err, ShouldBeNil)
				So(resp, ShouldBeValidStartResponse, sid)
			})
		})
		Convey("Without subject id", func() {
			Convey("Should return invalid argument gRPC error", func() {
				resp, err := s.client.Start(context.Background(), &mnemosyne.StartRequest{})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(ErrMissingSubjectID))
			})
		})
	}))
}

func TestRPCServer_Get_postgresStore(t *testing.T) {
	var (
		sid string
		at  *mnemosyne.AccessToken
	)
	Convey("Get", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosyne.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return corresponding session", func() {
					resp, err := s.client.Get(context.Background(), &mnemosyne.GetRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidGetResponse, sid)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Get(context.Background(), &mnemosyne.GetRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(ErrMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				access := mnemosyne.DecodeAccessTokenString("0000000000test")
				resp, err := s.client.Get(context.Background(), &mnemosyne.GetRequest{
					AccessToken: &access,
				})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.NotFound, "mnemosyned: session (get) does not exists")
			})
		})
	}))
}

func TestRPCServer_Exists_postgresStore(t *testing.T) {
	var (
		sid string
		at  *mnemosyne.AccessToken
	)
	Convey("Exists", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosyne.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return true", func() {
					resp, err := s.client.Exists(context.Background(), &mnemosyne.ExistsRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp.Exists, ShouldBeTrue)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return false", func() {
					resp, err := s.client.Exists(context.Background(), &mnemosyne.ExistsRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(ErrMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return false", func() {
				access := mnemosyne.DecodeAccessTokenString("0000000000test")
				resp, err := s.client.Exists(context.Background(), &mnemosyne.ExistsRequest{
					AccessToken: &access,
				})

				So(err, ShouldBeNil)
				So(resp.Exists, ShouldBeFalse)
			})
		})
	}))
}
