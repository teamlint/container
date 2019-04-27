package inject

import (
	"net"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestContainer_Provide(t *testing.T) {
	t.Run("function", func(t *testing.T) {
		container, err := New(
			Provide(func() *http.Server {
				return &http.Server{}
			}),
			Provide(func(server *http.Server) *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "test",
				}
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("function with error", func(t *testing.T) {
		container, err := New(
			Provide(func() (*net.TCPAddr, error) {
				return &net.TCPAddr{
					Zone: "test",
				}, errors.New("build error")
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.EqualError(t, container.Populate(&addr), "*net.TCPAddr: build error")
	})

	t.Run("function with nil error", func(t *testing.T) {
		container, err := New(
			Provide(func() (*net.TCPAddr, error) {
				return &net.TCPAddr{
					Zone: "test",
				}, nil
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.NotNil(t, addr)
		require.Equal(t, "test", addr.Zone)
	})

	t.Run("function without arguments", func(t *testing.T) {
		_, err := New(
			Provide(func() {}),
		)

		// todo: improve error message
		require.EqualError(t, err, "could not compile container: provide failed: provider must be a function with returned value and optional error")
	})

	// todo: implement struct provide
	t.Run("struct", func(t *testing.T) {
		type StructProvider struct {
			TCPAddr *net.TCPAddr `inject:""`
			UDPAddr *net.UDPAddr `inject:""`
		}

		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{Zone: "tcp"}
			}),
			Provide(func() *net.UDPAddr {
				return &net.UDPAddr{Zone: "udp"}
			}),
			Provide(&StructProvider{}),
		)

		require.NoError(t, err)

		var sp *StructProvider
		require.NoError(t, container.Populate(&sp))
		require.Equal(t, "tcp", sp.TCPAddr.Zone)
		require.Equal(t, "udp", sp.UDPAddr.Zone)
	})
}

func TestContainer_ProvideAs(t *testing.T) {
	t.Run("provide as", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "test",
				}
			}, As(new(net.Addr))),
		)

		require.NoError(t, err)

		var addr net.Addr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "test", addr.(*net.TCPAddr).Zone)
	})

	t.Run("provide as struct", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, As(http.Server{})),
		)

		require.EqualError(t, err, "could not compile container: provide failed: argument for As() must be pointer to interface type, got http.Server")
	})

	t.Run("provide as struct pointer", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, As(new(http.Server))),
		)

		require.EqualError(t, err, "could not compile container: provide failed: argument for As() must be pointer to interface type, got *http.Server")
	})

	t.Run("provide as not implemented interface", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{}
			}, As(new(http.Handler))),
		)

		require.EqualError(t, err, "could not compile container: provide failed: *net.TCPAddr not implement http.Handler interface")
	})
}

func TestContainer_Apply(t *testing.T) {
	t.Run("apply function", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) {
				addr.Zone = "two"
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "two", addr.Zone)
	})

	t.Run("apply without result", func(t *testing.T) {
		container, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) {
				addr.Zone = "two"
			}),
		)

		require.NoError(t, err)

		var addr *net.TCPAddr
		require.NoError(t, container.Populate(&addr))
		require.Equal(t, "two", addr.Zone)
	})

	t.Run("apply error", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) (err error) {
				return errors.New("boom")
			}),
		)

		require.EqualError(t, err, "could not compile container: apply error: boom")
	})

	t.Run("apply incorrect function", func(t *testing.T) {
		_, err := New(
			Provide(func() *net.TCPAddr {
				return &net.TCPAddr{
					Zone: "one",
				}
			}),
			Apply(func(addr *net.TCPAddr) (s string) {
				return "string"
			}),
		)

		require.EqualError(t, err, "could not compile container: modifier must be a function with optional error as result")
	})
}

//
// func TestContainer_ProvideName(t *testing.T) {
// 	t.Run("provide two named implementations as one interface", func(t *testing.T) {
// 		var container = &Container{}
//
// 		require.NoError(t, container.Provide(func() *net.TCPAddr {
// 			return &net.TCPAddr{}
// 		}, As(new(net.Addr)), Name("tcp")))
//
// 		require.NoError(t, container.Provide(func() *net.UDPAddr {
// 			return &net.UDPAddr{}
// 		}, As(new(net.Addr)), Name("udp")))
// 	})
//
// 	t.Run("provide two implementations as one interface without name", func(t *testing.T) {
// 		var container = &Container{}
//
// 		require.NoError(t, container.Provide(func() *net.TCPAddr {
// 			return &net.TCPAddr{}
// 		}, As(new(net.Addr))))
//
// 		require.NoError(t, container.Provide(func() *net.UDPAddr {
// 			return &net.UDPAddr{}
// 		}, As(new(net.Addr))))
// 	})
// }
