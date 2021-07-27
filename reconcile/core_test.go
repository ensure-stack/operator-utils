package reconcile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ensure-stack/operator-utils/testutils"
)

func TestNamespace(t *testing.T) {
	t.Run("creates", func(t *testing.T) {
		c := testutils.NewClient()
		ctx := context.Background()
		log := testutils.NewLogger(t)

		sa := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-sa",
				Namespace: "test-ns",
			},
		}

		notFound := errors.NewNotFound(schema.GroupResource{}, "")
		c.
			On("Get", mock.Anything, client.ObjectKey{
				Name:      "test-sa",
				Namespace: "test-ns",
			}, mock.Anything).
			Return(notFound)

		c.
			On("Create", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		_, err := Namespace(ctx, log, c, sa)
		require.NoError(t, err)
		c.AssertCalled(
			t, "Create", mock.Anything, sa, mock.Anything)
	})
}

func TestServiceAccount(t *testing.T) {
	t.Run("creates", func(t *testing.T) {
		c := testutils.NewClient()
		ctx := context.Background()
		log := testutils.NewLogger(t)

		sa := &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-sa",
				Namespace: "test-ns",
			},
		}

		notFound := errors.NewNotFound(schema.GroupResource{}, "")
		c.
			On("Get", mock.Anything, client.ObjectKey{
				Name:      "test-sa",
				Namespace: "test-ns",
			}, mock.Anything).
			Return(notFound)

		c.
			On("Create", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		_, err := ServiceAccount(ctx, log, c, sa)
		require.NoError(t, err)
		c.AssertCalled(
			t, "Create", mock.Anything, sa, mock.Anything)
	})

	t.Run("updates", func(t *testing.T) {
		t.Run("Secrets", func(t *testing.T) {
			c := testutils.NewClient()
			ctx := context.Background()
			log := testutils.NewLogger(t)

			desiredSA := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sa",
					Namespace: "test-ns",
				},
				Secrets: []corev1.ObjectReference{
					{Name: ""},
				},
			}

			existingSA := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sa",
					Namespace: "test-ns",
				},
			}

			c.
				On("Get", mock.Anything, client.ObjectKey{
					Name:      "test-sa",
					Namespace: "test-ns",
				}, mock.Anything).
				Run(func(args mock.Arguments) {
					arg := args.Get(2).(*corev1.ServiceAccount)
					existingSA.DeepCopyInto(arg)
				}).
				Return(nil)

			c.
				On("Update", mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			_, err := ServiceAccount(ctx, log, c, desiredSA)
			require.NoError(t, err)
		})

		t.Run("ImagePullSecrets", func(t *testing.T) {
			c := testutils.NewClient()
			ctx := context.Background()
			log := testutils.NewLogger(t)

			desiredSA := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sa",
					Namespace: "test-ns",
				},
				ImagePullSecrets: []corev1.LocalObjectReference{
					{Name: ""},
				},
			}

			existingSA := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sa",
					Namespace: "test-ns",
				},
			}

			c.
				On("Get", mock.Anything, client.ObjectKey{
					Name:      "test-sa",
					Namespace: "test-ns",
				}, mock.Anything).
				Run(func(args mock.Arguments) {
					arg := args.Get(2).(*corev1.ServiceAccount)
					existingSA.DeepCopyInto(arg)
				}).
				Return(nil)

			c.
				On("Update", mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			_, err := ServiceAccount(ctx, log, c, desiredSA)
			require.NoError(t, err)
		})

		t.Run("AutomountServiceAccountToken", func(t *testing.T) {
			c := testutils.NewClient()
			ctx := context.Background()
			log := testutils.NewLogger(t)

			desiredSA := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sa",
					Namespace: "test-ns",
				},
				AutomountServiceAccountToken: pointer.BoolPtr(false),
			}

			existingSA := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sa",
					Namespace: "test-ns",
				},
			}

			c.
				On("Get", mock.Anything, client.ObjectKey{
					Name:      "test-sa",
					Namespace: "test-ns",
				}, mock.Anything).
				Run(func(args mock.Arguments) {
					arg := args.Get(2).(*corev1.ServiceAccount)
					existingSA.DeepCopyInto(arg)
				}).
				Return(nil)

			c.
				On("Update", mock.Anything, mock.Anything, mock.Anything).
				Return(nil)

			_, err := ServiceAccount(ctx, log, c, desiredSA)
			require.NoError(t, err)
		})
	})
}
