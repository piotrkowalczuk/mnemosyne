package mnemosyned

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestSessionManager_mockedStore(t *testing.T) {
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
			token := "supertoken"

			itSuccess := func() {
				Convey("Not return any error", func() {
					So(err, ShouldBeNil)
				})
				Convey("Return session with same bag", func() {
					So(res.Session.Bag, ShouldResemble, req.Session.Bag)
				})
				Convey("Return session with same subject accessToken", func() {
					So(res.Session, ShouldNotBeNil)
					So(res.Session.SubjectId, ShouldEqual, req.Session.SubjectId)
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

			Convey("With subject accessToken and bag", func() {
				expireAt, err := ptypes.TimestampProto(time.Now())
				So(err, ShouldBeNil)

				req = &mnemosynerpc.StartRequest{Session: &mnemosynerpc.Session{SubjectId: subjectID, Bag: bag}}
				session = &mnemosynerpc.Session{AccessToken: token, SubjectId: subjectID, Bag: bag, ExpireAt: expireAt}

				Convey("Without storage error", func() {
					suite.store.On("Start", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(session, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					So(err, ShouldBeNil)
					Convey("Should", itSuccess)
				})
				Convey("With storage postgres error", func() {
					expectedErr = pq.Error{Message: "fake postgres error"}
					suite.store.On("Start", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(nil, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					Convey("Should return grpc error with code 13", func() {
						So(err, ShouldBeGRPCError(ShouldEqual), codes.Unknown, expectedErr.Error())
					})
					Convey("Should return an nil response", func() {
						So(res, ShouldBeNil)
					})
				})
			})
			Convey("With subject and without bag", func() {
				expireAt, err := ptypes.TimestampProto(time.Now())
				So(err, ShouldBeNil)

				req = &mnemosynerpc.StartRequest{Session: &mnemosynerpc.Session{SubjectId: subjectID}}
				session = &mnemosynerpc.Session{AccessToken: token, SubjectId: subjectID, ExpireAt: expireAt}
				suite.store.On("Start", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Once().
					Return(session, expectedErr)

				res, err = suite.service.Start(context.Background(), req)

				So(err, ShouldBeNil)
				Convey("Should", itSuccess)
			})
			Convey("Without subject and with bag", func() {
				req = &mnemosynerpc.StartRequest{Session: &mnemosynerpc.Session{Bag: bag}}
				expectedErr = errors.New("session cannot be started, subject accessToken is missing")
				suite.store.On("Start", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
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

func TestSessionManager_Start_postgresStore(t *testing.T) {
	factor := 2

	Convey("Start", t, func() {
		Convey("With single node", WithE2ESuite(t, func(s *e2eSuite) {
			Convey("With subject id", func() {
				sid := "entity:1"
				Convey("Should return newly created session", func() {
					resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
						Session: &mnemosynerpc.Session{SubjectId: sid},
					})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidStartResponse, sid)
				})
			})
			Convey("Without subject id", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, status.Convert(errMissingSession).Message())
				})
			})
		}))
		Convey("With cluster", WithE2ESuites(t, factor, func(s e2eSuites) {
			Convey("With subject id", func() {
				Convey("Should return newly created session", func() {
					for i := 0; i < factor; i++ {
						sid := fmt.Sprintf("entity:%d", i)
						resp, err := s[i].client.Start(context.Background(), &mnemosynerpc.StartRequest{
							Session: &mnemosynerpc.Session{SubjectId: sid},
						})
						So(err, ShouldBeNil)
						So(resp, ShouldBeValidStartResponse, sid)
					}
				})
			})
			Convey("Without subject id", func() {
				Convey("Should return invalid argument gRPC error", func() {
					for i := 0; i < factor; i++ {
						resp, err := s[i].client.Start(context.Background(), &mnemosynerpc.StartRequest{})

						So(resp, ShouldBeNil)
						So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, status.Convert(errMissingSession).Message())
					}
				})
			})
		}))
	})
}

func TestSessionManager_Get_postgresStore(t *testing.T) {
	factor := 2
	var (
		subjectID   string
		accessToken string
	)
	Convey("Get", t, func() {
		Convey("With single node", WithE2ESuite(t, func(s *e2eSuite) {
			Convey("With existing session", func() {
				subjectID = "entity:1"
				resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
					Session: &mnemosynerpc.Session{SubjectId: subjectID},
				})

				So(err, ShouldBeNil)
				So(resp, ShouldBeValidStartResponse, subjectID)

				accessToken = resp.Session.AccessToken
				Convey("With proper access token", func() {
					Convey("Should return corresponding session", func() {
						resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{
							AccessToken: accessToken,
						})

						So(err, ShouldBeNil)
						So(resp, ShouldBeValidGetResponse, subjectID)
					})
				})
				Convey("Without access token", func() {
					Convey("Should return invalid argument gRPC error", func() {
						resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{})

						So(resp, ShouldBeNil)
						So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, status.Convert(errMissingAccessToken).Message())
					})
				})
			})
			Convey("With unknown access token", func() {
				Convey("Should return not found gRPC error", func() {
					resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{
						AccessToken: "0000000000test",
					})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.NotFound, "mnemosyned: "+storage.ErrSessionNotFound.Error())
				})
			})
		}))
		Convey("With cluster", WithE2ESuites(t, factor, func(s e2eSuites) {
			Convey("With existing session", func() {
				var tokens []string
				for i := 0; i < factor; i++ {
					subjectID := fmt.Sprintf("entity:%d", i)
					resp, err := s[i].client.Start(context.Background(), &mnemosynerpc.StartRequest{
						Session: &mnemosynerpc.Session{SubjectId: subjectID},
					})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidStartResponse, subjectID)

					tokens = append(tokens, resp.Session.AccessToken)
				}

				for i := 0; i < factor; i++ {
					Convey(fmt.Sprintf("As node#%d", i), func() {
						for j, tok := range tokens {
							msg := fmt.Sprintf("With someones else access token: %d", j)
							if i == j {
								msg = "With it's own access token"
							}
							Convey(msg, func() {
								Convey("Should return corresponding session", func() {
									resp, err := s[i].client.Get(context.Background(), &mnemosynerpc.GetRequest{
										AccessToken: tok,
									})

									So(err, ShouldBeNil)
									So(resp, ShouldBeValidGetResponse, subjectID)
									Printf("session retrieved from '%s' through '%s'\n", s[j].daemon.Addr().String(), s[i].daemon.Addr().String())
								})
							})
						}
						Convey("Without access token", func() {
							Convey("Should return invalid argument gRPC error", func() {
								resp, err := s[i].client.Get(context.Background(), &mnemosynerpc.GetRequest{})

								So(resp, ShouldBeNil)
								So(err, ShouldBeGRPCError(ShouldEndWith), codes.InvalidArgument, "missing access token")
							})
						})
					})
				}
			})
			Convey("With non-existing session", func() {
				for i := 0; i < factor; i++ {
					Convey(fmt.Sprintf("As node#%d", i), func() {
						Convey("Should return corresponding session", func() {
							resp, err := s[i].client.Get(context.Background(), &mnemosynerpc.GetRequest{
								AccessToken: "non-existing-token",
							})

							So(resp, ShouldBeNil)
							So(err, ShouldBeGRPCError(ShouldEndWith), codes.NotFound, "session not found")
						})
					})
				}
			})
			//accessToken = resp.Session.AccessToken
			//Convey("With proper access token", func() {
			//	Convey("Should return corresponding session", func() {
			//		resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{
			//			AccessToken: accessToken,
			//		})
			//
			//		So(err, ShouldBeNil)
			//		So(resp, ShouldBeValidGetResponse, subjectID)
			//	})
			//})
			//Convey("Without access token", func() {
			//	Convey("Should return invalid argument gRPC error", func() {
			//		resp, err := s.client.Get(context.Background(), &mnemosynerpc.GetRequest{})
			//
			//		So(resp, ShouldBeNil)
			//		So(err, ShouldBeGRPCError(ShouldEqual),codes.InvalidArgument, grpc.ErrorDesc(erringAccessToken))
			//	})
			//})
		}))
	})
}

