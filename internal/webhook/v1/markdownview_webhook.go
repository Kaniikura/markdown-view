/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	viewv1 "github.com/Kaniikura/markdown-view/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// nolint:unused
// log is for logging in this package.
var markdownviewlog = logf.Log.WithName("markdownview-resource")

// SetupMarkdownViewWebhookWithManager registers the webhook for MarkdownView in the manager.
func SetupMarkdownViewWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&viewv1.MarkdownView{}).
		WithValidator(&MarkdownViewCustomValidator{}).
		WithDefaulter(&MarkdownViewCustomDefaulter{}).
		Complete()
}

func logValidationError(markdownview *viewv1.MarkdownView, errs field.ErrorList) error {
	err := apierrors.NewInvalid(schema.GroupKind{Group: "view.kaniikura.github.io", Kind: "MarkdownView"}, markdownview.Name, errs)
	markdownviewlog.Error(err, "validation error", "name", markdownview.Name)
	return err
}

// +kubebuilder:webhook:path=/mutate-view-kaniikura-github-io-v1-markdownview,mutating=true,failurePolicy=fail,sideEffects=None,groups=view.kaniikura.github.io,resources=markdownviews,verbs=create;update,versions=v1,name=mmarkdownview-v1.kb.io,admissionReviewVersions=v1

// MarkdownViewCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind MarkdownView when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type MarkdownViewCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &MarkdownViewCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind MarkdownView.
func (d *MarkdownViewCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	markdownview, ok := obj.(*viewv1.MarkdownView)
	if !ok {
		return fmt.Errorf("expected an MarkdownView object but got %T", obj)
	}
	markdownviewlog.Info("Defaulting for MarkdownView", "name", markdownview.GetName())

	if len(markdownview.Spec.ViewerImage) == 0 {
		markdownview.Spec.ViewerImage = "peaceiris/mdbook:latest"
	}

	return nil
}

// +kubebuilder:webhook:path=/validate-view-kaniikura-github-io-v1-markdownview,mutating=false,failurePolicy=fail,sideEffects=None,groups=view.kaniikura.github.io,resources=markdownviews,verbs=create;update;delete,versions=v1,name=vmarkdownview-v1.kb.io,admissionReviewVersions=v1

// MarkdownViewCustomValidator struct is responsible for validating the MarkdownView resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type MarkdownViewCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &MarkdownViewCustomValidator{}

func (v *MarkdownViewCustomValidator) validate(obj *viewv1.MarkdownView) (admission.Warnings, error) {
	var errs field.ErrorList

	if obj.Spec.Replicas < 1 || obj.Spec.Replicas > 5 {
		errs = append(errs, field.Invalid(field.NewPath("spec", "replicas"), obj.Spec.Replicas, "replicas must be in the range of 1 to 5."))
	}

	if _, ok := obj.Spec.Markdowns["SUMMARY.md"]; !ok {
		errs = append(errs, field.Required(field.NewPath("spec", "markdowns"), "markdowns must have SUMMARY.md."))
	}

	if len(errs) > 0 {
		return nil, logValidationError(obj, errs)
	}

	return nil, nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type MarkdownView.
func (v *MarkdownViewCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	markdownview, ok := obj.(*viewv1.MarkdownView)
	if !ok {
		return nil, fmt.Errorf("expected a MarkdownView object but got %T", obj)
	}
	markdownviewlog.Info("Validation for MarkdownView upon creation", "name", markdownview.GetName())

	return v.validate(markdownview)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type MarkdownView.
func (v *MarkdownViewCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	markdownview, ok := newObj.(*viewv1.MarkdownView)
	if !ok {
		return nil, fmt.Errorf("expected a MarkdownView object for the newObj but got %T", newObj)
	}
	markdownviewlog.Info("Validation for MarkdownView upon update", "name", markdownview.GetName())

	return v.validate(markdownview)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type MarkdownView.
func (v *MarkdownViewCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	markdownview, ok := obj.(*viewv1.MarkdownView)
	if !ok {
		return nil, fmt.Errorf("expected a MarkdownView object but got %T", obj)
	}
	markdownviewlog.Info("Validation for MarkdownView upon deletion", "name", markdownview.GetName())

	// No specific validation needed for deletion, but logging is added for traceability.
	return nil, nil
}
