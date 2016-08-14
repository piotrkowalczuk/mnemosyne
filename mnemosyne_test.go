package mnemosyne

import (
	"reflect"
	"testing"

	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func TestMnemosyne(t *testing.T) {
	ctx := context.Background()
	subjectID := "1"
	subjectClient := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36"
	bag := map[string]string{"username": "johnsnow@gmail.com"}

	addr, closer := mnemosyned.TestDaemon(t, mnemosyned.TestDaemonOpts{
		StoragePostgresAddress: testPostgresAddress,
	})
	defer closer.Close()

	m, err := New(MnemosyneOpts{
		Addresses: []string{addr.String()},
	})
	defer m.Close()

	if err != nil {
		t.Fatalf("unexpected client initialization error: %s", err.Error())
	}
	ses, err := m.Start(ctx, subjectID, subjectClient, bag)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	ses, err = m.Get(ctx, ses.AccessToken)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	meta := metadata.Pairs(mnemosynerpc.AccessTokenMetadataKey, ses.AccessToken)
	ctx = metadata.NewContext(context.Background(), meta)
	ses, err = m.FromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if ses.AccessToken == "" {
		t.Error("access token cannot be nil")
	}
	if len(ses.AccessToken) == 0 {
		t.Error("access token should not be empty")
	}
	if ses.ExpireAt.Seconds == 0 || ses.ExpireAt.Nanos == 0 {
		t.Error("expire at should not be zero value")
	}
	if ses.SubjectId != subjectID {
		t.Errorf("wrong subject id, expected %s but got %s", ses.SubjectId, subjectID)
	}
	if ses.SubjectClient != subjectClient {
		t.Errorf("wrong subject client, expected %s but got %s", ses.SubjectClient, subjectClient)
	}
	if !reflect.DeepEqual(ses.Bag, bag) {
		t.Errorf("wrong bag, expected:\n%s\nbut got:\n %s", ses.Bag, bag)
	}
}

func TestNew_missingAddresses(t *testing.T) {
	_, err := New(MnemosyneOpts{})
	if err == nil {
		t.Fatal("error expected, got nil")
	}
}
