package policies

import (
	"reflect"
	"testing"
)

// TestSjsonPath covers the escaping applied to configuration paths before they
// reach sjson. Segments are literal keys: Kubernetes annotation and label keys
// such as "eks.amazonaws.com/role-arn" carry dots of their own and must address
// one key rather than a chain of nested objects.
func TestSjsonPath(t *testing.T) {
	tests := []struct {
		name string
		path []string
		want string
	}{
		{"empty", nil, ""},
		{"single plain segment", []string{"namespace"}, "namespace"},
		{"nested plain path", []string{"spec", "template", "spec"}, "spec.template.spec"},
		{"numeric segment keeps array index meaning", []string{"spec", "containers", "0", "image"}, "spec.containers.0.image"},
		{"empty segment preserved", []string{"metadata", "", "name"}, "metadata..name"},
		{
			"dotted annotation key",
			[]string{"metadata", "annotations", "eks.amazonaws.com/role-arn"},
			`metadata.annotations.eks\.amazonaws\.com/role-arn`,
		},
		{
			"dotted label key",
			[]string{"metadata", "labels", "kubernetes.io/service-name"},
			`metadata.labels.kubernetes\.io/service-name`,
		},
		{"wildcard star escaped", []string{"a*b"}, `a\*b`},
		{"wildcard question escaped", []string{"a?b"}, `a\?b`},
		{"pipe escaped", []string{"a|b"}, `a\|b`},
		{"hash escaped", []string{"a#b"}, `a\#b`},
		{"at escaped", []string{"a@b"}, `a\@b`},
		{"colon escaped", []string{":2313"}, `\:2313`},
		{"backslash escaped", []string{`a\b`}, `a\\b`},
		{"multiple specials in one segment", []string{"a.b*c"}, `a\.b\*c`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sjsonPath(tt.path); got != tt.want {
				t.Errorf("sjsonPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

// TestApplyConfigRemovesLiteralKeys verifies that a remove action deletes the
// exact key named by its path, leaving the rest of the configuration intact.
func TestApplyConfigRemovesLiteralKeys(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]interface{}
		path   []string
		want   map[string]interface{}
	}{
		{
			name:   "plain key",
			config: map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]interface{}{"roleArn": "v", "keep": "k"}}},
			path:   []string{"configuration", "metadata", "annotations", "roleArn"},
			want:   map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]interface{}{"keep": "k"}}},
		},
		{
			name:   "dotted annotation key",
			config: map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]interface{}{"eks.amazonaws.com/role-arn": "v", "keep": "k"}}},
			path:   []string{"configuration", "metadata", "annotations", "eks.amazonaws.com/role-arn"},
			want:   map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]interface{}{"keep": "k"}}},
		},
		{
			name:   "dotted label key",
			config: map[string]interface{}{"metadata": map[string]interface{}{"labels": map[string]interface{}{"kubernetes.io/service-name": "v", "app": "web"}}},
			path:   []string{"configuration", "metadata", "labels", "kubernetes.io/service-name"},
			want:   map[string]interface{}{"metadata": map[string]interface{}{"labels": map[string]interface{}{"app": "web"}}},
		},
		{
			name:   "dotted key does not create nested objects",
			config: map[string]interface{}{"annotations": map[string]interface{}{"io.kompose.service": "v"}},
			path:   []string{"configuration", "annotations", "io.kompose.service"},
			want:   map[string]interface{}{"annotations": map[string]interface{}{}},
		},
		{
			name:   "nested plain path still resolves",
			config: map[string]interface{}{"spec": map[string]interface{}{"template": map[string]interface{}{"spec": map[string]interface{}{"serviceAccountName": "sa", "nodeName": "n"}}}},
			path:   []string{"configuration", "spec", "template", "spec", "serviceAccountName"},
			want:   map[string]interface{}{"spec": map[string]interface{}{"template": map[string]interface{}{"spec": map[string]interface{}{"nodeName": "n"}}}},
		},
		{
			name:   "numeric segment still removes array element",
			config: map[string]interface{}{"spec": map[string]interface{}{"containers": []interface{}{map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"}}}},
			path:   []string{"configuration", "spec", "containers", "0"},
			want:   map[string]interface{}{"spec": map[string]interface{}{"containers": []interface{}{map[string]interface{}{"name": "b"}}}},
		},
		{
			name:   "absent key is a no-op",
			config: map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]interface{}{"keep": "k"}}},
			path:   []string{"configuration", "metadata", "annotations", "missing.key"},
			want:   map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]interface{}{"keep": "k"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyConfigRemoves(tt.config, []PolicyAction{newComponentRemoveAction("id", tt.path)})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyConfigRemoves() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
