/*
* Tencent is pleased to support the open source community by making 蓝鲸 available.
* Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
* Licensed under the MIT License (the ",License"); you may not use this file except
* in compliance with the License. You may obtain a copy of the License at
* http://opensource.org/licenses/MIT
* Unless required by applicable law or agreed to in writing, software distributed under
* the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
* either express or implied. See the License for the specific language governing permissions and
* limitations under the License.
 */

package godriver

import (
	"context"

	"configcenter/src/storage/mongobyc"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var _ mongobyc.Transaction = (*transaction)(nil)

type transaction struct {
	mongocli       *client
	innerSession   mongo.Session
	collectionMaps map[collectionName]mongobyc.CollectionInterface
}

func newSessionTransaction(mongocli *client, clientSession mongo.Session) *transaction {
	return &transaction{
		mongocli:       mongocli,
		innerSession:   clientSession,
		collectionMaps: map[collectionName]mongobyc.CollectionInterface{},
	}
}

func (t *transaction) StartTransaction() error {
	return t.innerSession.StartTransaction()
}
func (t *transaction) AbortTransaction() error {
	return t.innerSession.AbortTransaction(context.TODO())
}
func (t *transaction) CommitTransaction() error {
	return t.innerSession.CommitTransaction(context.TODO())
}

func (t *transaction) Collection(collName string) mongobyc.CollectionInterface {

	target, ok := t.collectionMaps[collectionName(collName)]
	if !ok {

		target = newCollection(t.mongocli, collName)
		t.collectionMaps[collectionName(collName)] = target
	}

	return target
}

func (t *transaction) Close() error {

	for _, coll := range t.collectionMaps {
		switch target := coll.(type) {
		case *collection:
			if err := target.Close(); nil != err {
				return err
			}
		}
	}
	t.collectionMaps = map[collectionName]mongobyc.CollectionInterface{}
	return nil
}