func TestSessionManager_Context_postgresStore(t *testing.T) {
	var (
		subjectID   string
		accessToken string
	)
	Convey("Context", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			subjectID = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				Session: &mnemosynerpc.Session{SubjectId: subjectID},
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, subjectID)

			accessToken = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return corresponding session", func() {
					meta := metadata.Pairs(mnemosyne.AccessTokenMetadataKey, string(accessToken))
					ctx := metadata.NewOutgoingContext(context.Background(), meta)
					resp, err := s.client.Context(ctx, &empty.Empty{})

					So(err, ShouldBeNil)
					So(resp, ShouldBeValidContextResponse, subjectID)
				})
			})
			Convey("Without metadata", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Context(context.Background(), &empty.Empty{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, "mnemosyned: missing access token in metadata")
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"some-key": "some-value"}))
					resp, err := s.client.Context(ctx, &empty.Empty{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, "mnemosyned: missing access token in metadata")
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				meta := metadata.Pairs(mnemosyne.AccessTokenMetadataKey, "0000000000test")
				ctx := metadata.NewOutgoingContext(context.Background(), meta)
				resp, err := s.client.Context(ctx, &empty.Empty{})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError(ShouldEqual), codes.NotFound, "mnemosyned: "+storage.ErrSessionNotFound.Error())
			})
		})
	}))
}

