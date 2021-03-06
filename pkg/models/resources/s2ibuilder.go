/*

 Copyright 2019 The KubeSphere Authors.

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

package resources

import (
	"github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"kubesphere.io/kubesphere/pkg/informers"
	"kubesphere.io/kubesphere/pkg/server/params"
	"sort"
)

type s2iBuilderSearcher struct {
}

func (*s2iBuilderSearcher) get(namespace, name string) (interface{}, error) {
	return informers.S2iSharedInformerFactory().Devops().V1alpha1().S2iBuilders().Lister().S2iBuilders(namespace).Get(name)
}

// exactly Match
func (*s2iBuilderSearcher) match(kv map[string]string, item *v1alpha1.S2iBuilder) bool {
	for k, v := range kv {
		if !match(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

// Fuzzy searchInNamespace
func (*s2iBuilderSearcher) fuzzy(kv map[string]string, item *v1alpha1.S2iBuilder) bool {
	for k, v := range kv {
		if !fuzzy(k, v, item.ObjectMeta) {
			return false
		}
	}
	return true
}

func (*s2iBuilderSearcher) compare(a, b *v1alpha1.S2iBuilder, orderBy string) bool {
	return compare(a.ObjectMeta, b.ObjectMeta, orderBy)
}

func (s *s2iBuilderSearcher) search(namespace string, conditions *params.Conditions, orderBy string, reverse bool) ([]interface{}, error) {
	s2iBuilders, err := informers.S2iSharedInformerFactory().Devops().V1alpha1().S2iBuilders().Lister().S2iBuilders(namespace).List(labels.Everything())

	if err != nil {
		return nil, err
	}

	result := make([]*v1alpha1.S2iBuilder, 0)

	if len(conditions.Match) == 0 && len(conditions.Fuzzy) == 0 {
		result = s2iBuilders
	} else {
		for _, item := range s2iBuilders {
			if s.match(conditions.Match, item) && s.fuzzy(conditions.Fuzzy, item) {
				result = append(result, item)
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if reverse {
			tmp := i
			i = j
			j = tmp
		}
		return s.compare(result[i], result[j], orderBy)
	})

	r := make([]interface{}, 0)
	for _, i := range result {
		r = append(r, i)
	}
	return r, nil
}
