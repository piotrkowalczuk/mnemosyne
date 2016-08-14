package mnemosyned

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestRPCServer_mockedStore(t *testing.T) {
	var (
		err         error
		suite       *integrationSuite
		expectedErr error
		subjectID   string
		bag         map[string]string
		session     *mnemosynerpc.Session
	)

	suite = &integrationSuite{}
	suite.setup(t)

	Convey("RPCServer", t, func() {
		Convey("Start", func() {
			var (
				req *mnemosynerpc.StartRequest
				res *mnemosynerpc.StartResponse
			)

			expectedErr = nil
			subjectID = "subject_id"
			bag = map[string]string{"key": "value"}
			token := mnemosynerpc.NewAccessToken("key", "hash")

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

				req = &mnemosynerpc.StartRequest{SubjectId: subjectID, Bag: bag}
				session = &mnemosynerpc.Session{AccessToken: token, SubjectId: subjectID, Bag: bag, ExpireAt: expireAt}

				Convey("Without storage error", func() {
					suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(session, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					So(err, ShouldBeNil)
					Convey("Should", itSuccess)
				})
				Convey("With storage postgres error", func() {
					expectedErr = pq.Error{Message: "fake postgres error"}
					suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(nil, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					Convey("Should return grpc error with code 13", func() {
						So(err, ShouldBeGRPCError, codes.Unknown, expectedErr.Error())
					})
					Convey("Should return an nil response", func() {
						So(res, ShouldBeNil)
					})
				})
			})
			Convey("With subject and without bag", func() {
				expireAt, err := ptypes.TimestampProto(time.Now())
				So(err, ShouldBeNil)

				req = &mnemosynerpc.StartRequest{SubjectId: subjectID}
				session = &mnemosynerpc.Session{AccessToken: token, SubjectId: subjectID, ExpireAt: expireAt}
				suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Once().
					Return(session, expectedErr)

				res, err = suite.service.Start(context.Background(), req)

				So(err, ShouldBeNil)
				Convey("Should", itSuccess)
			})
			Convey("Without subject and with bag", func() {
				req = &mnemosynerpc.StartRequest{Bag: bag}
				expectedErr = errors.New("session cannot be started, subject id is missing")
				suite.store.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
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
				resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
					SubjectId: sid,
				})

				So(err, ShouldBeNil)
				So(resp, ShouldBeValidStartResponse, sid)
			})
		})
		Convey("Without subject id", func() {
			Convey("Should return invalid argument gRPC error", func() {
				resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(errMissingSubjectID))
			})
		})
	}))
}

func TestRPCServer_Get_postgresStore(t *testing.T) {
	var (
		sid string
		at  string
	)
	Convey("Get", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return corresponding session", func() {
					resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidGetResponse, sid)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(errMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{
					AccessToken: "0000000000test",
				})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.NotFound, "mnemosyned: session (get) does not exists")
			})
		})
	}))
}

func TestRPCServer_Context_postgresStore(t *testing.T) {
	var (
		sid string
		at  string
	)
	Convey("Context", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return corresponding session", func() {
					meta := metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, string(at))
					ctx := metadata.NewContext(context.Background(), meta)
					resp, err := s.client.Context(ctx, &empty.Empty{})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidContextResponse, sid)
				})
			})
			Convey("Without metadata", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Context(context.Background(), &empty.Empty{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, "mnemosyned: missing access token in metadata")
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"some-key": "some-value"}))
					resp, err := s.client.Context(ctx, &empty.Empty{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, "mnemosyned: missing access token in metadata")
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				meta := metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, "0000000000test")
				ctx := metadata.NewContext(context.Background(), meta)
				resp, err := s.client.Context(ctx, &empty.Empty{})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.NotFound, "mnemosyned: session (context) does not exists")
			})
		})
	}))
}

func TestRPCServer_Exists_postgresStore(t *testing.T) {
	var (
		sid string
		at  string
	)
	Convey("Exists", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return true", func() {
					resp, err := s.client.Exists(context.Background(), &mnemosynerpc.ExistsRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp.Exists, ShouldBeTrue)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Exists(context.Background(), &mnemosynerpc.ExistsRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(errMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return false", func() {
				resp, err := s.client.Exists(context.Background(), &mnemosynerpc.ExistsRequest{
					AccessToken: "0000000000test",
				})

				So(err, ShouldBeNil)
				So(resp.Exists, ShouldBeFalse)
			})
		})
	}))
}

