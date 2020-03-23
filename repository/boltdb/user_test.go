package boltdb

import (
	"testing"

	"github.com/ImageWare/TLSential/model"
	"github.com/boltdb/bolt"
)

const TestDBPath = ".test.db"

func TestKvstore_BoltStore(t *testing.T) {
	db, err := bolt.Open(TestDBPath, 0666, nil)
	if err != nil {
		t.Fatalf("Error opening test db: %s", err.Error())
	}
	defer db.Close()

	ur, err := NewUserRepository(db)
	if err != nil {
		t.Fatalf("Error on NewUserRepository: %s", err.Error())
	}

	t.Run("USERS", func(t *testing.T) {
		t.Run("Get Nonexistant", func(t *testing.T) {
			u, err := ur.GetUser("test")
			if err != nil {
				t.Fatal(err)
			}

			if u != nil {
				t.Error("Expected nil user returned.")
			}
		})

		u1 := &model.User{Name: "test-name", Hash: "test-hash", Role: "admin"}
		u2 := &model.User{Name: "test-name2", Hash: "test-hash", Role: "admin"}
		u3 := &model.User{Name: "test-name3", Hash: "test-hash", Role: "admin"}

		t.Run("Save User", func(t *testing.T) {
			err := ur.SaveUser(u1)
			if err != nil {
				t.Fatal(err)
			}
			err = ur.SaveUser(u2)
			if err != nil {
				t.Fatal(err)
			}
			err = ur.SaveUser(u3)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Get Existing", func(t *testing.T) {
			u, err := ur.GetUser("test-name")
			if err != nil {
				t.Fatal(err)
			}

			if u == nil {
				t.Error("Unexpected nil user returned.")
			}

			// expected, got
			assertEqualUser(t, u1, u)
		})

		t.Run("Get All Users", func(t *testing.T) {
			u, err := ur.GetAllUsers()
			if err != nil {
				t.Fatal(err)
			}
			if len(u) != 3 {
				t.Error("Expected 3 items at this point")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := ur.DeleteUser("test-user")
			if err != nil {
				t.Fatal(err)
			}
			u, err := ur.GetUser("test-user")
			if err != nil {
				t.Fatal(err)
			}
			if u != nil {
				t.Error("Expected nil returned after delete.")
			}

			// Should also work on non existent ident.
			err = ur.DeleteUser("test-user4")
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Delete All Users", func(t *testing.T) {
			err := ur.DeleteAllUsers()
			if err != nil {
				t.Fatal(err)
			}
			// Test by getting something that was there
			u, err := ur.GetUser("test-user2")
			if err != nil {
				t.Fatal(err)
			}
			if u != nil {
				t.Error("Expected nil returned after delete all.")
			}
		})
	})

}

func assertEqualUser(t *testing.T, expect *model.User, got *model.User) {
	t.Helper()
	if expect.Name != got.Name {
		t.Errorf("Mismatched Names:\n\tgot %s \n\twant %s", got.Name, expect.Name)
	}
	if expect.Hash != got.Hash {
		t.Errorf("Mismatched Hashes:\n\tgot %s \n\twant %s", got.Hash, expect.Hash)
	}
	if expect.Role != got.Role {
		t.Errorf("Mismatched Roles:\n\tgot %s \n\twant %s", got.Role, expect.Role)
	}
}

func testEq(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
