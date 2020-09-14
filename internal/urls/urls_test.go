package urls

import "testing"

const (
	urlExample1 = "https://www.reddit.com/r/golang/comments/is81vi/just_because_you_can_doesnt_mean_you_should"
	urlExample2 = "https://www.reddit.com/r/golang/comments/isg3a7/bitmealum_an_e2e_encrypted_email_alternative"
	urlExample3 = "https://www.reddit.com/r/golang/comments/is8eo9/godiagrams_create_architecture_diagrams_with_go"
)

func TestShorten(t *testing.T) {
	type args struct {
		longURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "given long URL example #1, it should return an 8 character long deterministic hash",
			args: args{
				longURL: urlExample1,
			},
			want: "043fdd58",
		},
		{
			name: "given long URL example #2, it should return an 8 character long deterministic hash",
			args: args{
				longURL: urlExample2,
			},
			want: "de3d8e77",
		},
		{
			name: "given long URL example #3, it should return an 8 character long deterministic hash",
			args: args{
				longURL: urlExample3,
			},
			want: "13555926",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Shorten(tt.args.longURL); got != tt.want {
				t.Errorf("Shorten() = %v, want %v", got, tt.want)
			}
		})
	}
}
