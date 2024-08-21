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
	//"bytes"
	"context"

	//"encoding/json"
	//"fmt"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	//v1 "k8s.io/client-go/applyconfigurations/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)


// TestMonitorReconciler reconciles a TestMonitor object
type TestMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}
var existService = make(map[string]string,0)


// +kubebuilder:rbac:groups=devops.pixocial.io,resources=testmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.pixocial.io,resources=testmonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=devops.pixocial.io,resources=testmonitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TestMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *TestMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	 _ = log.FromContext(ctx)
	 //env := "beta"
	 extexNs := []string{"kube-system","observable"}
	service:= &corev1.ServiceList{}
	for _, svc := range service.Items{
		if !isExcludedNamespace(svc.Namespace,extexNs) {
			existService[svc.GetName()+svc.GetNamespace()] = svc.Spec.ClusterIP
		}
	}


	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}
func isExcludedNamespace(ns string, nsList []string)bool{
	for _, name := range nsList {
		if name == ns {
			return true
		}
	}
	return false
}
func GetPrivateZOneList(){

}

func (r *TestMonitorReconciler) createSvcHandeler(e event.CreateEvent) bool {
	ctx := context.TODO() // 根据需要提供 context
	//log.FromContext(ctx).Info("Service created", "name",e)

	obj := e.Object
	// 尝试将 obj 断言为 corev1.Service
	v1svc, _ := obj.(*corev1.Service)

	log.FromContext(ctx).Info("type","type",reflect.TypeOf(e))

	//log.FromContext(ctx).Info("type","type",v1svc.Name)
	//log.FromContext(ctx).Info("type","type",v1svc.Spec)
	//log.FromContext(ctx).Info("type","type",v1svc.Namespace)
	if _,ok := existService[v1svc.GetName()+v1svc.GetNamespace()];ok{
		log.FromContext(ctx).Info("已存在")
	}else{
		log.FromContext(ctx).Info("需要新增")
		existService[v1svc.GetName()+v1svc.GetNamespace()] = v1svc.Spec.ClusterIP
	}
	// 实现你的创建逻辑
	return true
}

func (r *TestMonitorReconciler) deleteSvcHandeler(e event.DeleteEvent) bool {
	log.FromContext(context.TODO()).Info("Service deleted", "name",e)

	// 实现你的删除逻辑
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Service{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: r.createSvcHandeler,
			DeleteFunc: r.deleteSvcHandeler,
		}).
		//Watches(&corev1.Service{},handler.Funcs{
		//	CreateFunc: r.createSvcHandeler,
		//	DeleteFunc: r.createSvcHandeler,
	    //}).
		Complete(r)
}

