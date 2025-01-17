package criteria

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/pomerium/pomerium/pkg/grpc/session"
	"github.com/pomerium/pomerium/pkg/grpc/user"
)

func TestEmails(t *testing.T) {
	t.Run("no session", func(t *testing.T) {
		res, err := evaluate(t, `
allow:
  and:
    - email: test@example.com
`, []dataBrokerRecord{}, Input{Session: InputSession{ID: "SESSION_ID"}})
		require.NoError(t, err)
		require.Equal(t, false, res["allow"])
		require.Equal(t, false, res["deny"])
	})
	t.Run("by email", func(t *testing.T) {
		res, err := evaluate(t, `
allow:
  and:
    - email: test@example.com
`,
			[]dataBrokerRecord{
				&session.Session{
					Id:     "SESSION_ID",
					UserId: "USER_ID",
				},
				&user.User{
					Id:    "USER_ID",
					Email: "test@example.com",
				},
			},
			Input{Session: InputSession{ID: "SESSION_ID"}})
		require.NoError(t, err)
		require.Equal(t, true, res["allow"])
		require.Equal(t, false, res["deny"])
	})
	t.Run("by impersonate email", func(t *testing.T) {
		res, err := evaluate(t, `
allow:
  and:
    - email: test2@example.com
`,
			[]dataBrokerRecord{
				&session.Session{
					Id:               "SESSION_ID",
					UserId:           "USER_ID",
					ImpersonateEmail: proto.String("test2@example.com"),
				},
				&user.User{
					Id:    "USER_ID",
					Email: "test1@example.com",
				},
			},
			Input{Session: InputSession{ID: "SESSION_ID"}})
		require.NoError(t, err)
		require.Equal(t, true, res["allow"])
		require.Equal(t, false, res["deny"])
	})
}
