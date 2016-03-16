// +build unit,!postgres

package main

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/protot"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestRPCServer(t *testing.T) {
	var (
		err         error
		suite       *integrationSuite
		storage     *storageMock
		expectedErr error
		subjectID   string
		bag         map[string]string
		session     *mnemosyne.Session
		token       *mnemosyne.AccessToken
	)

	storage = &storageMock{}
	suite = newIntegrationSuite(storage)
	if err = suite.serve(grpc.WithInsecure()); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
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
				req = &mnemosyne.StartRequest{SubjectId: subjectID, Bag: bag}
				session = &mnemosyne.Session{AccessToken: token, SubjectId: subjectID, Bag: bag, ExpireAt: protot.Now()}

				Convey("Without storage error", func() {
					storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Once().
						Return(session, expectedErr)

					res, err = suite.service.Start(context.Background(), req)

					Convey("Should", itSuccess)
				})
				Convey("With storage postgres error", func() {
					expectedErr = pq.Error{Message: "fake postgres error"}
					storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
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
				req = &mnemosyne.StartRequest{SubjectId: subjectID}
				session = &mnemosyne.Session{AccessToken: token, SubjectId: subjectID, ExpireAt: protot.Now()}
				storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Once().
					Return(session, expectedErr)

				res, err = suite.service.Start(context.Background(), req)

				Convey("Should", itSuccess)
			})
			Convey("Without subject and with bag", func() {
				req = &mnemosyne.StartRequest{Bag: bag}
				expectedErr = errors.New("mnemosyned: session cannot be started, subject id is missing")
				storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
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

	suite.teardown()
}
