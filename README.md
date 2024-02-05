## Prerequisite

###  Install kubebuilder

https://book.kubebuilder.io/quick-start.html
```
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

mkdir -p ~/projects/guestbook
cd ~/projects/guestbook
kubebuilder init --domain my.domain --repo my.domain/guestbook
```

In this repo
> kubebuilder init --plugins go/v4 --domain a2ush.dev --repo my.domain/guideline

## Create API

```
$ kubebuilder create api --group sample.a2ush.dev --version v1 --kind SampleObject
INFO Create Resource [y/n]                        
y
INFO Create Controller [y/n]                      
y
...
```

## Create your CRD

Edit `api/v1/sampleobject_types.go`.
```
// SampleObjectSpec defines the desired state of SampleObject
type SampleObjectSpec struct {
	Filename string `json:"filename"`
	// +kubebuilder:default=BLANK
	// +optional
	Reason string `json:"reason"`
}

// SampleObjectStatus defines the observed state of SampleObject
type SampleObjectStatus struct {
	Filename string `json:"filename"`
	// +optional
	Reason string `json:"reason"`
}
```

And you can run `make manifests` for creating CRD to `config/crd/bases/sample.a2ush.dev.a2ush.dev_sampleobjects.yaml`.

## create controller

controllers keep to maintain a resource's "status" in the statue of "desire".

Edit `internal/controller/sampleobject_controller.go`.
```
func (r *SampleObjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

    // Catch events related to your custom resource
	instance := &samplea2ushdevv1.SampleObject{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		r.Recorder.Event(instance, "Normal", "Deleted", fmt.Sprintf("Deleted resource %s", req.NamespacedName.String()))
		return reconcile.Result{}, nil
	}
	desire := instance.DeepCopy()

	return ctrl.Result{}, nil
}
```

### Add new method

```
func (r *SampleObjectReconciler) changeStatus(desire, instance *samplea2ushdevv1.SampleObject) error {
	desire.Status.Filename = instance.Spec.Filename
	desire.Status.Reason = instance.Spec.Reason
	err := r.Status().Update(context.TODO(), desire)

	return err
}
```

### Record kubernetes event

`cmd/main.go`
```
	if err = (&controller.SampleObjectReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("sampleobject-controller"), // Add
```

Edit `internal/controller/sampleobject_controller.go`.
```
r.Recorder.Event(instance, "Normal", "Deleted", fmt.Sprintf("Deleted resource %s", req.NamespacedName.String()))
```

note: if you forget to edit main.go, you may get the following error
```
2024-02-05T14:32:50+09:00	INFO	Observed a panic in reconciler: runtime error: invalid memory address or nil pointer dereference
```