func TestSessionManager_Exists_postgresStore(t *testing.T) {
	var (
		subjectID   string
		accessToken string
	)
	Convey("Exists", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			subjectID = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				Session: &mnemosynerpc.Session{SubjectId: subjectID},
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, subjectID)

			accessToken = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return true", func() {
					res, err := s.client.Exists(context.Background(), &mnemosynerpc.ExistsRequest{
						AccessToken: accessToken,
					})

					So(err, ShouldBeNil)
					So(res.GetValue(), ShouldBeTrue)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					res, err := s.client.Exists(context.Background(), &mnemosynerpc.ExistsRequest{})

					So(res, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, status.Convert(errMissingAccessToken).Message())
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return false", func() {
				res, err := s.client.Exists(context.Background(), &mnemosynerpc.ExistsRequest{
					AccessToken: "0000000000test",
				})

				So(err, ShouldBeNil)
				So(res.GetValue(), ShouldBeFalse)
			})
		})
	}))
}

func TestSessionManager_Abandon_postgresStore(t *testing.T) {
	var (
		subjectID   string
		accessToken string
	)
	Convey("Abandon", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			subjectID = "entity:1"
			resp, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				Session: &mnemosynerpc.Session{SubjectId: subjectID},
			})

			So(err, ShouldBeNil)
			So(resp, ShouldBeValidStartResponse, subjectID)

			accessToken = resp.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return true", func() {
					res, err := s.client.Abandon(context.Background(), &mnemosynerpc.AbandonRequest{
						AccessToken: accessToken,
					})

					So(err, ShouldBeNil)
					So(res.GetValue(), ShouldBeTrue)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.Abandon(context.Background(), &mnemosynerpc.AbandonRequest{})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, status.Convert(errMissingAccessToken).Message())
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				resp, err := s.client.Abandon(context.Background(), &mnemosynerpc.AbandonRequest{
					AccessToken: "0000000000test",
				})

				So(resp, ShouldBeNil)
				So(err, ShouldBeGRPCError(ShouldEqual), codes.NotFound, "mnemosyned: "+storage.ErrSessionNotFound.Error())
			})
		})
	}))
}

func TestSessionManager_Delete_postgresStore(t *testing.T) {
	var (
		subjectID   string
		accessToken string
	)
	Convey("Delete", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			subjectID = "entity:1"
			res, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				Session: &mnemosynerpc.Session{SubjectId: subjectID},
			})

			So(err, ShouldBeNil)
			So(res, ShouldBeValidStartResponse, subjectID)

			accessToken = res.Session.AccessToken
			Convey("With proper access token", func() {
				Convey("Should return that one record affected", func() {
					res, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{
						AccessToken: accessToken,
					})

					So(err, ShouldBeNil)
					So(res.GetValue(), ShouldEqual, 1)
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					res, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{})

					So(res, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, "mnemosyned: none of expected arguments was provided")
				})
			})
			Convey("With date ranges far in the future", func() {
				Convey("Should return that no records was affected", func() {
					from, err := ptypes.TimestampProto(time.Now().AddDate(100, 0, 0))
					So(err, ShouldBeNil)
					to, err := ptypes.TimestampProto(time.Now().AddDate(101, 0, 0))
					So(err, ShouldBeNil)

					res, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{
						AccessToken:  accessToken,
						ExpireAtFrom: from,
						ExpireAtTo:   to,
					})

					So(err, ShouldBeNil)
					So(res.GetValue(), ShouldEqual, 0)
				})
			})
			Convey("With date ranges that match existing session", func() {
				Convey("Should return that one record was affected", func() {
					from, err := ptypes.TimestampProto(time.Now().AddDate(-100, 0, 0))
					So(err, ShouldBeNil)
					to, err := ptypes.TimestampProto(time.Now().AddDate(100, 0, 0))
					So(err, ShouldBeNil)

					res, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{
						AccessToken:  accessToken,
						ExpireAtFrom: from,
						ExpireAtTo:   to,
					})

					So(err, ShouldBeNil)
					So(res.GetValue(), ShouldEqual, 1)
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return that even single record was affected", func() {
				res, err := s.client.Delete(context.Background(), &mnemosynerpc.DeleteRequest{
					AccessToken: "0000000000test",
				})

				So(err, ShouldBeNil)
				So(res.GetValue(), ShouldEqual, 0)
			})
		})
	}))
}

