package main

import (
	"errors"

	"github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/protot"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var _ = Describe("RPCServer", func() {
	var (
		err         error
		suite       *integrationSuite
		storage     *storageMock
		expectedErr error
		subjectID   string
		bag         map[string]string
		session     *mnemosyne.Session
		token       *mnemosyne.Token
	)
	BeforeSuite(func() {
		storage = &storageMock{}
		suite = newIntegrationSuite(storage)
		suite.serve(grpc.WithInsecure())
	})
	AfterSuite(func() {
		Expect(suite.teardown()).ToNot(HaveOccurred())
	})
	BeforeEach(func() {
		expectedErr = nil
		subjectID = "subject_id"
		bag = map[string]string{"key": "value"}
		tk := mnemosyne.EncodeTokenString("key", "hash")
		token = &tk
	})
	Describe("Start", func() {
		var (
			req *mnemosyne.StartRequest
			res *mnemosyne.StartResponse
		)

		itSuccess := func() {
			It("should not return any error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
			It("should return session with same bag", func() {
				Expect(res.Session.Bag).To(Equal(req.Bag))
			})
			It("should return session with same subject id", func() {
				Expect(res.Session.SubjectId).To(Equal(req.SubjectId))
			})
			It("should return session with expire at timestamp", func() {
				AssertTimestamp(res.Session.ExpireAt)
			})
			It("should return session with token", func() {
				AssertToken(res.Session.Token)
			})
		}
		JustBeforeEach(func() {
			res, err = suite.service.Start(context.Background(), req)
		})
		Context("with subject id and bag", func() {
			BeforeEach(func() {
				req = &mnemosyne.StartRequest{SubjectId: subjectID, Bag: bag}
				session = &mnemosyne.Session{Token: token, SubjectId: subjectID, Bag: bag, ExpireAt: protot.Now()}
			})
			Context("without storage error", func() {
				BeforeEach(func() {
					storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Return(session, expectedErr).
						Once()
				})
				itSuccess()
			})
			Context("with storage postgres error", func() {
				BeforeEach(func() {
					expectedErr = pq.Error{Message: "fake postgres error"}
					storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
						Return(nil, expectedErr).
						Once()
				})
				It("should return grpc error with code 13", func() {
					AssertGRPCError(err, codes.Internal, expectedErr.Error())
				})
				It("should return an nil response", func() {
					Expect(res).To(BeNil())
				})
			})
		})
		Context("with subject and without bag", func() {
			BeforeEach(func() {
				req = &mnemosyne.StartRequest{SubjectId: subjectID}
				session = &mnemosyne.Session{Token: token, SubjectId: subjectID, ExpireAt: protot.Now()}
				storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Return(session, expectedErr).
					Once()
			})
			itSuccess()
		})
		Context("without subject and with bag", func() {
			BeforeEach(func() {
				req = &mnemosyne.StartRequest{Bag: bag}
				expectedErr = errors.New("mnemosyned: session cannot be started, subject id is missing")
				storage.On("Start", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).
					Return(session, expectedErr).
					Once()
			})
			It("should return an error", func() {
				Expect(err).ToNot(Equal(expectedErr))
			})
			It("should return an nil response", func() {
				Expect(res).To(BeNil())
			})
		})
	})
})