func TestRPCServer_Abandon_postgresStore(t *testing.T) {
	var (
		sid string
		at  string
	)
	Convey("Abandon", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return true", func() {
					resp, err := s.client.Abandon(context.Background(), &mnemosynerpc.AbandonRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp.Abandoned, ShouldBeTrue)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Abandon(context.Background(), &mnemosynerpc.AbandonRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(errMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				resp, err := s.client.Abandon(context.Background(), &mnemosynerpc.AbandonRequest{
					AccessToken: "0000000000test",
				})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.NotFound, grpc.ErrorDesc(errSessionNotFound))
			})
		})
	}))
}

func TestRPCServer_Delete_postgresStore(t *testing.T) {
	var (
		sid string
		at  string
	)
	Convey("Delete", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return that one record affected", func() {
					resp, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{
						AccessToken: at,
					})

					So(err, ShouldBeNil)
					So(resp.Count, ShouldEqual, 1)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, "mnemosyned: none of expected arguments was provided")
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return that even single record was affected", func() {
				resp, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{
					AccessToken: "0000000000test",
				})

				So(err, ShouldBeNil)
				So(resp.Count, ShouldEqual, 0)
			})
		})
	}))
}

func TestRPCServer_SetValue_postgresStore(t *testing.T) {
	var (
		sid string
		at  string
	)
	Convey("SetValue", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			sid = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				SubjectId: sid,
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, sid)

			at = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return that one record affected", func() {
					resp, err := s.client.SetValue(context.Background(), &mnemosynerpc.SetValueRequest{
						AccessToken: at,
						Key:         "key",
						Value:       "value",
					})

					So(err, ShouldBeNil)
					So(resp.Bag, ShouldContainKey, "key")
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.SetValue(context.Background(), &mnemosynerpc.SetValueRequest{
						Key:   "key",
						Value: "value",
					})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError, codes.InvalidArgument, grpc.ErrorDesc(errMissingAccessToken))
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				resp, err := s.client.SetValue(context.Background(), &mnemosynerpc.SetValueRequest{
					AccessToken: "0000000000test",
					Key:         "key",
					Value:       "value",
				})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError, codes.NotFound, grpc.ErrorDesc(errSessionNotFound))
			})
		})
	}))
}

func TestRPCServer_List_postgresStore(t *testing.T) {
	var (
		sid string
	)
	nb := 20
	Convey("List", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("Having multiple sessions active", func() {
			for i := 0; i < nb; i++ {
				resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
					SubjectId: strconv.Itoa(i),
				})
				So(err, ShouldBeNil)
				So(resp, ShouldBeValidStartResponse, sid)
			}
			Convey("With empty request", func() {
				Convey("Should return last 10 sessions", func() {
					resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{})

					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(len(resp.Sessions), ShouldEqual, 10)
				})
			})
			Convey("With limit set", func() {
				Convey("Should return specified numer of sessions", func() {
					resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Limit: int64(nb),
					})

					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(len(resp.Sessions), ShouldEqual, nb)
				})
			})
			Convey("With offset higher than overall number of sessions", func() {
				Convey("Should return empty collection", func() {
					resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Offset: int64(nb),
					})

					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(len(resp.Sessions), ShouldEqual, 0)
				})
			})
			Convey("With expire at to set in the past", func() {
				Convey("Should return empty collection", func() {
					past, err := ptypes.TimestampProto(time.Now().Add(-5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						ExpireAtTo: past,
					})

					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(len(resp.Sessions), ShouldEqual, 0)
				})
			})
			Convey("With expire at from set in the future", func() {
				Convey("Should return empty collection", func() {
					future, err := ptypes.TimestampProto(time.Now().Add(5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						ExpireAtFrom: future,
					})

					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(len(resp.Sessions), ShouldEqual, 0)
				})
			})
			Convey("With time range set very wide and maximum offset", func() {
				Convey("Should return all possible sessions", func() {
					from, err := ptypes.TimestampProto(time.Now().Add(-5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					to, err := ptypes.TimestampProto(time.Now().Add(5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Limit:        int64(nb),
						ExpireAtFrom: from,
						ExpireAtTo:   to,
					})

					So(err, ShouldBeNil)
					So(resp, ShouldNotBeNil)
					So(len(resp.Sessions), ShouldEqual, nb)
				})
			})
		})
		Convey("Without single session active", func() {
			Convey("Should return empty collection", func() {
				resp, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
					Limit: 100,
				})

				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(len(resp.Sessions), ShouldEqual, 0)
			})
		})
	}))
}
