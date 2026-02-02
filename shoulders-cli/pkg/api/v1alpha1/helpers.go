package v1alpha1

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func TypeMeta(kind string) v1.TypeMeta {
	return v1.TypeMeta{
		APIVersion: Group + "/" + Version,
		Kind:       kind,
	}
}

func ObjectMeta(name, namespace string) v1.ObjectMeta {
	meta := v1.ObjectMeta{Name: name}
	if namespace != "" {
		meta.Namespace = namespace
	}
	return meta
}
