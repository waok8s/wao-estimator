package controllers

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	v1beta1 "github.com/Nedopro2022/wao-estimator/api/v1beta1"
)

func TestGetFieldValue(t *testing.T) {
	type args struct {
		f    v1beta1.Field
		node *corev1.Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"default_empty", args{
			v1beta1.Field{Default: ""},
			nil},
			""},
		{"default", args{
			v1beta1.Field{Default: "a"},
			nil},
			"a"},
		{"label", args{
			v1beta1.Field{Default: "a", Override: &v1beta1.FieldRef{Label: pointer.StringPtr("xxx")}},
			&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "n0", Labels: map[string]string{"xxx": "b"}}}},
			"b"},
		{"label_wrong", args{
			v1beta1.Field{Default: "a", Override: &v1beta1.FieldRef{Label: pointer.StringPtr("xxx")}},
			&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "n0", Labels: map[string]string{"yyy": "b"}}}},
			"a"},
		{"label_nothing", args{
			v1beta1.Field{Default: "a", Override: &v1beta1.FieldRef{Label: pointer.StringPtr("xxx")}},
			&corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "n0", Labels: map[string]string{}}}},
			"a"},
		{"label_empty_node", args{
			v1beta1.Field{Default: "a", Override: &v1beta1.FieldRef{Label: pointer.StringPtr("xxx")}},
			&corev1.Node{}},
			"a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFieldValue(tt.args.f, tt.args.node); got != tt.want {
				t.Errorf("GetFieldValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