func TestSessionManager_SetValue_postgresStore(t *testing.T) {
	var (
		subjectID   string
		accessToken string
	)
	Convey("SetValue", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("With existing session", func() {
			subjectID = "entity:1"
			res, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
				Session: &mnemosynerpc.Session{SubjectId: subjectID},
			})

			So(err, ShouldBeNil)
			So(res, ShouldBeValidStartResponse, subjectID)

			accessToken = res.GetSession().GetAccessToken()
			Convey("With proper access token", func() {
				Convey("Should return that one record affected", func() {
					res, err := s.client.SetValue(context.Background(), &mnemosynerpc.SetValueRequest{
						AccessToken: accessToken,
						Key:         "key",
						Value:       "value",
					})

					So(err, ShouldBeNil)
					So(res.GetBag(), ShouldContainKey, "key")
				})
			})
			Convey("Without access token", func() {
				Convey("Should return invalid argument gRPC error", func() {
					resp, err := s.client.SetValue(context.Background(), &mnemosynerpc.SetValueRequest{
						Key:   "key",
						Value: "value",
					})

					So(resp, ShouldBeNil)
					So(err, ShouldBeGRPCError(ShouldEqual), codes.InvalidArgument, status.Convert(errMissingAccessToken).Message())
				})
			})
		})
		Convey("With unknown access token", func() {
			Convey("Should return not found gRPC error", func() {
				res, err := s.client.SetValue(context.Background(), &mnemosynerpc.SetValueRequest{
					AccessToken: "0000000000test",
					Key:         "key",
					Value:       "value",
				})

				So(res, ShouldBeNil)
				So(err, ShouldBeGRPCError(ShouldEqual), codes.NotFound, "mnemosyned: "+storage.ErrSessionNotFound.Error())
			})
		})
	}))
}

func TestSessionManager_List_postgresStore(t *testing.T) {
	var (
		subjectID string
	)
	nb := 20
	Convey("sessionManagerList", t, WithE2ESuite(t, func(s *e2eSuite) {
		Convey("Having multiple sessions active", func() {
			for i := 0; i < nb; i++ {
				res, err := s.client.Start(context.Background(), &mnemosynerpc.StartRequest{
					Session: &mnemosynerpc.Session{SubjectId: strconv.Itoa(i)},
				})
				So(err, ShouldBeNil)
				So(res, ShouldBeValidStartResponse, subjectID)
			}
			Convey("With empty request", func() {
				Convey("Should return last 10 sessions", func() {
					res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{})

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(len(res.Sessions), ShouldEqual, 10)
				})
			})
			Convey("With limit set", func() {
				Convey("Should return specified numer of sessions", func() {
					res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Limit: int64(nb),
					})

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(len(res.Sessions), ShouldEqual, nb)
				})
			})
			Convey("With offset higher than overall number of sessions", func() {
				Convey("Should return empty collection", func() {
					res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Offset: int64(nb),
					})

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(len(res.Sessions), ShouldEqual, 0)
				})
			})
			Convey("With expire at to set in the past", func() {
				Convey("Should return empty collection", func() {
					past, err := ptypes.TimestampProto(time.Now().Add(-5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Query: &mnemosynerpc.Query{
							ExpireAtTo: past,
						},
					})

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(len(res.Sessions), ShouldEqual, 0)
				})
			})
			Convey("With expire at from set in the future", func() {
				Convey("Should return empty collection", func() {
					future, err := ptypes.TimestampProto(time.Now().Add(5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Query: &mnemosynerpc.Query{
							ExpireAtFrom: future,
						},
					})

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(len(res.GetSessions()), ShouldEqual, 0)
				})
			})
			Convey("With time range set very wide and maximum offset", func() {
				Convey("Should return all possible sessions", func() {
					from, err := ptypes.TimestampProto(time.Now().Add(-5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					to, err := ptypes.TimestampProto(time.Now().Add(5 * time.Hour).UTC())
					So(err, ShouldBeNil)

					res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
						Limit: int64(nb),
						Query: &mnemosynerpc.Query{
							ExpireAtFrom: from,
							ExpireAtTo:   to,
						},
					})

					So(err, ShouldBeNil)
					So(res, ShouldNotBeNil)
					So(len(res.Sessions), ShouldEqual, nb)
				})
			})
		})
		Convey("Without single session active", func() {
			Convey("Should return empty collection", func() {
				res, err := s.client.List(context.Background(), &mnemosynerpc.ListRequest{
					Limit: 100,
				})

				So(err, ShouldBeNil)
				So(res, ShouldNotBeNil)
				So(len(res.Sessions), ShouldEqual, 0)
			})
		})
	}))
}
