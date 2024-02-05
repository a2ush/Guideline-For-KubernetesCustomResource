/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	samplea2ushdevv1 "my.domain/guideline/api/v1"
)

// SampleObjectReconciler reconciles a SampleObject object
type SampleObjectReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=sample.a2ush.dev.a2ush.dev,resources=sampleobjects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sample.a2ush.dev.a2ush.dev,resources=sampleobjects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sample.a2ush.dev.a2ush.dev,resources=sampleobjects/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SampleObject object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *SampleObjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrllog.FromContext(ctx)

	// Catch events related to your custom resource
	instance := &samplea2ushdevv1.SampleObject{}
	err := r.Get(ctx, req.NamespacedName, instance)

	//対象の custom resource が etcd から削除されている場合
	if err != nil && errors.IsNotFound(err) {

		// file を削除する処理を書く
		r.Recorder.Event(instance, "Normal", "Deleted", fmt.Sprintf("Deleted resource %s", req.NamespacedName.String()))
		return reconcile.Result{}, nil
	}

	desire := instance.DeepCopy()
	fmt.Printf("Resource. Filename/Reason: %s/%s \n", instance.Spec.Filename, instance.Spec.Reason)

	// filename が変更された場合
	if desire.Status.Filename != instance.Spec.Filename {
		// Status.Filename "xxx" means the filename already exists.
		if desire.Status.Filename != "" {

			//最初に新しいファイルを作成する
			//create file
			//その後に既存のファイルを削除する
			// delete file

			// 検知された SampleObject に紐づく Filename を新しい filename に変更する（オブジェクトの Status を変更する）
			if err = r.changeStatus(desire, instance); err != nil {
				return reconcile.Result{}, err
			}
			r.Recorder.Event(instance, "Normal", "Updated", fmt.Sprintf("Updated resource. Filename/Reason: %s/%s", desire.Status.Filename, desire.Status.Reason))

			// Status.Filename "xxxx" means the file already exists.it's the first created time!
		} else {

			// ここに File を新しく作る処理を書く（controller がファイルを作るお話）
			if err = r.changeStatus(desire, instance); err != nil {
				return reconcile.Result{}, err
			}
			r.Recorder.Event(instance, "Normal", "Created", fmt.Sprintf("Created resource. Filename/Reason: %s/%s", desire.Status.Filename, desire.Status.Reason))

		}
	}

	// reason が変更された場合
	if desire.Status.Reason != instance.Spec.Reason {
		if err = r.changeStatus(desire, instance); err != nil {
			return reconcile.Result{}, err
		}
		r.Recorder.Event(instance, "Normal", "Updated", fmt.Sprintf("Updated resource. Filename/Reason: %s/%s", desire.Status.Filename, desire.Status.Reason))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SampleObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplea2ushdevv1.SampleObject{}).
		Complete(r)
}

func (r *SampleObjectReconciler) changeStatus(desire, instance *samplea2ushdevv1.SampleObject) error {

	desire.Status.Filename = instance.Spec.Filename
	desire.Status.Reason = instance.Spec.Reason
	if err := r.Status().Update(context.TODO(), desire); err != nil {
		return err
	}
	return nil
}
