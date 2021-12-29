package repository

import (
	"context"
	"database/sql"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/config"
	"github.com/ToTranMinhNhut/S3_FriendManagementAPI_NhutTo/internal/models"
	"github.com/stretchr/testify/require"
)

// loadSqlTestFile will be loading a mock testdata
func loadSqlTestFile(t *testing.T, db *sql.DB, sqlfile string) {
	// Read sql file
	b, err := ioutil.ReadFile(sqlfile)
	require.NoError(t, err)

	_, err = db.Exec(string(b))
	require.NoError(t, err)
}

func TestRepository_CreateFriend(t *testing.T) {
	tcs := map[string]struct {
		userId   int
		friendId int
		expError error
	}{
		"success with adding input of userIds": {
			userId:   100,
			friendId: 104,
		},
		"query by an unknown input userIds": {
			userId:   100,
			friendId: 99,
			expError: errors.New("models: unable to insert into friends: pq: insert or update on table \"friends\" violates foreign key constraint \"friends_friend_id_fkey\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			err = repo.CreateFriend(ctx, tc.userId, tc.friendId)
			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepository_IsExistedFriend(t *testing.T) {
	tcs := map[string]struct {
		userId    int
		friendId  int
		expResult bool
		expError  error
	}{
		"success with adding input of userIds": {
			userId:    100,
			friendId:  102,
			expResult: true,
		},
		"query by an unknown input userIds": {
			userId:    100,
			friendId:  99,
			expResult: false,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.IsExistedFriend(ctx, tc.userId, tc.friendId)

			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestRepository_IsBlockedUser(t *testing.T) {
	tcs := map[string]struct {
		userId    int
		friendId  int
		expResult bool
		expError  error
	}{
		"success with adding input of userIds": {
			userId:    100,
			friendId:  103,
			expResult: true,
		},
		"query by an unknown input userIds": {
			userId:    100,
			friendId:  99,
			expResult: false,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.IsBlockedUser(ctx, tc.userId, tc.friendId)

			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestRepository_GetFriendsByID(t *testing.T) {
	tcs := map[string]struct {
		userId    int
		expResult models.FriendSlice
		expError  error
	}{
		"success with adding input of userId": {
			userId: 100,
			expResult: models.FriendSlice{
				&models.Friend{UserID: 100, FriendID: 102},
			},
		},
		"query by an unknown input userId": {
			userId:    99,
			expResult: nil,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.GetFriendsByID(ctx, tc.userId)

			require.NoError(t, err)
			require.Equal(t, len(tc.expResult), len(result))
			for i, ss := range tc.expResult {
				require.Equal(t, ss, result[i])
			}
		})
	}
}

func TestRepository_GetUserBlocksByID(t *testing.T) {
	tcs := map[string]struct {
		userId    int
		expResult models.UserBlockSlice
		expError  error
	}{
		"success with adding input of userId": {
			userId: 100,
			expResult: models.UserBlockSlice{
				&models.UserBlock{RequestorID: 100, TargetID: 103},
				&models.UserBlock{RequestorID: 100, TargetID: 104},
			},
		},
		"query by an unknown input userId": {
			userId:    99,
			expResult: nil,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.GetUserBlocksByID(ctx, tc.userId)

			require.NoError(t, err)
			require.Equal(t, len(tc.expResult), len(result))
			for i, ss := range tc.expResult {
				require.Equal(t, ss, result[i])
			}
		})
	}
}

func TestRepository_CreateSubscription(t *testing.T) {
	tcs := map[string]struct {
		requestorId int
		targetId    int
		expError    error
	}{
		"success with adding input of userIds": {
			requestorId: 102,
			targetId:    103,
		},
		"query by an unknown input userIds": {
			requestorId: 99,
			targetId:    100,
			expError:    errors.New("models: unable to insert into subscriptions: pq: insert or update on table \"subscriptions\" violates foreign key constraint \"subscriptions_subscription_requestor_id_fkey\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			err = repo.CreateSubscription(ctx, tc.requestorId, tc.targetId)

			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepository_GetRecipientEmails(t *testing.T) {
	tcs := map[string]struct {
		senderId  int
		expResult models.UserSlice
		expError  error
	}{
		"success with adding input of userId": {
			senderId: 100,
			expResult: models.UserSlice{
				{Email: "common@example.com"},
			},
		},
		"query by an unknown input userId": {
			senderId: 99,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.GetRecipientEmails(ctx, tc.senderId)

			require.NoError(t, err)
			require.Equal(t, len(tc.expResult), len(result))
			for i, ss := range tc.expResult {
				require.Equal(t, ss.Email, result[i].Email)
			}
		})
	}
}

func TestRepository_CreateUserBlock(t *testing.T) {
	tcs := map[string]struct {
		requestorId int
		targetId    int
		expError    error
	}{
		"success with adding input of userIds": {
			requestorId: 100,
			targetId:    101,
		},
		"query by an unknown input userIds": {
			requestorId: 99,
			targetId:    101,
			expError:    errors.New("models: unable to insert into user_blocks: pq: insert or update on table \"user_blocks\" violates foreign key constraint \"user_blocks_requestor_id_fkey\""),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			err = repo.CreateUserBlock(ctx, tc.requestorId, tc.targetId)

			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRepository_IsSubscribedFriend(t *testing.T) {
	tcs := map[string]struct {
		requestorId int
		targetId    int
		expResult   bool
		expError    error
	}{
		"success with adding input of userIds": {
			requestorId: 101,
			targetId:    103,
			expResult:   true,
		},
		"query by an unknown input userIds": {
			requestorId: 100,
			targetId:    99,
			expResult:   false,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.IsSubscribedUser(ctx, tc.requestorId, tc.targetId)

			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestRepository_GetUserIDByEmail(t *testing.T) {
	tcs := map[string]struct {
		email     string
		expResult int
		expError  error
	}{
		"success with adding input of email": {
			email:     "john@example.com",
			expResult: 100,
		},
		"query by an unknown input email": {
			email:    "test@example.com",
			expError: errors.New("sql: no rows in result set"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.GetUserIDByEmail(ctx, tc.email)

			if tc.expError != nil {
				require.EqualError(t, err, tc.expError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestRepository_GetEmailsByUserIDs(t *testing.T) {
	tcs := map[string]struct {
		userIds   []int
		expResult []string
		expError  error
	}{
		"success with adding input of userIds": {
			userIds:   []int{100, 101},
			expResult: []string{"john@example.com", "andy@example.com"},
		},
		"query by an unknown input userIds": {
			userIds:  []int{99},
			expError: errors.New("sql: no rows in result set"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.GetEmailsByUserIDs(ctx, tc.userIds)

			require.NoError(t, err)
			require.Equal(t, len(tc.expResult), len(result))
			for i, ss := range tc.expResult {
				require.Equal(t, ss, result[i])
			}
		})
	}
}

func TestRepository_GetUsers(t *testing.T) {
	tcs := map[string]struct {
		expResult models.UserSlice
		expError  error
	}{
		"successfully get all users": {
			expResult: models.UserSlice{
				&models.User{Name: "john", Email: "john@example.com"},
				&models.User{Name: "andy", Email: "andy@example.com"},
				&models.User{Name: "common", Email: "common@example.com"},
				&models.User{Name: "lisa", Email: "lisa@example.com"},
				&models.User{Name: "kate", Email: "kate@example.com"},
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			ctx := context.Background()
			db, err := config.NewDatabase()
			require.NoError(t, err)
			repo := NewDBRepo(db)

			// load testdata
			loadSqlTestFile(t, db, "testdata/friends.sql")
			result, err := repo.GetUsers(ctx)
			require.NoError(t, err)
			require.Equal(t, len(tc.expResult), len(result))
			for i, ss := range tc.expResult {
				require.Equal(t, ss.Email, result[i].Email)
			}
		})
	}
}
