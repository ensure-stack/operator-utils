package reconcile

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Namespace reconciles a /v1, Kind=Namespace.
func Namespace(
	ctx context.Context,
	log logr.Logger,
	c client.Client,
	desiredNamespace *corev1.Namespace,
) (
	actualNamespace *corev1.Namespace,
	err error,
) {
	key := client.ObjectKey{
		Name:      desiredNamespace.Name,
		Namespace: desiredNamespace.Namespace,
	}

	// Lookup current version of the object
	actualNamespace = &corev1.Namespace{}
	err = c.Get(ctx, key, actualNamespace)
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("getting Namespace: %w", err)
	}

	if errors.IsNotFound(err) {
		// Namespace needs to be created
		log.V(1).Info("creating", "Namespace", key.String())
		if err = c.Create(ctx, desiredNamespace); err != nil {
			return nil, fmt.Errorf("creating Namespace: %w", err)
		}
		return desiredNamespace, nil
	}

	return actualNamespace, nil
}

// ServiceAccount reconciles a /v1, Kind=ServiceAccount.
func ServiceAccount(
	ctx context.Context,
	log logr.Logger,
	c client.Client,
	desiredServiceAccount *corev1.ServiceAccount,
) (
	actualServiceAccount *corev1.ServiceAccount,
	err error,
) {
	key := client.ObjectKey{
		Name:      desiredServiceAccount.Name,
		Namespace: desiredServiceAccount.Namespace,
	}

	// Lookup current version of the object
	actualServiceAccount = &corev1.ServiceAccount{}
	err = c.Get(ctx, key, actualServiceAccount)
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("getting ServiceAccount: %w", err)
	}

	if errors.IsNotFound(err) {
		// ServiceAccount needs to be created
		log.V(1).Info("creating", "ServiceAccount", key.String())
		if err = c.Create(ctx, desiredServiceAccount); err != nil {
			return nil, fmt.Errorf("creating ServiceAccount: %w", err)
		}
		return desiredServiceAccount, nil
	}

	// Check ServiceAccount for update
	var needsUpdate bool
	if !equality.Semantic.DeepEqual(
		actualServiceAccount.Secrets, desiredServiceAccount.Secrets) {
		actualServiceAccount.Secrets = desiredServiceAccount.Secrets
		needsUpdate = true
	}
	if !equality.Semantic.DeepEqual(
		actualServiceAccount.ImagePullSecrets, desiredServiceAccount.ImagePullSecrets) {
		actualServiceAccount.ImagePullSecrets = desiredServiceAccount.ImagePullSecrets
		needsUpdate = true
	}

	if desiredServiceAccount.AutomountServiceAccountToken == nil &&
		actualServiceAccount.AutomountServiceAccountToken != nil {
		actualServiceAccount.AutomountServiceAccountToken = nil
		needsUpdate = true
	}
	if desiredServiceAccount.AutomountServiceAccountToken != nil &&
		actualServiceAccount.AutomountServiceAccountToken == nil ||
		desiredServiceAccount.AutomountServiceAccountToken != nil &&
			*actualServiceAccount.AutomountServiceAccountToken !=
				*desiredServiceAccount.AutomountServiceAccountToken {
		actualServiceAccount.AutomountServiceAccountToken =
			desiredServiceAccount.AutomountServiceAccountToken
		needsUpdate = true
	}

	if needsUpdate {
		log.V(1).Info("updating", "ServiceAccount", key.String())
		if err = c.Update(ctx, actualServiceAccount); err != nil {
			return nil, fmt.Errorf("updating ServiceAccount: %w", err)
		}
	}

	return actualServiceAccount, nil
}
