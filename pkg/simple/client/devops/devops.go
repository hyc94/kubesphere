/*
Copyright 2018 The KubeSphere Authors.
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

package devops

import (
	"fmt"
	"k8s.io/klog"
	"kubesphere.io/kubesphere/pkg/gojenkins"
	"sync"
)

const (
	JenkinsAllUserRoleName = "kubesphere-user"
)

type DevopsClient struct {
	jenkinsClient *gojenkins.Jenkins
}

// 创建DevopsClient
func NewDevopsClient(options *DevopsOptions) (*DevopsClient, error) {
	var d DevopsClient

	// 创建jenkins实例
	jenkins := gojenkins.CreateJenkins(nil, options.Host, options.MaxConnections, options.Username, options.Password)
	// 检查jenkins是否可以连接
	jenkins, err := jenkins.Init()
	if err != nil {
		klog.Errorf("failed to connecto to jenkins role, %+v", err)
		return nil, err
	}

	d.jenkinsClient = jenkins

	// 初始化jenkins
	err = d.initializeJenkins()
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return &d, nil
}

func NewDevopsClientOrDie(options *DevopsOptions) *DevopsClient {
	jenkins := gojenkins.CreateJenkins(nil, options.Host, options.MaxConnections, options.Username, options.Password)
	jenkins, err := jenkins.Init()
	if err != nil {
		klog.Errorf("failed to connecto to jenkins role, %+v", err)
		panic(err)
	}

	d := &DevopsClient{
		jenkinsClient: jenkins,
	}

	err = d.initializeJenkins()
	if err != nil {
		klog.Error(err)
		panic(err)
	}

	return d
}

func (c *DevopsClient) Jenkins() *gojenkins.Jenkins {
	return c.jenkinsClient
}

var mutex = sync.Mutex{}

// 初始化jenkins，主要是在Global Role和Project Role中新增kubesphere通用角色
func (c *DevopsClient) initializeJenkins() error {
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()

	if c.jenkinsClient == nil {
		return fmt.Errorf("jenkins intialization failed")
	}

	// 获取jenkins的Global Roles中是否存在kubesphere-user
	// 调用接口http://ks-jenkins.kubesphere-devops-system.svc/role-strategy/strategy/getRole/
	globalRole, err := c.jenkinsClient.GetGlobalRole(JenkinsAllUserRoleName)
	if err != nil {
		klog.Error(err)
		return err
	}

	// Jenkins uninitialized, create global role
	if globalRole == nil {
		// 在Global Roles新增kubesphere-user
		// 调用接口http://ks-jenkins.kubesphere-devops-system.svc/role-strategy/strategy/addRole
		_, err := c.jenkinsClient.AddGlobalRole(JenkinsAllUserRoleName, gojenkins.GlobalPermissionIds{GlobalRead: true}, true)
		if err != nil {
			klog.Error(err)
			return err
		}
	}

	// 在Project Roles新增kubesphere-user：AddProjectRole()
	// 调用接口http://ks-jenkins.kubesphere-devops-system.svc/role-strategy/strategy/addRole
	_, err = c.jenkinsClient.AddProjectRole(JenkinsAllUserRoleName, "\\n\\s*\\r", gojenkins.ProjectPermissionIds{SCMTag: true}, true)
	if err != nil {
		klog.Error(err)
		return err
	}

	return nil
}
