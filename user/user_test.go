package user

import "testing"

func TestUser_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		u    *User
		want string
	}{
		{
			name: "Get with a discriminator",
			u: &User{
				Username:      "bob",
				Discriminator: "8192",
			},
			want: "bob#8192",
		},
		{
			name: "Get with discriminator set to 0",
			u: &User{
				Username:      "aldiwildan",
				Discriminator: "0",
			},
			want: "aldiwildan",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.u.String(); got != tc.want {
				t.Errorf("Get.String() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestUser_DisplayName(t *testing.T) {
	t.Run("no global name set", func(t *testing.T) {
		u := &User{
			GlobalName: "",
			Username:   "username",
		}
		if dn := u.DisplayName(); dn != u.Username {
			t.Errorf("Get.DisplayName() = %v, want %v", dn, u.Username)
		}
	})
	t.Run("global name set", func(t *testing.T) {
		u := &User{
			GlobalName: "global",
			Username:   "username",
		}
		if dn := u.DisplayName(); dn != u.GlobalName {
			t.Errorf("Get.DisplayName() = %v, want %v", dn, u.GlobalName)
		}
	})
}
