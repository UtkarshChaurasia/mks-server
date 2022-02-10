// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2022 Satyam Bhardwaj <sabhardw@redhat.com>
// SPDX-FileCopyrightText: 2022 Utkarsh Chaurasia <uchauras@redhat.com>
// SPDX-FileCopyrightText: 2022 Avinal Kumar <avinkuma@redhat.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mkstaskrun

import (
	"fmt"
	"time"

	"github.com/MiniTeks/mks-server/pkg/apis/mkscontroller/v1alpha1"
	clientset "github.com/MiniTeks/mks-server/pkg/client/clientset/versioned"
	informer "github.com/MiniTeks/mks-server/pkg/client/informers/externalversions/mkscontroller/v1alpha1"
	"github.com/MiniTeks/mks-server/pkg/db"
	"github.com/MiniTeks/mks-server/pkg/tconfig"
	"github.com/go-redis/redis/v8"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var rClient *redis.Client

type Controller struct {
	kubeclientset  kubernetes.Interface
	mksclientset   clientset.Interface
	mksTaskRunSync cache.InformerSynced
	queue          workqueue.RateLimitingInterface
}

func NewController(kubeclientset kubernetes.Interface,
	mksclientset clientset.Interface,
	mksinformer informer.MksTaskRunInformer, redisClient *redis.Client) *Controller {
	rClient = redisClient
	controller := &Controller{
		kubeclientset:  kubeclientset,
		mksclientset:   mksclientset,
		mksTaskRunSync: mksinformer.Informer().HasSynced,
		queue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mks-taskrun-controller"),
	}

	mksinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addController,
		UpdateFunc: controller.updateController,
		DeleteFunc: controller.deleteController,
	})

	return controller
}

func (c *Controller) addController(obj interface{}) {
	fmt.Println("MksTaskRun has been created")

	tp := &tconfig.TektonParam{}
	tcl, err := tp.Client()
	if err != nil {
		fmt.Errorf("Cannot connect to Tekton client: %w", err)
		return
	}
	var crtobj = obj.(*v1alpha1.MksTaskRun)
	ttr, err := Create(tcl, crtobj, metav1.CreateOptions{}, crtobj.Namespace)
	if err != nil {
		db.Increment(rClient, "mksTaskRunfailed")
		fmt.Errorf("Cannot create Tekton TaskRun: %w", err)
		return
	} else {
		db.Increment(rClient, "mksTaskRuncreated")
	}

	fmt.Printf("Successfully created Tekton TaskRun: %s\n", ttr.Name)
	c.queue.Add(obj)
}

func (c *Controller) updateController(oldObj, newObj interface{}) {
	fmt.Println("MksTaskRun has been updated")
}

func (c *Controller) deleteController(obj interface{}) {

	tp := &tconfig.TektonParam{}
	tcl, err := tp.Client()
	if err != nil {
		fmt.Errorf("Cannot connect to Tekton client: %w", err)
		return
	}

	var delobj = obj.(*v1alpha1.MksTaskRun)
	delerr := Delete(tcl, delobj.Name, metav1.DeleteOptions{}, delobj.Namespace)
	if delerr != nil {
		fmt.Errorf("Cannot delete MksTaskRun: %v", delerr)
		return
	} else {
		fmt.Println("MksTaskRun has been deleted")
		db.Increment(rClient, "mksTaskRundeleted")
	}
	c.queue.Add(obj)
}

func (c *Controller) Run(ch <-chan struct{}) {
	fmt.Println("starting controller")
	if !cache.WaitForCacheSync(ch, c.mksTaskRunSync) {
		fmt.Print("waiting for cache to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func (c *Controller) worker() {
	for c.processItem() {

	}
}

func (c *Controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	return true
}
