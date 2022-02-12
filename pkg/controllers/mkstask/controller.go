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

package mkstask

import (
	"fmt"
	"log"
	"time"

	"github.com/MiniTeks/mks-server/pkg/apis/mkscontroller/v1alpha1"
	examplecomclientset "github.com/MiniTeks/mks-server/pkg/client/clientset/versioned"
	appsinformers "github.com/MiniTeks/mks-server/pkg/client/informers/externalversions/mkscontroller/v1alpha1"
	appslisters "github.com/MiniTeks/mks-server/pkg/client/listers/mkscontroller/v1alpha1"
	"github.com/MiniTeks/mks-server/pkg/db"
	"github.com/MiniTeks/mks-server/pkg/tconfig"
	"github.com/go-redis/redis/v8"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

var rClient *redis.Client

type controller struct {
	clientset          examplecomclientset.Clientset
	mksTaskLister      appslisters.MksTaskLister
	mksTaskCacheSynced cache.InformerSynced
	queue              workqueue.RateLimitingInterface
}

func NewController(clientset examplecomclientset.Clientset, mksTaskInformer appsinformers.MksTaskInformer, redisClient *redis.Client) *controller {
	rClient = redisClient
	c := &controller{
		clientset:          clientset,
		mksTaskLister:      mksTaskInformer.Lister(),
		mksTaskCacheSynced: mksTaskInformer.Informer().HasSynced,
		queue:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mks-task-controller"),
	}

	mksTaskInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			UpdateFunc: c.handleUpdate,
			DeleteFunc: c.handleDel,
		},
	)
	return c
}

func (c *controller) Run(ch <-chan struct{}) {
	fmt.Println("starting controller")
	if !cache.WaitForCacheSync(ch, c.mksTaskCacheSynced) {
		fmt.Print("waiting for cache to be synced\n")
	}

	go wait.Until(c.worker, 1*time.Second, ch)

	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	return true
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("\n Add handler was called")
	fmt.Println(obj)
	tp := &tconfig.TektonParam{}
	cs, er := tp.Client()
	if er != nil {
		log.Fatalf("Cannot get tekton client: %s", er.Error())
	}
	addObj := obj.(*v1alpha1.MksTask)
	tsk, err := Create(cs, addObj, metav1.CreateOptions{}, addObj.GetObjectMeta().GetNamespace())
	if err != nil {
		db.Increment(rClient, "MKSTASKFAILED")
		fmt.Errorf("Couldn't create tekton task: %s", err.Error())
	} else {
		db.Increment(rClient, "MKSTASKCREATED")
		fmt.Println("tekton task created")
		fmt.Printf("uid %s", tsk.UID)
	}

	c.queue.Add(obj)
}

func (c *controller) handleUpdate(old, obj interface{}) {
	fmt.Println("\n Update handler was called")
	tp := &tconfig.TektonParam{}
	cs, er := tp.Client()
	if er != nil {
		log.Fatalf("Cannot get tekton client: %s", er.Error())
	}
	updObj := obj.(*v1alpha1.MksTask)
	tsk, err := Update(cs, updObj, metav1.UpdateOptions{}, updObj.GetObjectMeta().GetNamespace())
	if err != nil {
		fmt.Errorf("Couldn't update tekton task: %s", err.Error())
	} else {
		fmt.Println("tekton task updated")
		fmt.Printf("uid %s", tsk.UID)
	}
	c.queue.Add(obj)
}

func (c *controller) handleDel(obj interface{}) {
	fmt.Println("\n Delete handler was called")
	tp := &tconfig.TektonParam{}
	cs, er := tp.Client()
	if er != nil {
		log.Fatalf("Cannot get tekton client: %s", er.Error())
	}
	delObj := obj.(*v1alpha1.MksTask)
	err := Delete(cs, delObj.Name, metav1.DeleteOptions{}, delObj.GetObjectMeta().GetNamespace())
	if err != nil {
		log.Fatalf("Cannot delete tekton task!!: %s", err.Error())
	}
	db.Increment(rClient, "MKSTASKDELETED")
	c.queue.Add(obj)
}
